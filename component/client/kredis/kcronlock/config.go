package kcronlock

import (
	"git.4321.sh/feige/flygo/core/kapp"
)

type Config struct {
	Prefix string
}

func DefaultConfig() *Config {
	return &Config{
		Prefix: "kl:" + kapp.Name() + ":", //前缀默认kl:appname:
	}
}
