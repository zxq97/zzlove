package social

import (
	"zzlove/pb/social"
)

func toFollowRequest(uid, touid int64) *social_svc.FollowRequest {
	return &social_svc.FollowRequest{
		FollowItem: toFollowItem(uid, touid),
	}
}

func toFollowItem(uid, touid int64) *social_svc.RelationItem {
	return &social_svc.RelationItem{
		Uid:   uid,
		ToUid: touid,
	}
}

func toListRequest(uid, cursor, offset int64) *social_svc.ListRequest {
	return &social_svc.ListRequest{
		Uid:    uid,
		Cursor: cursor,
		Offset: offset,
	}
}

func toCountRequest(uid int64) *social_svc.CountRequest {
	return &social_svc.CountRequest{
		Uid: uid,
	}
}

func toRelationRequest(uid int64, uids []int64) *social_svc.RelationRequest {
	return &social_svc.RelationRequest{
		Uid:  uid,
		Uids: uids,
	}
}

func toBlackRequest(uid, touid int64) *social_svc.BlackRequest {
	return &social_svc.BlackRequest{
		BlackItem: toFollowItem(uid, touid),
	}
}
