package teadb

type SQLAction = int

const (
	SQLInsert SQLAction = 1
	SQLSelect SQLAction = 2
	SQLDelete SQLAction = 3
	SQLUpdate SQLAction = 4
)
