package paginator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFactory(t *testing.T) {

	FirstCallOption := func(target *int) Option {
		return func(p Paginator, first bool) {
			if first {
				*target += 1
			}
		}
	}

	OtherCallOption := func(target *int) Option {
		return func(p Paginator, first bool) {
			if !first {
				*target += 1
			}
		}
	}

	Convey("TestFirstOption", t, func() {
		calls := 0
		factory := NewFactory(FirstCallOption(&calls))
		p := factory()
		So(calls, ShouldEqual, 1)
		p.Normalize()
		p.Normalize()
		So(calls, ShouldEqual, 1)
	})
	Convey("TestOtherOption", t, func() {
		ss := 0
		factory := NewFactory(OtherCallOption(&ss))
		p := factory()
		So(ss, ShouldEqual, 0)
		p.Normalize()
		p.Normalize()
		So(ss, ShouldEqual, 2)
	})

}
