package util

import (
	"fmt"
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

func getCacheImg(imgName string) (img image.Image, err error) {
	file, err := os.Open(CacheRoot + imgName)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

// xxx.jpg.pure.100x200
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

func WriteCache(imgName, category, imgArg string, img image.Image) {
	cacheName := genCacheName(imgName, category, imgArg)
	cacheFile, err := os.Create(CacheRoot + cacheName)
	if err != nil {
		fmt.Printf("WriteCache err:%v", err)
		return
	}
	defer cacheFile.Close()
	jpeg.Encode(cacheFile, img, nil)
	imgCache.Add(cacheName)
}

func genCacheName(imgName, category, imgArg string) string {
	return fmt.Sprintf("%s.%s.%s", imgName, category, imgArg)
}

func FindInCache(imgName, category, imgArg string) image.Image {
	cacheName := genCacheName(imgName, category, imgArg)
	if !imgCache.Contains(cacheName) {
		return nil
	}

	img, err := getCacheImg(cacheName)
	if err != nil {
		return nil
	}
	return img
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
