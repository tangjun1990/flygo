package khttp

import (
	"git.4321.sh/feige/flygo/core/kcfg"
	"git.4321.sh/feige/flygo/core/klog"
)

type Container struct {
	name   string
	config *Config
	logger *klog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger,
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := kcfg.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", klog.FieldErr(err), klog.FieldKey(key))
		return c
	}
	c.name = key
	c.config.ServiceName = key
	return c
}

// LoadConfig 手动加载配置
func LoadConfig(name string, config *Config) *Container {
	c := DefaultContainer()
	c.name = name
	c.config = config
	return c
}

func (c *Container) Build(options ...Option) *Component {

	if options == nil {
		options = make([]Option, 0)
	}

	if c.config.HookLog {
		options = append(options, WithHook(logHook))
	}

	for _, option := range options {
		option(c)
	}
	return newComponent(c.config.ServiceName, c.config, c.logger)
}
