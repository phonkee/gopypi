package metadata

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseQuery(t *testing.T) {

	Convey("ParseInvalidQuery", t, func() {
		action := NewAction()
		ParseQuery("q=string&other=unknown", action)

		So(len(action.GetQueryParamNames()), ShouldEqual, 2)

	})

}