package teadb

type SQLCond struct {
	Expr   string
	Params map[string]interface{}
}
