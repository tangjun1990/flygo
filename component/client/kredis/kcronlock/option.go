package kcronlock

import (
	"github.com/tangjun1990/flygo/component/client/kredis"
)

func WithClient(client *kredis.Component) Option {
	return func(c *Container) {
		c.client = client
	}
}

func WithPrefix(prefix string) Option {
	return func(c *Container) {
		c.config.Prefix = prefix
	}
}
