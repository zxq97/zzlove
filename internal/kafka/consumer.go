package kafka

import (
	"time"
	"zzlove/global"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

func NewConsumer(broker, topics []string, group string) (*cluster.Consumer, error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = -1
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Group.Return.Notifications = true

	consumer, err := cluster.NewConsumer(broker, group, topics, config)
	if err != nil {
		global.ExcLogger.Printf("NewConsumer broker %v topics %v group %v err %v", broker, topics, group, err)
		return nil, err
	}
	return consumer, nil
}
