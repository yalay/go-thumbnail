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

// http://xxx.com/t?n=22.jpg&s=100x200&c=pure
func imageHandler(context *gin.Context) {
	imgName := context.Query("n")
	size := context.Query("s")
	category := context.Query("c")

	cacheImg := util.FindInCache(imgName, category, size)
	if cacheImg != nil {
		rspImgWriter(cacheImg, context)
		return
	}

	srcImg, err := util.LoadImage(imgName, category)
	if err != nil {
		fmt.Printf("[GIN] LoadImage error:%v\n", err)
		context.String(http.StatusForbidden, "LoadImage error:%v", err)
		return
	}

	dstWidth, dstHeight := util.ParseImgArg(size)
	var dstImg image.Image
	if dstHeight == 0 || dstWidth == 0 {
		dstImg = srcImg
	} else {
		thumbImg := util.Thumbnail(dstWidth, dstHeight, srcImg)
		dstImg = util.CropImg(thumbImg, int(dstWidth), int(dstHeight))
		go util.WriteCache(imgName, category, size, dstImg)
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
	router.GET("/t", imageHandler)
	router.Run(":6789")
}
