package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"
	"zzlove/internal/constant"
	"zzlove/internal/generate"

	"github.com/gin-gonic/gin"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
)

func InitLogger(apiLog, excLog *log.Logger) {
	apiLogger = apiLog
	excLogger = excLog
}

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				excLogger.Println(err, string(debug.Stack()))
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}

func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func Access(track string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), constant.TrackKey, track)
		reqid, ok := c.Get(constant.ReqIDKey)
		if !ok {
			reqid = generate.UUID()
		}
		ctx = context.WithValue(ctx, constant.ReqIDKey, reqid)
		c.Request = c.Request.WithContext(ctx)
		now := time.Now()
		c.Next()
		apiLogger.Println(c.Request.URL, reqid, time.Since(now))
	}
}
