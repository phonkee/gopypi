package core

import (
	"github.com/urfave/cli"
	"fmt"
)

var (
	CliApp     *cli.App
	configflag cli.StringFlag
)

/*
getconfig returns config along with error
*/
func getconfig(c *cli.Context) (result Config, err error) {
	if result, err = NewConfigFromFilename(c.String("config")); err != nil {
		err = exitError("Config returned error: %s", err)
		return
	}
	// validate configuration
	if err = result.Validate(); err != nil {
		err = exitError("Config returned error: %s", err)
		return
	}
	return
}

/*
exitError return ExitError by formatting
 */
func exitError(format string, args ...interface{}) *cli.ExitError {
	message := fmt.Sprintf(format, args...)
	return cli.NewExitError(message, 1)
}

func MigrateAction(c *cli.Context) (err error) {
	var cfg Config
	if cfg, err = getconfig(c); err != nil {
		return exitError("Migrate returned error: %s", err)
	}

	Migrate(cfg)

	return
}

func RunserverAction(c *cli.Context) (err error) {
	var cfg Config
	if cfg, err = getconfig(c); err != nil {
		return
	}

	var s Server
	if s, err = New(cfg); err != nil {
		return
	}
	if err = s.ListenAndServe(); err != nil {
		return exitError("Listen and serve returned error: %s", err)
	}

	return nil
}

func init() {

	configflag = cli.StringFlag{
		Name:   "config, c",
		Value:  "gopypi.conf",
		Usage:  "Config filename",
		EnvVar: "GOPYPI_CONFIG",
	}

	CommandRunserver := cli.Command{
		Name:    "runserver",
		Aliases: []string{"run"},
		Usage:   "Runs gopypi server (along with background tasks)",
		Flags: []cli.Flag{
			configflag,
		},
		Action: RunserverAction,
	}

	CommandMigrate := cli.Command{
		Name:     "migrate",
		Usage:    "Migrates database",
		Flags: []cli.Flag{
			configflag,
		},

		Action: MigrateAction,
	}

	CommandMakeConfig := cli.Command{
		Name:  "makeconfig",
		Usage: "Interactive build configuration file",
		Action: func(c *cli.Context) (err error) {
			cg := NewConfGen()
			if err = cg.Run(); err != nil {
				return
			}

			return
		},
	}

	// instantiate cli app
	CliApp = cli.NewApp()
	CliApp.Usage = "private pypi server implementation"
	CliApp.Version = VERSION
	CliApp.Commands = []cli.Command{
		{
			Name:  "createadmin",
			Usage: "create new admin user",
			Flags: []cli.Flag{
				configflag,
			},
			Action: func(c *cli.Context) (err error) {
				var cfg Config
				if cfg, err = getconfig(c); err != nil {
					return
				}

				var cmd Command
				cmd = &CreateAdminCommand{Config: cfg}

				if err = cmd.Run(); err != nil {
					return exitError("Createadmin returned error: %s", err)
				}

				return nil
			},
		},
		{
			Name:  "changepassword",
			Usage: "Change password for given user",
			Flags: []cli.Flag{
				configflag,
			},
			Action: func(c *cli.Context) (err error) {
				var cfg Config
				if cfg, err = getconfig(c); err != nil {
					return
				}

				var cmd Command
				cmd = &ChangePasswordCommand{Config: cfg}

				if err = cmd.Run(); err != nil {
					return
				}

				return nil
			},
		},
		CommandMakeConfig,
		CommandMigrate,
		CommandRunserver,
		{
			Name:  "cleanupdownloadstats",
			Usage: "Cleans up download stats",
			Flags: []cli.Flag{
				configflag,
			},
			Action: func(c *cli.Context) (err error) {
				var cfg Config
				if cfg, err = getconfig(c); err != nil {
					return
				}

				cfg.Manager().DownloadStats().Cleanup()


				return nil
			},
		},
	}
}
