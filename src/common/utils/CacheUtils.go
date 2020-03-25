package utils

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"sync"
)

type ICacheConfig interface {
	Alias() string
}

// memory [interval means the gc time. The cache will check at each time interval, whether item has expired.]
type CacheMemoryConfig struct {
	Interval int32 `json:"interval"`
}

func (c CacheMemoryConfig) Alias() string { return "memory" }

// file
type CacheFileConfig struct {
	CachePath      string `json:"CachePath"`
	FileSuffix     string `json:"FileSuffix"`
	DirectoryLevel string `json:"DirectoryLevel"`
	EmbedExpiry    string `json:"EmbedExpiry"`
}

func (c CacheFileConfig) Alias() string { return "file" }

// redis
type CacheRedisConfig struct {
	Key      string `json:"key"`
	Conn     string `json:"conn"`
	DbNum    string `json:"dbNum"`
	Password string `json:"password"`
}

func (c CacheRedisConfig) Alias() string { return "redis" }

// memcache
type CacheMemcacheConfig struct {
	Conn string `json:"conn"`
}

func (c CacheMemcacheConfig) Alias() string { return "memcache" }

var (
	cacheOnce      sync.Once
	CacheForMemory cache.Cache
)

func init() {
	cacheOnce.Do(func() {
		var err error
		memoryConfig := CacheMemoryConfig{Interval: 30}
		CacheForMemory = cache.NewMemoryCache()
		err = CacheForMemory.StartAndGC(JsonString(memoryConfig))
		if err != nil {
			logs.Error("CacheForMemory.StartAndGC failed", err)
		}
	})
}
