package metadata

import (
	"encoding/json"
	"reflect"
	"strings"
)

/*
Action interface that describes single action (or http method)
*/
type Action interface {

	// Description sets description for action
	Description(string) Action

	// GetDescription returns description
	GetDescription() string

	// Field adds or retrieves fields (recursively)
	Field(...string) Field

	// HasField returns whether field is set
	HasField(names ...string) bool

	// GetFieldNames returns names of all fields
	GetFieldNames() []string

	// From inspects given value and makes appropriate steps
	From(v interface{}) Action

	// GetData returns dynamic data (for json etc..)
	GetData() map[string]interface{}

	// MarshalJSON satisfies json marshaller
	MarshalJSON() ([]byte, error)

	// Debug enables debugging
	Debug() Action

	// isDebug returns whether debugging is enabled
	isDebug() bool

	// ParseQueryParam parses
	ParseQueryParam(query string) Action

	// QueryParam returns or adds new query param
	QueryParam(name string) Field

	// RemoveQueryParam removes query param
	RemoveQueryParam(name string) Action

	// GetQueryParamNames returns names of all available query params
	GetQueryParamNames() []string
}

/*
NewAction creates fresh new action
*/
func NewAction() Action {
	return &action{
		query: newStructField(),
	}
}

/*
action is implementation of Action interface
*/
type action struct {

	// description of action
	description string

	// field
	field Field

	// debug
	debug bool

	// query params as Field
	query Field
}

/*
Description sets action description
*/
func (a *action) Description(description string) Action {
	a.description = strings.TrimSpace(description)
	return a
}

/*
GetDescription returns action description
*/
func (a *action) GetDescription() string {
	return a.description
}

/*
Field adds or retrieves field
*/
func (a *action) Field(names ...string) Field {

	if a.field == nil {
		a.field = newStructField()
	}

	return a.field.Field(names...)
}

/*
HasField returns whether field is set
*/
func (a *action) HasField(names ...string) bool {

	if len(names) == 0 {
		panic("please provide top level field")
	}

	if a.field == nil {
		return false
	}

	return a.field.HasField(names...)
}

/*
GetFieldNames returns names of all fields
*/
func (a *action) GetFieldNames() []string {

	if a.field == nil {
		return []string{}
	}

	return a.field.GetFieldNames()
}

/*
Read target structure and add fields
*/
func (a *action) From(target interface{}) Action {
	typ := reflect.TypeOf(target)
	for {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		} else {
			break
		}
	}

	if !(typ.Kind() == reflect.Struct || typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice) {
		panic("Metadata.Action.From supports only structs/array/slice")
	}

	a.field = newField().From(target)
	return a
}

/*
GetData returns data for json marshalling etc..
*/
func (a *action) GetData() (result map[string]interface{}) {

	result = map[string]interface{}{}

	if a.field != nil {
		data := a.field.GetData()

		// no need to have this information
		delete(data, "label")
		delete(data, "description")
		delete(data, "required")

		result["body"] = data
	}
	if a.description != "" {
		result["description"] = a.description
	}

	// check if we have query params
	if a.query.NumFields() > 0 {
		qpmap := map[string]Field{}

		for _, qpName := range a.query.GetFieldNames() {
			qpmap[qpName] = a.query.Field(qpName)
		}

		result["query"] = qpmap
	}

	return
}

/*
MarshalJSON returns json representation of metadata
*/
func (a *action) MarshalJSON() (result []byte, err error) {
	result, err = json.Marshal(a.GetData())
	return
}

/*
Debug enables debugging
*/
func (a *action) Debug() Action {
	a.debug = true
	return a
}

/*
isDebug returns whether debugging is enabled
*/
func (a *action) isDebug() bool {
	return a.debug
}

/*
ParseQueryParam parses url query and sets QueryParams
*/
func (a *action) ParseQueryParam(query string) Action {
	var (
		err error
	)
	if err = ParseQuery(query, a); err != nil {
		loggerError(a.isDebug(), "parse query error: %v", query)
	}
	return a
}

/*
QueryParam returns or adds new query param
*/
func (a *action) QueryParam(name string) Field {
	return a.query.Field(name)
}

/*
RemoveQueryParam removes query param by name
*/
func (a *action) RemoveQueryParam(name string) Action {
	a.query.RemoveField(name)
	return a
}

/*
GetQueryParamNames returns names of all available query params
*/
func (a *action) GetQueryParamNames() []string {
	return a.query.GetFieldNames()
}
