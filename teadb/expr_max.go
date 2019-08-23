package teadb

type MaxExpr struct {
	Field string
}

func NewMaxExpr(field string) *MaxExpr {
	return &MaxExpr{
		Field: field,
	}
}
