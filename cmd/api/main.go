package main

import (
	"net/http"
	"zzlove/client/article"
	"zzlove/client/kafka"
	"zzlove/client/social"
	"zzlove/client/user"
	"zzlove/conf"
	"zzlove/global"
	"zzlove/internal/constant"
	"zzlove/internal/middleware"
	"zzlove/internal/rpc"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := conf.LoadYaml(conf.ApiConfPath)
	if err != nil {
		panic(err)
	}
	articleConf, err := conf.LoadYaml(conf.ArticleConfPath)
	if err != nil {
		panic(err)
	}
	//commentConf, err := conf.LoadYaml(conf.CommentConfPath)
	//if err != nil {
	//	panic(err)
	//}
	//feedConf, err := conf.LoadYaml(conf.FeedConfPath)
	//if err != nil {
	//	panic(err)
	//}
	//searchConf, err := conf.LoadYaml(conf.SearchConfPath)
	//if err != nil {
	//	panic(err)
	//}
	socialConf, err := conf.LoadYaml(conf.SocialConfPath)
	if err != nil {
		panic(err)
	}
	userConf, err := conf.LoadYaml(conf.UserConfPath)
	if err != nil {
		panic(err)
	}

	global.ApiLogger, err = conf.InitLog(config.LogPath.Api)
	if err != nil {
		panic(err)
	}
	global.ExcLogger, err = conf.InitLog(config.LogPath.Exc)
	if err != nil {
		panic(err)
	}
	global.DbgLogger, err = conf.InitLog(config.LogPath.Debug)
	if err != nil {
		panic(err)
	}

	articleConn, err := rpc.NewGrpcConn(articleConf)
	if err != nil {
		panic(err)
	}
	article.InitClient(articleConn)

	socialConn, err := rpc.NewGrpcConn(socialConf)
	if err != nil {
		panic(err)
	}
	social.InitClient(socialConn)

	userConn, err := rpc.NewGrpcConn(userConf)
	if err != nil {
		panic(err)
	}
	user.InitClient(userConn)

	err = kafka.InitClient(config.Kafka)
	if err != nil {
		panic(err)
	}

	route := gin.New()
	route.Use(
		middleware.Access(constant.TrackKey),
		middleware.Timeout(constant.DefaultTimeout),
		middleware.Recover(),
	)

	Router(route)

	pprof.Register(route)
	s := &http.Server{
		Addr:           config.Svc.Bind,
		Handler:        route,
		ReadTimeout:    constant.DefaultIOTimeout,
		WriteTimeout:   constant.DefaultIOTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
