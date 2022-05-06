package main

import (
	"net/http"
	_ "net/http/pprof"
	"zzlove/client/article"
	"zzlove/client/social"
	"zzlove/client/user"
	"zzlove/conf"
	"zzlove/global"
	"zzlove/internal/concurrent"
	"zzlove/internal/kafka"
	"zzlove/internal/rpc"
	"zzlove/server/async"
)

func main() {
	config, err := conf.LoadYaml(conf.ASyncConfPath)
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

	concurrent.Go(func() {
		async.InitConsumer(config.Kafka.Addr, []string{kafka.UserActionTopic}, "async_svc")
	})

	_ = http.ListenAndServe(config.Svc.Bind, nil)
}
