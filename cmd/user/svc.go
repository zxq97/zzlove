package main

import (
	"context"
	"zzlove/pb/user"
	"zzlove/server/user"
)

type UserSvc struct {
}

func (UserSvc) GetUserinfo(ctx context.Context, req *user_svc.UserInfoRequest) (*user_svc.UserInfoResponse, error) {
	userInfo, err := user.GetUserinfo(ctx, req.Uid)
	if err != nil || userInfo == nil {
		return &user_svc.UserInfoResponse{}, err
	}
	return &user_svc.UserInfoResponse{
		Userinfo: userInfo.ToUserinfo(),
	}, nil
}

func (UserSvc) GetBatchUserinfo(ctx context.Context, req *user_svc.UserInfoBatchRequest) (*user_svc.UserInfoBatchResponse, error) {
	userMap, err := user.GetBatchUserinfo(ctx, req.Uids)
	if err != nil || userMap == nil {
		return &user_svc.UserInfoBatchResponse{}, err
	}
	infoMap := make(map[int64]*user_svc.UserInfo, len(userMap))
	for k, v := range userMap {
		infoMap[k] = v.ToUserinfo()
	}
	return &user_svc.UserInfoBatchResponse{
		Userinfos: infoMap,
	}, nil
}

func (UserSvc) GetHistoryBrowse(ctx context.Context, req *user_svc.ListRequest) (*user_svc.ListResponse, error) {
	uids, nextCur, err := user.GetHistoryBrowse(ctx, req.Uid, req.Cursor, req.Offset)
	if err != nil || uids == nil {
		return nil, err
	}
	return &user_svc.ListResponse{
		TargetIds:  uids,
		NextCursor: nextCur,
	}, nil
}

func (UserSvc) GetCollectionList(ctx context.Context, req *user_svc.ListRequest) (*user_svc.ListResponse, error) {
	uids, nextCur, err := user.GetCollectionList(ctx, req.Uid, req.Cursor, req.Offset)
	if err != nil || uids == nil {
		return nil, err
	}
	return &user_svc.ListResponse{
		TargetIds:  uids,
		NextCursor: nextCur,
	}, nil
}

func (UserSvc) Collection(ctx context.Context, req *user_svc.CollectionRequest) (*user_svc.EmptyResponse, error) {
	err := user.Collection(ctx, req.CollectionInfo.Uid, req.CollectionInfo.ToUid)
	return nil, err
}

func (UserSvc) CancelCollection(ctx context.Context, req *user_svc.CollectionRequest) (*user_svc.EmptyResponse, error) {
	err := user.CancelCollection(ctx, req.CollectionInfo.Uid, req.CollectionInfo.ToUid)
	return nil, err
}

func (UserSvc) AddBrowse(ctx context.Context, req *user_svc.AddBrowseRequest) (*user_svc.EmptyResponse, error) {
	err := user.AddBrowse(ctx, req.BrowseInfo.Uid, req.BrowseInfo.ToUid)
	return nil, err
}

func (UserSvc) CreateUser(ctx context.Context, req *user_svc.CreateUserRequest) (*user_svc.EmptyResponse, error) {
	err := user.CreateUser(ctx, req.Userinfo.Uid, req.Userinfo.Gender, req.Userinfo.Nickname, req.Userinfo.Introduction)
	return nil, err
}
