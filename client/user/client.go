package user

import (
	"context"
	"zzlove/internal/model"
	"zzlove/pb/user"

	"google.golang.org/grpc"
)

var (
	client user_svc.UserClient
)

func InitClient(conn *grpc.ClientConn) {
	client = user_svc.NewUserClient(conn)
}

func GetUserinfo(ctx context.Context, uid int64) (*model.User, error) {
	res, err := client.GetUserinfo(ctx, toUserinfoRequest(uid))
	if err != nil || res == nil {
		return nil, err
	}
	return toUser(res.Userinfo), nil
}

func GetBatchUserinfo(ctx context.Context, uids []int64) (map[int64]*model.User, error) {
	res, err := client.GetBatchUserinfo(ctx, toBatchUserinfoRequest(uids))
	if err != nil || res == nil {
		return nil, err
	}
	userMap := make(map[int64]*model.User, len(uids))
	for k, v := range res.Userinfos {
		userMap[k] = toUser(v)
	}
	return userMap, nil
}

func GetHistoryBrowse(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	res, err := client.GetHistoryBrowse(ctx, toListRequest(uid, cursor, offset))
	if err != nil || res == nil {
		return nil, 0, err
	}
	return res.TargetIds, res.NextCursor, nil
}

func GetCollectionList(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	res, err := client.GetCollectionList(ctx, toListRequest(uid, cursor, offset))
	if err != nil || res == nil {
		return nil, 0, err
	}
	return res.TargetIds, res.NextCursor, nil
}

func Collection(ctx context.Context, uid, touid int64) error {
	_, err := client.Collection(ctx, toCollectionRequest(uid, touid))
	return err
}

func CancelCollection(ctx context.Context, uid, touid int64) error {
	_, err := client.CancelCollection(ctx, toCollectionRequest(uid, touid))
	return err
}

func AddBrowse(ctx context.Context, uid, touid int64) error {
	_, err := client.AddBrowse(ctx, toAddBrowseRequest(uid, touid))
	return err
}

func CreateUser(ctx context.Context, uid int64, nickname, introduction string, gender int32) error {
	_, err := client.CreateUser(ctx, toCreateUserRequest(uid, nickname, introduction, gender))
	return err
}
