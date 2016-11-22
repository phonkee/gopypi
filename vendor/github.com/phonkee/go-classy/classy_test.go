package classy

import (
	"testing"

	"net/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClassy(t *testing.T) {

	Convey("test register view", t, func() {

		//router := mux.NewRouter()
		//products := router.PathPrefix("/product").Subrouter()
		//
		//classy := New(TestProductView{})
		//classy.Register(products)
	})

}

// Test classes
type TestProductView struct {
	ViewSet
}

func (t TestProductView) Create(w http.ResponseWriter, r *http.Request) {

}
