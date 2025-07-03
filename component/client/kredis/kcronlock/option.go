package kcronlock

import (
	"git.4321.sh/feige/flygo/component/client/kredis"
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
