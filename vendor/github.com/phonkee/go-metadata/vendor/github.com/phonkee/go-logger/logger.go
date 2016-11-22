package logger

import (
	"fmt"
	"strings"

	"github.com/mgutz/ansi"
)

// Simple Logger
type Logger func(enabled bool, format string, args ...interface{})

/*
Returns named logger that depends on enabled flag.
*/
func named(name string, color string) Logger {

	name = strings.TrimSpace(name)

	return func(enabled bool, format string, args ...interface{}) {

		// if not enabled skip
		if !enabled {
			return
		}

		// add newline
		format = strings.TrimRight(format, " \n\t")
		message := fmt.Sprintf(format, args...)

		if name != "" {
			message = fmt.Sprintf("%v %v", ansi.Color(name, color), message)
		}

		println(message)
	}
}

/*
Debug returns debug logger with given name
*/
func Debug(name string) Logger {
	return named(name + ".debug", "green")
}

/*
Info returns debug logger with given name
*/
func Info(name string) Logger {
	return named(name + ".info", "green")
}

/*
Error returns debug logger with given name
*/
func Error(name string) Logger {
	return named(name + ".error", "red")
}

/*
Warning returns debug logger with given name
*/
func Warning(name string) Logger {
	return named(name + ".warning", "red+h")
}
