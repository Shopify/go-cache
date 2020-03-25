package cache

type Encoder interface {
	Encode(*Item) ([]byte, error)
}

type Decoder interface {
	Decode([]byte) (*Item, error)
}

type Encoding interface {
	Encoder
	Decoder
}

var DefaultEncoding = GobEncoding
