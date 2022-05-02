package main

import (
	"log"
	"net"
	"net/http"
	"zzlove/conf"
	"zzlove/internal/rpc"
	"zzlove/pb/user"
	"zzlove/server/user"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
	dbgLogger *log.Logger
)

func main() {
	config, err := conf.LoadYaml(conf.UserConfPath)
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

	user.InitLogger(apiLogger, excLogger, dbgLogger)
	err = user.InitService(config)
	if err != nil {
		panic(err)
	}

	rpc.InitLogger(apiLogger, excLogger)
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
