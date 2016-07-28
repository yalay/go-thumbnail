package util

import (
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"time"
)

var (
	imgCache = NewSet()
)

func init() {
	go run()
}

func WriteCache(imgUrl string, img image.Image) {
	cacheName := genCacheName(imgUrl)
	cacheFile, err := os.Create(CacheRoot + cacheName)
	if err != nil {
		Logln("WriteCache err:" + err.Error())
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

	cacheBuff, _ := ioutil.ReadFile(CacheRoot + cacheName)
	return cacheBuff
}

// %2FPure%2F22.jpg100x100
func loadCache() {
	newImgCache := NewSet()
	if _, err := os.Stat(CacheRoot); err != nil {
		os.MkdirAll(CacheRoot, os.ModePerm)
	}

	files, _ := ioutil.ReadDir(CacheRoot)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		newImgCache.Add(file.Name())
	}
	imgCache = newImgCache
}

func genCacheName(imgUrl string) string {
	return Md5Sum(imgUrl)
}

func run() {
	loadCache()
	timer := time.NewTimer(time.Hour)
	for {
		select {
		case <-timer.C:
			loadCache()
			timer.Reset(time.Hour)
		}
	}
}
