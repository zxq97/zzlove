package main

import (
	"zzlove/cmd/api/article"
	"zzlove/cmd/api/social"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	articleGroup := r.Group("/article")
	articleGroup.GET("/info", article.HandleInfo)
	articleGroup.POST("/publish", article.HandlePublish)

	//commentGroup := r.Group("/comment")

	//feedGroup := r.Group("/feed")

	socialGroup := r.Group("/social")
	socialGroup.GET("/follow", social.HandleFollow)
	socialGroup.GET("/unfollow", social.HandleUnfollow)
	socialGroup.GET("/black", social.HandleBlack)
	socialGroup.GET("/cancel_black", social.HandleCancelBlack)
	socialGroup.GET("/follow_list", social.HandleFollowList)
	socialGroup.GET("/follower_list", social.HandleFollowerList)
	socialGroup.GET("/black_list", social.HandleBlackList)
	socialGroup.GET("/follow_count", social.HandleFollowCount)

	//userGroup := r.Group("/user")
}
