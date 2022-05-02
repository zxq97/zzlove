package rpc

import (
	"context"
	"time"
	"zzlove/internal/constant"
	"zzlove/internal/generate"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func defaultTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, d)
	}
	return ctx, cancel
}

func withOutgoing(ctx context.Context) context.Context {
	rawid := ctx.Value(constant.ReqIDKey)
	if rawid != nil {
		if reqid, ok := rawid.(string); ok {
			return metadata.AppendToOutgoingContext(ctx, constant.ReqIDKey, reqid)
		}
	}
	return metadata.AppendToOutgoingContext(ctx, constant.ReqIDKey, generate.UUID())
}

func getIncomingReqID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		rs := md.Get(constant.ReqIDKey)
		if len(rs) > 0 {
			return rs[0]
		}
	}
	return ""
}

func getReqInfo(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if ok {
		return p.Addr.String()
	}
	return ""
}
