package cache

import (
	"bytes"
	"encoding/gob"
)

var GobEncoding Encoding = NewGobEncoding()

func NewGobEncoding() *gobEncoding {
	return &gobEncoding{}
}

type gobEncoding struct {}

func (e *gobEncoding) Encode(data interface{}) ([]byte, error) {
	encoded := &bytes.Buffer{}
	enc := gob.NewEncoder(encoded)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return encoded.Bytes(), nil
}

func (e *gobEncoding) Decode(b []byte, data interface{}) error {
	dec := gob.NewDecoder(bytes.NewReader(b))
	return dec.Decode(data)
}
