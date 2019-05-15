package teahooks

var eventFunctions = map[Event][]Handler{}

type Handler func()

// add handlers
func On(event Event, f Handler) {
	locker.Lock()
	defer locker.Unlock()

	funcList, ok := eventFunctions[event]
	if ok {
		funcList = append(funcList, f)
	} else {
		funcList = []Handler{f}
	}
	eventFunctions[event] = funcList
}

// call handlers
func Call(event Event) {
	locker.Lock()
	defer locker.Unlock()

	funcList, ok := eventFunctions[event]
	if !ok {
		return
	}
	for _, handler := range funcList {
		handler()
	}
}
