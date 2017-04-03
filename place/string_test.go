package place

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSplit(t *testing.T) {
	Convey("拆分字符串", t, func() {
		a := Split("a b  c")
		So(a, ShouldContain, "c")
		So(len(a), ShouldEqual, 3)
		So(Split("a b c"), ShouldContain, "c")
		So(Split("a b,c"), ShouldContain, "c")
		a = Split("a b ; , c")
		So(a, ShouldContain, "c")
		So(len(a), ShouldEqual, 3)
		a = Split(" c;")
		So(a, ShouldContain, "c")
		So(len(a), ShouldEqual, 1)
	})
}
