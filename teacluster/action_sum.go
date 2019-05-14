package teacluster

import (
	"github.com/iwind/TeaGo/logs"
)

// cluster -> master|node
type SumAction struct {
	Action
}

func (this *SumAction) Name() string {
	return "sum"
}

func (this *SumAction) OnSuccess(success *SuccessAction) error {
	logs.Println("sum:", success.Data)

	// write to local file
	//file := files.NewFile(Tea.ConfigFile("cluster.sum"))
	//file.WriteString()

	return nil
}

func (this *SumAction) OnFail(fail *FailAction) error {
	// TODO retry later
	return nil
}

func (this *SumAction) TypeId() int8 {
	return 9
}
