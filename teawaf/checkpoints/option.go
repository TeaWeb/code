package checkpoints

// attach option
type Option struct {
	Name        string
	Code        string
	Value       string // default value
	IsRequired  bool
	Size        int
	Comment     string
	Placeholder string
	RightLabel  string
	MaxLength   int
	Validate    func(value string) (ok bool, message string)
}

func NewOption(name string, code string) *Option {
	return &Option{
		Name: name,
		Code: code,
	}
}
