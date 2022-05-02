package main

import (
	"zzlove/cmd/api/article"
	"zzlove/cmd/api/social"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	articleGroup := r.Group("/article")
	articleGroup.GET("/info", article.HandleInfo)

	//commentGroup := r.Group("/comment")
	//
	//feedGroup := r.Group("/feed")

	socialGroup := r.Group("/social")
	socialGroup.GET("/follow", social.HandleFollow)

	//userGroup := r.Group("/user")
}
