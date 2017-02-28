package main

import (
	"bytes"
	"conf"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"

	"github.com/gin-gonic/gin"
)

// http://127.0.0.1:6789/pure/22.jpg
func imageHandler(context *gin.Context) {
	imgPath := context.Param("path")
	size := context.Query("s")

	if doSkip(imgPath, context) {
		context.Status(http.StatusOK)
		return
	}

	if size != "" {
		imgPath = imgPath + "?s=" + size
	} else {
		if conf.IsSpider(context.Request.UserAgent()) {
			imgPath = imgPath + "?s=" + conf.GetExtImgSize()
		} else {
			if !conf.IsAllowedRefer(context.Request.Referer()) {
				if util.DoAd(context) {
					adImgPath := conf.GetRandomAdPath()
					if adImgPath != "" {
						imgPath = adImgPath
					}
				} else {
					imgPath = imgPath + "?s=" + conf.GetExtImgSize()
				}
			}
		}
	}

	rspImg(imgPath, context)
	return
}

// 原图本地获取或者跳转到img服务器
func rspImg(imgPath string, context *gin.Context) {
	imgUrl, err := url.Parse(imgPath)
	if err != nil || imgUrl == nil {
		context.Status(http.StatusNotFound)
		return
	}

	if len(imgUrl.Host) > 0 {
		context.Redirect(http.StatusFound, imgPath)
		return
	}

	// cache
	cacheBuff := util.FindInCache(imgUrl.String())
	if len(cacheBuff) > 0 {
		rspCacheControl(cacheBuff, context)
		return
	}

	buff := &bytes.Buffer{}
	thumbImg := getThumbnailImg(imgUrl)
	if thumbImg == nil {
		context.Status(http.StatusNotFound)
		return
	}

	jpeg.Encode(buff, thumbImg, nil)
	rspCacheControl(buff.Bytes(), context)
}

func getThumbnailImg(imgUrl *url.URL) image.Image {
	srcImg, err := util.LoadImage(imgUrl.Path)
	if err != nil {
		log.Println("[GIN] LoadImage error:" + err.Error())
		return nil
	}

	imgValues, _ := url.ParseQuery(imgUrl.RawQuery)
	size := imgValues.Get("s")
	if size == "" {
		return srcImg
	} else {
		dstWidth, dstHeight := util.ParseImgArg(size)
		if dstHeight == 0 && dstWidth == 0 {
			return srcImg
		}

		var thumbImg image.Image
		if dstHeight == 0 || dstWidth == 0 {
			thumbImg = util.ThumbnailSimple(dstWidth, dstHeight, srcImg)
		} else {
			thumbImg = util.ThumbnailCrop(dstWidth, dstHeight, srcImg)
		}

		go util.WriteCache(imgUrl.String(), thumbImg)
		return thumbImg
	}
}

func doSkip(imgPath string, context *gin.Context) bool {
	// 忽略favicon
	if strings.HasSuffix(imgPath, "favicon.ico") {
		return true
	}
	return false
}

func rspCacheControl(data []byte, context *gin.Context) {
	eTag := util.Md5Sum(string(data))
	reqTag := context.Request.Header.Get("If-None-Match")
	if reqTag != "" && reqTag == eTag {
		context.Header("ETag", eTag)
		context.Status(http.StatusNotModified)
	} else {
		cacheSince := time.Now().Format(http.TimeFormat)
		cacheUntil := time.Now().AddDate(0, 0, 1).Format(http.TimeFormat)

		context.Header("ETag", eTag)
		context.Header("Cache-Control", "max-age=86400")
		context.Header("Last-Modified", cacheSince)
		context.Header("Expires", cacheUntil)
		context.Data(http.StatusOK, "image/jpeg", data)
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover:%v\n", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(util.Counter(), gin.Logger(), gin.Recovery())
	router.GET("/*path", imageHandler)
	router.Run(":" + strconv.Itoa(conf.GetListenPort()))
}
