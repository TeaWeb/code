package actions

type ActionDefinition struct {
	Name        string
	Code        ActionString
	Description string
	Instance    ActionInterface
}
