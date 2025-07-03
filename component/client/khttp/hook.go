package khttp

import (
	"fmt"
	"strings"
	"time"

	"git.4321.sh/feige/flygo/core/klog"
)

// 拦截器
type handler func(*KResponse)

type Hook func(c *Config, logger *klog.Component) func(next handler) handler

// 日志拦截器
func logHook(c *Config, logger *klog.Component) func(next handler) handler {
	return func(next handler) handler {
		return func(kr *KResponse) {
			fields := make([]klog.Field, 0)
			start := time.Now()
			next(kr)

			cost := time.Since(start)

			agent := ""
			if values, ok := kr.Request().Header["User-Agent"]; ok && len(values) > 0 {
				agent = values[0]
			}

			fields = append(fields,
				klog.FieldClientKind(),
				klog.FieldHttpMethod(kr.Req.Method),
				klog.FieldDuration(cost), // 耗时
				klog.FieldNetPeerIp(strings.Split(kr.Request().RemoteAddr, ":")[0]), // 请求方式
				klog.FieldHttpMethod(kr.Request().Method),                           // 请求方式
				klog.FieldHttpHost(kr.Request().Host),                               // 实际请求域名
				klog.FieldHttpPath(kr.Request().URL.Path),                           // 地址
				klog.FieldHttpTarget(kr.Request().RequestURI),                       // target
				klog.FieldHttpUserAgent(agent),                                      // user agent
				klog.FieldHttpStatusCode(kr.StatusCode()),
			)

			// request 日志
			if c.HookReq {
				for k, v := range kr.Req.Header {
					fields = append(fields, klog.Any("http.request.header."+strings.ReplaceAll(strings.ToLower(k), "-", "_"), strings.Join(v, ",")))
				}
				if kr.ReqBody != nil {
					fields = append(fields, klog.String("http.request.body", string(kr.ReqBody)))
				}
			}

			// response 日志
			if c.HookRsp {
				// 返回过长，则不记录
				bodyStr := kr.String()
				size := len(bodyStr)
				if size > c.MaxResContentSize {
					bodyStr = bodyStr[:c.MaxResContentSize] + " ..."
				}

				for k, v := range kr.ResHeader() {
					fields = append(fields, klog.Any("http.response.header."+strings.ReplaceAll(strings.ToLower(k), "-", "_"), strings.Join(v, ",")))
				}

				fields = append(fields,
					klog.Int("http.response_content_length", size),
					klog.String("http.response.body", bodyStr),
				)
			}

			uri := strings.Split(kr.Request().URL.String(), "?")[0]

			// 慢接口日志
			if c.SlowLogThreshold > time.Duration(0) && c.SlowLogThreshold < cost {
				logger.WithCtx(kr.Request().Context()).Warn(fmt.Sprintf("[khttp][%s]request slow", uri), fields...)
			}

			if kr.Error() != nil {
				fields = append(fields, klog.FieldErr(kr.Error()))
				logger.WithCtx(kr.Request().Context()).Error(fmt.Sprintf("[khttp][%s]request error : %s", uri, kr.Error()), fields...)
			} else {
				logger.WithCtx(kr.Request().Context()).Info(fmt.Sprintf("[khttp][%s]request success", uri), fields...)
			}

			return
		}
	}
}
