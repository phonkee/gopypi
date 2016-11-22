/*
@TODO: add more documentation to this code, since it is more advanced code and can bring misunderstandnis.
*/
package metadata

import (
	"reflect"
	"time"
)

const (
	FIELD_INVALID = "invalid"
	FIELD_ARRAY = "array"
	FIELD_STRUCT = "struct"
	FIELD_MAP = "map"
	FIELD_INTEGER = "integer"
	FIELD_UNSIGNED_INGETER = "unsigned"
	FIELD_STRING = "string"
	FIELD_BOOLEAN = "boolean"
	FIELD_FLOAT = "float"
	FIELD_DATETIME = "datetime"
)

// field type func returns Field by reflect value
type FieldTypeFunc func(reflect.Type) Field

var (
	kinds = map[reflect.Kind]FieldTypeFunc{}

	// mapping of custom types
	types = map[reflect.Type]FieldTypeFunc{}

	AVAILABLE_FIELDS = []string{FIELD_INVALID, FIELD_ARRAY, FIELD_STRUCT, FIELD_MAP, FIELD_INTEGER,
		FIELD_UNSIGNED_INGETER,
		FIELD_STRING, FIELD_BOOLEAN, FIELD_FLOAT, FIELD_DATETIME}
)

func init() {
	// register kinds
	RegisterKind(ftf(FIELD_INTEGER), reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
	RegisterKind(ftf(FIELD_UNSIGNED_INGETER), reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)
	RegisterKind(ftf(FIELD_STRING), reflect.String)
	RegisterKind(ftf(FIELD_BOOLEAN), reflect.Bool)
	RegisterKind(ftf(FIELD_FLOAT), reflect.Float32, reflect.Float64)

	// special types
	RegisterKind(ftfStruct, reflect.Struct)
	RegisterKind(ftfArray, reflect.Array, reflect.Slice)
	RegisterKind(ftfMap, reflect.Map)
	RegisterType(ftf(FIELD_DATETIME), time.Now())
}

// register kinds
func RegisterKind(f FieldTypeFunc, kind ...reflect.Kind) {
	for _, k := range kind {
		kinds[k] = f
	}
}

// RegisterType provides register for custom types (e.g. time.Time)
func RegisterType(f FieldTypeFunc, values ...interface{}) (err error) {
	for _, val := range values {
		typ := reflect.TypeOf(val)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		types[typ] = f
	}
	return
}

// returns field by kind
func getFieldByKind(typ reflect.Type) (field Field) {

	typn := typ
	if typn.Kind() == reflect.Ptr {
		typn = typ.Elem()
	}

	if fn, ok := kinds[typn.Kind()]; ok {
		return fn(typ)
	}

	// if something is not implemented
	// @TODO: shouldn't we panic here?
	return newField().Type(FIELD_INVALID)
}

// returns field by value
func getField(typ reflect.Type) Field {

	orig := typ

	required := true

	for {
		if typ.Kind() == reflect.Ptr {
			required = false
			typ = typ.Elem()
		} else {
			break
		}
	}

	// custom types
	if fn, ok := types[typ]; ok {
		return fn(orig)
	}

	return getFieldByKind(orig).Required(required)
}

// default ftf imlpementation
func ftf(fieldtype string) FieldTypeFunc {
	return func(typ reflect.Type) (result Field) {
		result = newField().Type(fieldtype)
		if typ.Kind() == reflect.Ptr {
			result.Required(false)
		}
		return
	}
}

// field type function for struct
func ftfStruct(typ reflect.Type) (result Field) {
	result = newStructField()

	required := true

	for {
		if typ.Kind() == reflect.Ptr {
			required = false
			typ = typ.Elem()
		} else {
			break
		}
	}

	result.Required(required)

	for i := 0; i < typ.NumField(); i++ {
		ft := typ.Field(i)
		tag, _ := ParseTag(ft.Tag.Get("json"))
		if tag == "-" {
			continue
		}
		name := ft.Name
		if tag != "" {
			name = tag
		}

		// what to do with this??
		result.addField(name, getField(ft.Type))
	}

	return
}

// array type
func ftfArray(typ reflect.Type) (result Field) {
	return newField().Type(FIELD_ARRAY).addField("value", getField(typ.Elem()))
}

// map type
func ftfMap(typ reflect.Type) (result Field) {
	return newField().Type(FIELD_MAP).addField("key", getField(typ.Key())).addField("value", getField(typ.Elem()))
	return
}
