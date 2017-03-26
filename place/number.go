package place
import (
	"bytes"
	"encoding/binary"
)

func Int2Bytes(num int32) []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, num)
	return buf.Bytes()
}

func Bytes2Int(bs []byte) int32 {
	buf := bytes.NewBuffer(bs)
	var number int32
	binary.Read(buf, binary.BigEndian, &number)
	return number
}
