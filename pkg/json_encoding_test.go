package cache

import (
	"testing"
)

func Test_jsonEncoding(t *testing.T) {
	testEncoding(t, JsonEncoding)
}
