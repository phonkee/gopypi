package metadata

import (
	"net/url"
	"strings"
)

/*
ParseQuery parses query and udpates action
query has following format

	var1=type&var2=type

Where type is type of Field

Example:
	q=string&page=integer

*/
func ParseQuery(query string, action Action) (err error) {

	var (
		values url.Values
	)

	if values, err = url.ParseQuery(query); err != nil {
		return
	}

	for k := range values {
		ft := strings.TrimSpace(values.Get(k))

		if !stringListContains(AVAILABLE_FIELDS, ft) {
			loggerWarning(action.isDebug(), "Type not found %v, using it anyway", ft)
		}

		// add query param
		action.QueryParam(k).Type(ft)
	}

	return
}
