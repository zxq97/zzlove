package async

import (
	"encoding/json"
	"time"
	"zzlove/global"
	"zzlove/internal/kafka"

	"github.com/Shopify/sarama"
)

func InitConsumer(broker, topics []string, group string) {
	consumer, err := kafka.NewConsumer(broker, topics, group)
	if err != nil || consumer == nil {
		global.ExcLogger.Printf("InitConsumer broker %v topics %v group %v err %v", broker, topics, group, err)
		panic(err)
	}
	defer consumer.Close()

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				process(msg)
				consumer.MarkOffset(msg, "")
			}
		case err = <-consumer.Errors():
			global.ExcLogger.Println("InitConsumer consumer err", err)
		case nft := <-consumer.Notifications():
			global.ApiLogger.Println("InitConsumer consumer Notifications", nft)
		}
	}
}

func process(msg *sarama.ConsumerMessage) {
	val := msg.Value
	global.ApiLogger.Printf("process msg key %v val %v", string(msg.Key), string(val))
	event := &kafka.KafkaMessage{}
	err := json.Unmarshal(val, event)
	if err != nil {
		global.ExcLogger.Printf("process key %v val %v unmarshal err %v", string(msg.Key), string(val), err)
		return
	}
	now := time.Now()
	ctx, cancel := kafka.ConsumerContext(event.MsgID)
	defer cancel()
	switch event.Event {
	case kafka.EventPublish:
		publishArticle(ctx, event.Info.UID, event.Info.ArticleID)
	case kafka.EventFollow:
		follow(ctx, event.Info.UID, event.Info.ToUID)
	case kafka.EventUnfollow:
		unfollow(ctx, event.Info.UID, event.Info.ToUID)
	case kafka.EventBlack:
		black(ctx, event.Info.UID, event.Info.ToUID)
	}
	global.ApiLogger.Printf("process key %v val %v time %v", string(msg.Key), string(val), time.Since(now))
}
