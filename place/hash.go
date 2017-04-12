package place

import (
	"bufio"
	"crypto/sha256"
	"os"
)

func Hash(file string) (sum, head []byte, err error) {
	var (
		f     *os.File
		count int
	)
	if f, err = os.Open(file); err != nil {
		return
	}
	reader := bufio.NewReader(f)
	buffer := make([]byte, 65536) // 缓存64K
	hash := sha256.New()
	isHead := true
	for {
		if count, err = reader.Read(buffer); err != nil {
			break
		}
		if isHead {
			isHead = false
			head = buffer[:261]
		}
		hash.Write(buffer[:count])
	}
	sum = hash.Sum(nil)
	return sum, head, nil
}
