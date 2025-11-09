package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

type CacheData struct {
	MessageID string    `json:"message_id"`
	SentAt    time.Time `json:"sent_at"`
}

func NewCache(host, port string) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Cache{client: client}, nil
}

func (c *Cache) Set(ctx context.Context, id int64, msgID string, sentAt time.Time) error {
	data := CacheData{MessageID: msgID, SentAt: sentAt}

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("sent_message:%d", id)
	return c.client.Set(ctx, key, json, 7*24*time.Hour).Err()
}

func (c *Cache) Get(ctx context.Context, id int64) (*CacheData, error) {
	key := fmt.Sprintf("sent_message:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err

	}
	var data CacheData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}
func (c *Cache) Close() error {
	return c.client.Close()
}
