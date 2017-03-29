package place

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBytesPrefix(t *testing.T) {
	Convey("增加字符串前缀", t, func() {
		So([]byte("aa123"), ShouldResemble, BytesPrefix("aa", []byte("123")))
	})
}
