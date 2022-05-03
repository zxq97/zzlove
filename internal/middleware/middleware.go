package middleware

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"
	"zzlove/global"
	"zzlove/internal/constant"
	"zzlove/internal/generate"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.ExcLogger.Println(err, string(debug.Stack()))
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
		global.ApiLogger.Println(c.Request.URL, reqid, time.Since(now))
	}
}
