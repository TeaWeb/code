package actions

type ActionDefinition struct {
	Name     string
	Code     ActionString
	Instance ActionInterface
}
