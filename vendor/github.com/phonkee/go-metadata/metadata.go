package metadata

import (
	"encoding/json"
	"strings"
)

/*
MetaData interface that describes heep endpoint methods
*/
type Metadata interface {
	// action returns either existing or newly created action
	Action(method string) Action

	// RemoveAction removes action
	RemoveAction(method string) Metadata

	// Name sets name to metadata
	Name(string) Metadata

	// GetName returns name of metadata
	GetName() string

	// Description sets description of metadata
	Description(string) Metadata

	// GetDescription returns description
	GetDescription() string

	// GetData returns dynamic data (for json etc..)
	GetData() map[string]interface{}

	// MarshalJSON satisfies json marshaller
	MarshalJSON() ([]byte, error)

	// Debug enables debugging for metadata
	Debug() Metadata

	// isDebug returns whether debugging is enabled
	isDebug() bool
}

/*
New returns new metadata instance
*/
func New(label ...string) (result Metadata) {
	result = &metadata{
		actions: map[string]Action{},
	}

	// global debug
	if debugEnabled {
		result.Debug()
	}

	if len(label) > 0 {
		result.Name(label[0])
	}

	return
}

/*
Metadata interface implementation
*/
type metadata struct {
	// map to actions
	actions map[string]Action

	// name of metadata
	name string

	// description of metadata
	description string

	// debug enabled
	debug bool
}

/*
Actions either returns existing action, or it creates new one
*/
func (m *metadata) Action(method string) Action {

	method = cleanMethod(method)

	// check if action exists
	if _, ok := m.actions[method]; !ok {
		na := NewAction()

		// if debug is enabled, action must too
		if m.isDebug() {
			na.Debug()
		}

		m.actions[method] = na
	}

	return m.actions[method]
}

/*
RemoveAction removes action from metadata
*/
func (m *metadata) RemoveAction(method string) Metadata {
	delete(m.actions, cleanMethod(method))
	return m
}

/*
Name sets name of metadata
*/
func (m *metadata) Name(name string) Metadata {
	m.name = strings.TrimSpace(name)
	return m
}

/*
GetName returns name of metadata
*/
func (m *metadata) GetName() string {
	return m.name
}

/*
Description sets description of metadata
*/
func (m *metadata) Description(description string) Metadata {
	m.description = strings.TrimSpace(description)
	return m
}

/*
GetDescription returns description of metadata
*/
func (m *metadata) GetDescription() string {
	return m.description
}

/*
GetData returns data for json marshalling etc..
*/
func (m *metadata) GetData() (result map[string]interface{}) {
	result = map[string]interface{}{}

	if m.name != "" {
		result["name"] = m.name
	}

	if m.description != "" {
		result["description"] = m.description
	}

	if len(m.actions) > 0 {
		result["actions"] = m.actions
	}

	return
}

/*
MarshalJSON returns json representation of metadata
*/
func (m *metadata) MarshalJSON() (result []byte, err error) {
	result, err = json.Marshal(m.GetData())
	return
}

/*
Debug enables debugging
*/
func (m *metadata) Debug() Metadata {
	m.debug = true
	return m
}

/*
isDebug returns whether debugging is enabled
*/
func (m *metadata) isDebug() bool {
	return m.debug
}
