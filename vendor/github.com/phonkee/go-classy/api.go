/*
Api provides simple functions that return registrar.

The api is fairly simple and can be used in many different ways.

Example:

	classy.Name("api:{{view}}).Register(
        router,
        classy.New(&ProductDetailView{}),
	)

    classy.Debug().Use(middleware1, middleware2).Name("api:{{view}}").Register(
        router,
        classy.New(&ProductDetailView{}).Use(middleware3),
        classy.New(&ProductApproveView{}).Path("/approve").Name("approve"),
    )

    classy.Register(
        router,
        classy.New(&ProductDetailView{}).Path("/product/").Name(),
        classy.New(&ProductApproveView{}).Path("/product/approve").Debug(),
    )
*/
package classy

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var (
	debug bool
)

/*
Debug function is basically Regsitrar with debug enabled
*/
func Debug() {
	debug = true
}

/*
Name function returns Regsitrar without given name
*/
func Name(name string) Registrar {
	return newRegistrar().Name(name)
}

/*
Path function creates Regsitrar with set Path
*/
func Path(path string) Registrar {
	return newRegistrar().Path(path)
}

/*
Register function is basically Regsitrar without any options
*/
func Register(router *mux.Router, views ...Classy) Registrar {
	return newRegistrar().Register(router, views...)
}

/*
Use function is basically Regsitrar with given middlewares
*/
func Use(middleware ...alice.Constructor) Registrar {
	return newRegistrar().Use(middleware...)
}

/*
New instantiates new classy view, that is able to introspect view for http methods.
*/
func New(view Viewer) Classy {
	return newClassy(view)
}
