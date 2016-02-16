package util

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/url"
	"os"
	"time"
)

var (
	imgCache = NewSet()
)

func init() {
	go run()
}

func WriteCache(imgPath, imgArg string, img image.Image) {
	cacheName := genCacheName(imgPath, imgArg)
	cacheFile, err := os.Create(CacheRoot + cacheName)
	if err != nil {
		fmt.Printf("WriteCache err:%v", err)
		return
	}
	defer cacheFile.Close()
	jpeg.Encode(cacheFile, img, nil)
	imgCache.Add(cacheName)
}

func FindInCache(imgPath, imgArg string) ([]byte, error) {
	cacheName := genCacheName(imgPath, imgArg)
	if !imgCache.Contains(cacheName) {
		return nil, nil
	}

	return ioutil.ReadFile(CacheRoot + cacheName)
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

func genCacheName(imgPath, imgArg string) string {
	return fmt.Sprintf("%s", url.QueryEscape(imgPath+imgArg))
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
