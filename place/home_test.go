package place

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHome(t *testing.T) {
	Convey("获取当前用户目录", t, func() {
		path, _ := Home()
		So(path, ShouldEqual, "/home/ender")
	})
}
