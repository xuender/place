package place

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNowFormat(t *testing.T) {
	Convey("当前时间格式化", t, func() {
		t, _ := time.Parse("2006-01-02", "2017-04-02")
		So(TimeFormat(t, "yyyy/mm"), ShouldEqual, "2017/04")
		So(TimeFormat(t, "yy/m/dd"), ShouldEqual, "17/4/02")
		So(TimeFormat(t, "yy/m/d"), ShouldEqual, "17/4/2")
	})
}
