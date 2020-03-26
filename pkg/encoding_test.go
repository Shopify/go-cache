package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Foo string
}

func testEncoding(t *testing.T, encoding Encoding) {
	t.Run("empty struct", func(t *testing.T) {
		item := struct{}{}
		enc, err := encoding.Encode(item)
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec struct{}
		err = encoding.Decode(enc, &dec)
		require.NoError(t, err)
		require.EqualValues(t, item, dec)
	})

	t.Run("test struct", func(t *testing.T) {
		item := testStruct{Foo: "bar"}
		enc, err := encoding.Encode(item)
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec testStruct
		err = encoding.Decode(enc, &dec)
		require.NoError(t, err)
		require.EqualValues(t, item, dec)
	})

	t.Run("time struct", func(t *testing.T) {
		item, err := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 15:04:05 -0700")
		require.NoError(t, err)

		enc, err := encoding.Encode(item)
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec time.Time
		err = encoding.Decode(enc, &dec)
		require.NoError(t, err)
		require.EqualValues(t, item, dec)
	})

	t.Run("integer", func(t *testing.T) {
		item := 123
		enc, err := encoding.Encode(item)
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec int
		err = encoding.Decode(enc, &dec)
		require.NoError(t, err)
		require.EqualValues(t, item, dec)
	})

	t.Run("float", func(t *testing.T) {
		item := 1.23
		enc, err := encoding.Encode(item)
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec float64
		err = encoding.Decode(enc, &dec)
		require.NoError(t, err)
		require.EqualValues(t, item, dec)
	})

	t.Run("string", func(t *testing.T) {
		item := "123"
		enc, err := encoding.Encode(item)
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec string
		err = encoding.Decode(enc, &dec)
		require.NoError(t, err)
		require.EqualValues(t, item, dec)
	})

	t.Run("nil", func(t *testing.T) {
		var item interface{}
		enc, err := encoding.Encode(&item)
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec interface{}
		err = encoding.Decode(enc, &dec)
		require.NoError(t, err)
		require.EqualValues(t, item, dec)
	})

	t.Run("not a pointer", func(t *testing.T) {
		enc, err := encoding.Encode("123")
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec string
		err = encoding.Decode(enc, dec)
		require.EqualError(t, err, "argument to Get() must be a pointer")
		require.Zero(t, dec)
	})

	t.Run("wrong type", func(t *testing.T) {
		enc, err := encoding.Encode("123")
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec int
		err = encoding.Decode(enc, &dec)
		require.Error(t, err)
		require.Zero(t, dec)
	})
}
