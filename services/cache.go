package services

import (
	"encoding/gob"
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache go-cache
var Cache *cache.Cache

func init() {
	Cache = cache.New(10*time.Minute, 5*time.Minute)
	gobRegister()
}

func gobRegister() {
	gob.Register(&CourseCourseCategoryList{})
}
