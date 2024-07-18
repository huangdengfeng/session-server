package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"session-server/entity/errs"
	"time"
)

type RedisDao struct {
	redisClient redis.UniversalClient
}

func NewRedisDao(redisClient redis.UniversalClient) *RedisDao {
	return &RedisDao{redisClient: redisClient}
}

func (r *RedisDao) Exists(ctx context.Context, sessionId string) (bool, error) {
	exists, err := r.redisClient.Exists(ctx, sessionId).Result()
	if err != nil {
		log.Errorf("expire error %s", err)
		return false, errs.RedisError.Newf(err)
	}
	return exists == 1, nil
}

func (r *RedisDao) Expire(ctx context.Context, key string, duration time.Duration) error {
	if err := r.redisClient.Expire(ctx, key, duration).Err(); err != nil {
		log.Errorf("expire error %s", err)
		return errs.RedisError.Newf(err)
	}
	return nil
}

func (r *RedisDao) Hset(ctx context.Context, key string, fv ...any) error {
	if len(fv)%2 != 0 {
		return errs.BasArgs.Newf("hset fv must even number")
	}
	if len(fv) < 2 {
		return errs.BasArgs.Newf("hset fv must >= 2")
	}
	if err := r.redisClient.HSet(ctx, key, fv).Err(); err != nil {
		log.Errorf("hset error %s", err)
		return errs.RedisError.Newf(err)
	}
	return nil
}

func (r *RedisDao) Hget(ctx context.Context, key string, field string) ([]byte, error) {
	value, err := r.redisClient.HGet(ctx, key, field).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("hget error %s", err)
		return nil, errs.RedisError.Newf(err)
	}
	return value, nil
}

func (r *RedisDao) HgetAll(ctx context.Context, key string) (map[string][]byte, error) {
	result, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("hgetAll error %s", err)
		return nil, errs.RedisError.Newf(err)
	}
	byteMap := make(map[string][]byte)
	for k, v := range result {
		byteMap[k] = []byte(v)
	}
	return byteMap, nil
}

func (r *RedisDao) Del(ctx context.Context, key string) error {
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		log.Errorf("del error %s", err)
		return errs.RedisError.Newf(err)
	}
	return nil
}

func (r *RedisDao) Hdel(ctx context.Context, key string, field string) error {
	if err := r.redisClient.HDel(ctx, key, field).Err(); err != nil {
		log.Errorf("hdel error %s", err)
		return errs.RedisError.Newf(err)
	}
	return nil
}
