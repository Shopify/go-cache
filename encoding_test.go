package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Results:
//
// encoding   format  operation   ns/op
// json       int     encode      187
// json       float64 encode      270
// json       string  encode      192
// json       int     decode      274
// json       float64 decode      356
// json       string  decode      322
// literal    int     encode      75.3
// literal    float64 encode      187
// literal    string  encode      45.7
// literal    int     decode      105
// literal    float64 decode      124
// literal    string  decode      92.8
// gob        int     encode      1004
// gob        float64 encode      1034
// gob        string  encode      1000
// gob        int     decode      480
// gob        float64 decode      486
// gob        string  decode      491
//
// Average of operations and types:
// encoding	 ns/op
// literal	 105
// json		 267
// gob       749
//
func BenchmarkEncoding(b *testing.B) {
	benchmarks := map[string]Encoding{
		"gob":     GobEncoding,
		"json":    JsonEncoding,
		"literal": NewLiteralEncoding(nil),
	}
	for name, encoding := range benchmarks {
		b.Run(name, func(b *testing.B) {
			b.Run("int", func(b *testing.B) {
				b.Run("encode", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						_, _ = encoding.Encode(123)
					}
				})

				b.Run("decode", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						var val int
						_ = encoding.Decode([]byte("123"), &val)
					}
				})
			})

			b.Run("float64", func(b *testing.B) {
				b.Run("encode", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						_, _ = encoding.Encode(12.3)
					}
				})

				b.Run("decode", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						var val float64
						_ = encoding.Decode([]byte("12.3"), &val)
					}
				})
			})

			b.Run("string", func(b *testing.B) {
				b.Run("encode", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						_, _ = encoding.Encode("123")
					}
				})

				b.Run("decode", func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						var val string
						_ = encoding.Decode([]byte("123"), &val)
					}
				})
			})
		})
	}
}

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
		enc, err := encoding.Encode("12.3")
		require.NoError(t, err)
		require.NotNil(t, enc)

		var dec int
		err = encoding.Decode(enc, &dec)
		require.Error(t, err)
		require.Zero(t, dec)
	})
}
