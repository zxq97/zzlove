package article

import (
	"net/http"
	"zzlove/cmd/api/env"
	"zzlove/server/article"
	"zzlove/server/user"

	"github.com/gin-gonic/gin"
)

func HandleInfo(ctx *gin.Context) {
	articleID := ctx.GetInt64("article_id")
	articleInfo, err := article.GetArticle(ctx.Request.Context(), articleID)
	if err != nil || articleInfo == nil {
		env.ExcLogger.Printf("ctx %v GetArticle articleid %v err %v", ctx, articleID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	uid := articleInfo.UID
	userInfo, err := user.GetUserinfo(ctx.Request.Context(), uid)
	if err != nil || userInfo == nil {
		env.ExcLogger.Printf("ctx %v GetUserinfo uid %v err %v", ctx, uid, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"article": articleInfo,
		"user":    userInfo,
	})
}
