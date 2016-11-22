package xmlrpc

import (
	"bytes"
	"go/types"
	"strconv"

	"strings"
)

/*
New implementation of parameter
*/
type Param interface {
	// Name returns code name
	Name() string

	// Type returns type
	Type() string

	// SetType sets type (useful in case of named types)
	SetType(string)

	// Writes Field
	FromEtree(element string, resultvar string, errvar string) string

	// Writes param to element
	ToEtree(element string, resultvar string, errvar string) string

	// returns all imports for given param
	Imports() []string
}

/*
getParam returns appropriate param based on given variable
*/
func getParam(variable *types.Var) Param {
	switch x := variable.Type().(type) {
	case *types.Basic:
		bitSize := 0
		unsigned := false
		switch x.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
			bitSize = 0
			unsigned = false
			switch x.Kind() {
			case types.Int8:
				bitSize = 8
			case types.Int16:
				bitSize = 16
			case types.Int32:
				bitSize = 32
			case types.Int64:
				bitSize = 64
			}
			return newIntParam(variable.Name(), bitSize, unsigned)
		case types.String:
			return newStringParam(variable.Name())
		case types.Bool:
			return newBoolParam(variable.Name())
		case types.Float32, types.Float64:
			bitsize := 32
			if x.Kind() == types.Float64 {
				bitsize = 64
			}
			return newFloatParam(variable.Name(), bitsize)
		}
	case *types.Struct:
		return newStructParam(variable)
	case *types.Array:
		Exit("No support for array")
	case *types.Slice:
		// support for []byte (transforms to base64)
		if x.Elem().String() == "byte" {
			return newBytesParam(variable.Name())
		}
		// all other slices
		v := types.NewVar(0, variable.Pkg(), variable.Name(), x.Elem())
		sliceElemParam := getParam(v)
		return newSliceParam(variable.Name(), x.Elem().String(), sliceElemParam)
	case *types.Named:
		// first we check for error
		if variable.Type().String() == "error" {
			return newErrorParam("err")
		}

		// create new variable for named type and call getParam
		v := types.NewVar(variable.Pos(), variable.Pkg(), variable.Name(), x.Underlying())
		result := getParam(v)

		t := strings.TrimLeft(x.String(), ".")
		result.SetType(t)

		return result
	default:
		// pass
	}
	Exit("not supported param: %v", variable.Type().String())

	return nil
}

/*
newBoolParam returns boolParam instance (Param implementation for type bool)
*/
func newBoolParam(name string) Param {
	return &boolParam{
		name: name,
		typ:  "bool",
	}
}

/*
boolParam - Param implementation of boolean values
*/
type boolParam struct {
	name string
	typ  string
}

func (b *boolParam) Imports() []string  { return []string{} }
func (p *boolParam) Name() string       { return p.name }
func (p *boolParam) Type() string       { return p.typ }
func (p *boolParam) SetType(typ string) { p.typ = typ }
func (p *boolParam) FromEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}
	RenderTemplateInto(&buf, `
	var {{.Varname}} {{.Type}}
	if {{.Varname}}, {{.ErrorVar}} = xmlrpc.XPathValueGetBool({{.Element}}, "{{.Name}}"); {{.ErrorVar}} != nil {
		return
	}
	`, map[string]interface{}{
		"Element":  element,
		"ErrorVar": errvar,
		"Type":     p.Type(),
		"Varname":  resultvar,
		"Name":     p.name,
	})

	return buf.String()
}
func (p *boolParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `
		{{.Temp}} := "0"
		if {{.Varname}} {
			{{.Temp}} = "1"
		}
		{{.Element}}.CreateElement("boolean").SetText({{.Temp}})`, map[string]interface{}{
		"Element":  element,
		"Varname":  resultvar,
		"ErrorVar": errvar,
		"Temp":     GenerateVariableName("boolstr"),
	})

	return buf.String()
}

/*
newBytesParam returns new bytesParam
*/
func newBytesParam(name string) Param {
	return &bytesParam{
		name: name,
		typ:  "[]byte",
	}
}

/*
bytesParam is Param imlpementation for []byte
*/
type bytesParam struct {
	name string
	typ  string
}

func (p *bytesParam) Imports() []string  { return []string{"encoding/base64"} }
func (p *bytesParam) Name() string       { return p.name }
func (p *bytesParam) Type() string       { return p.typ }
func (p *bytesParam) SetType(typ string) { p.typ = typ }
func (p *bytesParam) FromEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}
	RenderTemplateInto(&buf, `var {{.Varname}} {{.Type}}
	if {{.Varname}}, {{.ErrorVar}} = xmlrpc.XPathValueGetBytes({{.Element}}, "{{.Name}}"); {{.ErrorVar}} != nil {
		return
	}
	`, map[string]interface{}{
		"Element":  element,
		"ErrorVar": errvar,
		"Type":     p.Type(),
		"Varname":  resultvar,
		"Name":     p.name,
	})

	return buf.String()
}
func (p *bytesParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}
	RenderTemplateInto(&buf, `
		{{.Temp}} := base64.StdEncoding.EncodeToString({{.Varname}})
		{{.Element}}.CreateElement("base64").SetText({{.Temp}})
	`, map[string]interface{}{
		"Element":  element,
		"Varname":  resultvar,
		"ErrorVar": errvar,
		"Temp":     GenerateVariableName("temp_variable"),
	})

	return buf.String()
}

/*
newFloatParam returns new floatParam
*/
func newFloatParam(name string, bitsize int) Param {
	return &floatParam{
		name:    name,
		bitsize: bitsize,
	}
}

/*
floatParam is Param imlpementation for float variables (32, 64)
*/
type floatParam struct {
	name    string
	bitsize int
	typ     string
}

func (b *floatParam) Imports() []string { return []string{} }
func (p *floatParam) Name() string      { return p.name }
func (p *floatParam) Type() string {
	if p.typ != "" {
		return "float" + strconv.Itoa(p.bitsize)
	}
	return p.typ
}
func (p *floatParam) SetType(typ string) { p.typ = typ }

func (p *floatParam) FromEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}
	RenderTemplateInto(&buf, `
	var {{.Varname}} {{.Type}}
	if {{.Varname}}, {{.ErrorVar}} = xmlrpc.XPathValueGetFloat{{.BitSize}}({{.Element}}, "{{.Name}}"); {{.ErrorVar}} != nil {
		return
	}
	`, map[string]interface{}{
		"Element":  element,
		"ErrorVar": errvar,
		"Type":     p.Type(),
		"Varname":  resultvar,
		"Name":     p.name,
		"BitSize":  p.bitsize,
	})

	return buf.String()
}
func (p *floatParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}
	RenderTemplateInto(&buf, `{{if eq .BitSize 32}}{{.Element}}.CreateElement("double").SetText(strconv.FormatFloat(float64({{.Varname}}), 'f', -1, 32))
		{{else}}{{.Element}}.CreateElement("double").SetText(strconv.FormatFloat({{.Varname}}, 'f', -1, 64)){{end}}`, map[string]interface{}{
		"BitSize":  p.bitsize,
		"Element":  element,
		"Varname":  resultvar,
		"ErrorVar": errvar,
		"Temp":     GenerateVariableName(""),
	})

	return buf.String()
}

/*
newIntParam returns new intParam (Param) instance
*/
func newIntParam(name string, bitSize int, unsigned bool) Param {
	return &intParam{
		name:     name,
		bitSize:  bitSize,
		unsigned: unsigned,
	}
}

/*
intParam is Param implementation covering all integers
*/
type intParam struct {
	bitSize  int
	name     string
	typ      string
	unsigned bool
}

func (b *intParam) Imports() []string { return []string{} }

/*
Name returns name of param
*/
func (i *intParam) Name() string {
	return i.name
}

/*
Type returns type of param
*/
func (i *intParam) Type() string {
	if i.typ != "" {
		return i.typ
	}

	result := "int"

	if i.unsigned {
		result = "u" + result
	}

	if i.bitSize > 0 {
		result += strconv.Itoa(i.bitSize)
	}

	return result
}

func (p *intParam) SetType(typ string) { p.typ = typ }

func (i *intParam) getParseFunc() string {
	return "XPathValueGetInt"

}

func (i *intParam) FromEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `
	var {{.Varname}} {{.Type}}
	if {{.Varname}}, {{.ErrorVar}} = xmlrpc.{{.ParseFunc}}({{.Element}}, "{{.Name}}"); {{.ErrorVar}} != nil {
		return
	}`, map[string]interface{}{
		"Element":   element,
		"ErrorVar":  errvar,
		"Type":      i.Type(),
		"Varname":   resultvar,
		"Name":      i.Name(),
		"ParseFunc": i.getParseFunc(),
	})

	return buf.String()
}

func (i *intParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `{{.Element}}.CreateElement("int").SetText(strconv.Itoa(int({{.ResultVar}})))`,
		map[string]interface{}{
			"Element":   element,
			"ResultVar": resultvar,
			"ErrorVar":  errvar,
		},
	)

	return buf.String()
}

func newStructParam(variable *types.Var) Param {
	strukt := variable.Type().(*types.Struct)

	result := &structParam{
		name:   variable.Name(),
		typ:    variable.Type().String(),
		params: make([]Param, 0, strukt.NumFields()),
	}

	for i := 0; i < strukt.NumFields(); i++ {
		result.params = append(result.params, getParam(strukt.Field(i)))
	}

	return result
}

/*
structParam is implementation of Param for struct variables
*/
type structParam struct {
	name   string
	typ    string
	params []Param
}

func (b *structParam) Imports() []string  { return []string{} }
func (p *structParam) Name() string       { return p.name }
func (p *structParam) Type() string       { return p.typ }
func (p *structParam) SetType(typ string) { p.typ = typ }
func (p *structParam) FromEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `
	// rendering struct
	{{.ResultVar}} := {{.Type}}{}

	{{$temp := GenerateVariableName "name_elem" }}
	{{$nameVar := GenerateVariableName "name" }}
	{{$valueVar := GenerateVariableName "value" }}

	// Lets iterate over given members.
	// @TODO: we should check first "struct" if not provided it's probably error
	for _, member := range {{.Element}}.FindElements("struct/members") {
		var {{$temp}} *etree.Element
		if {{$temp}} = member.FindElement("name"); {{$temp}} == nil {
			return errors.New("no name provided")
		}

		{{$nameVar}} := {{$temp}}.Text()

		var {{$valueVar}} *etree.Element
		if {{$valueVar}} = member.FindElement("value"); {{$valueVar}} == nil {
			return errors.New("no name provided")
		}

		// switch over param names (over all params)
		switch {{$nameVar}} {
			{{range $index,$param := .Params}}
				case "{{$param.Name}}": {{$paramTmp := GenerateVariableName }}
				{{$name := $param.Name }}
				{{$param.FromEtree $valueVar $paramTmp "err" }}

				// Assign to variable (for pointer support we can provide it here
				{{$.ResultVar}}.{{$name}} = {{$paramTmp}}{{end}}
		}
	}
	`, map[string]interface{}{
		"Type":      p.Type(),
		"ResultVar": resultvar,
		"Element":   element,
		"Params":    p.params,
	})

	return buf.String()
}
func (p *structParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `
		{{.StructVar}} := {{.Element}}.CreateElement("struct")
		// iterate over struct members
		{{range .Params}}
			{{$MemberVar:= GenerateVariableName "member"}}
			{{$MemberVar}} := {{$.StructVar}}.CreateElement("member")

			// first create "name" xml element with member name
			{{$MemberVar}}.CreateElement("name").SetText("{{.Name}}")

			{{$TempValueVar := GenerateVariableName "value"}}
			{{$TempValueVar}} := {{$MemberVar}}.CreateElement("value")

			// make shortcut to struct member {{$StructItemVar := GenerateVariableName "struct_var"}}
			{{$StructItemVar}} := {{$.ResultVar}}.{{.Name}}

			// set value
			{{.ToEtree $TempValueVar $StructItemVar $.ErrorVar }}
		{{end}}
	`,
		map[string]interface{}{
			"Element":   element,
			"ErrorVar":  errvar,
			"Params":    p.params,
			"ResultVar": resultvar,
			"StructVar": GenerateVariableName("struct"),
		},
	)

	return buf.String()
}

func newSliceParam(name string, typ string, obj Param) Param {
	return &sliceParam{
		name:   name,
		typ:    typ,
		object: obj,
	}
}

type sliceParam struct {
	name   string
	typ    string
	object Param
}

func (b *sliceParam) Imports() []string  { return []string{} }
func (p *sliceParam) Name() string       { return p.name }
func (p *sliceParam) Type() string       { return "[]" + p.object.Type() }
func (p *sliceParam) SetType(typ string) { p.typ = typ }

func (p *sliceParam) FromEtree(element string, resultvar string, errvar string) string {

	buf := bytes.Buffer{}
	RenderTemplateInto(&buf, `
	// This is slice implementation of {{.ResultVar}}
	{{.ResultVar}} := []{{.Type}}{}

	{{$memberVar := GenerateVariableName "member"}}

	// Lets iterate over given members.
	for _, {{$memberVar}} := range {{.Element}}.FindElements("array/data/value") {
		{{$targetName := GenerateVariableName "value"}}
		{{.Object.FromEtree $memberVar $targetName .ErrVar }}
		{{.ResultVar}} = append({{.ResultVar}}, {{$targetName}})
	}
	`, map[string]interface{}{
		"Element":   element,
		"ErrVar":    errvar,
		"ResultVar": resultvar,
		"Type":      p.Type(),
		"Object":    p.object,
	})

	return buf.String()
}
func (p *sliceParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `{{.Temp}} := {{.Element}}.CreateElement("array").CreateElement("data")
		for _, {{.TempItem}} := range {{.ResultVar}} {
			{{.TempValueVar}} := {{.Temp}}.CreateElement("value")
			{{.Object.ToEtree .TempValueVar .TempItem .ErrorVar }}
		}
	`, map[string]interface{}{
		"Element":      element,
		"ErrorVar":     errvar,
		"Object":       p.object,
		"ResultVar":    resultvar,
		"Temp":         GenerateVariableName("array_data"),
		"TempItem":     GenerateVariableName("item"),
		"TempValueVar": GenerateVariableName("value"),
	})

	return buf.String()
}

/*
newErrorParam returns new errorParam (Param implementation for error)
*/
func newErrorParam(name string) Param {
	return &errorParam{
		name: name,
		typ:  "error",
	}
}

/*
errorParam Param implementation for error
*/
type errorParam struct {
	name string
	typ  string
}

func (b *errorParam) Imports() []string  { return []string{} }
func (p *errorParam) Name() string       { return p.name }
func (p *errorParam) Type() string       { return p.typ }
func (p *errorParam) SetType(typ string) { p.typ = typ }

func (p *errorParam) FromEtree(element string, resultvar string, errvar string) string {
	return ""
}
func (p *errorParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `

		// move this code to error.
		{{.Fault}} := {{.Element}}.CreateElement("fault")

		{{.ErrorCode}} := 500

		// Try to cast error to xmlrpc.Error (with code added)
		if {{.CastVar}}, {{.OkVar}} := {{.ResultVar}}.(xmlrpc.Error); {{.OkVar}} {
			{{.ErrorCode}} = {{.CastVar}}.Code()
		}

		{{.StructVar}} := {{.Fault}}.CreateElement("value").CreateElement("struct")

		{{$member := GenerateVariableName "member"}}
		{{$member}} := {{.StructVar}}.CreateElement("member")
		{{$member}}.CreateElement("name").SetText("faultCode")
		{{$member}}.CreateElement("value").CreateElement("int").SetText(strconv.Itoa({{.ErrorCode}}))

		{{$member := GenerateVariableName "member"}}
		{{$member}} := {{.StructVar}}.CreateElement("member")
		{{$member}}.CreateElement("name").SetText("faultString")
		{{$member}}.CreateElement("value").CreateElement("string").SetText( {{.ResultVar}}.Error())

	`,
		map[string]interface{}{
			"CastVar":   GenerateVariableName("code"),
			"StructVar": GenerateVariableName("struct"),
			"OkVar":     GenerateVariableName("ok"),
			"ErrorCode": GenerateVariableName("code"),
			"Element":   element,
			"Fault":     GenerateVariableName("fault"),
			"ResultVar": resultvar,
			"ErrorVar":  errvar,
		})

	return buf.String()
}

/*
newStringParam returns new strinParam
*/
func newStringParam(name string) Param {
	return &stringParam{
		name: name,
		typ:  "string",
	}
}

/*
stringParam is Param imlpementation for string variables
*/
type stringParam struct {
	name string
	typ  string
}

func (b *stringParam) Imports() []string  { return []string{} }
func (p *stringParam) Name() string       { return p.name }
func (p *stringParam) Type() string       { return p.typ }
func (p *stringParam) SetType(typ string) { p.typ = typ }

func (p *stringParam) FromEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}
	RenderTemplateInto(&buf, `
	var {{.Varname}} {{.Type}}
	if {{.Varname}}, {{.ErrorVar}} = xmlrpc.XPathValueGetString({{.Element}}, "{{.Name}}"); {{.ErrorVar}} != nil {
		return
	}
	`, map[string]interface{}{
		"Element":  element,
		"ErrorVar": errvar,
		"Type":     p.Type(),
		"Varname":  resultvar,
		"Name":     p.name,
	})

	return buf.String()
}
func (p *stringParam) ToEtree(element string, resultvar string, errvar string) string {
	buf := bytes.Buffer{}

	RenderTemplateInto(&buf, `{{.Element}}.CreateElement("string").SetText({{.Varname}})`, map[string]interface{}{
		"Element":  element,
		"Varname":  resultvar,
		"ErrorVar": errvar,
		"Temp":     GenerateVariableName(""),
	})

	return buf.String()
}
