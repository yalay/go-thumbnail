package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strings"
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
		context.String(http.StatusNotFound, "path error")
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
		fmt.Printf("[GIN] LoadFile error:%v\n", err)
		context.String(http.StatusNotFound, "LoadFile error:%v", err)
	} else {
		context.Data(http.StatusOK, "image/jpeg", imgBuff)
	}

	return
}

// 只有指定Refer才允许crop
func rspThumbnailImg(imgPath, size string, context *gin.Context) {
	cacheBuff := util.FindInCache(imgPath, size)
	if len(cacheBuff) > 0 {
		context.Data(http.StatusOK, "image/jpeg", cacheBuff)
		return
	}

	var thumbImg image.Image
	referUrl := context.Request.Referer()
	if !strings.Contains(referUrl, util.AllowedRefer) {
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
	context.Data(http.StatusOK, "image/jpeg", buff.Bytes())

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
		fmt.Printf("[GIN] LoadImage error:%v\n", err)
		return nil
	}

	if doCrop {
		return util.ThumbnailCrop(dstWidth, dstHeight, srcImg)
	} else {
		return util.ThumbnailSimple(dstWidth, dstHeight, srcImg)
	}

}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/*path", imageHandler)
	router.Run(":" + util.ServePort)
}
