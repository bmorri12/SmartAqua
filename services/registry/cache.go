package main

import (
	"flag"
	"github.com/bmorri12/SmartAqua/pkg/cache"
)

const (
	flagCacheSize    = "cacheSize"
	defaultCacheSize = 102400
)

var (
	confCacheSize = flag.Int(flagCacheSize, defaultCacheSize, "maximum size of cache")
)

var MemCache cache.Cache

func getCache() cache.Cache {
	if MemCache == nil {
		MemCache = cache.NewMemCache(*confCacheSize)
	}
	return MemCache
}
