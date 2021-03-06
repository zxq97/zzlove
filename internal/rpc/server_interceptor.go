package rpc

import (
	"context"
	"runtime/debug"
	"time"
	"zzlove/global"
	"zzlove/internal/constant"

	"google.golang.org/grpc"
)

func recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	reqid := getIncomingReqID(ctx)
	if reqid != "" {
		ctx = context.WithValue(ctx, constant.ReqIDKey, reqid)
	}
	clientip := getReqInfo(ctx)
	if clientip != "" {
		ctx = context.WithValue(ctx, constant.ClientIPKey, clientip)
	}
	now := time.Now()
	defer func() {
		global.ApiLogger.Println(time.Since(now))
		if err := recover(); err != nil {
			global.ExcLogger.Println(info.FullMethod, err, string(debug.Stack()))
		}
	}()
	return handler(ctx, req)
}
