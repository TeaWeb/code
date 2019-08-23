package teadb

type AvgExpr struct {
	Field string
}

func NewAvgExpr(field string) *AvgExpr {
	return &AvgExpr{
		Field: field,
	}
}
