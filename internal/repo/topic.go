package repo

import (
	"context"
	"fmt"
	"mcontext/internal/model"

	"github.com/go-redis/redis/v8"
)

type TopicRepo interface {
	LoadAllTopicData(ctx context.Context, topicDatas []model.TopicData) error
	RemoveAllTopicData(ctx context.Context) error

	// TopicData:{辩题}
	GetPoint(ctx context.Context, topic string, role model.Role) (string, error)
}

type TopicRepoImpl struct {
	rdb *redis.Client
}

// LoadAllTopicData 将所有固定的 topicData 加载到 redis 中
func (r *TopicRepoImpl) LoadAllTopicData(ctx context.Context, topicDataList []model.TopicData) error {
	for _, topicData := range topicDataList {
		if err := r.rdb.HSet(ctx, "TopicData:"+topicData.Topic, map[string]interface{}{
			"TopicID":  topicData.TopicID,
			"Topic":    topicData.Topic,
			"ProPoint": topicData.ProPoint,
			"ConPoint": topicData.ConPoint,
		}).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (r *TopicRepoImpl) RemoveAllTopicData(ctx context.Context) error {
	var cursor uint64
	var keys []string
	var err error
	prefix := "TopicData:"

	// 使用 SCAN 命令遍历所有匹配的键
	for {
		keys, cursor, err = r.rdb.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		// 删除匹配的键
		if len(keys) > 0 {
			if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete keys: %w", err)
			}
		}

		// 如果 cursor 为 0，表示已经遍历完所有匹配的键
		if cursor == 0 {
			break
		}
	}

	return nil
}

// GetPoint 根据辩题和立场读取立论
func (r *TopicRepoImpl) GetPoint(ctx context.Context, topic string, role model.Role) (string, error) {
	return r.rdb.HGet(ctx, "TopicData:"+topic, role.CacheName()).Result()
}

func NewTopicRepo(rdb *redis.Client) TopicRepo {
	return &TopicRepoImpl{
		rdb: rdb,
	}
}
