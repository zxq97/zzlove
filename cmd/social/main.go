package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
	social_svc.RegisterSocialServer(svc, SocialSvc{})
	_, err = er.Register()
	if err != nil {
		panic(err)
	}

	errch := make(chan error)
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = http.ListenAndServe(config.Svc.Bind, nil)
		errch <- err
	}()

	lis, err := net.Listen("tcp", config.Svc.Addr)
	if err != nil {
		panic(err)
	}

	go func() {
		err = svc.Serve(lis)
		errch <- err
	}()

	select {
	case <-sign:
		global.ApiLogger.Println("receive signal done")
		er.Stop()
	case <-errch:
		global.ExcLogger.Println("receive err done", err)
		er.Stop()
	}
}
