package service

import (
	"encoding/json"
	"mcontext/internal/model"
	"mcontext/internal/repo"
	"mcontext/internal/state"
	"os"
)

func WarmUp() error {
	topicDataFile, err := os.ReadFile("data/topic_datas.json")
	if err != nil {
		return err
	}

	var topicDatas []model.TopicData
	err = json.Unmarshal(topicDataFile, &topicDatas)
	if err != nil {
		return err
	}

	for _, topicData := range topicDatas {
		repo.Rdb.HSet(repo.Ctx, "TopicData:"+topicData.Topic, map[string]interface{}{
			"TopicID":  topicData.TopicID,
			"Topic":    topicData.Topic,
			"ProPoint": topicData.ProPoint,
			"ConPoint": topicData.ConPoint,
		})
	}

	roundFile, err := os.ReadFile("data/round")
	if err != nil {
		return err
	}

	repo.Rdb.Set(repo.Ctx, "NextDebateTag", string(roundFile), 0)
	repo.Rdb.SAdd(repo.Ctx, "ActiveDebateMemory", "")

	state.SetAvailable()

	return nil
}

func CoolDown() error {
	state.SetUnavailable()

	debateTags, _ := repo.Rdb.SMembers(repo.Ctx, "ActiveDebateMemory").Result()
	for _, debateTag := range debateTags {
		saveDebateMemoryToFile(debateTag)
	}
	return nil
}
