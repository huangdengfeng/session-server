package logic

import (
	"context"
	"session-server/entity/errs"
	"session-server/entity/session"
	"session-server/repo/cache"
	"slices"
	"strconv"
	"time"
)

type SessionService struct {
	redisDao *cache.RedisDao
}

const (
	MaxInactiveIntervalKey = "_maxInactiveInterval"
	CreateTimeKey          = "_createTime"
	LastAccessTimeKey      = "_lastAccessTimeKey"
)

var InternalKey = []string{MaxInactiveIntervalKey, CreateTimeKey, LastAccessTimeKey}

func NewSessionService(redisDao *cache.RedisDao) *SessionService {
	return &SessionService{redisDao: redisDao}
}

// Create 创建session
//
//	@Description: 创建session
//	@param ctx
//	@param maxInactiveInterval 超时时间单位秒
//	@param attributes 属性，可以为空
//	@return string sessionId
//	@return error
func (s *SessionService) Create(ctx context.Context, maxInactiveInterval int64, attributes map[string][]byte) (string, error) {
	if attributes == nil {
		attributes = make(map[string][]byte)
	}
	for k, _ := range attributes {
		if slices.Contains(InternalKey, k) {
			return "", errs.AttrKeyLimit.Newf(k)
		}
	}

	attributes[MaxInactiveIntervalKey] = []byte(strconv.FormatInt(maxInactiveInterval, 10))
	attributes[CreateTimeKey] = []byte(strconv.FormatInt(time.Now().Unix(), 10))
	fv := make([]any, 0, len(attributes)*2+len(InternalKey))
	for k, v := range attributes {
		fv = append(fv, k, v)
	}
	sessionId := session.New()
	if err := s.redisDao.Hset(ctx, sessionId, fv...); err != nil {
		return "", err
	}
	if err := s.redisDao.Expire(ctx, sessionId, time.Duration(maxInactiveInterval)*time.Second); err != nil {
		return "", err
	}
	return sessionId, nil
}

func (s *SessionService) Invalidate(ctx context.Context, sessionId string) error {
	return s.redisDao.Del(ctx, sessionId)
}

func (s *SessionService) GetAttribute(ctx context.Context, sessionId string, field string) ([]byte, error) {
	maxInactiveInterval, err := s.GetMaxInactiveInterval(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	value, err := s.redisDao.Hget(ctx, sessionId, field)
	if err != nil {
		return nil, err
	}
	if err := s.redisDao.Expire(ctx, sessionId, time.Duration(maxInactiveInterval)*time.Second); err != nil {
		return nil, err
	}
	return value, nil
}

func (s *SessionService) SetAttribute(ctx context.Context, sessionId string, field string, value []byte) error {
	maxInactiveInterval, err := s.GetMaxInactiveInterval(ctx, sessionId)
	if err != nil {
		return err
	}
	if err := s.redisDao.Hset(ctx, sessionId, field, value); err != nil {
		return err
	}
	if err := s.redisDao.Expire(ctx, sessionId, time.Duration(maxInactiveInterval)*time.Second); err != nil {
		return err
	}
	return nil
}

func (s *SessionService) RemoveAttribute(ctx context.Context, sessionId string, field string) error {
	maxInactiveInterval, err := s.GetMaxInactiveInterval(ctx, sessionId)
	if err != nil {
		return err
	}
	if err := s.redisDao.Hdel(ctx, sessionId, field); err != nil {
		return err
	}
	if err := s.redisDao.Expire(ctx, sessionId, time.Duration(maxInactiveInterval)*time.Second); err != nil {
		return err
	}
	return nil
}

func (s *SessionService) GetMaxInactiveInterval(ctx context.Context, sessionId string) (int64, error) {
	maxInactiveInterval, err := s.redisDao.Hget(ctx, sessionId, MaxInactiveIntervalKey)
	if err != nil {
		return 0, err
	}
	if len(maxInactiveInterval) == 0 {
		return 0, errs.SessionInvalid
	}
	expireSeconds, err := strconv.Atoi(string(maxInactiveInterval))
	if err != nil {
		return 0, errs.RedisError.Newf(err)
	}
	return int64(expireSeconds), nil
}

func (s *SessionService) GetAllAttribute(ctx context.Context, sessionId string) (map[string][]byte, error) {
	maxInactiveInterval, err := s.GetMaxInactiveInterval(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	value, err := s.redisDao.HgetAll(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if err := s.redisDao.Expire(ctx, sessionId, time.Duration(maxInactiveInterval)*time.Second); err != nil {
		return nil, err
	}
	for _, v := range InternalKey {
		delete(value, v)
	}
	return value, nil
}
