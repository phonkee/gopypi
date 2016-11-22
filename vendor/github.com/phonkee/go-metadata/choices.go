package metadata

import "encoding/json"

/*
Choices interface
*/
type Choices interface {
	// Add adds new choice
	Add(value interface{}, display string) Choices

	// returns choices count
	Count() int

	// Support for marshalling
	MarshalJSON() ([]byte, error)
}

/*
Returns new choices
*/
func newChoices() Choices {
	return &choices{
		choices: []Choice{},
	}
}

// choice type
type Choice struct {
	Value   interface{} `json:"value"`
	Display string      `json:"display"`
}

/*
choices is Choices implementation
*/
type choices struct {
	choices []Choice
}

/*
Add single choice
*/
func (c *choices) Add(value interface{}, display string) Choices {
	c.choices = append(c.choices, Choice{Value: value, Display: display})
	return c
}

/*
Count returns count of choices
*/
func (c *choices) Count() int {
	return len(c.choices)
}

/*
MarshalJSON returns json representation of choices
*/
func (c *choices) MarshalJSON() (result []byte, err error) {
	result, err = json.Marshal(c.choices)
	return
}
