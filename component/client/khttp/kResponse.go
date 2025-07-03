package khttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 返回的 response
// todo: 封闭结构体内属性
type KResponse struct {
	Err      error
	Req      *http.Request
	Resp     *http.Response
	ReqBody  []byte
	RespBody []byte
}

func (kr *KResponse) Error() error {
	return kr.Err
}

// TimeOut 请求是否超时
func (kr *KResponse) TimeOut() bool {
	switch kr.Err.(type) {
	case *url.Error:
		return kr.Err.(*url.Error).Timeout()
	default:
		return false
	}
}

// StatusCode http code
func (kr *KResponse) StatusCode() int {
	// 服务端 error code
	if kr.Error() == nil && kr.Response() != nil {
		return kr.Response().StatusCode
	}

	// 客户端 error code
	if kr.TimeOut() {
		return http.StatusRequestTimeout
	}
	return http.StatusBadRequest
}

func (kr *KResponse) SetError(err error) {
	kr.Err = err
}

func (kr *KResponse) Request() *http.Request {
	return kr.Req
}

func (kr *KResponse) Response() *http.Response {
	return kr.Resp
}

func (kr *KResponse) ResHeader() (h http.Header) {
	if kr.Response() != nil {
		h = kr.Response().Header
	}
	return
}

//
func (kr *KResponse) ToBytes() ([]byte, error) {
	if kr.Err != nil {
		return nil, kr.Err
	}

	if kr.RespBody != nil {
		return kr.RespBody, nil
	}

	// 重新读取 body
	defer kr.Resp.Body.Close()
	respBody, err := ioutil.ReadAll(kr.Resp.Body)
	if err != nil {
		kr.Err = err
		return nil, err
	}
	kr.RespBody = respBody

	return kr.RespBody, nil
}

func (kr *KResponse) Bytes() []byte {
	data, _ := kr.ToBytes()
	return data
}

func (kr *KResponse) ToString() (string, error) {
	data, err := kr.ToBytes()
	return string(data), err
}

func (kr *KResponse) String() string {
	data, _ := kr.ToBytes()
	return string(data)
}

func (kr *KResponse) AssignTo(v interface{}) error {
	data, err := kr.ToBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
