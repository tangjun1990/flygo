package kcronlock

import (
	"context"
	"sync"
	"time"

	kredis2 "github.com/tangjun1990/flygo/component/client/kredis"

	"github.com/tangjun1990/flygo/core/klog"
)

type redisLock struct {
	mutex  sync.RWMutex
	client *kredis2.Component
	key    string
	locker *kredis2.Lock
	logger *klog.Component
}

func newRedisLock(client *kredis2.Component, key string, logger *klog.Component) *redisLock {
	return &redisLock{
		mutex:  sync.RWMutex{},
		client: client,
		key:    key,
		locker: nil,
		logger: logger,
	}
}

func (c *redisLock) Lock(ctx context.Context, ttl time.Duration) error {
	locker := c.client.LockClient()
	lock, err := locker.Obtain(ctx, c.key, ttl, kredis2.WithLockOptionRetryStrategy(kredis2.LinearBackoffRetry(ttl)))
	if err != nil {
		return err
	}
	c.mutex.Lock()
	c.locker = lock
	c.mutex.Unlock()
	return nil
}

func (c *redisLock) Unlock(ctx context.Context) error {
	c.mutex.RLock()
	locker := c.locker
	c.mutex.RUnlock()
	if locker == nil {
		return nil
	}

	err := c.locker.Release(ctx)
	if err != nil {
		c.logger.WithCtx(ctx).Warn("cron unlock warning", klog.FieldErr(err))
	}
	return nil
}

func (c *redisLock) Refresh(ctx context.Context, ttl time.Duration) error {
	c.mutex.RLock()
	locker := c.locker
	c.mutex.RUnlock()
	if locker == nil {
		return nil
	}

	return locker.Refresh(ctx, ttl)
}
