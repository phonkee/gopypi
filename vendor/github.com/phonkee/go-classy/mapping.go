package classy

type Mapping interface {

	// Add adds mapping from http mehod to struct method names
	Add(httpmethod string, names ...string) Mapping

	// Debug enables debugging on mapping
	Debug() Mapping

	// returns if debug is enabled
	isDebug() bool

	// GetMap returns map of struct method to http method
	Get() map[string]string

	// set name
	Name(string) Mapping

	// returns name
	GetName() string

	// Renames field
	Rename(from, to string) Mapping
}

/*
Do not forget on this

AVAILABLE_METHODS

*/
func NewMapping(items ...[]string) (result Mapping) {
	result = newMapping()

	for _, item := range items {
		if len(item) == 0 {
			continue
		}

		if len(item) == 1 {
			result.Add(item[0], item[0])
		} else if len(item) > 1 {
			result.Add(item[0], item[1:]...)
		}
	}

	return
}

/*
newMapping creates new mapping
*/
func newMapping() Mapping {
	result := mapping{
		dispatcher: make(map[string][]string),
		name:       "",
	}
	return result
}

/*
Mapping implementation
*/
type mapping struct {
	debug bool

	// mapping from http method to other methods
	dispatcher map[string][]string

	name string
}

/*
Add adds aliases
*/
func (m mapping) Add(httpmethod string, names ...string) Mapping {

	// check if http method exists
	if !stringListContains(AVAILABLE_METHODS, httpmethod) {
		loggerError(m.isDebug(), "http method %v not in available methods %v", httpmethod, AVAILABLE_METHODS)
		return m
	}

	if _, ok := m.dispatcher[httpmethod]; !ok {
		m.dispatcher[httpmethod] = []string{}
	}

	m.dispatcher[httpmethod] = append(m.dispatcher[httpmethod], names...)

	return m
}

/*
Debug enables debugging on mapping
*/
func (m mapping) Debug() Mapping {
	m.debug = true
	return m
}

/*
debugEnabled returns whether debug is enabled
*/
func (m mapping) isDebug() bool {
	return m.debug
}

/*
Get returns actual mapping from struct method to http method.
*/
func (m mapping) Get() (result map[string]string) {
	result = make(map[string]string)

	var (
		ok    bool
		value string
	)

	// hm - http methods
	// sms - structmethods
	for hm, sms := range m.dispatcher {
		sms = removeDuplicates(sms)
		if len(sms) == 0 {
			continue
		}

		// iterate over all available struct methods
		for _, sm := range sms {

			// struct method already in result for other http method
			if value, ok = result[sm]; ok && value != hm {
				loggerError(m.isDebug(), "Method %v already registered to %v", sm, value)
				continue
			}
			result[sm] = hm
		}
	}

	return
}

func (m mapping) GetName() string {
	return m.name
}

func (m mapping) Name(name string) Mapping {
	m.name = name
	return m
}

func (m mapping) Rename(from, to string) Mapping {

	for http, structmethods := range m.dispatcher {

		for i, sm := range structmethods {
			if sm == from {
				m.dispatcher[http][i] = to
			}
		}
	}

	return m
}
