package util

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"time"
)

const (
	CachePath = "./cache/"
)

var (
	imgCache = NewSet()
)

func init() {
	go run()
}

// xxx.jpg.100x200
func loadCache() {
	newImgCache := NewSet()
	files, _ := ioutil.ReadDir(CachePath)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		newImgCache.Add(file.Name())
	}
	imgCache = newImgCache
}

func WriteCache(imgName, imgArg string, img image.Image) {
	cacheName := genCacheName(imgName, imgArg)
	cacheFile, err := os.Create(CachePath + cacheName)
	if err != nil {
		fmt.Printf("WriteCache err:%v", err)
		return
	}
	defer cacheFile.Close()
	jpeg.Encode(cacheFile, img, nil)
	imgCache.Add(cacheName)
}

func genCacheName(imgName, imgArg string) string {
	return fmt.Sprintf("%s.%s", imgName, imgArg)
}

func FindInCache(imgName, imgArg string) image.Image {
	cacheName := genCacheName(imgName, imgArg)
	if !imgCache.Contains(cacheName) {
		return nil
	}

	img, err := LoadImage(CachePath + cacheName)
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
