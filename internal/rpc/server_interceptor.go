package rpc

import (
	"context"
	"runtime/debug"
	"time"
	"zzlove/internal/constant"

	"google.golang.org/grpc"
)

func recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	reqid := getIncomingReqID(ctx)
	if reqid != "" {
		ctx = context.WithValue(ctx, constant.ReqIDKey, reqid)
	}
	now := time.Now()
	defer func() {
		apiLogger.Println(time.Since(now))
		if err := recover(); err != nil {
			excLogger.Println(info.FullMethod, err, string(debug.Stack()))
		}
	}()
	return handler(ctx, req)
}
