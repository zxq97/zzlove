package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func timeout(d time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, res interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var cancel context.CancelFunc
		ctx = withOutgoing(ctx)
		ctx, cancel = defaultTimeout(ctx, d)
		if cancel != nil {
			defer cancel()
		}
		return invoker(ctx, method, req, res, cc, opts...)
	}
}
