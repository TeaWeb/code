package checkpoints

// check point definition
type CheckPointDefinition struct {
	Name        string
	Description string
	Prefix      string
	HasParams   bool // has sub params
	Instance    CheckPointInterface
}
