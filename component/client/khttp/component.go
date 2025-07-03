package khttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"git.4321.sh/feige/flygo/core/klog"
)

const PackageName = "core.khttp"

type Component struct {
	name       string
	config     *Config
	httpClient *http.Client
	logger     *klog.Component
}

func newComponent(name string, config *Config, logger *klog.Component) *Component {
	h := &http.Client{
		Timeout: config.Timeout,
	}

	return &Component{
		name:       name,
		config:     config,
		logger:     logger.With(klog.FieldComponent(PackageName)).With(klog.FieldComponentName(name)),
		httpClient: h,
	}
}

func (c *Component) ConfigKey() string {
	return c.config.ServiceName
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Start() error {
	return nil
}

func (c *Component) Stop() error {
	return nil
}

// Client 获取 KHttp 组件 Client
func (c *Component) Client() *http.Client {
	return c.httpClient
}

const (
	ContentTypeJson = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

type KHeader map[string]string
type KBody map[string]string

// 从 header 中获取 contentType
func (h KHeader) getContentType() string {
	contentType, ok := h["Content-Type"]
	if ok {
		return contentType
	} else {
		return ContentTypeForm
	}
}

// Get get 请求
func (c *Component) Get(ctx context.Context, path string, vs ...interface{}) (resp *KResponse) {
	header, query := getHeaderBody(vs)
	// 对url参数 和 query 进行合并
	path, err := addQueryValuesIntoPath(path, query)
	if err != nil {
		resp.SetError(err)
		return
	}
	resp = c.getHttpRequest(ctx, http.MethodGet, path, nil, header)
	if resp.Error() != nil {
		return
	}
	c.send(resp)
	return
}

func (c *Component) PostJson(ctx context.Context, path string, body map[string]interface{}, vs ...interface{}) (r *KResponse) {
	header, _ := getHeaderBody(vs)
	header["Content-Type"] = ContentTypeJson
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return
	}

	return c.Post(ctx, path, bodyJson, header)
}

func (c *Component) PostForm(ctx context.Context, path string, body map[string]string, vs ...interface{}) (r *KResponse) {
	header, _ := getHeaderBody(vs)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	values := url.Values{}
	for k, v := range body {
		values.Add(k, v)
	}
	return c.Post(ctx, path, values, header)
}

// Post 请求  body->[]byte/string/url.Values
func (c *Component) Post(ctx context.Context, path string, body interface{}, vs ...interface{}) (kResp *KResponse) {
	kResp = &KResponse{}
	header, _ := getHeaderBody(vs)
	var b []byte
	// 对 body 类型推断
	switch v := body.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	case url.Values:
		b = []byte(v.Encode())
	default:
		kResp.SetError(errors.New("body type is error"))
		return
	}

	kResp = c.getHttpRequest(ctx, http.MethodPost, path, b, header)
	if kResp.Error() != nil {
		return
	}
	c.send(kResp)
	return
}

// 获取 request
func (c *Component) getHttpRequest(ctx context.Context, method, path string, body []byte, header KHeader) (kResp *KResponse) {
	kResp = &KResponse{}
	if !isPullPath(path) {
		host := c.config.Host
		path = host + path
	}
	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(body))
	if err != nil {
		kResp.SetError(err)
		return
	}
	// set header
	if header != nil {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
	kResp.Req = req
	kResp.ReqBody = body
	return
}

// 重试
func (c *Component) send(resp *KResponse) {
	// 执行发送请求
	c.do(resp)

	retry := c.config.Retry
	// retry
	if resp.Response() == nil && retry > 0 {
		delay := time.NewTimer(c.config.RetryInterval)
		defer delay.Stop()
		for resp.Response() == nil && resp.Error() != nil && retry > 0 {
			delay.Reset(c.config.RetryInterval)
			select {
			case <-resp.Req.Context().Done():
				return
			case <-delay.C:
				// 重试
				c.do(resp)
				retry--
			}
		}
	}
	return
}

// 1. 发送请求
// 2. 加载拦截器
func (c *Component) do(kResponse *KResponse) {
	// 发送请求
	var h handler = func(kr *KResponse) {
		r, e := c.httpClient.Do(kr.Req)
		if e != nil {
			kr.SetError(e)
			return
		}
		// 处理 Resp.body
		b, _ := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		kr.RespBody = b
		kr.Resp = r
		return
	}

	// 拦截器 数组中越靠前的，越里层
	for _, hook := range c.config.hooks {
		h = hook(c.config, c.logger)(h)
	}
	h(kResponse)
	return
}

// 将 query 中的参数合并到 url 中
func addQueryValuesIntoPath(path string, query map[string]string) (newPath string, err error) {
	target, err := url.Parse(path)
	if err != nil {
		return
	}
	urlValues := target.Query()
	for key, value := range query {
		urlValues.Add(key, value)
	}
	target.RawQuery = urlValues.Encode()
	return target.String(), nil
}

// 解析 v interface 中的 body  header
func getHeaderBody(vs []interface{}) (header KHeader, body KBody) {
	header = make(KHeader)
	body = make(KBody)
	for _, v := range vs {
		switch v.(type) {
		case KHeader:
			header = v.(KHeader)
		case KBody:
			body = v.(KBody)
		default:
			continue
		}
	}
	return
}

// 判断 path 是否是 完整地址
func isPullPath(path string) bool {
	reg := regexp.MustCompile(`^(http|https):\/\/`)
	strArr := reg.FindAllString(path, -1)
	return len(strArr) > 0
}
