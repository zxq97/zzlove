package user

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	redisCli redis.Cmdable
	mcCli    *memcache.Client
	dbCli    *gorm.DB
	slaveCli *gorm.DB

	infoLogger  *log.Logger
	excLogger   *log.Logger
	debugLogger *log.Logger
)
