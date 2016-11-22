package classy

import (
	"testing"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/phonkee/go-response"
	. "github.com/smartystreets/goconvey/convey"
)

/*
test middleware
*/
func middleware(http.Handler) http.Handler {
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
	})
}

type TestView struct{ GenericView }

func (t *TestView) GET(http.ResponseWriter, *http.Request) response.Response { return nil }
func (t *TestView) POST(string, *http.Request)                               {}

func TestRegistrar(t *testing.T) {

	Convey("Test New and debug flag", t, func() {
		r1 := newRegistrar()
		So(r1.isDebug(), ShouldBeFalse)

		r2 := newRegistrar().Debug()
		So(r2.isDebug(), ShouldBeTrue)
	})

	Convey("Test Name", t, func() {
		name := "somename"
		r1 := newRegistrar().Name(name)
		So(r1.(registrar).name, ShouldEqual, name)
	})

	Convey("Test Register", t, func() {
		Path("/api").Debug().Register(
			mux.NewRouter(),
			New(&TestView{}),
		)
	})
}
