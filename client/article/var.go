package article

import (
	"log"
	"zzlove/rpc/article"
)

var (
	client article_svc.ArticleService

	infoLogger  *log.Logger
	excLogger   *log.Logger
	debugLogger *log.Logger
)

func InitLogger(infoLog, excLog, dbgLog *log.Logger) {
	infoLogger = infoLog
	excLogger = excLog
	debugLogger = dbgLog
}
