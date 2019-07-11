package tealogs

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/logs"
	"sync"
)

// 存储策略相关
var storageMap = map[string]StorageInterface{} // policyId => StorageInterface
var storageNamesMap = map[string]string{}      // policyId => policy name
var storageLocker = sync.Mutex{}

// 通过策略ID查找存储
func FindPolicyStorage(policyId string) StorageInterface {
	storageLocker.Lock()
	defer storageLocker.Unlock()

	storage, ok := storageMap[policyId]
	if ok {
		return storage
	}

	policy := teaconfigs.NewAccessLogStoragePolicyFromId(policyId)
	if policy == nil || !policy.On {
		storageMap[policyId] = nil
		return nil
	}

	storageNamesMap[policyId] = policy.Name

	storage = DecodePolicyStorage(policy)
	if storage != nil {
		err := storage.Start()
		if err != nil {
			logs.Println("access log storage '"+policyId+"/"+FindPolicyName(policyId)+"' start failed:", err.Error())
			storage = nil
		}
	}
	storageMap[policyId] = storage

	return storage
}

// 清除策略相关信息
func ResetPolicyStorage(policyId string) {
	storageLocker.Lock()
	storage, ok := storageMap[policyId]
	if ok {
		delete(storageMap, policyId)
		delete(storageNamesMap, policyId)
	}
	storageLocker.Unlock()

	if storage != nil {
		storage.Close()
	}
}

// 解析策略中的存储对象
func DecodePolicyStorage(policy *teaconfigs.AccessLogStoragePolicy) StorageInterface {
	if policy == nil {
		return nil
	}

	var instance StorageInterface = nil
	switch policy.Type {
	case StorageTypeFile:
		instance = new(FileStorage)
	case StorageTypeES:
		instance = new(ESStorage)
	case StorageTypeMySQL:
		instance = new(MySQLStorage)
	case StorageTypeTCP:
		instance = new(TCPStorage)
	case StorageTypeCommand:
		instance = new(CommandStorage)
	}
	if instance == nil {
		return nil
	}

	teautils.MapToObjectJSON(policy.Options, instance)

	return instance
}

// 查找策略名称
func FindPolicyName(policyId string) string {
	storageLocker.Lock()
	name, _ := storageNamesMap[policyId]
	storageLocker.Unlock()

	return name
}
