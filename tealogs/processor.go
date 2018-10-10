package tealogs

type Processor interface {
	Process(accessLog *AccessLog)
}
