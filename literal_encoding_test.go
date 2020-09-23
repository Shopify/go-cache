package cache

import (
	"fmt"
	"testing"
)

func Test_literalEncoding(t *testing.T) {
	for _, e := range []Encoding{JsonEncoding, GobEncoding} {
		t.Run(fmt.Sprintf("%T", e), func(t *testing.T) {
			testEncoding(t, NewLiteralEncoding(e))
		})
	}
}
