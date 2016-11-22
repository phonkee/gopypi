package classy

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMapping(t *testing.T) {
	Convey("Test simple mapping", t, func() {
		m := NewMapping(
			[]string{"GET", "Get"},
			[]string{"POST"},
			[]string{"OPTIONS", "OPTIONS", "Metadata"},
		)

		result := m.Get()

		So(result["Get"], ShouldEqual, "GET")
		So(result["POST"], ShouldEqual, "POST")
		So(result["OPTIONS"], ShouldEqual, "OPTIONS")
		So(result["Metadata"], ShouldEqual, "OPTIONS")
	})
}
