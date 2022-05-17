package rpc

import (
	"fmt"
	"time"
	"zzlove/conf"
	"zzlove/internal/constant"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

func NewGrpcConn(config *conf.Conf) (*grpc.ClientConn, error) {
	er := newEtcdDiscover(config.Etcd.Addr, time.Duration(config.Etcd.TTL)*time.Second, config.Svc.Name)
	resolver.Register(er)
	initBreaker(config.Svc.Name)
	conn, err := grpc.Dial(
		er.Scheme()+":///",
		grpc.WithInsecure(),
		grpc.WithResolvers(er),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				timeout(constant.DefaultTimeout),
				demote(config.Svc.Name),
			)))
	return conn, err
}

func NewGrpcServer(config *conf.Conf) (*grpc.Server, *EtcdRegister) {
	er := newEtcdRegister(config.Etcd.Addr, time.Duration(config.Etcd.TTL)*time.Second, config.Svc.Name+"_"+config.Svc.Addr, config.Svc.Addr)
	opt := []grpc.ServerOption{
		grpc.UnaryInterceptor(recovery),
	}
	svc := grpc.NewServer(opt...)
	return svc, er
}
