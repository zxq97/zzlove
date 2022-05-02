package constant

import "time"

const (
	ReqIDKey = "req_id"
	TrackKey = "track"

	ApiTrackKey     = "api"
	ArticleTrackKey = "article"
	AsyncTrackKey   = "async"
	ChatTrackKey    = "char"
	CommentTrackKey = "comment"
	FeedTrackKey    = "feed"
	SearchTrackKey  = "search"
	SocialTrackKey  = "social"
	UserTrackKey    = "user"

	DefaultTimeout = time.Second
	DefaultTicker  = time.Minute

	EtcdLeaseTTL = 10

	EtcdScheme = "etcd"
)
