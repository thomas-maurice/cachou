package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(opts *redis.Options) *RedisStorage {
	return &RedisStorage{
		client: redis.NewClient(opts),
	}
}

func (s *RedisStorage) key(objType, objUID []byte) string {
	return fmt.Sprint("%s#%s", string(objType), string(objUID))
}

func (s *RedisStorage) Put(objType []byte, objUID []byte, objData []byte) error {
	return s.client.Set(context.Background(), string(s.key(objType, objUID)), objData, 0).Err()
}

func (s *RedisStorage) Get(objType []byte, objUID []byte) ([]byte, error) {
	res := s.client.Get(context.Background(), string(s.key(objType, objUID)))
	if res.Err() != nil {
		return nil, res.Err()
	}

	b, err := res.Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return b, err
}

func (s *RedisStorage) Del(objType []byte, objUID []byte) error {
	return s.client.Del(context.Background(), string(s.key(objType, objUID))).Err()
}
