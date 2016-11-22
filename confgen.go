package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// template used for confgen
const CONF_TEMPLATE = `
[core]
listen = '{{.listen}}'
secret_key = '{{.secret_key}}'

[database]
driver = '{{.driver}}'
dsn = '{{.dsn}}'

[packages]
directory = '{{.packages_dir}}'

[download_stats]
archive_weekly = 4
archive_monthly = 4
`

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("test").Parse(CONF_TEMPLATE))
}

/*
NewConfGen returns ConfGen instance
*/
func NewConfGen() *confgen {
	return &confgen{}
}

/*
confgen structure
*/
type confgen struct{}

/*
Run runs confgenerator
*/
func (c *confgen) Run() (err error) {

	data := map[string]interface{}{
		"filename":     TerminalGetStringValue("Please enter config filename", "gopypi.conf"),
		"listen":       TerminalGetStringValue("Please enter listen", "0.0.0.0:80"),
		"driver":       TerminalGetStringValueChoices("Please enter driver", []string{"postgres"}, "postgres"),
		"dsn":          TerminalGetStringValue("Please enter database dsn"),
		"packages_dir": TerminalGetStringValue("Please enter target packages directory", ".packages"),
		"secret_key":   fmt.Sprintf("%x", GenerateSalt(CONFGEN_SECRET_KEY_LENGTH)),
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, data); err != nil {
		return
	}

	targetFilename := data["filename"].(string)

	// write file
	if _, err := os.Stat(targetFilename); os.IsNotExist(err) {
		ioutil.WriteFile(targetFilename, []byte(strings.TrimLeft(buf.String(), " ")), 0666)
	} else {
		return fmt.Errorf("File %v already exists.", targetFilename)
	}

	return
}
