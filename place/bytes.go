package place

import (
	"bytes"
)

func BytesPrefix(prefix string, bs []byte) []byte {
	buf := bytes.NewBuffer([]byte(prefix))
	buf.Write(bs)
	return buf.Bytes()
}
