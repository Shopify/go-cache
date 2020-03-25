package cache

import (
	"encoding/gob"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_gobEncoding(t *testing.T) {
	type test struct {
		Foo string
	}
	gob.Register(test{})

	tests := map[string]Item{
		"empty":      {},
		"expiration": {Expiration: time.Unix(10000, 0)},
		"struct":     {Data: test{Foo: "bar"}},
		"integer":    {Data: 123},
		"float":      {Data: 1.2},
		"string":     {Data: "123"},
		"nil":        {Data: nil},
	}

	for name, item := range tests {
		t.Run(name, func(t *testing.T) {
			enc, err := GobEncoding.Encode(&item)
			assert.NoError(t, err)
			assert.NotNil(t, enc)

			dec, err := GobEncoding.Decode(enc)
			assert.NoError(t, err)
			assert.NotNil(t, dec)

			assert.EqualValues(t, item, *dec)
		})
	}
}
