package teadb

type MinExpr struct {
	Field string
}

func NewMinExpr(field string) *MinExpr {
	return &MinExpr{
		Field: field,
	}
}
