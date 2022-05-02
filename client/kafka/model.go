package kafka

type InfoMessage struct {
	UID       int64 `json:"uid"`
	ToUID     int64 `json:"to_uid"`
	FType     int32 `json:"f_type"`
	ArticleID int64 `json:"article_id"`
}

type ExtraMessage struct {
}

type KafkaMessage struct {
	Event string       `json:"event"`
	Info  InfoMessage  `json:"info"`
	Extra ExtraMessage `json:"extra"`
}
