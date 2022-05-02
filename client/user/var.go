package user

import (
	"log"
	"zzlove/rpc/user"
)

var (
	client user_svc.UserService

	infoLogger  *log.Logger
	excLogger   *log.Logger
	debugLogger *log.Logger
)

func InitLogger(infoLog, excLog, dbgLog *log.Logger) {
	infoLogger = infoLog
	excLogger = excLog
	debugLogger = dbgLog
}
