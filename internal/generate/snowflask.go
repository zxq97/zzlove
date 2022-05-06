package generate

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	err  error
)

func InitSnowFlask() error {
	node, err = snowflake.NewNode(time.Now().UnixNano() % 1024)
	return err
}

func SnowFlask() int64 {
	return node.Generate().Int64()
}
