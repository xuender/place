package place

import (
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMime(t *testing.T) {
	Convey("获取文件mime", t, func() {
		Convey("基础格式", func() {
			bs, _ := ioutil.ReadFile("../samples/sample.gif")
			k := Mime(bs, "")
			So(k.MIME.Type, ShouldEqual, "image")
		})
		Convey("扩展格式", func() {
			bs, _ := ioutil.ReadFile("/bin/bash")
			k := Mime(bs, "/bin/bash")
			So(k.MIME.Subtype, ShouldEqual, "x-executable")
			So(k.MIME.Type, ShouldEqual, "application")
		})
	})
}
