package cache

import "github.com/Shopify/go-encoding"

// This file provides backwards compatibility for when these encodings were provided by this package.
// Can all be removed in v2

// Deprecated: Use encoding.ValueEncoder
type Encoder encoding.ValueEncoder

// Deprecated: Use encoding.ValueDecoder
type Decoder encoding.ValueDecoder

// Deprecated: Use encoding.ValueEncoding
type Encoding encoding.ValueEncoding

var (
	gobEncoding     = encoding.NewValueEncoding(encoding.GobEncoding)
	DefaultEncoding = gobEncoding

	// GobEncoding is a Gob encoding
	// Deprecated: Use encoding.GobEncoding
	GobEncoding = gobEncoding

	// JsonEncoding is a JSON encoding
	// Deprecated: Use encoding.JSONEncoding
	// nolint:golint
	JsonEncoding = encoding.JSONEncoding

	// NewLiteralEncoding is an encoding that will try its best to store the data as is,
	// but fallback on another encoder if not possible.
	// Deprecated: Use encoding.NewLiteralEncoding
	NewLiteralEncoding = encoding.NewLiteralEncoding
)
