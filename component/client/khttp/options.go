package khttp

import (
	"git.4321.sh/feige/flygo/core/klog"
)

type Option func(c *Container)

func WithName(name string) Option {
	return func(c *Container) {
		c.name = name
	}
}

func WithConfig(conf *Config) Option {
	return func(c *Container) {
		c.config = conf
	}
}

func WithLogger(logger klog.Component) Option {
	return func(c *Container) {
		c.logger = logger.With(klog.FieldComponentName(PackageName))
	}
}

func WithHook(is ...Hook) Option {
	return func(c *Container) {
		if c.config.hooks == nil {
			c.config.hooks = make([]Hook, 0)
		}
		c.config.hooks = append(c.config.hooks, is...)
	}
}
