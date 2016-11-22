package template

import (
	"fmt"
	"html/template"
	"io"
)

// FuncMap is a convenience type that mirrors the FuncMap type in html/template
type FuncMap template.FuncMap

// HTML is another convenience type that mirrors the HTML type in html/template
// (http://golang.org/src/html/template/content.go?h=HTML#L120)
type HTML string

// AssetFunc is the function that go-bindata generates to look up a file
// by name
type AssetFunc func(string) ([]byte, error)

// Must is a helper that wraps a call to a function returning
// (*Template, error) and panics if the error is non-nil. It is intended for
// use in variable initializations such as
//	var t = template.Must(template.New("name").Parse("templates/my.tmpl"))
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(fmt.Sprintf("template error: %s", err))
	}
	if t == nil {
		panic(fmt.Sprintf("template was nil"))
	}
	return t
}

// Template is a wrapper around a Template (from html/template). It reads
// template file contents from a function instead of the filesystem.
type Template struct {
	AssetFunc AssetFunc
	tmpl      *template.Template
}

// New creates a new Template with the given name. It stores
// the given Asset() function for use later.
// Example usage:
//  tmpl := template.New("mytmpl", Asset) //Asset is the function that go-bindata generated for you
//
func New(name string, fn AssetFunc) *Template {
	return &Template{fn, template.New(name)}
}

// Name gets the name that was passed in the New function
func (t *Template) Name() string {
	return t.tmpl.Name()
}

// Funcs is a proxy to the underlying template's Funcs function
func (t *Template) Funcs(funcMap FuncMap) *Template {
	return t.replaceTmpl(t.tmpl.Funcs(template.FuncMap(funcMap)))
}

//Delims is a proxy to the underlying template's Delims function
func (t *Template) Delims(left, right string) *Template {
	return t.replaceTmpl(t.tmpl.Delims(left, right))
}

// Parse looks up the filename in the underlying Asset store,
// then calls the underlying template's Parse function with the result.
// returns an error if the file wasn't found or the Parse call failed
func (t *Template) Parse(filename string) (*Template, error) {
	tmplBytes, err := t.file(filename)
	if err != nil {
		return nil, err
	}
	newTmpl, err := t.tmpl.Parse(string(tmplBytes))
	if err != nil {
		return nil, err
	}
	return t.replaceTmpl(newTmpl), nil
}

// ParseFiles looks up all of the filenames in the underlying Asset store,
// concatenates the file contents together, then calls the underlying template's
// Parse function with the result. returns an error if any of the files
// don't exist or the underlying Parse call failed.
func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	fileBytes := []byte{}
	for _, filename := range filenames {
		tmplBytes, err := t.file(filename)
		if err != nil {
			return nil, err
		}
		fileBytes = append(fileBytes, tmplBytes...)
	}
	newTmpl, err := t.tmpl.Parse(string(fileBytes))
	if err != nil {
		return nil, err
	}
	return t.replaceTmpl(newTmpl), nil
}

// Execute is a proxy to the underlying template's Execute function
func (t *Template) Execute(w io.Writer, data interface{}) error {
	return t.tmpl.Execute(w, data)
}

// replaceTmpl is a convenience function to replace t.tmpl with the given tmpl
func (t *Template) replaceTmpl(tmpl *template.Template) *Template {
	t.tmpl = tmpl
	return t
}

// file is a convenience function to look up fileName using t.AssetFunc, then
// return the contents or an error if the file doesn't exist
func (t *Template) file(fileName string) ([]byte, error) {
	tmplBytes, err := t.AssetFunc(fileName)
	if err != nil {
		return nil, err
	}
	return tmplBytes, nil
}
