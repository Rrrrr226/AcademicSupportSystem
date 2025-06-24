package rds

import (
	"bytes"
	"strings"
	"sync"
)

var keyPool = sync.Pool{
	New: func() interface{} {
		b := &bytes.Buffer{}
		b.Grow(32)
		return b
	},
}

func Key(business string, indexes ...string) string {
	buf := keyPool.Get().(*bytes.Buffer)

	buf.WriteString("/pinnacle/")
	buf.WriteString(business)
	buf.WriteString("/")
	buf.WriteString(strings.Join(indexes, "/"))

	key := buf.String()
	buf.Reset()
	keyPool.Put(buf)
	return key
}
