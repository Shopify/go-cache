package cache

type Encoder interface {
	Encode(data interface{}) ([]byte, error)
}

type Decoder interface {
	Decode(b []byte, data interface{}) error
}

type Encoding interface {
	Encoder
	Decoder
}

var DefaultEncoding = GobEncoding
