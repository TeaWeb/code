package agents

import (
	"bytes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"os/exec"
	"runtime"
	"strings"
)

// Script文件数据源
type ScriptSource struct {
	Path       string           `yaml:"path" json:"path"`
	ScriptType string           `yaml:"scriptType" json:"scriptType"` // 脚本类型，可以为path, code
	Script     string           `yaml:"script" json:"script"`         // 脚本代码
	Env        []*EnvVariable   `yaml:"env" json:"env"`               // 环境变量设置
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

// 格式化脚本
func (this *ScriptSource) FormattedScript() string {
	script := this.Script
	script = strings.Replace(script, "\r", "", -1)
	return script
}

// 保存到本地
func (this *ScriptSource) Generate(id string) (path string, err error) {
	if runtime.GOOS == "windows" {
		path = Tea.ConfigFile("agents/source." + id + ".bat")
	} else {
		path = Tea.ConfigFile("agents/source." + id + ".script")
	}
	shFile := files.NewFile(path)
	if !shFile.Exists() {
		err = shFile.WriteString(this.FormattedScript())
		if err != nil {
			return
		}
		err = shFile.Chmod(0777)
		if err != nil {
			return
		}
	}
	return
}

// 执行
func (this *ScriptSource) Execute(params map[string]string) (value interface{}, err error) {
	// 脚本
	if this.ScriptType == "code" {
		path, err := this.Generate(stringutil.Rand(16))
		if err != nil {
			return nil, err
		}
		this.Path = path

		defer func() {
			err := files.NewFile(this.Path).Delete()
			if err != nil {
				logs.Error(err)
			}
		}()
	}

	if len(this.Path) == 0 {
		return nil, errors.New("path or script should not be empty")
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
