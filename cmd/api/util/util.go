package util

import (
	"encoding/json"
	"zzlove/global"
	"zzlove/internal/cast"
	"zzlove/internal/concurrent"
	"zzlove/internal/kafka"
)

func SendMessage(key int64, msg *kafka.KafkaMessage) {
	b, err := json.Marshal(msg)
	if err != nil {
		global.ExcLogger.Printf("sendMessage json msg %#v err %v", msg, err)
		return
	}
	concurrent.Go(func() {
		err = kafka.SendMessage(kafka.UserActionTopic, []byte(cast.FormatInt(key)), b)
		if err != nil {
			global.ExcLogger.Printf("sendMessage msg %#v err %v", msg, err)
		}
	})
}
