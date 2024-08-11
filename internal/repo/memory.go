package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mcontext/internal/model"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type MemoryRepo interface {
	SetNextDebateTag(ctx context.Context, tag int) error
	GetNextDebateTag(ctx context.Context) (int, error)
	IncrGetNextDebateTag(ctx context.Context) (int, error)
	RemoveNextDebateTag(ctx context.Context) error

	AddActiveDebateMemoryTag(ctx context.Context, tag int) error
	RemoveActiveDebateMemoryTag(ctx context.Context, tag int) error
	GetActiveDebateMemoryTags(ctx context.Context) ([]int, error)
	RemoveActiveDebateMemoryTags(ctx context.Context) error
	IsInActiveDebateMemoryTags(ctx context.Context, tag int) (bool, error)

	GetDebateMemory(ctx context.Context, tag int) (*model.DebateMemory, error)

	GetDebateMemoryDialogs(ctx context.Context, tag int) ([]model.Dialog, error)
	AddDebateMemoryDialog(ctx context.Context, tag int, dialog model.Dialog) error
	RemoveDebateMemoryDialog(ctx context.Context, tag int) error

	GetDebateMemoryBase(ctx context.Context, tag int) (*model.DebateMemory, error)
	SetDebateMemoryBase(ctx context.Context, tag int, debateMemory *model.DebateMemory) error
	RemoveDebateMemoryBase(ctx context.Context, tag int) error
}

type MemoryRepoImpl struct {
	rdb       *redis.Client
	topicRepo TopicRepo
}

// SetNextDebateTag 设置 NextDebateTag 的值
func (r *MemoryRepoImpl) SetNextDebateTag(ctx context.Context, tag int) error {
	return r.rdb.Set(ctx, "NextDebateTag", tag, 0).Err()
}

// GetNextDebateTag 获取 NextDebateTag 的值
func (r *MemoryRepoImpl) GetNextDebateTag(ctx context.Context) (int, error) {
	res, err := r.rdb.Get(ctx, "NextDebateTag").Result()
	if err != nil {
		return -1, fmt.Errorf("get next debate tag failed: %w", err)
	}

	return strconv.Atoi(res)
}

// IncrGetNextDebateTag 原子获取并加一 NextDebateTag 值，返回加一后的数值
func (r *MemoryRepoImpl) IncrGetNextDebateTag(ctx context.Context) (int, error) {
	res, err := r.rdb.Incr(ctx, "NextDebateTag").Result()
	return int(res), err
}

// RemoveNextDebateTag 删除 NextDebateTag
func (r *MemoryRepoImpl) RemoveNextDebateTag(ctx context.Context) error {
	return r.rdb.Del(ctx, "NextDebateTag").Err()
}

// AddActiveDebateMemoryTag 添加 tag 这个元素到 ActiveDebateMemoryTags set 中
func (r *MemoryRepoImpl) AddActiveDebateMemoryTag(ctx context.Context, tag int) error {
	return r.rdb.SAdd(ctx, "ActiveDebateMemoryTags", tag).Err()
}

// RemoveActiveDebateMemoryTag 从 ActiveDebateMemoryTags set 中删除 tag 这个元素
func (r *MemoryRepoImpl) RemoveActiveDebateMemoryTag(ctx context.Context, tag int) error {
	return r.rdb.SRem(ctx, "ActiveDebateMemoryTags", tag).Err()
}

// GetActiveDebateMemoryTags 获取 ActiveDebateMemoryTags 这个 set 中所有的 tag
func (r *MemoryRepoImpl) GetActiveDebateMemoryTags(ctx context.Context) ([]int, error) {
	tags, err := r.rdb.SMembers(ctx, "ActiveDebateMemoryTags").Result()
	if err != nil {
		return nil, fmt.Errorf("get active debate memory tags failed: %w", err)
	}

	tagIntList := make([]int, len(tags))
	for i, tag := range tags {
		num, err := strconv.Atoi(tag)
		if err != nil {
			return nil, fmt.Errorf("parse active debate memory tags failed: %w", err)
		}
		tagIntList[i] = num
	}

	return tagIntList, nil
}

// RemoveActiveDebateMemoryTags 删除 ActiveDebateMemoryTags 这个 set
func (r *MemoryRepoImpl) RemoveActiveDebateMemoryTags(ctx context.Context) error {
	return r.rdb.Del(ctx, "ActiveDebateMemoryTags").Err()
}

func (r *MemoryRepoImpl) IsInActiveDebateMemoryTags(ctx context.Context, tag int) (bool, error) {
	return r.rdb.SIsMember(ctx, "ActiveDebateMemoryTags", tag).Result()
}

// GetDebateMemory 根据 tag 读取相关信息，组装成一个 DebateMemory
func (r *MemoryRepoImpl) GetDebateMemory(ctx context.Context, tag int) (*model.DebateMemory, error) {
	debateMemory, err := r.GetDebateMemoryBase(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("get debate memory base failed: %w", err)
	}

	debateMemory.Point, err = r.topicRepo.GetPoint(ctx, debateMemory.Topic, debateMemory.Role)
	if err != nil {
		return nil, fmt.Errorf("get debate memory point failed: %w", err)
	}

	dialogs, err := r.GetDebateMemoryDialogs(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("get debate memory dialogs failed: %w", err)
	}

	debateMemory.Dialogs = dialogs
	return debateMemory, nil
}

// GetDebateMemoryDialogs 读取某个具体的 DebateMemory 的对话上下文
func (r *MemoryRepoImpl) GetDebateMemoryDialogs(ctx context.Context, tag int) ([]model.Dialog, error) {
	dialogsData, err := r.rdb.LRange(ctx, "DebateMemory:dialogs:"+strconv.Itoa(tag), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("get debate memory dialogs failed: %w", err)
	}

	dialogs := make([]model.Dialog, len(dialogsData))
	for i, dialogStr := range dialogsData {
		// 这里的内存拷贝怎么办？
		err = json.Unmarshal([]byte(dialogStr), &dialogs[i])
		if err != nil {
			return nil, err
		}
	}

	return dialogs, nil
}

// AddDebateMemoryDialog 给具体的某个 DebateMemory 增加一段对话
func (r *MemoryRepoImpl) AddDebateMemoryDialog(ctx context.Context, tag int, dialog model.Dialog) error {
	dialogStr, err := json.Marshal(dialog)
	if err != nil {
		return fmt.Errorf("marshal dialog failed: %w", err)
	}

	return r.rdb.RPush(ctx, "DebateMemory:dialogs:"+strconv.Itoa(tag), string(dialogStr)).Err()
}

// RemoveDebateMemoryDialog 删除 DebateMemory 的对话上下文
func (r *MemoryRepoImpl) RemoveDebateMemoryDialog(ctx context.Context, tag int) error {
	key := "DebateMemory:dialogs:" + strconv.Itoa(tag)
	log.Printf("Remove %s\n", key)
	return r.rdb.Del(ctx, key).Err()
}

// GetDebateMemoryBase 获取某个具体的 DebateMemory 的固定部分
func (r *MemoryRepoImpl) GetDebateMemoryBase(ctx context.Context, tag int) (*model.DebateMemory, error) {
	// 读取基础部分：topic、role
	key := "DebateMemory:base:" + strconv.Itoa(tag)
	baseDataStr, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("get debate memory base failed: %w", err)
	}

	debateMemory := model.DebateMemory{}
	err = json.Unmarshal([]byte(baseDataStr), &debateMemory)
	// 结构体中存有 tag
	return &debateMemory, fmt.Errorf("unmarshal base data string err: %w", err)
}

// SetDebateMemoryBase 设置某个具体的 DebateMemory 的固定部分
func (r *MemoryRepoImpl) SetDebateMemoryBase(ctx context.Context, tag int, debateMemory *model.DebateMemory) error {
	debateMemoryStr, err := json.Marshal(debateMemory)
	if err != nil {
		return fmt.Errorf("marshal debate memory base failed: %w", err)
	}

	return r.rdb.Set(ctx, "DebateMemory:base:"+strconv.Itoa(tag), string(debateMemoryStr), 0).Err()
}

// RemoveDebateMemoryBase 删除某个具体的 DebateMemory 的固定部分
func (r *MemoryRepoImpl) RemoveDebateMemoryBase(ctx context.Context, tag int) error {
	key := "DebateMemory:base:" + strconv.Itoa(tag)
	log.Printf("Remove %s\n", key)
	return r.rdb.Del(ctx, key).Err()
}

func NewMemoryRepo(rdb *redis.Client, topicRepo TopicRepo) MemoryRepo {
	return &MemoryRepoImpl{rdb: rdb, topicRepo: topicRepo}
}
