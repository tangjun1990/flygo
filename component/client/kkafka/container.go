package kkafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/flygo/core/klog"
)

type Option func(c *Container)

type Container struct {
	config *config
	name   string
	logger *klog.Component
}

// DefaultContainer ...
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger.With(klog.FieldComponent(PackageName)),
	}
}

// Load ...
func Load(key string) *Container {
	c := DefaultContainer()
	if err := kcfg.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", klog.FieldErr(err), klog.FieldKey(key))
		return c
	}

	c.logger = c.logger.With(klog.FieldComponentName(key))
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	c.logger = c.logger.With(klog.FieldAddr(fmt.Sprintf("%s", c.config.Brokers)))
	return &Component{
		config:    c.config,
		logger:    c.logger,
		consumers: make(map[string]sarama.Consumer),
		producers: make(map[string]sarama.SyncProducer),
	}
}
