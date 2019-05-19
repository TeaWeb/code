package teautils

type MemoryGridCompressOpt struct {
	Level int
}

func NewMemoryGridCompressOpt(level int) *MemoryGridCompressOpt {
	return &MemoryGridCompressOpt{
		Level: level,
	}
}
