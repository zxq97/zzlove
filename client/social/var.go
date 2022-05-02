package social

import (
	"log"
	"zzlove/rpc/social"
)

var (
	client social_svc.SocialService

	infoLogger  *log.Logger
	excLogger   *log.Logger
	debugLogger *log.Logger
)

func InitLogger(infoLog, excLog, dbgLog *log.Logger) {
	infoLogger = infoLog
	excLogger = excLog
	debugLogger = dbgLog
}
