package util

import (
	"common"
	"conf"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

var imgCache *cache

type cache struct {
	sync.RWMutex
	imgNames common.Set
}

func init() {
	imgCache = &cache{
		imgNames: common.NewSet(),
	}
	loadCache()
}

func (c *cache) Add(imgName string) {
	c.Lock()
	c.imgNames.Add(imgName)
	c.Unlock()
}

func (c *cache) Contains(imgName string) bool {
	c.RLock()
	isExist := c.imgNames.Contains(imgName)
	c.RUnlock()
	return isExist
}

func WriteCache(imgUrl string, img image.Image) {
	cacheName := genCacheName(imgUrl)
	cacheDir := conf.GetCacheDir()
	cacheFile, err := os.Create(path.Join(cacheDir, cacheName))
	if err != nil {
		log.Println("WriteCache err:" + err.Error())
		return
	}
	defer cacheFile.Close()
	jpeg.Encode(cacheFile, img, nil)
	imgCache.Add(cacheName)
}

func FindInCache(imgUrl string) []byte {
	cacheName := genCacheName(imgUrl)
	if !imgCache.Contains(cacheName) {
		return nil
	}
	cacheDir := conf.GetCacheDir()
	cacheBuff, _ := ioutil.ReadFile(path.Join(cacheDir, cacheName))
	return cacheBuff
}

// %2FPure%2F22.jpg100x100
func loadCache() {
	cacheDir := conf.GetCacheDir()
	if _, err := os.Stat(cacheDir); err != nil {
		os.MkdirAll(cacheDir, os.ModePerm)
	}

	files, _ := ioutil.ReadDir(cacheDir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		imgCache.Add(file.Name())
	}
}

func genCacheName(imgUrl string) string {
	return Md5Sum(imgUrl)
}
