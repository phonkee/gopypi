package signals

import (
	"testing"

	"context"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConnect(t *testing.T) {

	Convey("Test connect", t, func() {
		x := func(ctx context.Context) {}

		ss := New()
		ss.Connect(ReceiverFunc(x))

		sig, _ := ss.(*signal)
		So(len(sig.receivers), ShouldEqual, 1)
	})

	Convey("Test connect with id", t, func() {
		x := func(ctx context.Context) {}

		ss := New()
		ss.Connect(ReceiverFunc(x), "signal")

		sig, _ := ss.(*signal)
		So(len(sig.receivers), ShouldEqual, 1)

		ss.Connect(ReceiverFunc(x), "signal")
		So(len(sig.receivers), ShouldEqual, 1)

		ss.Connect(ReceiverFunc(x), "other")
		So(len(sig.receivers), ShouldEqual, 2)

		ss.Connect(ReceiverFunc(x), "other")
		So(len(sig.receivers), ShouldEqual, 2)

	})

	Convey("Test Dispatch (wait)", t, func() {
		var called bool
		sx := New()
		sx.Connect(ReceiverFunc(func(ctx context.Context) {
			called = true
		}))
		sx.Dispatch(context.Background(), true)
		So(called, ShouldBeTrue)
	})

}
