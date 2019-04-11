package teacache

import "github.com/iwind/TeaGo/maps"

// 所有的缓存配置
func AllCacheTypes() []maps.Map {
	return []maps.Map{
		{
			"name": "内存",
			"type": "memory",
		},
		{
			"name": "文件",
			"type": "file",
		},
		{
			"name": "Redis",
			"type": "redis",
		},
		{
			"name": "LevelDB",
			"type": "leveldb",
		},
	}
}

// 类型名称
func TypeName(typeCode string) string {
	for _, m := range AllCacheTypes() {
		if m.GetString("type") == typeCode {
			return m.GetString("name")
		}
	}
	return ""
}
