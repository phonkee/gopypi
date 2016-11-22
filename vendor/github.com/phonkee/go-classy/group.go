package classy

import (
	"strings"

	"github.com/justinas/alice"
)

/*
Group creates new group of multiple classy views
*/
func Group(path string, views ...Classy) (result Classy) {
	result = group{
		chain: alice.New(),
		views: views,
	}
	result = result.Path(path)

	return
}

/*
group is Group implementation, which basically groups multiple classy news and provides groupping
*/
type group struct {
	chain alice.Chain
	debug bool
	name  string
	path  string
	views []Classy
}

/*
Debug sets debugging on group
*/
func (g group) Debug() Classy {
	g.debug = true
	return g
}

/*
IsDebug returns whether debugging is enabled for group
*/
func (g group) isDebug() bool {
	return g.debug
}

/*
Name sets name of group
*/
func (g group) Name(name string) Classy {
	g.name = strings.TrimSpace(name)
	return g
}

/*
GetName returns name of group
*/
func (g group) GetName() string {
	return g.name
}

/*
Path sets path to group, all views in group  paths will be prepended with this path
*/
func (g group) Path(path string) Classy {
	g.path = strings.TrimSpace(path)
	return g
}

/*
Add Middleware
*/
func (g group) Use(m ...alice.Constructor) Classy {
	g.chain = g.chain.Append(m...)
	return g
}

/*
getBoundMethods collects bound methods from all views and adds path, updates name (if available)
*/
func (g group) getBoundMethods() (result []BoundMethod) {
	result = []BoundMethod{}

	// range over views and get all bound methods
	for _, view := range g.views {

		// set debug if available
		if g.isDebug() {
			view.Debug()
		}

		vbms := view.getBoundMethods()

		// iterate over bound methods of view and update information from group (path, name)
		for _, vbm := range vbms {

			// add name if available
			if g.name != "" {
				vbm.Name = makeName(g.name, vbm.Name)
			}

			// add path if set
			if g.path != "" {
				vbm.Path = g.path + vbm.Path
			}

			// add group middlewares
			vbm.Handlerfunc = g.chain.ThenFunc(vbm.Handlerfunc).ServeHTTP

			result = append(result, vbm)
		}
	}

	return
}
