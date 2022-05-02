package rpc

import (
	"fmt"
	"log"
	"time"
	"zzlove/conf"
	"zzlove/internal/constant"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
)

func InitLogger(apiLog, excLog *log.Logger) {
	apiLogger = apiLog
	excLogger = excLog
}

func NewGrpcConn(config *conf.Conf) (*grpc.ClientConn, error) {
	er := newEtcdDiscover(config.Etcd.Addr, time.Duration(config.Etcd.TTL)*time.Second, config.Svc.Name)
	resolver.Register(er)
	conn, err := grpc.Dial(
		er.Scheme()+":///",
		grpc.WithInsecure(),
		grpc.WithResolvers(er),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithUnaryInterceptor(timeout(constant.DefaultTimeout)))
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
