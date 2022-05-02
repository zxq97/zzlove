package kafka

import "log"

const (
	UserActionTopic = "user_action"

	EventPublish  = "publish"
	EventFollow   = "follow"
	EventUnfollow = "unfollow"
	EventBlack    = "black"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
	dbgLogger *log.Logger
)

func InitLogger(apiLog, excLog, dbgLog *log.Logger) {
	apiLogger = apiLog
	excLogger = excLog
	dbgLogger = dbgLog
}
