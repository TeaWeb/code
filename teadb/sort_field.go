package teadb

type SortType = string

const (
	SortAsc  = "asc"
	SortDesc = "desc"
)

type SortField struct {
	Name string
	Type SortType
}

func (this *SortField) IsAsc() bool {
	return this.Type == SortAsc
}

func (this *SortField) IsDesc() bool {
	return this.Type == SortDesc
}
