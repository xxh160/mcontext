package service

import (
	"encoding/json"
	"mcontext/internal/model"
	"mcontext/internal/repo"
	"os"
	"strconv"
	"strings"
)

func InitMemory(topic string, role int, question string) (model.DebateMemory, error) {
	debateTagStr, _ := repo.Rdb.Get(repo.Ctx, "NextDebateTag").Result()
	debateTag, _ := strconv.Atoi(debateTagStr)
	repo.Rdb.Incr(repo.Ctx, "NextDebateTag")

	var point string
	if role == 1 {
		point, _ = repo.Rdb.HGet(repo.Ctx, "TopicData:"+topic, "ProPoint").Result()
	} else {
		point, _ = repo.Rdb.HGet(repo.Ctx, "TopicData:"+topic, "ConPoint").Result()
	}

	firstDialog := model.Dialog{Question: question, Answer: point}
	repo.Rdb.Set(repo.Ctx, "DebateMemory:base:"+strconv.Itoa(debateTag), topic+"|"+strconv.Itoa(role), 0)
	repo.Rdb.RPush(repo.Ctx, "DebateMemory:dialogs:"+strconv.Itoa(debateTag), string(jsonMarshal(firstDialog)))
	repo.Rdb.SAdd(repo.Ctx, "ActiveDebateMemory", strconv.Itoa(debateTag))

	debateMemory := model.DebateMemory{
		DebateTag: debateTag,
		Topic:     topic,
		Role:      role,
		Point:     point,
		Dialogs:   []model.Dialog{firstDialog},
	}
	return debateMemory, nil
}

func GetMemory(debateTag int) (model.DebateMemory, error) {
	baseData, _ := repo.Rdb.Get(repo.Ctx, "DebateMemory:base:"+strconv.Itoa(debateTag)).Result()
	parts := strings.Split(baseData, "|")
	topic := parts[0]
	role, _ := strconv.Atoi(parts[1])

	var point string
	if role == 1 {
		point, _ = repo.Rdb.HGet(repo.Ctx, "TopicData:"+topic, "ProPoint").Result()
	} else {
		point, _ = repo.Rdb.HGet(repo.Ctx, "TopicData:"+topic, "ConPoint").Result()
	}

	dialogsData, _ := repo.Rdb.LRange(repo.Ctx, "DebateMemory:dialogs:"+strconv.Itoa(debateTag), 0, -1).Result()
	dialogs := make([]model.Dialog, len(dialogsData))
	for i, d := range dialogsData {
		json.Unmarshal([]byte(d), &dialogs[i])
	}

	debateMemory := model.DebateMemory{
		DebateTag: debateTag,
		Topic:     topic,
		Role:      role,
		Point:     point,
		Dialogs:   dialogs,
	}
	return debateMemory, nil
}

func UpdateMemory(debateTag int, dialog model.Dialog, last bool) error {
	repo.Rdb.RPush(repo.Ctx, "DebateMemory:dialogs:"+strconv.Itoa(debateTag), string(jsonMarshal(dialog)))

	if last {
		go saveDebateMemoryToFile(strconv.Itoa(debateTag))
	}
	return nil
}

func saveDebateMemoryToFile(debateTag string) {
	baseData, _ := repo.Rdb.Get(repo.Ctx, "DebateMemory:base:"+debateTag).Result()
	parts := strings.Split(baseData, "|")
	topic := parts[0]
	role, _ := strconv.Atoi(parts[1])

	var point string
	if role == 1 {
		point, _ = repo.Rdb.HGet(repo.Ctx, "TopicData:"+topic, "ProPoint").Result()
	} else {
		point, _ = repo.Rdb.HGet(repo.Ctx, "TopicData:"+topic, "ConPoint").Result()
	}

	dialogsData, _ := repo.Rdb.LRange(repo.Ctx, "DebateMemory:dialogs:"+debateTag, 0, -1).Result()
	dialogs := make([]model.Dialog, len(dialogsData))
	for i, d := range dialogsData {
		json.Unmarshal([]byte(d), &dialogs[i])
	}

	debateTagInt, _ := strconv.Atoi(debateTag)

	debateMemory := model.DebateMemory{
		DebateTag: debateTagInt,
		Topic:     topic,
		Role:      role,
		Point:     point,
		Dialogs:   dialogs,
	}

	file, _ := json.MarshalIndent(debateMemory, "", " ")
	_ = os.WriteFile("data/debate_memory_"+debateTag+".json", file, 0644)

	repo.Rdb.Del(repo.Ctx, "DebateMemory:base:"+debateTag)
	repo.Rdb.Del(repo.Ctx, "DebateMemory:dialogs:"+debateTag)
	repo.Rdb.SRem(repo.Ctx, "ActiveDebateMemory", debateTag)
}

func jsonMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
