package storage

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisStorage(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

	redisStorage := NewRedisStorage(&redis.Options{
		Addr: s.Addr(),
	})

	defer redisStorage.Close()

	genericTest(t, redisStorage)
}
