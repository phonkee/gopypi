package xmlrpc

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"bytes"

	"crypto/md5"

	"sync/atomic"

	"io"

	"github.com/fatih/color"
)

var (
	idcounter uint64 = 0
)

/*
GenerateVariableName generates unique variable name
*/
func GenerateVariableName(prefix... string) string {
	id := atomic.AddUint64(&idcounter, 1)

	p := "v_"
	if len(prefix) > 0 && strings.TrimSpace(prefix[0]) != ""{
		p = strings.TrimSpace(prefix[0]) + "_"
	}

	return p + strconv.Itoa(int(id))
}

/*
Exit prints error to stdout and exits with error
*/
func Exit(err interface{}, args...interface{}) {
	msg := ""
	switch err := err.(type) {
	case string:
		msg = fmt.Sprintf(err, args...)
	case error:
		msg = err.Error()
	default:
		msg = fmt.Sprintf("%v", err)

	}

	color.Red("Error: " + msg)
	os.Exit(1)
}

func RenderTemplateInto(w io.Writer, tpl string, data map[string]interface{}, funcmaps ...template.FuncMap) {
	fmt.Fprintln(w, RenderTemplate(tpl, data, funcmaps...))
	return
}

/*
RenderTemplate renders text template with data
*/
func RenderTemplate(tpl string, data map[string]interface{}, funcmaps ...template.FuncMap) string {
	var (
		buf bytes.Buffer
		err error
		t   *template.Template
	)

	funcmap := template.FuncMap{
		"GenerateVariableName": GenerateVariableName,
		"IsBlank":              IsBlank,
		"IsNotBlank":           IsNotBlank,
	}

	// update with given funcmaps
	for _, item := range funcmaps {
		for key, value := range item {
			funcmap[key] = value
		}
	}

	if t, err = template.New(MD5(tpl)).Funcs(funcmap).Parse(tpl); err != nil {
		Exit("cannot parse template: " + err.Error())
	}

	if err = t.Execute(&buf, data); err != nil {
		Exit("cannot execute template: " + err.Error())
	}

	return buf.String()
}

/*
IsBlank check trimmed value equal to ""
*/
func IsBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

/*
IsNotBlank is opposite of IsBlank :-D
*/
func IsNotBlank(value string) bool {
	return !IsBlank(value)
}

/*
MD5 helper
*/
func MD5(value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}

/*
getAvailableMethodsVariable returns available methods global variable name
*/
func getAvailableMethodsVariable(service string) string {
	return fmt.Sprintf("availableMethodsFor%v", service)
}

func getStructName(service, method, what string) string {
	return "__" + service + method + what
}

func getRequestStructName(service, method string) string {
	return getStructName(service, method, "Request")
}

func getResponseStructName(service, method string) string {
	return getStructName(service, method, "Response")
}

func getAvailableMethods(methods []*rpcMethod) string {

	parts := make([]string, 0, len(methods))

	for _, method := range methods {
		parts = append(parts, strconv.Quote(method.Method))
	}

	return strings.Join(parts, ", ")
}
