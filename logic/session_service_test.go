package logic

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"os"
	"session-server/repo/cache"
	"testing"
)

var redisDao *cache.RedisDao

func setup() {
	// embedded cache
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{mr.Addr()}})
	redisDao = cache.NewRedisDao(client)
}
func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	// 运行测试
	exitCode := m.Run()
	// 退出测试
	teardown()
	os.Exit(exitCode)
}

func TestCreate(t *testing.T) {
	service := NewSessionService(redisDao)
	sessionId, err := service.Create(context.Background(), 1800, map[string][]byte{
		"uid": []byte("1000"),
	})
	fmt.Println("sessionId:", sessionId)
	assert.NoError(t, err)
}
