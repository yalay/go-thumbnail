package main

import (
	"bytes"
	"fmt"
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

	cacheBuff := util.FindInCache(imgPath, size)
	if len(cacheBuff) > 0 {
		context.Data(http.StatusOK, "image/jpeg", cacheBuff)
		return
	}

	if size != "" {
		rspThumbnailImg(imgPath, size, context)
		return
	}

	referUrl := context.Request.Referer()
	if strings.Contains(referUrl, util.AllowedRefer) {
		rspOriginImg(imgPath, context)
	} else {
		rspWaterMarkImg(imgPath, context)
	}

	return
}

// 原图获取跳转到img服务器
func rspOriginImg(imgPath string, context *gin.Context) {
	context.Redirect(http.StatusFound, util.RedirectUrl+imgPath)
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
		context.String(http.StatusNotFound, "LoadImage error:%v", err)
		return
	}

	thumbImg := util.Thumbnail(dstWidth, dstHeight, srcImg)
	dstImg := util.CropImg(thumbImg, int(dstWidth), int(dstHeight))
	buff := &bytes.Buffer{}
	jpeg.Encode(buff, dstImg, nil)
	context.Data(http.StatusOK, "image/jpeg", buff.Bytes())

	go util.WriteCache(imgPath, size, dstImg)
	return
}

func rspWaterMarkImg(imgPath string, context *gin.Context) {
	cacheBuff := util.FindInCache(imgPath, util.WaterSize)
	if len(cacheBuff) > 0 {
		context.Data(http.StatusOK, "image/jpeg", cacheBuff)
		return
	}

	waterImg, err := util.WaterMark(imgPath)
	if err != nil {
		context.String(http.StatusNotFound, "Water mark error:%v", err)
	} else {
		waterBuff := &bytes.Buffer{}
		jpeg.Encode(waterBuff, waterImg, nil)
		context.Data(http.StatusOK, "image/jpeg", waterBuff.Bytes())

		go util.WriteCache(imgPath, util.WaterSize, waterImg)
	}

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
	router.GET("/*path", imageHandler)
	router.Run(":" + util.ServePort)
}
