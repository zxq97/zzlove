package article

import (
	"fmt"
	"log"
	"zzlove/conf"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jinzhu/gorm"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
)

func InitLogger(apiLog, excLog *log.Logger) {
	apiLogger = apiLog
	excLogger = excLog
}

type ArticleDAL struct {
	mc    *memcache.Client
	db    *gorm.DB
	slave *gorm.DB
}

func NewArticleDAL(config *conf.Conf) (*ArticleDAL, error) {
	mcCli := conf.GetMC(config.MC.Addr)
	dbCli, err := conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	if err != nil {
		return nil, err
	}
	slaveCli, err := conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Slave.User, config.Slave.Password, config.Slave.Host, config.Slave.Port, config.Slave.DB))
	if err != nil {
		return nil, err
	}
	return &ArticleDAL{
		mc:    mcCli,
		db:    dbCli,
		slave: slaveCli,
	}, nil
}
