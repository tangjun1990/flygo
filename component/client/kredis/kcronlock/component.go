package kcronlock

import (
	"sync"

	"github.com/tangjun1990/flygo/component/client/kredis"
	"github.com/tangjun1990/flygo/component/task/kcron"

	"github.com/tangjun1990/flygo/core/klog"
)

type Component struct {
	name   string
	config *Config
	logger *klog.Component
	client *kredis.Component
	mutuex sync.RWMutex
}

func newComponent(name string, config *Config, logger *klog.Component, client *kredis.Component) *Component {
	reg := &Component{
		name:   name,
		logger: logger,
		config: config,
		client: client,
	}
	return reg
}

func (c *Component) NewLock(key string) kcron.Lock {
	return newRedisLock(c.client, c.config.Prefix+key, c.logger)
}
