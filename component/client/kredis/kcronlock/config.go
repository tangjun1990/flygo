package kcronlock

import (
	"github.com/tangjun1990/flygo/core/kapp"
)

type Config struct {
	Prefix string
}

func DefaultConfig() *Config {
	return &Config{
		Prefix: "kl:" + kapp.Name() + ":", //前缀默认kl:appname:
	}
}
