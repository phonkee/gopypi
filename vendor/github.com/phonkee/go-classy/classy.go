package classy

import (
	"strings"

	"reflect"

	"net/http"

	"fmt"

	"github.com/justinas/alice"
	"github.com/phonkee/go-response"
)

// support for custom http handlers that return response.Response
type ResponseHandlerFunc func(http.ResponseWriter, *http.Request) response.Response

/*
Classy is struct that wraps classy view and can register views to gorilla mux
*/
type Classy interface {

	// Debug enables debugging of classy view
	Debug() Classy

	// isDebug returns whether debug is enabled
	isDebug() bool

	// Name sets name of given classy view (route will be registered under this name)
	Name(name string) Classy

	// GetName returns name
	GetName() string

	// Path sets path (optional)
	Path(path string) Classy

	// Use adds middlewares for given classy view
	Use(middlewares ...alice.Constructor) Classy

	// GetMethods returns bound methods
	getBoundMethods() []BoundMethod
}

/*
newClassy returns new classy view
*/
func newClassy(view Viewer) Classy {
	return classy{
		chain:      alice.New(),
		path:       "/",
		structname: getStructName(view),
		view:       view,
	}.Name(getViewName(view, true))
}

/*
Implementation of classy
*/
type classy struct {
	chain      alice.Chain
	debug      bool
	name       string
	path       string
	structname string
	view       Viewer
}

/*
Debug enables debugging of classy View
*/
func (c classy) Debug() Classy {
	c.debug = true
	return c
}

/*
isDebug returns whether debugging is enabled
*/
func (c classy) isDebug() bool {
	return c.debug
}

/*
GetBoundMethods scans struct for all possible methods and maps them by Routes
*/
func (c classy) getBoundMethods() (result []BoundMethod) {

	result = make([]BoundMethod, 0)

	// methods that should be skipped
	ignoreMethods := []string{"Before", "Routes"}

	fooType := reflect.TypeOf(c.view)

	// inspect all struct methods and iterate over them.
	// then try to find whether
	for i := 0; i < fooType.NumMethod(); i++ {
		methodname := fooType.Method(i).Name
		method := reflect.ValueOf(c.view).MethodByName(string(methodname))

		// ignore methods
		if stringListContains(ignoreMethods, methodname) {
			continue
		}

		found := false
		// iterate over routes map

		for path, mapping := range c.view.Routes() {
			if c.isDebug() {
				mapping = mapping.Debug()
			}

			// get final mapping struct method => http method
			mappingmap := mapping.Get()

			joinedpath := joinPaths(c.path, path)

			// StructMethod is human readable name {struct}.{method}
			sm := fmt.Sprintf("%v.%v", c.structname, methodname)

			// method is found, now we need to check them
			if hm, ok := mappingmap[methodname]; ok {
				switch fn := method.Interface().(type) {
				case func(http.ResponseWriter, *http.Request):
					bm := BoundMethod{
						Handlerfunc:  c.getHandlerFunc(c.view, fn),
						Method:       hm,
						Name:         makeName(mapping.GetName(), c.name),
						Path:         joinedpath,
						StructMethod: sm,
					}
					result = append(result, bm)
					found = true
				case func(http.ResponseWriter, *http.Request) response.Response:
					bm := BoundMethod{
						Handlerfunc:  c.getHandlerFuncResponse(c.view, fn),
						Method:       hm,
						Name:         makeName(mapping.GetName(), c.name),
						Path:         joinedpath,
						StructMethod: sm,
					}
					result = append(result, bm)
					found = true
				default:
					loggerError(c.isDebug(), "Found method %+v with bad signature.", methodname)
				}
			}
		}

		if !found {
			// check type whether there is no typo
			switch method.Interface().(type) {
			case func(http.ResponseWriter, *http.Request):
				loggerWarning(c.isDebug(), "method %v is http.HandlerFunc but not found in method mapping(probably typo?)", methodname)
				continue
			case func(http.ResponseWriter, *http.Request) response.Response:
				loggerWarning(c.isDebug(), "method %v is ResponseHandlerFunc but not found in method mapping (probably typo?)", methodname)
				continue
			}

		}

	}

	return
}

func (c classy) getHandlerFunc(view Viewer, fn http.HandlerFunc) http.HandlerFunc {
	return c.chain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {

		// check if Before returns response
		response := view.Before(w, r)
		if response != nil {
			response.Write(w, r)
			return
		}

		// call method
		fn(w, r)
	}).ServeHTTP
}

func (c classy) getHandlerFuncResponse(view Viewer, fn ResponseHandlerFunc) http.HandlerFunc {
	return c.chain.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if Before returns response
		response := view.Before(w, r)
		if response != nil {
			response.Write(w, r)
			return
		}

		// call method
		response = fn(w, r)
		if response != nil {
			response.Write(w, r)
			return
		}
	}).ServeHTTP
}

/*
Name sets name of classy view
*/
func (c classy) Name(name string) Classy {
	c.name = strings.TrimSpace(name)
	return c
}

/*
GetName returns name of classy view
*/
func (c classy) GetName() string {
	return c.name
}

/*
Path sets path (optional)
*/
func (c classy) Path(path string) Classy {
	c.path = path
	return c
}

/*
Add Middleware
*/
func (c classy) Use(m ...alice.Constructor) Classy {
	c.chain = c.chain.Append(m...)
	return c
}

/*
BoundMethod is representation of struct method
*/
type BoundMethod struct {
	Handlerfunc  http.HandlerFunc
	Path         string
	Name         string
	Method       string
	StructMethod string
}

func (b BoundMethod) String() string {
	return fmt.Sprintf("path: %v method: %v name: %v method: %v", b.Path, b.Method, b.Name, b.StructMethod)
}
