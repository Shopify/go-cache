package cache

import (
	"encoding/gob"
	"testing"
	"time"
)

func Test_gobEncoding(t *testing.T) {
	gob.Register(struct{}{})
	gob.Register(testStruct{})
	gob.Register(time.Time{})

	testEncoding(t, GobEncoding)
}
