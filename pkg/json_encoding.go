package cache

import (
	"encoding/json"
)

var JsonEncoding Encoding = NewJsonEncoding()

func NewJsonEncoding() *jsonEncoding {
	return &jsonEncoding{}
}

type jsonEncoding struct {}

func (e *jsonEncoding) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (e *jsonEncoding) Decode(b []byte, data interface{}) error {
	return json.Unmarshal(b, data)
}
