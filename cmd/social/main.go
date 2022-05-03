package main

import (
	"net"
	"net/http"
	"zzlove/conf"
	"zzlove/global"
	"zzlove/internal/rpc"
	"zzlove/pb/social"
	"zzlove/server/social"
)

func main() {
	config, err := conf.LoadYaml(conf.SocialConfPath)
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

	err = social.InitService(config)
	if err != nil {
		panic(err)
	}

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
