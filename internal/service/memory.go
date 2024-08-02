package service

import (
	"context"
	"encoding/json"
	"log"
	"mcontext/internal/conf"
	"mcontext/internal/erro"
	"mcontext/internal/model"
	"mcontext/internal/repo"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MemoryService interface {
	Init(ctx context.Context) error
	Exit(ctx context.Context) error
	CreateMemory(ctx context.Context, topic string, role model.Role, question string) (*model.DebateMemory, error)
	GetMemory(ctx context.Context, debateTag int) (*model.DebateMemory, error)
	UpdateMemory(ctx context.Context, debateTag int, dialog model.Dialog, last bool) error
}

type MemoryServiceImpl struct {
	memoryRepo   repo.MemoryRepo
	topicService TopicService
}

// 初始化 NextDebateTag
func (s *MemoryServiceImpl) Init(ctx context.Context) error {
	log.Printf("MemoryService initing...\n")

	fileContent, err := os.ReadFile(conf.RoundPath)
	if err != nil {
		return err
	}

	content := strings.TrimSpace(string(fileContent))
	num, err := strconv.Atoi(content)
	if err != nil {
		return err
	}

	log.Printf("Round: %d\n", num)

	// 设置 NextDebateTag
	s.memoryRepo.SetNextDebateTag(ctx, num)
	// ActiveDebateMemoryTags set 惰性创建，没必要在这里创建空集合
	return nil
}

// 将 ActiveDebateMemoryTags 所代表的所有 DebateMemory 持久化到文件
// 将 NextDebateTag 数值持久化到文件 round
// 同时从 redis 中删除 NextDebateTag 和 ActiveDebateMemoryTags
func (s *MemoryServiceImpl) Exit(ctx context.Context) error {
	log.Printf("MemoryService exiting...\n")

	tags, err := s.memoryRepo.GetActiveDebateMemoryTags(ctx)
	if err != nil {
		return err
	}

	log.Printf("Active tags: %v\n", tags)

	for _, tag := range tags {
		if err = s.persistDelete(ctx, tag); err != nil {
			log.Printf("PersistDebateMemory error: %s\n", err)
			continue
		}
	}

	// 删除 ActiveDebateMemoryTags
	if err = s.memoryRepo.RemoveActiveDebateMemoryTags(ctx); err != nil {
		return err
	}

	// 读取并持久化 NextDebateTag
	tagNum, err := s.memoryRepo.GetNextDebateTag(ctx)
	if err != nil {
		return err
	}

	// 覆盖写
	if err = os.WriteFile(conf.RoundPath, []byte(strconv.Itoa(tagNum)), 0644); err != nil {
		return err
	}

	// 删除 NextDebateTag
	if err = s.memoryRepo.RemoveNextDebateTag(ctx); err != nil {
		return err
	}

	return nil
}

// 将 debateMemory 持久化到文件，并删除相关的 redis 变量
func (s *MemoryServiceImpl) persistDelete(ctx context.Context, tag int) error {
	log.Printf("Persist and delete: %d\n", tag)

	debateMemory, err := s.memoryRepo.GetDebateMemory(ctx, tag)
	if err != nil {
		return err
	}

	fileContent, err := json.MarshalIndent(debateMemory, "", "  ")
	if err != nil {
		return err
	}

	// 持久化到 data/debate_memory_{DebateTag}.json
	filename := filepath.Join(conf.DataDir, "debate_memory_"+strconv.Itoa(debateMemory.DebateTag)+".json")
	err = os.WriteFile(filename, fileContent, 0644)
	if err != nil {
		return err
	}

	log.Printf("Save to %s\n", filename)

	// 删除 debateMemory 相关的基础信息
	if err = s.memoryRepo.RemoveDebateMemoryBase(ctx, debateMemory.DebateTag); err != nil {
		return err
	}
	// 删除 debateMemory 相关的对话
	if err = s.memoryRepo.RemoveDebateMemoryDialog(ctx, debateMemory.DebateTag); err != nil {
		return err
	}
	// 从 redis ActiveDebateMemoryTags 中删除这个 debateMemoryTag
	if err = s.memoryRepo.RemoveActiveDetabeMemoryTag(ctx, debateMemory.DebateTag); err != nil {
		return err
	}

	return nil
}

func (s *MemoryServiceImpl) CreateMemory(ctx context.Context, topic string, role model.Role, question string) (*model.DebateMemory, error) {
	// 获取新 tag
	newTag, err := s.memoryRepo.IncrGetNextDebateTag(ctx)
	if err != nil {
		return nil, err
	}

	// 构造 DebateMemory 的 base 部分
	debateMemory := model.DebateMemory{
		DebateTag: newTag,
		Topic:     topic,
		Role:      role,
	}

	// 存储 base 部分到 redis
	if err = s.memoryRepo.SetDebateMemoryBase(ctx, newTag, &debateMemory); err != nil {
		return nil, err
	}

	// 根据辩题和立场得到立论
	debateMemory.Point, err = s.topicService.GetPoint(ctx, topic, role)
	if err != nil {
		return nil, err
	}

	// 构造出第一轮对话
	firstDialog := model.Dialog{Question: question, Answer: debateMemory.Point}

	// 存储第一轮对话
	if err = s.memoryRepo.AddDebateMemoryDialog(ctx, newTag, firstDialog); err != nil {
		return nil, err
	}

	debateMemory.Dialogs = []model.Dialog{firstDialog}

	// 将 tag 增加到 ActiveDebateMemoryTags set 中
	if err = s.memoryRepo.AddActiveDebateMemoryTag(ctx, newTag); err != nil {
		return nil, err
	}

	log.Printf("CreateMemory debateTag: %d\n", newTag)

	return &debateMemory, nil
}

func (s *MemoryServiceImpl) GetMemory(ctx context.Context, debateTag int) (*model.DebateMemory, error) {
	log.Printf("GetMemory debateTag: %d\n", debateTag)
	return s.memoryRepo.GetDebateMemory(ctx, debateTag)
}

func (s *MemoryServiceImpl) UpdateMemory(ctx context.Context, debateTag int, dialog model.Dialog, last bool) error {
	log.Printf("UpdateMemory debateTag: %d\n", debateTag)

	// 检查 debateTag 是否在 ActiveDebateMemoryTags set 中
	res, err := s.memoryRepo.IsInActiveDebateMemoryTags(ctx, debateTag)
	if err != nil {
		return err
	}

	// 不在 ActiveDebateMemoryTags set 中
	if !res {
		return &erro.CustomError{Msg: "DebateTag " + strconv.Itoa(debateTag) +
			" is not in ActiveDebateMemoryTags set"}
	}

	// 存储新的对话
	if err := s.memoryRepo.AddDebateMemoryDialog(ctx, debateTag, dialog); err != nil {
		return err
	}

	// 如果是最后一轮对话，则持久化到文件
	if last {
		go func() {
			s.persistDelete(ctx, debateTag)
		}()
	}

	return nil
}

func NewMemoryService(memoryRepo repo.MemoryRepo, topicService TopicService) MemoryService {
	return &MemoryServiceImpl{
		memoryRepo:   memoryRepo,
		topicService: topicService,
	}
}
