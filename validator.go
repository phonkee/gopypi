package core

import (
	"encoding/json"
	"strings"

	"fmt"

	"strconv"

	"github.com/asaskevich/govalidator"
)

func NewValidationResult() ValidationResult {
	return &validationResult{
		fieldErrors:   map[string][]ValidationError{},
		unboundErrors: []ValidationError{},
	}
}

type ValidationError struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error"`
}

type ValidationResult interface {
	// add field related error
	AddFieldError(field string, err error, code ...string) ValidationResult

	// add not field related error
	AddUnboundError(err error, code ...string) ValidationResult

	// returns whether result is valid
	IsValid() bool

	// marshals result to json
	MarshalJSON() ([]byte, error)

	// returns whether field has error
	HasFieldError(field string) bool
}

/*
Implementation of ValidationResult
*/
type validationResult struct {
	fieldErrors   map[string][]ValidationError
	unboundErrors []ValidationError
}

/*
Adds field error
*/
func (v *validationResult) AddFieldError(field string, err error, code ...string) ValidationResult {
	if _, ok := v.fieldErrors[field]; !ok {
		v.fieldErrors[field] = []ValidationError{}
	}

	ve := ValidationError{Error: err.Error()}
	if len(code) > 0 {
		ve.Code = code[0]
	}
	v.fieldErrors[field] = append(v.fieldErrors[field], ve)

	return v
}

/*
Add unbound error (not related to field)
*/
func (v *validationResult) AddUnboundError(err error, code ...string) ValidationResult {
	ve := ValidationError{Error: err.Error()}
	if len(code) > 0 {
		ve.Code = code[0]
	}
	v.unboundErrors = append(v.unboundErrors, ve)
	return v
}

/*
HasFieldError returns whether any errors are on result
*/
func (v *validationResult) HasFieldError(field string) bool {
	if value, ok := v.fieldErrors[field]; !ok {
		return false
	} else {
		return len(value) > 0
	}
}

/*
IsValid returns whether any errors are on result
*/
func (v *validationResult) IsValid() bool {
	return (len(v.fieldErrors) + len(v.unboundErrors)) == 0
}

/*
MarshalJSON marshals result to json
*/
func (v *validationResult) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{
		"fields":  v.fieldErrors,
		"unbound": v.unboundErrors,
	}

	return json.Marshal(result)
}

/*
Common validators
*/

/*
StringMinMaxValidator returns validator for length
*/
func StringMinMaxValidator(min int, max int) func(field string, value *string, vr ValidationResult) bool {
	minStr := strconv.Itoa(min)
	maxStr := strconv.Itoa(max)

	return func(field string, value *string, vr ValidationResult) bool {
		if !govalidator.StringLength(*value, minStr, maxStr) {
			vr.AddFieldError(field, fmt.Errorf("Value must have %v to %v characters", minStr, maxStr))
			return false
		}
		return true
	}
}

/*
Validate<field> function helpers

All these methods tah first the pointer to field (so it can manipulate the value e.g. TrimSpace) and also pointer to
ValidationResult to add errors
*/

var (
	usernameValidator = StringMinMaxValidator(USERNAME_MIN_LENGTH, USERNAME_MAX_LENGTH)

	// ValidatePassword validate password
	ValidatePassword = StringMinMaxValidator(PASSWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH)
)

/*
ValidateUsername validates username and adds error if available
*/
func ValidateUsername(field string, value *string, vr ValidationResult) bool {
	*value = strings.TrimSpace(*value)
	return usernameValidator(field, value, vr)
}

/*
ValidateEmail validates whether the value is valid email. Also it supports optional argument required.
*/
func ValidateEmail(field string, value *string, vr ValidationResult, required ...bool) bool {
	*value = strings.TrimSpace(*value)

	isRequired := IsEnabledOption(required)

	// not required not provided
	if !isRequired && *value == "" {
		return true
	}

	return govalidator.IsEmail(*value)
}
