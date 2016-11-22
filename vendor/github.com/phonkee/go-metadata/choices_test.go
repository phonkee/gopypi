package metadata

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChoices(t *testing.T) {
	Convey("Test Add Choice", t, func() {
		choices := newChoices()
		So(choices.Count(), ShouldEqual, 0)
		choices.Add("value", "display")
		So(choices.Count(), ShouldEqual, 1)
	})

	Convey("Test MarshalJSON", t, func() {
		choices := newChoices()

		_, err := choices.MarshalJSON()
		So(err, ShouldBeNil)
	})

}
