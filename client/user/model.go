package user

import (
	"zzlove/rpc/user"
	"zzlove/server/user"
)

func toUserinfoRequest(uid int64) *user_svc.UserInfoRequest {
	return &user_svc.UserInfoRequest{
		Uid: uid,
	}
}

func toBatchUserinfoRequest(uids []int64) *user_svc.UserInfoBatchRequest {
	return &user_svc.UserInfoBatchRequest{
		Uids: uids,
	}
}

func toListRequest(uid, cursor, offset int64) *user_svc.ListRequest {
	return &user_svc.ListRequest{
		Uid:    uid,
		Cursor: cursor,
		Offset: offset,
	}
}

func toCreateUserRequest(uid int64, nickname, introduction string, gender int32) *user_svc.CreateUserRequest {
	return &user_svc.CreateUserRequest{
		Userinfo: toUserItem(uid, nickname, introduction, gender),
	}
}

func toUserItem(uid int64, nickname, introduction string, gender int32) *user_svc.UserInfo {
	return &user_svc.UserInfo{
		Uid:          uid,
		Nickname:     nickname,
		Introduction: introduction,
		Gender:       gender,
	}
}

func toAddBrowseRequest(uid, toUID int64) *user_svc.AddBrowseRequest {
	return &user_svc.AddBrowseRequest{
		BrowseInfo: toCollectionInfo(uid, toUID),
	}
}

func toCollectionInfo(uid, targetID int64) *user_svc.RelationInfo {
	return &user_svc.RelationInfo{
		Uid:   uid,
		ToUid: targetID,
	}
}

func toCollectionRequest(uid, targetID int64) *user_svc.CollectionRequest {
	return &user_svc.CollectionRequest{
		CollectionInfo: toCollectionInfo(uid, targetID),
	}
}

func toUser(userinfo *user_svc.UserInfo) *user.User {
	return &user.User{
		UID:          userinfo.Uid,
		Nickname:     userinfo.Nickname,
		Introduction: userinfo.Introduction,
		Gender:       userinfo.Gender,
	}
}
