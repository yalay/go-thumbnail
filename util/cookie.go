package util

import (
	"conf"
	"log"
	"net/http"
	"net/url"
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
		refer := context.Request.Referer()
		if !conf.IsAllowedRefer(refer) {
			referUrl, _ := url.Parse(refer)
			adDuration := conf.GetAdDuration(referUrl.Host)
			userId := getUserId(context)
			if gAccessCount.get(userId) >= adDuration {
				gAccessCount.reset(userId)
				log.Println("[GIN] Ad. UserId:" + userId)
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

func getUserId(context *gin.Context) string {
	//return context.ClientIP() + context.Request.UserAgent()
	return context.ClientIP()
}
