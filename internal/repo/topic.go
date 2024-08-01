package repo

import (
	"context"
	"mcontext/internal/model"

	"github.com/go-redis/redis/v8"
)

type TopicRepo interface {
	LoadTopicDatas(ctx context.Context, topicDatas []model.TopicData) error
	RemoveTopicDatas(ctx context.Context) error

	// TopicData:{辩题}
	GetPoint(ctx context.Context, topic string, role model.Role) (string, error)
}

type TopicRepoImpl struct {
	rdb *redis.Client
}

// 将固定数据的 topicDatas 加载到 redis 中
func (r *TopicRepoImpl) LoadTopicDatas(ctx context.Context, topicDatas []model.TopicData) error {
	for _, topicData := range topicDatas {
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

func (r *TopicRepoImpl) RemoveTopicDatas(ctx context.Context) error {
	var cursor uint64
	var keys []string
	var err error
	prefix := "TopicData:"

	// 使用 SCAN 命令遍历所有匹配的键
	for {
		keys, cursor, err = r.rdb.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		// 删除匹配的键
		if len(keys) > 0 {
			if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		// 如果 cursor 为 0，表示已经遍历完所有匹配的键
		if cursor == 0 {
			break
		}
	}

	return nil
}

// 根据辩题和立场读取立论
func (r *TopicRepoImpl) GetPoint(ctx context.Context, topic string, role model.Role) (string, error) {
	return r.rdb.HGet(ctx, "TopicData:"+topic, role.String()).Result()
}

func NewTopicRepo(rdb *redis.Client) TopicRepo {
	return &TopicRepoImpl{
		rdb: rdb,
	}
}
