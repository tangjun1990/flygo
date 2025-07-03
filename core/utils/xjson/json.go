package xjson

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.Config{
	SortMapKeys:            true,
	UseNumber:              true,
	CaseSensitive:          true,
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(&v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, &v)
}
