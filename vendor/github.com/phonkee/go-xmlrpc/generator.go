/*
Generator generates two structures, one for parameters one for result values.

Result values must be two, first is actual return value and second is error.
If error is returned fault response will be generated.

*/
package xmlrpc

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"text/template"

	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/loader"
)

/*
Generator interface
*/
type Generator interface {
	// add service by its name
	AddService(name string) error

	// format returns formatted source code
	Format() []byte
}

/*
NewGenerator returns Generator implementation
*/
func NewGenerator(filename string) (Generator, error) {
	result := &generator{
		services: map[string][]*rpcMethod{},
	}

	result.addImports("github.com/phonkee/go-xmlrpc", "github.com/beevik/etree", "strconv")

	var err error

	// parse file
	if err = result.parseFile(filename); err != nil {
		Exit("this")
		return nil, err
	}
	return result, nil
}

/*
generator
*/
type generator struct {
	// buffer for accumulated output
	buf bytes.Buffer

	// parsed package
	pkg *types.Package

	imports []string

	// store methods
	services map[string][]*rpcMethod
}

func (g *generator) addImports(imports ...string) {
out:
	for _, imp := range imports {
		for _, existing := range g.imports {
			if imp == existing {
				continue out
			}
		}
	}

	g.imports = append(g.imports, imports...)
}

/*
parseFile parses file
*/
func (g *generator) parseFile(filename string) (err error) {

	fset := token.NewFileSet()

	kpath := "."

	pkgs, e := parser.ParseDir(fset, kpath, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, 0)
	if e != nil {
		Exit(e)
		return
	}

	astf := make([]*ast.File, 0)
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			astf = append(astf, f)
		}
	}
	var conf loader.Config
	conf.CreateFromFiles(".", astf...)

	prog, errLoad := conf.Load()

	if errLoad != nil {
		Exit(errLoad)
	}

	p := prog.Package(".")
	if p == nil {
		Exit("errors")
	}

	g.pkg = p.Pkg

	return nil

}

func (g *generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *generator) WriteTemplate(tpl string, data map[string]interface{}) {
	fm := template.FuncMap{
		"getAvailableMethodsVariable": getAvailableMethodsVariable,
		"getAvailableMethods":         getAvailableMethods,
	}
	g.Printf("%v", RenderTemplate(tpl, data, fm))
}

// format generates and returns the gofmt-ed contents of the Generator's buffer.
func (g *generator) Format() []byte {
	g.buf = bytes.Buffer{}

	// writeHeader writes header (package name, imports)
	g.writeHeader()

	// add additional imports for params
	for _, service := range g.services {
		for _, method := range service {
			for _, param := range method.Params {
				g.addImports(param.Imports()...)
			}
		}
	}

	g.WriteTemplate(`
	package {{.Package}}
	import (
		{{range .Imports}}"{{.}}"
		{{end}}
	)

	{{range $service, $methods := .Services}}
		{{ $availMethodsVarname := getAvailableMethodsVariable $service}}

		var (
			{{$availMethodsVarname}} = map[string]bool{
			{{range $methods }}"{{.Method}}": true,{{end}}
			}
		)

		/*
		MethodExists returns whether rpc method is available on service
		*/
		func (s *{{$service}}) MethodExists(method string) (ok bool) {
			_, ok = {{$availMethodsVarname}}[method]
			return
		}

		/*
		ListMethods returns list of all available methods for given service
		*/
		func (s *{{$service}}) ListMethods() []string {
			result := make([]string, 0, len({{$availMethodsVarname}}))
			for key := range {{$availMethodsVarname}} {
				result = append(result, key)
			}
			return result
		}

		/*
		Dispatch dispatches method on service, do not use this method directly.
		root is params *etree.Element (actually "methodCall/params"
		*/
		func (s *{{$service}}) Dispatch(method string, root *etree.Element) (doc *etree.Document, err error) {

			// call appropriate methods
			switch method { {{range $methods}}
			case "{{.Method}}":
				// Get parameters from xmlrpc request

				{{$resultVar := GenerateVariableName "result"}}

				{{ if .Result }}
					{{.FromEtree "root" $resultVar "err"}}
				{{ else }}
					{{.FromEtree "root" "" "err"}}
				{{ end }}

				// create *etree.Document
				doc = etree.NewDocument()
				doc.CreateProcInst("xml", "version=\"1.0\" encoding=\"UTF-8\"")
				{{$methodResponse := GenerateVariableName "methodResponse"}}{{$methodResponse}} := doc.CreateElement("methodResponse")
				if {{.ResultError.Name}} != nil {
					{{.ResultError.ToEtree $methodResponse .ResultError.Name "err" }}
				} else {
					// here is place where we need to hydrate results {{if .Result }} {{$tempParam := GenerateVariableName}}
						{{$tempParam}} := {{$methodResponse}}.CreateElement("params").CreateElement("param").CreateElement("value")
						{{.Result.ToEtree $tempParam $resultVar "err" }} {{end}}
				}{{end}}
			default:
				// method not found, this should not happened since we check whether method exists
				err = xmlrpc.ErrMethodNotFound
				return
			}
			return
		}
	{{end}}
	`, map[string]interface{}{
		"Services": g.services,
		"Package":  g.pkg.Name(),
		"Imports":  g.imports,
	})

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		fmt.Printf("warning: internal error: invalid Go generated: %s\n", err)
		fmt.Print("warning: compile the package to analyze the error\n")
		return g.buf.Bytes()
	}
	return src
}

/*
AddService adds service
*/
func (g *generator) AddService(name string) error {
	obj := g.pkg.Scope().Lookup(name)
	if obj == nil {
		return fmt.Errorf("Service %v unavailable.", name)
	}

	// if service not there create one
	if _, ok := g.services[name]; !ok {
		g.services[name] = []*rpcMethod{}
	}

	service := obj.Type()

	// prepare methodset
	mset := types.NewMethodSet(types.NewPointer(service))

	// iterate over methods
	for i := 0; i < mset.Len(); i++ {
		what := mset.At(i).Obj().(*types.Func)
		signature := what.Type().(*types.Signature)

		// add service method
		g.services[name] = append(g.services[name], newRPCMethod(name, what.Name(), signature))
	}

	return nil
}

/*
Write header
*/
func (g *generator) writeHeader() {
	g.Printf("// This file is autogenerated by xmlrpcgen\n")
	g.Printf("// do not change it directly!\n")
}
