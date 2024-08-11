package conf

import "path/filepath"

const (
	DataDir       = "data"
	RoundName     = "round"
	RoundPath     = DataDir + string(filepath.Separator) + RoundName
	TopicDataName = "topic_data.json"
	TopicDataPath = DataDir + string(filepath.Separator) + TopicDataName
	RedisMachine  = "127.0.0.1"
	RedisPort     = ":6379"
	RedisAddr     = RedisMachine + RedisPort
	ServerMachine = "0.0.0.0"
	ServerPort    = ":8080"
	ServerAddr    = ServerMachine + ServerPort
)
