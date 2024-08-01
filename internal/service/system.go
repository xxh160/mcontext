package service

import (
	"context"
	"mcontext/internal/util"
)

type SystemService interface {
	Reset(ctx context.Context) error
}

type SystemServiceImpl struct {
	topicService  TopicService
	memoryService MemoryService
}

// 重启整个系统，先让 memory exit（如果需要），然后重新 load TopicDatas，然后让 memory init
// Reset 执行期间是没有其余请求并发执行的
func (s *SystemServiceImpl) Reset(ctx context.Context) error {
	// 如果当前系统是可用状态，则让所有 DebateMemory 保存并清空
	if util.IsAvailabale() {
		util.SetUnavailable()
		if err := s.memoryService.Exit(ctx); err != nil {
			return err
		}
	}

	// 重新 load TopicDatas
	if err := s.topicService.LoadTopicDatas(ctx); err != nil {
		return err
	}

	// 让 memoryService 可用
	if err := s.memoryService.Init(ctx); err != nil {
		return err
	}

	// 让系统可用
	util.SetAvailable()

	return nil
}

func NewSystemService(topicService TopicService, memoryService MemoryService) SystemService {
	return &SystemServiceImpl{
		topicService:  topicService,
		memoryService: memoryService,
	}
}
