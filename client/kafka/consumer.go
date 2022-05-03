package kafka

import (
	"zzlove/global"

	"github.com/Shopify/sarama"
)

func NewKafkaConsumer(broker string) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		global.ExcLogger.Println("consumer connect err", err)
		return nil, err
	}
	return consumer, nil
}
