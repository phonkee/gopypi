package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSource(t *testing.T) {

	Convey("Test Default", t, func() {
		So(newSource(), ShouldNotBeNil)
		So(newSource().isDebug(), ShouldBeFalse)
		So(newSource().Debug().isDebug(), ShouldBeTrue)

		path := "/suggest"
		So(newSource(path).GetPath(), ShouldEqual, path)
	})

	Convey("Test Action", t, func() {
		s := newSource()
		So(s.GetAction(), ShouldNotBeNil)
	})

	Convey("Test Path", t, func() {
		So(newSource().Path("/hello?world=true").GetPath(), ShouldEqual, "/hello")
		So(newSource().Path("%ola").GetPath(), ShouldEqual, "")
	})

	Convey("Test Result", t, func() {

		action := NewAction()
		action.Field("hello")

		So(newSource().Action(action).Result("hello"), ShouldNotBeNil)
		So(newSource().Action(action).Result("nope"), ShouldNotBeNil)

		action.Field("world").Type(FIELD_ARRAY)

		ns := newSource()
		So(ns.Action(action).Result("world"), ShouldNotBeNil)

	})

	Convey("Test IsValid", t, func() {

		action := NewAction()
		action.Field("hello").Type(FIELD_ARRAY)
		So(newSource().Path("/user/suggest").Action(action).Result("hello").IsValid(), ShouldBeTrue)

	})
	Convey("Test GetData", t, func() {

		action := NewAction()
		action.Field("hello").Type(FIELD_ARRAY)

		data := newSource().Action(action).Result("hello").GetData()
		_ = data

	})

	Convey("Test GetData", t, func() {
		_, err := newSource().MarshalJSON()
		So(err, ShouldBeNil)

		action := NewAction()
		action.Field("hello").Type(FIELD_ARRAY)

		data := newSource().Path("/hello").Action(action).Result("hello").GetData()
		So(data["metadata"], ShouldNotBeNil)

	})

	Convey("Test MarshalJSON", t, func() {
		_, err := newSource().MarshalJSON()
		So(err, ShouldBeNil)
	})

	Convey("Test Value", t, func() {
		action := NewAction()
		action.Field("hello").Type(FIELD_ARRAY)
		value := "testvalueorsomething42"
		So(newSource().Value(value).GetValue(), ShouldEqual, value)

	})

	Convey("Test Display", t, func() {
		action := NewAction()
		action.Field("hello").Type(FIELD_ARRAY)
		value := "testvalueorsomething42"
		So(newSource().Display(value).GetDisplay(), ShouldEqual, value)

	})

}
