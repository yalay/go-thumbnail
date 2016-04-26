package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	cookieKey = "cnt"
)

func Counter() gin.HandlerFunc {
	return func(context *gin.Context) {
		value, err := context.Cookie(cookieKey)
		if err == nil && value != "" {
			count, _ := strconv.Atoi(value)
			count++
			cookieValue := strconv.Itoa(count)
			context.SetCookie(cookieKey, cookieValue, 0, "", "", false, false)
		} else {
			context.SetCookie(cookieKey, "0", 0, "", "", false, false)
		}
		context.Next()
	}
}
