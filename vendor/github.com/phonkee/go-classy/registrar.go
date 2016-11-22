/*
registrar provides object that handles all the registration. It is also used as shorthands in top level api.
*/
package classy

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/phonkee/go-response"
)

/*
Registrar interface

This interface is responsible for registering views.
*/
type Registrar interface {

	// Debug enables debug
	Debug() Registrar

	// returns if debug is enabled
	isDebug() bool

	// MethodNotAllowed sets response that will be returned
	MethodNotAllowed(response.Response) Registrar

	// Name sets name, this is a text/template that injects variable `view` with view name in snakecase
	Name(string) Registrar

	// Path sets path to registrar, this path will be prepended to all classy views paths
	Path(string) Registrar

	// Register registers all classy views to router
	// In the future we will probably set router parameter as interface{} and we will do type switch with concrete
	// impementations
	Register(router *mux.Router, views ...Classy) Registrar

	// Use sets middlewares to registrar (all classy views will use this middlewares)
	Use(middleware ...alice.Constructor) Registrar
}

/*
newRegistrar returns new default registrar
*/
func newRegistrar() Registrar {
	return registrar{
		chain: alice.New(),
		name:  "{name}",
		debug: debug,
	}
}

/*
Iplementation of Registrar interface
*/
type registrar struct {
	chain            alice.Chain
	debug            bool
	methodnotallowed response.Response
	name             string
	path             string
}

/*
Debug enables debug for registrar
*/
func (r registrar) Debug() Registrar {
	r.debug = true
	return r
}

/*
debugEnabled returns whether debug is enabled for registrar
*/
func (r registrar) isDebug() bool {
	return r.debug
}

/*
MethodNotAllowed sets response that will be returned in this situation
*/
func (r registrar) MethodNotAllowed(mna response.Response) Registrar {
	r.methodnotallowed = mna
	return r
}

/*
Debug enables debug for registrar
*/
func (r registrar) Name(name string) Registrar {
	r.name = name
	return r
}

/*
Path sets path to registrar that will be prepended to all classy views
*/
func (r registrar) Path(path string) Registrar {
	r.path = path
	return r
}

/*
Register registers classy views to given router
Currently the router is only gorilla mux router, but in the future I will implement type switch on router to provide
multiple implementations (third party routers)
*/
func (r registrar) Register(router *mux.Router, classyviews ...Classy) Registrar {
	if len(classyviews) == 0 {
		loggerError(r.isDebug(), "Register didn't received any classy views.")
	}

	// iterate over views and register them
	for _, classyview := range classyviews {

		// if registrar has enabled debug, enable debug in all views
		if r.isDebug() {
			classyview = classyview.Debug()
		}

		// range over bound methods and do following
		// 1. join paths
		// 2. resolve name
		for _, boundMethod := range classyview.getBoundMethods() {

			path := boundMethod.Path
			if r.path != "" {
				path = joinPaths(r.path, path)
			}

			name := makeName(r.name, boundMethod.Name)

			// register to router
			router.Handle(
				path,
				r.chain.ThenFunc(boundMethod.Handlerfunc),
			).Methods(boundMethod.Method).Name(name)

			// log info message
			loggerInfo(
				r.isDebug(),
				"Registered \"%v\", method: \"%v\", name: \"%v\", method: \"%v\"",
				path, boundMethod.Method, name, boundMethod.StructMethod,
			)
		}
	}

	return r
}

/*
Use sets middlewares to registrar
*/
func (r registrar) Use(middleware ...alice.Constructor) Registrar {
	r.chain = r.chain.Append(middleware...)
	return r
}
