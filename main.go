package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/http"
	"strings"
	"time"
	"util"

	"github.com/gin-gonic/gin"
)

// http://127.0.0.1:6789/pure/22.jpg
func imageHandler(context *gin.Context) {
	imgPath := context.Param("path")
	size := context.Query("s")

	if !strings.HasSuffix(imgPath, "jpg") &&
		!strings.HasSuffix(imgPath, "jpeg") &&
		!strings.HasSuffix(imgPath, "png") {
		context.Status(http.StatusNotFound)
		return
	}

	if size != "" {
		rspThumbnailImg(imgPath, size, context)
	} else {
		rspOriginImg(imgPath, context)
	}
	return
}

func rspOriginImg(imgPath string, context *gin.Context) {
	imgBuff, err := util.LoadFile(imgPath)
	if err != nil {
		util.Log("[GIN] LoadFile error:" + err.Error())
		context.Status(http.StatusNotFound)
	} else {
		rspCacheControl(imgBuff, context)
	}

	return
}

// 只有指定Refer才允许crop
func rspThumbnailImg(imgPath, size string, context *gin.Context) {
	cacheBuff := util.FindInCache(imgPath, size)
	if len(cacheBuff) > 0 {
		rspCacheControl(cacheBuff, context)
		return
	}

	var thumbImg image.Image
	referUrl := context.Request.Referer()
	if strings.Contains(referUrl, util.AllowedRefer) {
		thumbImg = getThumbnailImg(imgPath, size, true)
	} else {
		thumbImg = getThumbnailImg(imgPath, size, false)
	}
	if thumbImg == nil {
		context.String(http.StatusNotFound, "Thumbnail fail:%s-%s", imgPath, size)
		return
	}

	buff := &bytes.Buffer{}
	jpeg.Encode(buff, thumbImg, nil)
	rspCacheControl(buff.Bytes(), context)

	go util.WriteCache(imgPath, size, thumbImg)
	return
}

func getThumbnailImg(imgPath, size string, doCrop bool) image.Image {
	dstWidth, dstHeight := util.ParseImgArg(size)
	if dstHeight == 0 || dstWidth == 0 {
		return nil
	}

	srcImg, err := util.LoadImage(imgPath)
	if err != nil {
		util.Log("[GIN] LoadImage error:" + err.Error())
		return nil
	}

	if doCrop {
		return util.ThumbnailCrop(dstWidth, dstHeight, srcImg)
	} else {
		return util.ThumbnailSimple(dstWidth, dstHeight, srcImg)
	}
}

func rspCacheControl(data []byte, context *gin.Context) {
	eTag := string(util.Md5Sum(data))
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
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.LoggerWithWriter(util.GetLogBuf()), gin.Recovery())
	router.GET("/*path", imageHandler)
	router.Run(":" + util.ServePort)
}
