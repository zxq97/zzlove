package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"zzlove/conf"
	"zzlove/global"
	"zzlove/internal/rpc"
	"zzlove/pb/user"
	"zzlove/server/user"
)

func main() {
	config, err := conf.LoadYaml(conf.UserConfPath)
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

	err = user.InitService(config)
	if err != nil {
		panic(err)
	}

	svc, er := rpc.NewGrpcServer(config)
	defer er.Stop()
	user_svc.RegisterUserServer(svc, UserSvc{})
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
