package social

import (
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	redisCli redis.Cmdable
	dbCli    *gorm.DB
	slaveCli *gorm.DB

	infoLogger  *log.Logger
	excLogger   *log.Logger
	debugLogger *log.Logger
)
