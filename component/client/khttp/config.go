package khttp

import (
	"git.4321.sh/feige/flygo/core/utils/xtime"
	"time"

	"git.4321.sh/feige/flygo/core/kapp"
)

type Config struct {
	ClientName    string        // 本服务名称
	ServiceName   string        // 第三方服务名称
	Host          string        // 第三方 host
	Timeout       time.Duration // 超时时间 （秒）
	Retry         int           // 重试次数
	RetryInterval time.Duration // 重试时间间隔
	Mock          bool          // 是否开启mock
	MockGroupId   int           // mock 分组Id

	HookLog bool // 是否记录log
	HookReq bool // 是否记录请求参数
	HookRsp bool // 是否记录响应参数

	SlowLogThreshold  time.Duration // 服务慢日志，默认500ms
	MaxResContentSize int           //最大响应结果大小

	// 自定义字段
	Ext   map[string]string // 补充字段，用于自定义接口鉴权参数
	hooks []Hook            // 拦截器
}

func DefaultConfig() *Config {
	return &Config{
		ClientName:        kapp.Name(),
		Timeout:           0,
		Retry:             0,
		RetryInterval:     time.Millisecond * 100,
		HookLog:           true,
		HookReq:           true,
		HookRsp:           true,
		SlowLogThreshold:  xtime.Duration("500ms"),
		MaxResContentSize: 2048,
	}
}
