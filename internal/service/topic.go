package service

import (
	"context"
	"encoding/json"
	"mcontext/internal/conf"
	"mcontext/internal/model"
	"mcontext/internal/repo"
	"os"
)

type TopicService interface {
	LoadTopicDatas(ctx context.Context) error
	RemoveTopicDatas(ctx context.Context) error
	GetPoint(ctx context.Context, topic string, role model.Role) (string, error)
}

type TopicServiceImpl struct {
	repo repo.TopicRepo
}

func (s *TopicServiceImpl) LoadTopicDatas(ctx context.Context) error {
	fileContent, err := os.ReadFile(conf.TopicDatasPath)
	if err != nil {
		return err
	}

	var topicDatas []model.TopicData
	// 转化为 TopicData
	err = json.Unmarshal(fileContent, &topicDatas)
	if err != nil {
		return err
	}

	return s.repo.LoadTopicDatas(ctx, topicDatas)
}

func (s *TopicServiceImpl) RemoveTopicDatas(ctx context.Context) error {
	return s.repo.RemoveTopicDatas(ctx)
}

func (s *TopicServiceImpl) GetPoint(ctx context.Context, topic string, role model.Role) (string, error) {
	return s.repo.GetPoint(ctx, topic, role)
}

func NewTopicservice(repo repo.TopicRepo) TopicService {
	return &TopicServiceImpl{
		repo: repo,
	}
}
