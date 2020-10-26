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

	// Deprecated: Use encoding.GobEncoding
	GobEncoding = gobEncoding

	// Deprecated: Use encoding.JSONEncoding
	JsonEncoding = encoding.JSONEncoding

	// Deprecated: Use encoding.NewLiteralEncoding
	NewLiteralEncoding = encoding.NewLiteralEncoding
)
