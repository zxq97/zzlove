package social

import "time"

const (
	RelationFollow   = 1
	RelationFollower = 10
)

type Follow struct {
	UID       int64     `json:"uid"`
	FollowUID int64     `json:"follow_uid"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

type Follower struct {
	UID         int64     `json:"uid"`
	FollowerUID int64     `json:"follower_uid"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

type FollowCount struct {
	UID           int64 `json:"uid"`
	FollowCount   int64 `json:"follow_count"`
	FollowerCount int64 `json:"follower_count"`
}

type UserBlack struct {
	UID           int64     `json:"uid"`
	BlackTargetID int64     `json:"black_target_id"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

func (t *Follow) TableName() string {
	return "follow"
}

func (t *Follower) TableName() string {
	return "follower"
}

func (t *FollowCount) TableName() string {
	return "follow_count"
}

func (t *UserBlack) TableName() string {
	return "user_black"
}
