package comment

import "log"

var (
	infoLogger  *log.Logger
	excLogger   *log.Logger
	debugLogger *log.Logger
)

func InitLogger(infoLog, excLog, dbgLog *log.Logger) {
	infoLogger = infoLog
	excLogger = excLog
	debugLogger = dbgLog
}
