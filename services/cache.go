package services

import (
	"encoding/gob"
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache go-cache
var Cache *cache.Cache
var cacheDir = "./.cache/"

func init() {
	Cache = cache.New(20*time.Minute, 10*time.Minute)
	gobRegister()
}

func gobRegister() {
	gob.Register(&CourseCategoryList{})
}

// Cacher cache interface
type Cacher interface {
	getCacheKey() string
	getCache(string) (interface{}, bool)
	setCache(string) error
}

func cacheKey(c Cacher) string {
	return c.getCacheKey()
}

func getCache(c Cacher, fileName string) (interface{}, bool) {
	return c.getCache(fileName)
}

func setCache(c Cacher, fileName string) error {
	return c.setCache(fileName)
}
