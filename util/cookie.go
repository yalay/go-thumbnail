package util

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cookieKey = "cnt"
)

func Counter() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Request.Cookie(cookieKey)
		if err == nil {
			cnt, _ := strconv.Atoi(cookie.Value)
			cnt++
			cookie.Value = strconv.Itoa(cnt)
			http.SetCookie(context.Writer, cookie)
			Log(cookie.Value)
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
