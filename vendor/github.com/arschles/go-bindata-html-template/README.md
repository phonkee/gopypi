# go-bindata-html-template

[![GoDoc](https://godoc.org/github.com/arschles/go-bindata-html-template?status.svg)](https://godoc.org/github.com/arschles/go-bindata-html-template)

go-bindata-html-template is a wrapper for Go's built in
[`html/template`](godoc.org/html/template) package to work with template
contents embedded with the go-bindata tool instead of contents on the
filesystem See https://github.com/jteeuwen/go-bindata for more information
about embedding binary data with go-bindata.

It's compatible with a subset of the functionality in `html/template`.

Example usage (after running `go-bindata data/...` in your project directory):

```go
import (
  "github.com/arschles/go-bindata-html-template"
)

//...

func myHandler(res http.ResponseWriter, req *http.Request) {
  tmpl, err := template.New("mytmpl", Asset).Parse("data/templates/my.tmpl")
  if err != nil {
    log.Fatalf("error parsing template: %s", err)
  }
  err := tmpl.Execute(res)
  if err != nil {
    log.Fatalf("error executing template: %s", err)
  }
}
```
