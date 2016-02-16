package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"net/http"
	"util"

	"github.com/gin-gonic/gin"
)

// http://xxx.com/160214/22.jpg?s=100x200
func imageHandler(context *gin.Context) {
	imgPath := context.Param("path")
	size := context.Query("s")

	cacheBuff := util.FindInCache(imgPath, size)
	if len(cacheBuff) > 0 {
		// 用状态码201表示当前从缓存中读取的数据,便于日志直接查看
		context.Data(http.StatusCreated, "image/jpeg", cacheBuff)
		return
	}

	// 无size指定，默认为原图大小
	if size == "" {
		go rspOriginImg(imgPath, context)
	} else {
		go rspThumbnailImg(imgPath, size, context)
	}
	return
}

func rspOriginImg(imgPath string, context *gin.Context) {
	imgBuff, err := util.LoadFile(imgPath)
	if err != nil {
		fmt.Printf("[GIN] LoadFile error:%v\n", err)
		context.String(http.StatusNoContent, "LoadFile error:%v", err)
	} else {
		context.Data(http.StatusOK, "image/jpeg", imgBuff)
	}
	return
}

func rspThumbnailImg(imgPath, size string, context *gin.Context) {
	dstWidth, dstHeight := util.ParseImgArg(size)
	if dstHeight == 0 || dstWidth == 0 {
		context.String(http.StatusForbidden, "size forbidden")
		return
	}

	srcImg, err := util.LoadImage(imgPath)
	if err != nil {
		fmt.Printf("[GIN] LoadImage error:%v\n", err)
		context.String(http.StatusNoContent, "LoadImage error:%v", err)
		return
	}

	thumbImg := util.Thumbnail(dstWidth, dstHeight, srcImg)
	dstImg := util.CropImg(thumbImg, int(dstWidth), int(dstHeight))
	go util.WriteCache(imgPath, size, dstImg)

	buff := &bytes.Buffer{}
	jpeg.Encode(buff, dstImg, nil)
	context.Data(http.StatusOK, "image/jpeg", buff.Bytes())
	return
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/*path", imageHandler)
	router.Run(":6789")
}
