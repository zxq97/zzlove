package generate

import (
	"time"
	"github.com/sony/sonyflake"
)

var (
	sony *sonyflake.Sonyflake
)

func InitSonyFlake() {
	t, _ := time.Parse("2006-01-02 15:04:05","2022-01-01 00:00:00")
	set := sonyflake.Settings{
		StartTime: t,
	}
	sony = sonyflake.NewSonyflake(set)
}

func SonyFake() int64 {
	id, _ := sony.NextID()
	return int64(id)
}
