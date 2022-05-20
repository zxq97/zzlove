package main

import (
	"context"
	"zzlove/pb/social"
	"zzlove/server/social"
)

type SocialSvc struct {
}

func (SocialSvc) Follow(ctx context.Context, req *social_svc.FollowRequest) (*social_svc.EmptyResponse, error) {
	err := social.Follow(ctx, req.FollowItem.Uid, req.FollowItem.ToUid)
	return &social_svc.EmptyResponse{}, err
}

func (SocialSvc) Unfollow(ctx context.Context, req *social_svc.FollowRequest) (*social_svc.EmptyResponse, error) {
	err := social.Unfollow(ctx, req.FollowItem.Uid, req.FollowItem.ToUid)
	return &social_svc.EmptyResponse{}, err
}

func (SocialSvc) GetFollow(ctx context.Context, req *social_svc.ListRequest) (*social_svc.ListResponse, error) {
	uids, nextCur, err := social.GetFollow(ctx, req.Uid, req.Cursor, req.Offset)
	return &social_svc.ListResponse{
		Uids:       uids,
		NextCursor: nextCur,
	}, err
}

func (SocialSvc) GetFollower(ctx context.Context, req *social_svc.ListRequest) (*social_svc.ListResponse, error) {
	uids, nextCur, err := social.GetFollower(ctx, req.Uid, req.Cursor, req.Offset)
	return &social_svc.ListResponse{
		Uids:       uids,
		NextCursor: nextCur,
	}, err
}

func (SocialSvc) GetFollowCount(ctx context.Context, req *social_svc.CountRequest) (*social_svc.CountResponse, error) {
	fcnt, frcnt, err := social.GetFollowCount(ctx, req.Uid)
	return &social_svc.CountResponse{
		FollowCount:   fcnt,
		FollowerCount: frcnt,
	}, err
}

func (SocialSvc) GetRelations(ctx context.Context, req *social_svc.RelationRequest) (*social_svc.RelationResponse, error) {
	relationMap, err := social.GetRelations(ctx, req.Uid, req.Uids)
	return &social_svc.RelationResponse{
		Relation: relationMap,
	}, err
}

func (SocialSvc) Black(ctx context.Context, req *social_svc.BlackRequest) (*social_svc.EmptyResponse, error) {
	err := social.Black(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	return &social_svc.EmptyResponse{}, err
}

func (SocialSvc) CancelBlack(ctx context.Context, req *social_svc.BlackRequest) (*social_svc.EmptyResponse, error) {
	err := social.CancelBlack(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	return &social_svc.EmptyResponse{}, err
}

func (SocialSvc) CheckBlack(ctx context.Context, req *social_svc.BlackRequest) (*social_svc.BlackResponse, error) {
	ok, err := social.CheckBlack(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	return &social_svc.BlackResponse{
		IsBlack: ok,
	}, err
}

func (SocialSvc) CheckBatchBlack(ctx context.Context, req *social_svc.RelationRequest) (*social_svc.BlackBatchResponse, error) {
	blackMap, err := social.CheckBatchBlack(ctx, req.Uid, req.Uids)
	return &social_svc.BlackBatchResponse{
		Relation: blackMap,
	}, err
}

func (SocialSvc) GetBlackList(ctx context.Context, req *social_svc.ListRequest) (*social_svc.ListResponse, error) {
	uids, nextCur, err := social.GetBlackList(ctx, req.Uid, req.Cursor, req.Offset)
	return &social_svc.ListResponse{
		Uids:       uids,
		NextCursor: nextCur,
	}, err
}
