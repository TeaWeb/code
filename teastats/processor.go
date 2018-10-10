package teastats

import (
	"github.com/TeaWeb/code/teamongo"
	"sync"
	"github.com/TeaWeb/code/tealogs"
)

var collectionsMap = map[string]*teamongo.Collection{} // name => collection
var collectionsMutex = &sync.Mutex{}
var processors = []tealogs.Processor{
	new(DailyPVStat),
	new(HourlyPVStat),
	new(MonthlyPVStat),

	new(DailyRequestsStat),
	new(HourlyRequestsStat),
	new(MonthlyRequestsStat),

	new(DailyUVStat),
	new(HourlyUVStat),
	new(MonthlyUVStat),

	new(TopRegionStat),
	new(TopStateStat),
	new(TopOSStat),
	new(TopBrowserStat),
	new(TopRequestStat),
	new(TopCostStat),
}

type Processor struct {
}

func (this *Processor) Process(accessLog *tealogs.AccessLog) {
	for _, processor := range processors {
		processor.Process(accessLog)
	}
}

func findCollection(collectionName string, initFunc func()) *teamongo.Collection {
	collectionsMutex.Lock()
	defer collectionsMutex.Unlock()

	coll, found := collectionsMap[collectionName]
	if found {
		return coll
	}

	coll = teamongo.FindCollection(collectionName)
	collectionsMap[collectionName] = coll

	// 初始化
	if initFunc != nil {
		go initFunc()
	}

	return coll
}
