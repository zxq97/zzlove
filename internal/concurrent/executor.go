package concurrent

import (
	"log"
	"runtime/debug"
	"sync"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
)

func InitLogger(apiLog, excLog *log.Logger) {
	apiLogger = apiLog
	excLogger = excLog
}

type WaitGroup struct {
	wg sync.WaitGroup
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{wg: sync.WaitGroup{}}
}

func (ex *WaitGroup) Run(fn func()) {
	ex.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				excLogger.Println("executor", err, string(debug.Stack()))
			}
			ex.wg.Done()
		}()
		fn()
	}()
}

func (ex *WaitGroup) Wait() {
	ex.wg.Wait()
}

func Go(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				excLogger.Println("Go:", err, "\n", string(debug.Stack()))
			}
		}()
		fn()
	}()
}
