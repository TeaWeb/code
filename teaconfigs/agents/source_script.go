package agents

import (
	"bytes"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"os/exec"
)

// Script文件
type ScriptSource struct {
	Path       string           `yaml:"path" json:"path"`
	Env        []*EnvVariable   `yaml:"env" json:"env"` // 环境变量设置
	Cwd        string           `yaml:"cwd" json:"cwd"`
	DataFormat SourceDataFormat `yaml:"dataFormat" json:"dataFormat"` // 数据格式
}

// 获取新对象
func NewScriptSource() *ScriptSource {
	return &ScriptSource{
		Env: []*EnvVariable{},
	}
}

// 校验
func (this *ScriptSource) Validate() error {
	if len(this.Path) == 0 {
		return errors.New("path should not be empty")
	}

	return nil
}

// 名称
func (this *ScriptSource) Name() string {
	return "Shell脚本"
}

// 代号
func (this *ScriptSource) Code() string {
	return "script"
}

// 描述
func (this *ScriptSource) Description() string {
	return "通过执行本地的Shell脚本文件获取数据"
}

// 数据格式
func (this *ScriptSource) DataFormatCode() SourceDataFormat {
	return this.DataFormat
}

// 执行
func (this *ScriptSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.Path) == 0 {
		return nil, errors.New("path should not be empty")
	}

	cmd := exec.Command(this.Path)

	if len(this.Env) > 0 {
		for _, env := range this.Env {
			cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
		}
	}

	if len(params) > 0 {
		for key, value := range params {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}

	if len(this.Cwd) > 0 {
		cmd.Dir = this.Cwd
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		// do nothing
	}

	if stderr.Len() > 0 {
		logs.Println("error:", string(stderr.Bytes()))
	}

	return DecodeSource(stdout.Bytes(), this.DataFormat)
}

// 获取简要信息
func (this *ScriptSource) Summary() maps.Map {
	return maps.Map{
		"name":        this.Name(),
		"code":        this.Code(),
		"description": this.Description(),
	}
}

// 添加环境变量
func (this *ScriptSource) AddEnv(name, value string) {
	this.Env = append(this.Env, &EnvVariable{
		Name:  name,
		Value: value,
	})
}
