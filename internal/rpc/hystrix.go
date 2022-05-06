package rpc

import (
	"github.com/afex/hystrix-go/hystrix"
)

// fixme 之后替换
func initBreaker(name string) {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		RequestVolumeThreshold: 5,
		MaxConcurrentRequests:  10,
		Timeout:                3000,
		SleepWindow:            10000,
		ErrorPercentThreshold:  20,
	})
}
