package utils

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

// UnmarshalReader 将 r 中的 json 格式的数据, 解析到 v
func UnmarshalReader(r io.Reader, v interface{}) error {
	d := jsoniter.NewDecoder(r)
	return d.Decode(v)
}

// UnmarshalJSON 将 data 中的 json 格式的数据, 解析到 v
func UnmarshalJSON(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
