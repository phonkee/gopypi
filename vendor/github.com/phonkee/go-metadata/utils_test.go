package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtils(t *testing.T) {

	Convey("Test Parse Tag Blank", t, func() {
		tag := "-"
		t, _ := ParseTag(tag)
		So(t, ShouldEqual, "-")
	})

	Convey("Test Parse Tag", t, func() {
		tag := "name,omitempty"
		t, parsed := ParseTag(tag)
		So(t, ShouldEqual, "name")

		So(parsed.Contains("omitempty"), ShouldBeTrue)

	})

	Convey("Test Parse Tag No Options", t, func() {
		tag := "name"
		t, parsed := ParseTag(tag)
		So(t, ShouldEqual, "name")

		So(parsed.Contains("omitempty"), ShouldBeFalse)

	})

	Convey("Test Parse Tag multiple Options", t, func() {
		tag := "name,omitempty,other"
		t, parsed := ParseTag(tag)
		So(t, ShouldEqual, "name")

		So(parsed.Contains("omitempty"), ShouldBeTrue)
		So(parsed.Contains("other"), ShouldBeTrue)
		So(parsed.Contains("notthere"), ShouldBeFalse)

	})

	Convey("Test String List Contains", t, func() {
		So(stringListContains([]string{"yes", "no"}, "no"), ShouldBeTrue)
		So(stringListContains([]string{"yes", "no"}, "never"), ShouldBeFalse)
	})

}
