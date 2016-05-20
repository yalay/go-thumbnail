package util

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cookieKey   = "cnt"
	adMaxCounts = 24
	adHttpFlag  = http.StatusTemporaryRedirect
)

func Counter() gin.HandlerFunc {
	return func(context *gin.Context) {
		if ReferAllow(context.Request.Referer()) {
			return
		}

		cookie, err := context.Request.Cookie(cookieKey)
		if err == nil {
			if cookie.Path != "/" {
				cookie.Path = "/"
				cookie.Expires = time.Now().Add(4 * time.Hour)
			}
			cnt, _ := strconv.Atoi(cookie.Value)
			cnt = setAdStatus(cnt, context)
			cookie.Value = strconv.Itoa(cnt)
			http.SetCookie(context.Writer, cookie)
			Logln("[GIN] cookie value:" + cookie.Value)
		} else {
			http.SetCookie(context.Writer, &http.Cookie{
				Name:    cookieKey,
				Value:   "0",
				Path:    "/",
				MaxAge:  4 * 60 * 60,
				Expires: time.Now().Add(4 * time.Hour),
			})
		}
		context.Next()
	}
}

func setAdStatus(count int, context *gin.Context) int {
	if count < adMaxCounts {
		return count + 1
	}
	context.Status(adHttpFlag)
	return 0
}

func DoAd(context *gin.Context) bool {
	if context.Writer.Status() == adHttpFlag {
		return true
	}
	return false
}

func GetAdImgPath() string {
	return AdPath + "random.jpg"
}
