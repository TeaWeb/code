package scheduling

// 候选对象接口
type CandidateInterface interface {
	// 权重
	CandidateWeight() uint
	CandidateCodes() []string
}
