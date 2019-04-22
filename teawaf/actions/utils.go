package actions

var AllActions = []*ActionDefinition{
	{
		Name:     "阻止",
		Code:     ActionBlock,
		Instance: new(BlockAction),
	},
	{
		Name:     "允许通过",
		Code:     ActionAllow,
		Instance: new(AllowAction),
	},
	{
		Name:     "允许并记录日志",
		Code:     ActionLog,
		Instance: new(LogAction),
	},
}

func FindActionInstance(action ActionString) ActionInterface {
	for _, def := range AllActions {
		if def.Code == action {
			return def.Instance
		}
	}
	return nil
}

func FindActionName(action ActionString) string {
	for _, def := range AllActions {
		if def.Code == action {
			return def.Name
		}
	}
	return ""
}
