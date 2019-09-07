package teautils

import (
	"bytes"
	"io"
	"os"
	"sync"
)

// 文件Buffer
type FileBuffer struct {
	file string

	writer *os.File
	reader *os.File

	locker sync.Mutex
	buf    []byte

	maxFileSize     int
	currentFileSize int
}

// 获得FileBuffer对象
func NewFileBuffer(file string) *FileBuffer {
	return &FileBuffer{
		buf:  make([]byte, 1024),
		file: file,
	}
}

// 打开文件
func (this *FileBuffer) Open() error {
	writer, err := os.OpenFile(this.file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	this.writer = writer

	reader, err := os.OpenFile(this.file, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	this.reader = reader

	return nil
}

// 写入数据
func (this *FileBuffer) Write(data []byte) error {
	this.locker.Lock()
	defer this.locker.Unlock()

	if this.maxFileSize > 0 && this.currentFileSize > this.maxFileSize {
		_, err := this.reader.Seek(0, 0)
		if err != nil {
			return err
		}
		err = this.writer.Close()
		if err != nil {
			return err
		}
		this.writer, err = os.OpenFile(this.file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return err
		}
		this.currentFileSize = 0
	}

	n, err := this.writer.Write(data)
	if err != nil {
		return err
	}
	this.currentFileSize += n

	err = this.writer.Sync()
	if err != nil {
		return err
	}

	n, err = this.writer.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	this.currentFileSize += n

	err = this.writer.Sync()
	if err != nil {
		return err
	}

	return nil
}

// 读取数据
func (this *FileBuffer) Read() (data []byte, err error) {
	this.locker.Lock()
	defer this.locker.Unlock()

	for {
		n, err := this.reader.Read(this.buf)
		if n > 0 {
			b := this.buf[:n]
			nIndex := bytes.IndexByte(b, '\n')
			if nIndex > -1 {
				data = append(data, b[:nIndex]...)

				// 往前移
				_, err := this.reader.Seek(int64(nIndex-n)+1, 1)
				if err != nil {
					return nil, err
				}
				break
			} else {
				data = append(data, b...)
			}
		}
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, err
		}
	}

	return
}
