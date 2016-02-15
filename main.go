package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"util"

	"github.com/gin-gonic/gin"
)

// http://xxx.com/160214/22.jpg?s=100x200
func imageHandler(context *gin.Context) {
	imgPath := context.Param("path")
	size := context.Query("s")

	cacheImg := util.FindInCache(imgPath, size)
	if cacheImg != nil {
		rspImgWriter(cacheImg, context)
		return
	}

	srcImg, err := util.LoadImage(imgPath)
	if err != nil {
		fmt.Printf("[GIN] LoadImage error:%v\n", err)
		context.String(http.StatusNoContent, "LoadImage error:%v", err)
		return
	}

	// 无size指定，默认为原图大小
	var dstImg image.Image
	if size == "" {
		dstImg = srcImg
	} else {
		dstWidth, dstHeight := util.ParseImgArg(size)
		if dstHeight == 0 || dstWidth == 0 {
			context.String(http.StatusForbidden, "size forbidden")
			return
		}

		thumbImg := util.Thumbnail(dstWidth, dstHeight, srcImg)
		dstImg = util.CropImg(thumbImg, int(dstWidth), int(dstHeight))
		go util.WriteCache(imgPath, size, dstImg)
	}

	rspImgWriter(dstImg, context)
}

func rspImgWriter(dstImg image.Image, context *gin.Context) {
	buff := &bytes.Buffer{}
	jpeg.Encode(buff, dstImg, nil)
	context.Data(http.StatusOK, "image/jpeg", buff.Bytes())
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/*path", imageHandler)
	router.Run(":6789")
}
