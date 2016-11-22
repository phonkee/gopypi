/*
All interfaces used in classy module
*/
package classy

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/phonkee/go-response"
)

var (
	// List of available (supported) http methods. You can extend this with new methods
	AVAILABLE_METHODS = []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "TRACE", "HEAD", "DELETE"}
)

// Before func is called before any view is called. If Response is returned it's written and stopped execution
type BeforeFunc func(w http.ResponseWriter, r *http.Request) response.Response

/*
View

Interface for ClassyView, structs will be so-called class based view.
*/
type Viewer interface {
	Before(w http.ResponseWriter, r *http.Request) response.Response

	// Return used route map.
	Routes() map[string]Mapping
}

/*
GetFuncName returns function name (primarily for logging reasons)
*/
func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
