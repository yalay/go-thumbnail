package util

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	adMaxCounts = 16
	adHttpFlag  = http.StatusTemporaryRedirect
)

var (
	gAccessCount *AccessCount
)

type AccessCount struct {
	sync.RWMutex
	userCount map[string]int
}

func init() {
	gAccessCount = &AccessCount{
		userCount: make(map[string]int, 0),
	}
}

func (a *AccessCount) get(userId string) int {
	a.RLock()
	count := a.userCount[userId]
	a.RUnlock()
	return count
}

func (a *AccessCount) inc(userId string) {
	a.Lock()
	count := a.userCount[userId]
	a.userCount[userId] = count + 1
	a.Unlock()
}

func (a *AccessCount) reset(userId string) {
	a.Lock()
	a.userCount[userId] = 0
	a.Unlock()
}

func Counter() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !ReferAllow(context.Request.Referer()) {
			userId := getUserId(context)
			if gAccessCount.get(userId) >= adMaxCounts {
				gAccessCount.reset(userId)
				Logln("[GIN] Ad. UserId:" + userId)
				context.Status(adHttpFlag)
			} else {
				gAccessCount.inc(userId)
			}
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
	//return context.ClientIP() + context.Request.UserAgent()
	return context.ClientIP()
}
