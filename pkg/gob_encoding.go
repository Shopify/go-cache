package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/pkg/errors"
)

var GobEncoding Encoding = NewGobEncoding()

func NewGobEncoding() *gobEncoding {
	return &gobEncoding{}
}

type gobEncoding struct {}

func (e *gobEncoding) Encode(item *Item) ([]byte, error) {
	encoded := &bytes.Buffer{}
	enc := gob.NewEncoder(encoded)
	if err := enc.Encode(*item); err != nil {
		return nil, errors.Wrap(err, "unable to encode item")
	}

	return encoded.Bytes(), nil
}

func (e *gobEncoding) Decode(data []byte) (*Item, error) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	var item Item
	if err := dec.Decode(&item); err != nil {
		return nil, errors.Wrap(err, "unable to decode item")
	}

	return &item, nil

}
