package metadata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestField(t *testing.T) {

	Convey("Test Default Field", t, func() {
		f := newField()
		So(f, ShouldNotBeNil)

		So(f.isDebug(), ShouldBeFalse)
		So(f.Debug().isDebug(), ShouldBeTrue)
	})

	Convey("Test Label", t, func() {
		label := "testlabel"
		f := newField().Label(label)
		So(f.GetLabel(), ShouldEqual, label)
	})

	Convey("Test Description", t, func() {
		descrption := "testht"
		f := newField().Description(descrption)
		So(f.GetDescription(), ShouldEqual, descrption)
	})

	Convey("Test add Field", t, func() {
		f := newField()
		f.Debug()
		added := f.addField("hello", newField())
		So(added.isDebug(), ShouldBeTrue)

		// test debug
		So(newField().Debug().From(1).isDebug(), ShouldBeTrue)
	})

	Convey("Test Add Field", t, func() {
		f := newField()

		So(f.NumFields(), ShouldEqual, 0)

		name := "testfield"
		sub := f.Field(name)
		sub2 := f.Field(name)
		So(sub, ShouldEqual, sub2)
		So(f.NumFields(), ShouldEqual, 1)

		other := "other"
		otherfield := newField()
		f.addField(other, otherfield)

		So(f.Field(other), ShouldEqual, otherfield)

		So(func() { f.Field() }, ShouldPanic)

		resultuser := f.Field("result", "user")

		So(f.Field("result").Field("user"), ShouldEqual, resultuser)

	})

	Convey("Test Has Field", t, func() {
		f := newField()
		name := "testfield"

		So(func() { f.HasField() }, ShouldPanic)

		So(f.HasField(name), ShouldBeFalse)
		f.Field(name)
		So(f.HasField(name), ShouldBeTrue)

		f.Field("one", "two", "three")
		So(f.HasField("one", "two", "three"), ShouldBeTrue)
	})

	Convey("Test Fields", t, func() {
		f := newField()

		So(len(f.GetFieldNames()), ShouldEqual, 0)

		name := "some field"
		f.Field(name)

		So(len(f.GetFieldNames()), ShouldEqual, 1)

		f.Field(name)
		So(len(f.GetFieldNames()), ShouldEqual, 1)

		f.RemoveField(name)
		So(len(f.GetFieldNames()), ShouldEqual, 0)
	})

	Convey("Test GetData", t, func() {

		label := "lllabel"
		description := "dddesc"

		f := newStructField().Label(label).Description(description)

		data := f.GetData()

		So(data["label"], ShouldEqual, label)
		So(data["fields"], ShouldBeNil)

		f.Field("subfield")

		data = f.GetData()

		So(len(data["fields"].(map[string]Field)), ShouldEqual, 1)

		So(data["choices"], ShouldBeNil)

		f.Choices().Add("value", "display")
		data = f.GetData()

		So(data["choices"], ShouldHaveSameTypeAs, newChoices())

		nf := newField().From([]string{})
		So(nf.GetData()["value"], ShouldNotBeNil)

		nf = newField().From(map[string]string{})
		So(nf.GetData()["key"], ShouldNotBeNil)
		So(nf.GetData()["value"], ShouldNotBeNil)

		f.Source("/something")
		So(f.GetData()["source"], ShouldNotBeNil)

	})

	Convey("Test RemoveField", t, func() {
		f := newField()
		name := "testfield"
		So(f.HasField(name), ShouldBeFalse)
		f.Field(name)
		So(f.HasField(name), ShouldBeTrue)
		f.RemoveField(name)
		So(f.HasField(name), ShouldBeFalse)
	})

	Convey("Test Required", t, func() {
		f := newField()
		So(f.IsRequired(), ShouldBeFalse)

		f.Required(true)
		So(f.IsRequired(), ShouldBeTrue)
	})

	Convey("Test Type", t, func() {
		f := newField()
		So(f.GetType(), ShouldEqual, "")

		typ := "sometype"

		f.Type(typ)

		So(f.GetType(), ShouldEqual, typ)

	})

	Convey("Test Field.From", t, func() {
		f := newField()

		type TestStruct struct {
			First  string  `json:"first"`
			Second *string `json:"second"`
			Third  *struct {
			} `json:"third"`
		}

		f.From(TestStruct{})

		So(f.HasField("first"), ShouldBeTrue)
		So(f.IsRequired(), ShouldBeTrue)

		So(f.Field("second").IsRequired(), ShouldBeFalse)
		So(f.Field("third").IsRequired(), ShouldBeFalse)

		f.From(&TestStruct{})
		So(f.HasField("first"), ShouldBeTrue)
		So(f.IsRequired(), ShouldBeFalse)
		So(f.GetType(), ShouldEqual, FIELD_STRUCT)
	})

	Convey("Test MarshalJSON", t, func() {
		f := newField()
		_, err := f.MarshalJSON()
		So(err, ShouldBeNil)
	})

	Convey("Test Choices", t, func() {
		f := newField()
		So(f.Choices(), ShouldHaveSameTypeAs, newChoices())
	})

	Convey("Test Source", t, func() {

		action := NewAction()
		action.Field("result").Type(FIELD_ARRAY)

		mid := newField()

		src := mid.Source("/suggest").Action(action).Result("result")

		So(src, ShouldHaveSameTypeAs, newField().Source("/suggest"))
		So(newField().Source("/suggest").GetData(), ShouldNotBeNil)

	})

}
