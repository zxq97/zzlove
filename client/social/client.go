package social

import (
	"context"

	"google.golang.org/grpc"

	"zzlove/pb/social"
)

var (
	client social_svc.SocialClient
)

func InitClient(conn *grpc.ClientConn) {
	client = social_svc.NewSocialClient(conn)
}

func Follow(ctx context.Context, uid, touid int64) error {
	_, err := client.Follow(ctx, toFollowRequest(uid, touid))
	return err
}

func Unfollow(ctx context.Context, uid, touid int64) error {
	_, err := client.Unfollow(ctx, toFollowRequest(uid, touid))
	return err
}

func GetFollow(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	res, err := client.GetFollow(ctx, toListRequest(uid, cursor, offset))
	if err != nil || res == nil {
		return nil, 0, err
	}
	return res.Uids, res.NextCursor, nil
}

func GetFollower(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	res, err := client.GetFollower(ctx, toListRequest(uid, cursor, offset))
	if err != nil || res == nil {
		return nil, 0, err
	}
	return res.Uids, res.NextCursor, nil
}

func GetFollowCount(ctx context.Context, uid int64) (int64, int64, error) {
	res, err := client.GetFollowCount(ctx, toCountRequest(uid))
	if err != nil || res == nil {
		return 0, 0, err
	}
	return res.FollowCount, res.FollowerCount, nil
}

func GetRelations(ctx context.Context, uid int64, uids []int64) (map[int64]int32, error) {
	res, err := client.GetRelations(ctx, toRelationRequest(uid, uids))
	if err != nil || res == nil {
		return nil, err
	}
	return res.Relation, nil
}

func Black(ctx context.Context, uid, touid int64) error {
	_, err := client.Black(ctx, toBlackRequest(uid, touid))
	return err
}

func CancelBlack(ctx context.Context, uid, touid int64) error {
	_, err := client.CancelBlack(ctx, toBlackRequest(uid, touid))
	return err
}

func CheckBlack(ctx context.Context, uid, touid int64) (bool, error) {
	res, err := client.CheckBlack(ctx, toBlackRequest(uid, touid))
	if err != nil || res == nil {
		return false, err
	}
	return res.IsBlack, nil
}

func CheckBatchBlack(ctx context.Context, uid int64, uids []int64) (map[int64]bool, error) {
	res, err := client.CheckBatchBlack(ctx, toRelationRequest(uid, uids))
	if err != nil || res == nil {
		return nil, err
	}
	return res.Relation, nil
}

func GetBlackList(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	res, err := client.GetBlackList(ctx, toListRequest(uid, cursor, offset))
	if err != nil || res == nil {
		return nil, 0, err
	}
	return res.Uids, res.NextCursor, nil
}
