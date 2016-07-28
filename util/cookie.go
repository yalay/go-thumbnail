package util

import (
	"io/ioutil"
	"math/rand"
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

// 并发获取cookie缓存,用户id+[时间戳][计数]
var (
	cookieBuffMap = make(map[string]int, 0)
	cookieDelay   = 24 * time.Second
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

			userId := getUserId(context)
			cnt := getCookieValue(userId, cookie)
			cnt = setAdStatus(cnt, userId, context)
			cookie.Value = strconv.Itoa(cnt)
			http.SetCookie(context.Writer, cookie)
			//Logln("[GIN] userId:" + userId + " cookie value:" + cookie.Value)
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

func DoAd(context *gin.Context) bool {
	if context.Writer.Status() == adHttpFlag {
		return true
	}
	return false
}

func GetRandomAdPath() string {
	imgs, err := ioutil.ReadDir(ImgRoot + AdPath)
	if err != nil || len(imgs) == 0 {
		return ""
	}
	return AdPath + imgs[rand.Intn(len(imgs))].Name()
}

func getUserId(context *gin.Context) string {
	return context.ClientIP() + context.Request.UserAgent()
}

func getCookieValue(userId string, cookie *http.Cookie) int {
	if cnt, ok := cookieBuffMap[userId]; ok {
		return cnt
	} else {
		cookieCnt, _ := strconv.Atoi(cookie.Value)
		cookieBuffMap[userId] = cookieCnt
		go cookieBuffDelay(userId)
		return cookieCnt
	}
}

func setAdStatus(count int, userId string, context *gin.Context) int {
	if count < adMaxCounts {
		count++
	} else {
		count = 0
		Logln("[GIN] Ad. UserId:" + userId)
		context.Status(adHttpFlag)
	}

	if _, ok := cookieBuffMap[userId]; ok {
		cookieBuffMap[userId] = count
	}
	return count
}

func cookieBuffDelay(key string) {
	time.Sleep(cookieDelay)
	if _, ok := cookieBuffMap[key]; ok {
		delete(cookieBuffMap, key)
	}
}
