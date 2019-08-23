package teadb

type SumExpr struct {
	Field string
}

func NewSumExpr(field string) *SumExpr {
	return &SumExpr{
		Field: field,
	}
}
