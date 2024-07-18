package cache

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var redisDao *RedisDao

func setup() {
	// embedded cache
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{mr.Addr()}})
	redisDao = &RedisDao{client}
}
func teardown() {
	redisDao.redisClient.Shutdown(context.Background())
}

func TestMain(m *testing.M) {
	setup()
	// 运行测试
	exitCode := m.Run()
	// 退出测试
	teardown()
	os.Exit(exitCode)
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	key := "k1"
	f1 := "f1"
	v1 := []byte("v1")
	// 删除不存在的key
	err := redisDao.Del(ctx, key)
	assert.NoError(t, err)
	// 删除存在的key
	err = redisDao.Hset(ctx, key, f1, v1)
	assert.NoError(t, err)
	hget, err := redisDao.Hget(ctx, key, f1)
	assert.NoError(t, err)
	assert.Equal(t, hget, v1)
	err = redisDao.Del(ctx, key)
	assert.NoError(t, err)
	hget, err = redisDao.Hget(ctx, key, f1)
	assert.NoError(t, err)
	assert.Empty(t, hget)
}
