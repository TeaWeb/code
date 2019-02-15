package teaconfigs

// 看板
type Board struct {
	Charts []*BoardChart `yaml:"charts" json:"charts"`
}
