package logbuffer

import (
	"bufio"
	"os"
	"sync/atomic"
)

type File struct {
	name   string
	writer *os.File
	reader *os.File
	buffer *bufio.Reader
	size   int

	isWriting int32
}

func NewFile(name string) *File {
	return &File{
		name: name,
	}
}

func (this *File) Write(data []byte) (int, error) {
	atomic.StoreInt32(&this.isWriting, 1)
	if this.writer == nil {
		writer, err := os.OpenFile(this.name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			atomic.StoreInt32(&this.isWriting, 0)
			return 0, err
		}
		this.writer = writer
	}
	n, err := this.writer.Write(data)
	this.size += n

	n1, _ := this.writer.Write([]byte{'\n'})
	this.size += n1
	atomic.StoreInt32(&this.isWriting, 0)
	return n + n1, err
}

func (this *File) Read() (data []byte, err error) {
	v := atomic.LoadInt32(&this.isWriting)
	if v == 1 {
		return
	}
	if this.reader == nil {
		reader, err := os.OpenFile(this.name, os.O_RDONLY, 0777)
		if err != nil {
			return nil, err
		}
		this.reader = reader
		this.buffer = bufio.NewReader(reader)
	}
	line, err := this.buffer.ReadSlice('\n')
	if err == nil {
		if len(line) > 0 {
			line = line[:len(line)-1]
		}
	} else if err == bufio.ErrBufferFull {
		newLine := append([]byte{}, line...)
		for {
			line2, err := this.buffer.ReadSlice('\n')
			newLine = append(newLine, line2...)
			if err == bufio.ErrBufferFull {
				continue
			}
			break
		}
		line = newLine[:len(newLine)-1]
		err = nil
	}
	return line, err
}

func (this *File) Sync() error {
	return this.writer.Sync()
}

func (this *File) Size() int {
	return this.size
}

func (this *File) Close() error {
	var err error = nil
	if this.writer != nil {
		err1 := this.writer.Close()
		if err1 != nil {
			err = err1
		}
	}

	if this.reader != nil {
		err1 := this.reader.Close()
		if err1 != nil {
			err = err1
		}
	}
	return err
}

func (this *File) Delete() error {
	return os.Remove(this.name)
}
