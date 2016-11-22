package classy

import (
	"net/http"

	"github.com/phonkee/go-response"
)

/*
BaseView is blank implementation of basic view
*/
type BaseView struct{}

/*
Before blank implementation
*/
func (v BaseView) Before(w http.ResponseWriter, r *http.Request) response.Response {
	return nil
}

/*
GetRoutes blank implementation
*/
func (v *BaseView) Routes() map[string]Mapping {
	return map[string]Mapping{}
}

/*
View is predefined struct for list views
*/
type GenericView struct {
	BaseView
}

/*
Routes returns list of routes with predefined method maps
*/
func (l GenericView) Routes() map[string]Mapping {
	margs := [][]string{}
	for _, am := range AVAILABLE_METHODS {
		margs = append(margs, []string{am})
	}
	return map[string]Mapping{
		"/": NewMapping(margs...),
	}
}

/*
ListView is predefined struct for list views
*/
type ListView struct {
	BaseView
}

/*
Routes returns list of routes with predefined method maps
*/
func (l ListView) Routes() map[string]Mapping {
	return map[string]Mapping{
		"/": NewMapping(
			[]string{"GET", "List"},
			[]string{"OPTIONS", "Metadata"},
			[]string{"POST", "Create"},
		),
	}
}

/*
DetailView is predefined struct for detail
*/
type DetailView struct {
	BaseView
}

/*
Routes returns list of routes with predefined method maps
*/
func (d DetailView) Routes() (result map[string]Mapping) {
	result = map[string]Mapping{
		"/{pk:[0-9]+}/": NewMapping(
			[]string{"DELETE", "Delete"},
			[]string{"GET", "Retrieve"},
			[]string{"OPTIONS", "Metadata"},
			[]string{"POST", "Update"},
		),
	}
	return
}

/*
SlugDetailView is predefined struct for detail that handles id as slug
*/
type SlugDetailView struct {
	BaseView
}

/*
Routes returns list of routes with predefined method maps
*/
func (d SlugDetailView) Routes() map[string]Mapping {
	return map[string]Mapping{
		"/{slug}/": NewMapping(
			[]string{"DELETE", "Delete"},
			[]string{"GET", "Retrieve"},
			[]string{"OPTIONS", "Metadata"},
			[]string{"POST", "Update"},
		),
	}
}

/*
ViewSet is combination of list and detail views.
*/

type ViewSet struct {
	ListView
	DetailView
}

/*
Before is blank implementation for ViewSet
*/
func (v ViewSet) Before(w http.ResponseWriter, r *http.Request) response.Response {
	return nil
}

/*
Routes returns combination of list and detail routes
*/
func (v ViewSet) Routes() (result map[string]Mapping) {
	return JoinRoutes().
		Add(v.DetailView.Routes(), "{name}_detail", []string{"Metadata", "MetadataDetail"}).
		Add(v.ListView.Routes(), "{name}_list", []string{"Metadata", "MetadataList"}).
		Get()
}

type SlugViewSet struct {
	ListView
	SlugDetailView
}

/*
Before is blank implementation for ViewSet
*/
func (v SlugViewSet) Before(w http.ResponseWriter, r *http.Request) response.Response {
	return nil
}

/*
Routes returns combination of list and detail routes
*/
func (v SlugViewSet) Routes() map[string]Mapping {
	return JoinRoutes().
		Add(v.SlugDetailView.Routes(), "{name}_detail").
		Add(v.ListView.Routes(), "{name}_list").
		Get()
}
