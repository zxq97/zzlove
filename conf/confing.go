package conf

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
)

const (
	ApiConfPath     = "/home/work/zzlove/conf/yaml/api.yaml"
	ArticleConfPath = "/home/work/zzlove/conf/yaml/article.yaml"
	ASyncConfPath   = "/home/work/zzlove/conf/yaml/async.yaml"
	CommentConfPath = "/home/work/zzlove/conf/yaml/comment.yaml"
	SocialConfPath  = "/home/work/zzlove/conf/yaml/social.yaml"
	UserConfPath    = "/home/work/zzlove/conf/yaml/user.yaml"
	FeedConfPath    = "/home/work/zzlove/conf/yaml/feed.yaml"
	SearchConfPath  = "/home/work/zzlove/conf/yaml/search.yaml"

	MysqlAddr = "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True"
)

type MysqlConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RedisConf struct {
	Addr string `yaml:"addr"`
	DB   int    `yaml:"db"`
}

type RedisClusterConf struct {
	Addr []string `yaml:"addr"`
}

type MCConf struct {
	Addr []string `yaml:"addr"`
}

type SvcConf struct {
	Bind string `yaml:"bind"`
	Addr string `yaml:"addr"`
	Name string `yaml:"name"`
}

type HystrixConf struct {
	Name string `yaml:"name"`
	TTL  int    `yaml:"ttl"`
	MCR  int    `yaml:"mcr"`
	RVT  int    `yaml:"rvt"`
	SW   int    `yaml:"sw"`
	EPT  int    `yaml:"ept"`
}

type EtcdConf struct {
	Addr []string `yaml:"addr"`
	TTL  int      `yaml:"ttl"`
}

type KafkaConf struct {
	Addr []string `yaml:"addr"`
}

type LogConf struct {
	Api   string `yaml:"api"`
	Exc   string `yaml:"exc"`
	Debug string `yaml:"debug"`
}

type Conf struct {
	Mysql        MysqlConf        `yaml:"mysql"`
	Slave        MysqlConf        `yaml:"slave"`
	RedisCluster RedisClusterConf `yaml:"cluster"`
	MC           MCConf           `yaml:"mc"`
	Svc          SvcConf          `yaml:"svc"`
	Hystrix      HystrixConf      `yaml:"hystrix"`
	Etcd         EtcdConf         `yaml:"etcd"`
	Kafka        KafkaConf        `yaml:"kafka"`
	LogPath      LogConf          `yaml:"log_path"`
}

func LoadYaml(path string) (*Conf, error) {
	conf := new(Conf)
	y, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(y, conf)
	return conf, err
}

func GetMC(addr []string) *memcache.Client {
	return memcache.New(addr...)
}

func GetRedisCluster(addr []string) redis.Cmdable {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addr,
	})
}

func GetGorm(addr string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", addr)
	if err != nil {
		return nil, err
	}
	db.DB().SetConnMaxLifetime(time.Minute * 3)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	return db, nil
}

func InitLog(path string) (*log.Logger, error) {
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return log.New(fp, "", log.LstdFlags|log.Lshortfile), nil
}
