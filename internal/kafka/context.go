package kafka

import (
	"context"
	"zzlove/internal/constant"
	"zzlove/internal/generate"

	"google.golang.org/grpc/metadata"
)

func ConsumerContext(msgID string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultTimeout)
	if msgID == "" {
		msgID = generate.UUID()
	}
	ctx = context.WithValue(ctx, constant.MsgIDKey, msgID)
	ctx = metadata.AppendToOutgoingContext(ctx, constant.MsgIDKey, msgID)
	return ctx, cancel
}
