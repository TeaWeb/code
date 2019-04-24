package checkpoints

type Checkpoint struct {
}

func (this *Checkpoint) IsRequest() bool {
	return true
}

func (this *Checkpoint) ParamOptions() *ParamOptions {
	return nil
}
