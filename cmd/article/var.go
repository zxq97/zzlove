package main

import (
	"log"
	"zzlove/dal/article"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
	dbgLogger *log.Logger

	ArticleDAL *article.ArticleDAL
)
