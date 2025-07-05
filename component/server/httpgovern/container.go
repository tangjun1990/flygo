package httpgovern

import (
	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/flygo/core/kflag"
	"github.com/tangjun1990/flygo/core/klog"
)

type Container struct {
	config *Config
	name   string
	err    error
	logger *klog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: klog.FlygoLogger.With(klog.FieldComponent(PackageName)),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	c.logger = c.logger.With(klog.FieldComponentName(key))
	if err := kcfg.UnmarshalKey(key, &c.config); err != nil {
		c.err = err
		return c
	}

	if kflag.String("host") != "" {
		c.config.Host = kflag.String("host")
	}
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.name, c.config, c.logger)
}
