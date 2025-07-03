package kcronlock

import (
	"git.4321.sh/feige/flygo/component/client/kredis"
	"git.4321.sh/feige/flygo/core/kcfg"
	"git.4321.sh/feige/flygo/core/klog"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	logger *klog.Component
	client *kredis.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger.With(klog.FieldComponent(kredis.PackageName)),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := kcfg.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse Config error", klog.FieldErr(err), klog.FieldKey(key))
		return c
	}

	c.logger = c.logger.With(klog.FieldComponentName(key))
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	if c.client == nil {
		c.logger.Panic("client redis nil", klog.FieldKey("use WithClient method"))
	}
	return newComponent(c.name, c.config, c.logger, c.client)
}
