package model

type TopicData struct {
	TopicID  int    `json:"topicID"`
	Topic    string `json:"topic"`
	ProPoint string `json:"proPoint"`
	ConPoint string `json:"conPoint"`
}
