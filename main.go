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

	if doSkip(imgPath, context) {
		context.String(http.StatusOK, "skip")
		return
	}

	if size != "" {
		rspThumbnailImg(imgPath, size, context)
		return
	}

	rspOriginImg(imgPath, context)
	return
}

// 原图获取跳转到img服务器
func rspOriginImg(imgPath string, context *gin.Context) {
	referUrl := context.Request.Referer()
	if !strings.Contains(referUrl, util.AllowedRefer) {
		imgPath = imgPath + "?s="+ util.ExtImgSize
	}
	context.Redirect(http.StatusFound, util.RedirectUrl+imgPath)
	return
}

func rspThumbnailImg(imgPath, size string, context *gin.Context) {
	cacheBuff := util.FindInCache(imgPath, size)
	if len(cacheBuff) > 0 {
		context.Data(http.StatusOK, "image/jpeg", cacheBuff)
		return
	}

	thumbImg := getThumbnailImg(imgPath, size)
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

func rspWaterMarkImg(imgPath string, context *gin.Context) {
	cacheBuff := util.FindInCache(imgPath, util.WaterSize)
	if len(cacheBuff) > 0 {
		context.Data(http.StatusOK, "image/jpeg", cacheBuff)
		return
	}

	thumbImg := getThumbnailImg(imgPath, util.ExtImgSize)
	if thumbImg == nil {
		context.String(http.StatusNotFound, "Warter thumbnail fail:%s", imgPath)
		return
	}

	waterImg, err := util.WaterMark(thumbImg)
	if err != nil {
		context.String(http.StatusNotFound, "Water mark error:%v", err)
	} else {
		waterBuff := &bytes.Buffer{}
		jpeg.Encode(waterBuff, waterImg, nil)
		context.Data(http.StatusOK, "image/jpeg", waterBuff.Bytes())

		go util.WriteCache(imgPath, util.WaterSize, waterImg)
	}
}

func getThumbnailImg(imgPath, size string) image.Image {
	dstWidth, dstHeight := util.ParseImgArg(size)
	if dstHeight == 0 || dstWidth == 0 {
		return nil
	}

	srcImg, err := util.LoadImage(imgPath)
	if err != nil {
		fmt.Printf("[GIN] LoadImage error:%v\n", err)
		return nil
	}

	return util.ThumbnailCrop(dstWidth, dstHeight, srcImg)
}

func doSkip(imgPath string, context *gin.Context) bool {
	// 忽略favicon
	if strings.HasSuffix(imgPath, "favicon.ico") {
		return true
	}
	req := context.Request
	if req == nil {
		return true
	}

	ua := req.UserAgent()
	for _, spider := range util.Spiders {
		if strings.Contains(ua, spider) {
			fmt.Printf("[GIN] spider skip:%s\n", ua)
			return true
		}
	}

	return false
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(util.Counter())
	router.GET("/*path", imageHandler)
	router.Run(":" + util.ServePort)
}
