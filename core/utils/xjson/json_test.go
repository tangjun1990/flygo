package xjson

import (
	"fmt"
	"testing"
)

type test struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func TestMarshal(t *testing.T) {
	data := test{
		Name: "tangjun1",
		Desc: "<a>link</a>",
	}
	result, _ := Marshal(data)
	fmt.Println(string(result))
}

func TestUnmarshal(t *testing.T) {
	data := `{"name":"tangjun1","desc":"\u003ca\u003elink\u003c/a\u003e"}`
	var v test
	_ = Unmarshal([]byte(data), &v)
	fmt.Println(v)
}
