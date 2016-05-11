package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
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
		context.String(http.StatusOK, "skip")
		return
	}

	if !doModifiedSince(context) {
		context.String(http.StatusNotModified, "not modified")
		return
	}

	referUrl := context.Request.Referer()
	if size != "" {
		imgPath = imgPath + "?s=" + size
	} else {
		if !strings.Contains(referUrl, util.AllowedRefer) {
			imgPath = imgPath + "?s=" + util.ExtImgSize
		}
	}

	rspOriginImg(util.RedirectUrl+imgPath, context)
	return
}

// 原图本地获取或者跳转到img服务器
func rspOriginImg(imgPath string, context *gin.Context) {
	imgUrl, err := url.Parse(imgPath)
	if err != nil || imgUrl == nil {
		context.String(http.StatusNotFound, "404")
		return
	}

	if len(imgUrl.Host) > 0 {
		context.Redirect(http.StatusFound, imgPath)
		return
	}

	buff := &bytes.Buffer{}
	thumbImg := getThumbnailImg(imgUrl)
	jpeg.Encode(buff, thumbImg, nil)
	context.Data(http.StatusOK, "image/jpeg", buff.Bytes())

}

func getThumbnailImg(imgUrl *url.URL) image.Image {
	if imgUrl == nil {
		return nil
	}

	srcImg, err := util.LoadImage(imgUrl.Path)
	if err != nil {
		fmt.Printf("[GIN] LoadImage error:%v\n", err)
		return nil
	}

	imgValues, _ := url.ParseQuery(imgUrl.RawQuery)
	size := imgValues.Get("s")
	if size == "" {
		return srcImg
	} else {
		dstWidth, dstHeight := util.ParseImgArg(size)
		if dstHeight == 0 || dstWidth == 0 {
			return srcImg
		}
		return util.ThumbnailCrop(dstWidth, dstHeight, srcImg)
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

func doModifiedSince(context *gin.Context) bool {
	lastTime, err := time.Parse(http.TimeFormat, context.Request.Header.Get("If-Modified-Since"))
	if err != nil {
		return true
	}

	if lastTime.Add(4 * time.Hour).Before(time.Now()) {
		return true
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
