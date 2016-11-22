package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMetadata(t *testing.T) {

	Convey("Test Default Metadata", t, func() {
		md := New()
		So(len(md.(*metadata).actions), ShouldEqual, 0)

		// check debug
		So(md.isDebug(), ShouldBeFalse)
		So(md.Debug().isDebug(), ShouldBeTrue)

		So(debugEnabled, ShouldBeFalse)
		Debug()
		So(debugEnabled, ShouldBeTrue)
		So(New().isDebug(), ShouldBeTrue)

		debugEnabled = false
	})

	Convey("Test Add Action", t, func() {
		md := New()
		So(len(md.(*metadata).actions), ShouldEqual, 0)

		action := md.Action("GET")
		So(len(md.(*metadata).actions), ShouldEqual, 1)

		action2 := md.Action("GET")
		So(len(md.(*metadata).actions), ShouldEqual, 1)

		So(action, ShouldEqual, action2)

		// test with debug
		md = New().Debug()
		action = md.Action("GET")
		So(action.isDebug(), ShouldBeTrue)

	})

	Convey("Test Remove Action", t, func() {
		md := New()
		md.Action("GET")

		So(len(md.GetData()["actions"].(map[string]Action)), ShouldEqual, 1)

		md.RemoveAction("GET")
		So(md.GetData()["actions"], ShouldBeNil)

	})
	Convey("Test Name get/set", t, func() {
		md := New()

		So(md.GetName(), ShouldEqual, "")
		name := "some name"

		md.Name(name)
		So(md.GetName(), ShouldEqual, name)

		So(New(name).GetName(), ShouldEqual, name)

	})

	Convey("Test Description get/set", t, func() {
		md := New()

		So(md.GetDescription(), ShouldEqual, "")
		description := "description"

		md.Description(description)
		So(md.GetDescription(), ShouldEqual, description)
	})

	Convey("Test GetData/MarshalJSON", t, func() {

		name := "mdname"
		description := "mddescription"
		md := New().Name(name).Description(description)
		md.Action(ACTION_CREATE)

		data := md.GetData()
		So(data["name"], ShouldEqual, name)
		So(data["description"], ShouldEqual, description)

		_, err := md.MarshalJSON()
		So(err, ShouldBeNil)

	})

}
