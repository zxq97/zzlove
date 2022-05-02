package main

import (
	"net/http"
	"zzlove/client/article"
	"zzlove/client/kafka"
	"zzlove/client/social"
	"zzlove/client/user"
	"zzlove/cmd/api/env"
	"zzlove/conf"
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

	env.ApiLogger, err = conf.InitLog(config.LogPath.Api)
	if err != nil {
		panic(err)
	}
	env.ExcLogger, err = conf.InitLog(config.LogPath.Exc)
	if err != nil {
		panic(err)
	}
	env.DbgLogger, err = conf.InitLog(config.LogPath.Debug)
	if err != nil {
		panic(err)
	}

	rpc.InitLogger(env.ApiLogger, env.ExcLogger)

	articleConn, err := rpc.NewGrpcConn(articleConf)
	if err != nil {
		panic(err)
	}
	article.InitLogger(env.ApiLogger, env.ExcLogger, env.DbgLogger)
	article.InitClient(articleConn)

	socialConn, err := rpc.NewGrpcConn(socialConf)
	if err != nil {
		panic(err)
	}
	social.InitLogger(env.ApiLogger, env.ExcLogger, env.DbgLogger)
	social.InitClient(socialConn)

	userConn, err := rpc.NewGrpcConn(userConf)
	if err != nil {
		panic(err)
	}
	user.InitLogger(env.ApiLogger, env.ExcLogger, env.DbgLogger)
	user.InitClient(userConn)

	kafka.InitLogger(env.ApiLogger, env.ExcLogger, env.DbgLogger)
	err = kafka.InitClient(config.Kafka)
	if err != nil {
		panic(err)
	}

	middleware.InitLogger(env.ApiLogger, env.ExcLogger)
	env.Route = gin.New()
	env.Route.Use(
		middleware.Access(constant.TrackKey),
		middleware.Timeout(constant.DefaultTimeout),
		middleware.Recover(),
	)

	pprof.Register(env.Route)
	s := &http.Server{
		Addr:           config.Svc.Bind,
		Handler:        env.Route,
		ReadTimeout:    constant.DefaultIOTimeout,
		WriteTimeout:   constant.DefaultIOTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
