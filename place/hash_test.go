package place

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHash(t *testing.T) {
	Convey("Hash", t, func() {
		s, h, e := Hash("../LICENSE")
		So(e, ShouldEqual, nil)
		So(len(s), ShouldEqual, 32)
		So(len(h), ShouldEqual, 261)
	})
}
