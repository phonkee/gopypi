package metadata

import (
	"testing"

	"reflect"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFType(t *testing.T) {

	Convey("Test Simple", t, func() {

		f := getField(reflect.TypeOf(int(0)))
		So(f.GetType(), ShouldEqual, FIELD_INTEGER)

	})

}
