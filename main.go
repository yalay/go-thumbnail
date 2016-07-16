package main

import (
	"bytes"
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
		context.Status(http.StatusOK)
		return
	}

	if size != "" {
		imgPath = imgPath + "?s=" + size
	} else {
		if !util.ReferAllow(context.Request.Referer()) {
			if util.DoAd(context) {
				adImgPath := util.GetRandomAdPath()
				if adImgPath != "" {
					imgPath = adImgPath
				}
			} else {
				imgPath = imgPath + "?s=" + util.ExtImgSize
			}
		}
	}

	rspImg(util.RedirectUrl+imgPath, context)
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
		util.Logln("[GIN] LoadImage error:" + err.Error())
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

		thumbImg := util.ThumbnailCrop(dstWidth, dstHeight, srcImg)
		go util.WriteCache(imgUrl.String(), thumbImg)
		return thumbImg
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
			util.Logln("[GIN] spider skip:" + ua)
			return true
		}
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
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(util.Counter(), gin.LoggerWithWriter(util.GetLogBuf()), gin.Recovery())
	router.GET("/*path", imageHandler)
	router.Run(":" + util.ServePort)
}
