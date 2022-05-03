package social

import (
	"net/http"
	"zzlove/client/social"
	"zzlove/global"
	"zzlove/internal/cast"

	"github.com/gin-gonic/gin"
)

func HandleFollow(ctx *gin.Context) {
	uid := cast.ParseInt(ctx.Query("uid"), 0)
	touid := cast.ParseInt(ctx.Query("to_uid"), 0)

	if uid == touid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "uid touid same",
		})
		return
	}

	isBlack, err := social.CheckBlack(ctx.Request.Context(), uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v CheckBlack uid %v touid %v err %v", ctx, uid, touid, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	if isBlack {
		ctx.JSON(http.StatusOK, gin.H{
			"isblack": true,
		})
		return
	}

	err = social.Follow(ctx.Request.Context(), uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v Follow uid %v touid %v err %v", ctx, uid, touid, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
