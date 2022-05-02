package model

import (
	"time"
	"zzlove/pb/user"
)

type User struct {
	UID          int64  `json:"uid"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Gender       int32  `json:"gender"`
}

type Collection struct {
	UID      int64     `json:"uid"`
	TargetID int64     `json:"target_id"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

type BrowseHistory struct {
	UID   int64     `json:"uid"`
	ToUID int64     `json:"to_uid"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

func (u *User) ToUserinfo() *user_svc.UserInfo {
	return &user_svc.UserInfo{
		Uid:          u.UID,
		Nickname:     u.Nickname,
		Introduction: u.Introduction,
		Gender:       u.Gender,
	}
}

func (t *Collection) TableName() string {
	return "collection"
}

func (t *BrowseHistory) TableName() string {
	return "browse_history"
}
