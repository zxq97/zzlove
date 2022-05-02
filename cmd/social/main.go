package main

import (
	"log"
	"net"
	"net/http"
	"zzlove/conf"
	"zzlove/internal/rpc"
	"zzlove/pb/social"
	"zzlove/server/article"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
	dbgLogger *log.Logger
)

func main() {
	config, err := conf.LoadYaml(conf.SocialConfPath)
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

	article.InitLogger(apiLogger, excLogger, dbgLogger)
	err = article.InitService(config)
	if err != nil {
		panic(err)
	}

	rpc.InitLogger(apiLogger, excLogger)
	svc, er := rpc.NewGrpcServer(config)
	defer er.Stop()
	social_svc.RegisterSocialServer(svc, SocialSvc{})
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
