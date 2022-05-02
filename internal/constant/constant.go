package constant

import "time"

const (
	ClientIPKey = "client_ip"
	ReqIDKey    = "req_id"
	TrackKey    = "track"

	ApiTrackKey     = "api"
	ArticleTrackKey = "article"
	AsyncTrackKey   = "async"
	ChatTrackKey    = "char"
	CommentTrackKey = "comment"
	FeedTrackKey    = "feed"
	SearchTrackKey  = "search"
	SocialTrackKey  = "social"
	UserTrackKey    = "user"

	DefaultBatchCount = 1000

	DefaultTTL       = 20 * time.Second
	DefaultTimeout   = 3 * time.Second
	DefaultIOTimeout = 10 * time.Second
	DefaultTicker    = time.Minute

	EtcdLeaseTTL = 10

	EtcdScheme = "etcd"
)
