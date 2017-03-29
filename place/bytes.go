package place

import (
	"bytes"
	"encoding/binary"
)

func BytesPrefix(prefix string, bs []byte) []byte {
	buf := bytes.NewBuffer([]byte(prefix))
	buf.Write(bs)
	return buf.Bytes()
}

func Int2Bytes(num int) []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, int32(num))
	return buf.Bytes()
}

func Bytes2Int(bs []byte) int {
	buf := bytes.NewBuffer(bs)
	var number int32
	binary.Read(buf, binary.BigEndian, &number)
	return int(number)
}
