package main

import (
	"os"

	"github.com/phonkee/gopypi"

	"fmt"

	"github.com/urfave/cli"
)

func main() {

	// print logo before every command
	core.CliApp.Before = func(c *cli.Context) (err error) {
		// yeah print logo
		fmt.Println(core.LOGO)
		return
	}

	core.CliApp.Run(os.Args)
}
