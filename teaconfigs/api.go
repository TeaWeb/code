package teaconfigs

//  API定义
type API struct {
	Path           string          `yaml:"path" json:"path"`                     // TODO
	Address        string          `yaml:"address" json:"address"`               // TODO
	Methods        []string        `yaml:"methods" json:"methods"`               // TODO
	Params         []*APIParam     `yaml:"params" json:"params"`                 // TODO
	Name           string          `yaml:"name" json:"name"`                     // TODO
	Description    string          `yaml:"description" json:"description"`       // TODO
	Mock           []string        `yaml:"mock" json:"mock"`                     // TODO
	Author         string          `yaml:"author" json:"author"`                 // TODO
	Company        string          `yaml:"company" json:"company"`               // TODO
	IsAsynchronous bool            `yaml:"isAsynchronous" json:"isAsynchronous"` // TODO
	Timeout        float64         `yaml:"timeout" json:"timeout"`               // TODO
	MaxSize        uint            `yaml:"maxSize" json:"maxSize"`               // TODO
	Headers        []*HeaderConfig `yaml:"headers" json:"headers"`               // TODO
	TodoThings     []string        `yaml:"todo" json:"todo"`                     // TODO
	DoneThings     []string        `yaml:"done" json:"done"`                     // TODO
	Response       []byte          `yaml:"response" json:"response"`             // TODO
	Roles          []string        `yaml:"roles" json:"roles"`                   // TODO
	IsDeprecated   bool            `yaml:"isDeprecated" json:"isDeprecated"`     // TODO
	On             bool            `yaml:"on" json:"on"`                         // TODO
	Versions       []string        `yaml:"versions" json:"versions"`             // TODO
	ModifiedAt     int64           `yaml:"modifiedAt" json:"modifiedAt"`         // TODO
	Username       string          `yaml:"username" json:"username"`             // 添加API的用户名 TODO
}

func NewAPI() *API {
	return &API{}
}

func (this *API) Validate() error {
	return nil
}

func (this *API) AddParam(param *APIParam) {
	this.Params = append(this.Params, param)
}
