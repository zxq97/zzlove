package async

import (
	"context"
	"zzlove/client/social"
	"zzlove/global"
)

func follow(ctx context.Context, uid, touid int64) {

}

func unfollow(ctx context.Context, uid, touid int64) {

}

func black(ctx context.Context, uid, touid int64) {
	err := social.Unfollow(ctx, uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v black Unfollow uid %v touid %v err %v", ctx, uid, touid)
	}
}
