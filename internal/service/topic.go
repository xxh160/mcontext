package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mcontext/internal/conf"
	"mcontext/internal/model"
	"mcontext/internal/repo"
	"os"
)

type TopicService interface {
	LoadAllTopicData(ctx context.Context) error
	RemoveAllTopicData(ctx context.Context) error
	GetPoint(ctx context.Context, topic string, role model.Role) (string, error)
}

type TopicServiceImpl struct {
	repo repo.TopicRepo
}

func (s *TopicServiceImpl) LoadAllTopicData(ctx context.Context) error {
	log.Printf("TopicService loading...\n")

	fileContent, err := os.ReadFile(conf.TopicDataPath)
	if err != nil {
		return fmt.Errorf("read topicData error: %w", err)
	}

	var topicDataList []model.TopicData
	// 转化为 TopicData
	err = json.Unmarshal(fileContent, &topicDataList)
	if err != nil {
		return fmt.Errorf("unmarshal topicData error: %w", err)
	}

	return s.repo.LoadAllTopicData(ctx, topicDataList)
}

func (s *TopicServiceImpl) RemoveAllTopicData(ctx context.Context) error {
	log.Printf("TopicService removing...\n")
	return s.repo.RemoveAllTopicData(ctx)
}

func (s *TopicServiceImpl) GetPoint(ctx context.Context, topic string, role model.Role) (string, error) {
	return s.repo.GetPoint(ctx, topic, role)
}

func NewTopicService(repo repo.TopicRepo) TopicService {
	return &TopicServiceImpl{
		repo: repo,
	}
}
