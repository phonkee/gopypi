package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAction(t *testing.T) {

	Convey("Test Default Action", t, func() {
		a := NewAction()

		So(a.isDebug(), ShouldBeFalse)
		So(a.Debug().isDebug(), ShouldBeTrue)
	})

	Convey("Test Description", t, func() {
		a := NewAction()

		So(a.GetDescription(), ShouldEqual, "")
		description := "fieldname"
		So(a.Description(description).GetDescription(), ShouldEqual, description)

	})

	Convey("Test New Field", t, func() {
		a := NewAction()

		name := "fieldname"

		So(a.HasField(name), ShouldBeFalse)

		a.Field(name)
		So(a.HasField(name), ShouldBeTrue)

		So(func() {
			a.Field()
		}, ShouldPanic)

		result := a.Field("response", "user")
		So(a.Field("response").Field("user"), ShouldEqual, result)

	})

	Convey("Test Get Field Names", t, func() {
		a := NewAction()
		name := "fieldname"

		So(len(a.GetFieldNames()), ShouldEqual, 0)
		a.Field(name)
		So(len(a.GetFieldNames()), ShouldEqual, 1)

	})

	Convey("Test Action From", t, func() {
		a := NewAction()

		So(func() {
			a.From("")
		}, ShouldPanic)

		type TestStruct struct {
			A string `json:"a"`
			B string `json:"b"`
			C string `json:"-"`
		}

		a.From(TestStruct{})

		So(a.HasField("a"), ShouldBeTrue)
		So(a.HasField("b"), ShouldBeTrue)
		So(a.HasField("c"), ShouldBeFalse)

		a.From(&TestStruct{})

		So(a.HasField("a"), ShouldBeTrue)
		So(a.HasField("b"), ShouldBeTrue)
		So(a.HasField("c"), ShouldBeFalse)

	})

	Convey("Test Has Field", t, func() {
		a := NewAction()

		So(func() {
			a.HasField()
		}, ShouldPanic)

		a.Field("response").Field("user")

		So(a.HasField("response", "user"), ShouldBeTrue)

	})

	Convey("Test GetData", t, func() {
		a := NewAction()

		So(a.GetData()["type"], ShouldBeNil)

		name := "fieldname"
		a.Field(name)

		So(a.GetData()["body"], ShouldNotBeNil)
		So(a.GetData()["description"], ShouldBeNil)

		description := "desccrippp"
		a.Description(description).GetData()
	})

	Convey("Test Parse Query Params", t, func() {
		a := NewAction()

		a.ParseQueryParam("q=string&page=integer")
		So(len(a.GetQueryParamNames()), ShouldEqual, 2)
		So(a.GetData()["query"], ShouldNotBeNil)

		a = NewAction()
		a.ParseQueryParam("%")
		So(len(a.GetQueryParamNames()), ShouldEqual, 0)

	})

	Convey("Test Remove Query Param", t, func() {
		a := NewAction()

		a.ParseQueryParam("q=string&page=integer")
		So(len(a.GetQueryParamNames()), ShouldEqual, 2)

		a.RemoveQueryParam("q")
		So(len(a.GetQueryParamNames()), ShouldEqual, 1)


	})

}
