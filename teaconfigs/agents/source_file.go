package agents

import (
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// 数据文件
type FileSource struct {
	Path       string           `yaml:"path" json:"path"`
	DataFormat SourceDataFormat `yaml:"dataFormat" json:"dataFormat"` // 数据格式
}

// 获取新对象
func NewFileSource() *FileSource {
	return &FileSource{}
}

// 校验
func (this *FileSource) Validate() error {
	if len(this.Path) == 0 {
		return errors.New("path should not be empty")
	}

	return nil
}

// 名称
func (this *FileSource) Name() string {
	return "数据文件"
}

// 代号
func (this *FileSource) Code() string {
	return "file"
}

// 描述
func (this *FileSource) Description() string {
	return "通过读取本地文件获取数据"
}

// 数据格式
func (this *FileSource) DataFormatCode() SourceDataFormat {
	return this.DataFormat
}

// 执行
func (this *FileSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.Path) == 0 {
		return nil, errors.New("path should not be empty")
	}

	file := files.NewFile(this.Path)
	if !file.Exists() {
		return nil, errors.New("file does not exist")
	}

	data, err := file.ReadAll()
	if err != nil {
		return nil, err
	}
	return DecodeSource(data, this.DataFormat)
}

// 获取简要信息
func (this *FileSource) Summary() maps.Map {
	return maps.Map{
		"name":        this.Name(),
		"code":        this.Code(),
		"description": this.Description(),
	}
}
