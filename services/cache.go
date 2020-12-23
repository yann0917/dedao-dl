package services

import (
	"encoding/gob"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/yann0917/dedao-dl/utils"
)

// Cache go-cache
var Cache *cache.Cache
var cacheDir = "./.cache/"

var _ Cacher = (*CourseCategoryList)(nil)
var _ Cacher = (*CourseList)(nil)
var _ Cacher = (*CourseInfo)(nil)
var _ Cacher = (*ArticleList)(nil)
var _ Cacher = (*ArticleInfo)(nil)
var _ Cacher = (*ArticleDetail)(nil)
var _ Cacher = (*EbookDetail)(nil)
var _ Cacher = (*CourseList)(nil)

func init() {
	Cache = cache.New(2*time.Hour, 30*time.Minute)
	gobRegister()
}

func gobRegister() {
	gob.Register(&CourseCategoryList{})
	gob.Register(&CourseList{})
	gob.Register(&CourseInfo{})
	gob.Register(&ArticleList{})
	gob.Register(&ArticleInfo{})
	gob.Register(&ArticleDetail{})
	gob.Register(&EbookDetail{})
	gob.Register(&Token{})
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

// LoadCacheFile load local file
func LoadCacheFile(fname string) (err error) {
	err = Cache.LoadFile(cacheDir + fname)
	if err != nil {
		errors.Wrap(err, "load cache file error")
		_, err = utils.Mkdir(cacheDir)
		if err != nil {
			errors.Wrap(err, "make cache dir error")
		}
		return
	}
	return
}

// SaveCacheFile save cache
func SaveCacheFile(fname string) (err error) {
	err = Cache.SaveFile(cacheDir + fname)
	if err != nil {
		errors.Wrap(err, "save cache file error")
	}
	return
}
