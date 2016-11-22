package xmlrpc

import (
	"bytes"
	"go/types"
	"strings"
)

func newRPCMethod(service, method string, signature *types.Signature) *rpcMethod {
	result := &rpcMethod{
		Method:    method,
		Service:   service,
		Params:    []Param{},
		Signature: signature,
	}

	// iterate over params
	for i := 0; i < result.Signature.Params().Len(); i++ {
		result.Params = append(result.Params, getParam(result.Signature.Params().At(i)))
	}

	// if length is one only error is returned
	count := result.Signature.Results().Len()
	if count == 1 {
		resultType := result.Signature.Results().At(0).Type().String()
		if resultType != "error" {
			Exit("Service method %v.%v should return either (value, error) or just error, got %v", result.Service, result.Method, resultType)
		}
		result.ResultError = getParam(result.Signature.Results().At(0))
	} else if count == 2 {
		resultType := result.Signature.Results().At(1).Type().String()
		if resultType != "error" {
			Exit("Service method %v.%v should return either (value, error) or just error", result.Service, result.Method)
		}

		result.Result = getParam(result.Signature.Results().At(0))
		result.ResultError = getParam(result.Signature.Results().At(1))
	} else {
		Exit("Service %v method %v must return either 2 variables (result, error) or just error")
	}

	return result
}

/*
rpcMethod holds information about method
*/
type rpcMethod struct {

	// Method name
	Method string

	// Service name
	Service string

	// Method signature
	Signature *types.Signature

	// params
	Params []Param

	// Result specification
	Result Param

	// result error
	ResultError Param
}

func (r *rpcMethod) HasResult() bool {
	return r.Result != nil
}

/*
FromEtree writes code to get values from xml
*/
func (r *rpcMethod) FromEtree(element string, resultvar string, errorvar string) string {

	buf := bytes.Buffer{}

	methodParams := []string{}

	for i, param := range r.Params {

		newelem := GenerateVariableName()

		elemval := param.FromEtree(newelem, param.Name(), "err")

		RenderTemplateInto(&buf, `
			{{.Variable}} := {{.Root}}.FindElement("param[{{.Index}}]/value")
			if {{.Variable}} == nil {
				{{.ErrorVar}} = xmlrpc.Errorf(400, "could not find {{.Name}}")
				return
			}
			{{.Param}}

		`, map[string]interface{}{
			"Root":     element,
			"Variable": newelem,
			"Index":    i + 1,
			"Param":    elemval,
			"ErrorVar": errorvar,
			"Result":   r.Result,
			"Name":     param.Name(),
		})

		methodParams = append(methodParams, param.Name())
	}

	RenderTemplateInto(&buf, `
		// If following method call fails there are 2 possible reasons:
		// 1. you have either changed method signature or you deleted method. Please re-run "go generate"
		// 2. you have probably found a bug and you should file issue on github.
		// @TODO: add panic recovery that returns error with 500 code

		{{if .ResultVar}}
			var {{.ResultVar}} {{.Result.Type}}

			{{.ResultVar}}, {{.ErrorVar}} = s.{{.Method}}({{.Params}})
		{{else}}
			{{.ErrorVar}} = s.{{.Method}}({{.Params}})
		{{end}}`, map[string]interface{}{
		"Service":   r.Service,
		"Method":    r.Method,
		"Params":    strings.Join(methodParams, ", "),
		"ResultVar": resultvar,
		"ErrorVar":  errorvar,
		"Result":    r.Result,
	})

	return buf.String()
}
