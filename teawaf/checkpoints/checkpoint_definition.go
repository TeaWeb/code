package checkpoints

// check point definition
type CheckPointDefinition struct {
	Name     string
	Prefix   string
	Instance CheckPointInterface
}
