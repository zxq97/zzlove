package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"zzlove/conf"
	"zzlove/dal/article"
	"zzlove/internal/rpc"
	"zzlove/pb/article"
)

func main() {
	config, err := conf.LoadYaml(conf.ArticleConfPath)
	if err != nil {
		panic(err)
	}

	apiLogger, err = conf.InitLog(config.LogPath.Api)
	if err != nil {
		panic(err)
	}
	excLogger, err = conf.InitLog(config.LogPath.Exc)
	if err != nil {
		panic(err)
	}
	dbgLogger, err = conf.InitLog(config.LogPath.Debug)
	if err != nil {
		panic(err)
	}

	article.InitLogger(apiLogger, excLogger)
	ArticleDAL, err = article.NewArticleDAL(config)
	if err != nil {
		panic(err)
	}

	rpc.InitLogger(apiLogger, excLogger)
	svc, er := rpc.NewGrpcServer(config)
	defer er.Stop()
	article_svc.RegisterArticleServer(svc, ArcileSvc{})
	_, err = er.Register()
	if err != nil {
		panic(err)
	}

	go func() {
		_ = http.ListenAndServe(config.Svc.Bind, nil)
	}()

	lis, err := net.Listen("tcp", config.Svc.Addr)
	if err != nil {
		panic(err)
	}
	_ = svc.Serve(lis)
}
