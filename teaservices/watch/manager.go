package watch

import (
	"sync"
	"github.com/iwind/TeaGo/utils/string"
)

var taskMap = map[string]*Task{} // id => *task
var tasksLocker = sync.Mutex{}

func Add(task *Task) {
	tasksLocker.Lock()
	defer tasksLocker.Unlock()

	if len(task.Id) == 0 {
		task.Id = stringutil.Rand(16)
	}

	taskMap[task.Id] = task
}
