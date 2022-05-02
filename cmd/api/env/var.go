package env

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	ApiLogger *log.Logger
	ExcLogger *log.Logger
	DbgLogger *log.Logger

	Route *gin.Engine
)
