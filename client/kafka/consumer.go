package kafka

import (
	"github.com/Shopify/sarama"
)

func NewKafkaConsumer(broker string) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		excLogger.Println("consumer connect err", err)
		return nil, err
	}
	return consumer, nil
}
