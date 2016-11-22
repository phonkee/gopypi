package main

import (
	"os"

	"path"

	"fmt"
	"io/ioutil"
	"strings"

	"github.com/phonkee/go-xmlrpc"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "Filename",
		},
		cli.BoolFlag{
			Name: "debug",
		},
	}

	app.Action = func(c *cli.Context) error {

		var (
			err error
			gen xmlrpc.Generator
		)

		filename := c.String("file")

		// instantiate generator
		if gen, err = xmlrpc.NewGenerator(filename); err != nil {
			return err
		}

		for i := 0; i < c.NArg(); i++ {
			if err = gen.AddService(c.Args().Get(0)); err != nil {
				return err
			}
		}

		// add service

		// print
		result := gen.Format()

		if c.Bool("debug") {
			print(string(result))
		}

		name := strings.TrimSuffix(filename, path.Ext(path.Base(filename)))
		target := fmt.Sprintf("%v_xmlrpc.go", name)
		ioutil.WriteFile(target, []byte(result), 0666)

		return nil
	}
	app.Run(os.Args)
}
