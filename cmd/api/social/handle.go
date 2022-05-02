package social

import (
	"net/http"
	"zzlove/client/social"
	"zzlove/cmd/api/env"

	"github.com/gin-gonic/gin"
)

func HandleFollow(ctx *gin.Context) {
	uid := ctx.GetInt64("uid")
	touid := ctx.GetInt64("to_uid")

	if uid == touid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "uid touid same",
		})
		return
	}

	isBlack, err := social.CheckBlack(ctx.Request.Context(), uid, touid)
	if err != nil {
		env.ExcLogger.Printf("ctx %v CheckBlack uid %v touid %v err %v", ctx, uid, touid, err)
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
		env.ExcLogger.Printf("ctx %v Follow uid %v touid %v err %v", ctx, uid, touid, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
