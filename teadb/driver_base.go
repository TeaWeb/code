package teadb

import "errors"

var ErrorDBUnavailable = errors.New("database is not available")

type BaseDriver struct {
	isAvailable bool
}

func (this *BaseDriver) IsAvailable() bool {
	return this.isAvailable
}

func (this *BaseDriver) SetIsAvailable(b bool) {
	this.isAvailable = b
}
