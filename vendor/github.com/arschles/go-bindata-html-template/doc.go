// Package template allows standard html/template templates to be rendered from
// contents embedded with the go-bindata tool instead of the filesystem
//
// See https://github.com/jteeuwen/go-bindata for more information
// about embedding binary data with go-bindata.
//
// Usage example, after running
//  $ go-bindata data/...
// use:
//  tmpl, err := template.New("mytmpl", Asset).Parse("data/templates/my.tmpl")
//  if err != nil {
//    log.Fatalf("error parsing template: %s", err)
//  }
//  err := tmpl.Execute(os.Stdout)
//  if err != nil {
//    log.Fatalf("error executing template: %s", err)
//  }
package template
