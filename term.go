package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
)

/*
Exit exits with priting error message.
*/
func Exit(why interface{}, code ...int) {
	var message string

	exitcode := 1
	if len(code) > 0 {
		exitcode = code[0]
	}

	switch why := why.(type) {
	case string:
		message = why
	case fmt.Stringer:
		message = why.String()
	case error:
		message = why.Error()
	default:
		message = fmt.Sprintf("%v", why)
	}

	println(message)

	// now die
	os.Exit(exitcode)
}

/*
TerminalGetPasswordValue return password
*/
func TerminalGetPasswordValue(question string) (result string) {
	fmt.Printf("%v> ", question)

	var b []byte
	b, _ = gopass.GetPasswdMasked()
	result = string(b)
	return
}

/*
GetTermValue calls user input.
*/
func TerminalGetStringValue(question string, defs ...string) (result string) {

	def := ""
	if len(defs) > 0 {
		def = strings.TrimSpace(defs[0])
	}

	if def != "" {
		fmt.Printf("%v [default:%v]> ", question, def)
	} else {
		fmt.Printf("%v> ", question)
	}

	reader := bufio.NewReader(os.Stdin)
	if response, err := reader.ReadString('\n'); err != nil {
		panic(err)
	} else {
		result = strings.TrimSpace(response)
	}

	if result == "" {
		return def
	}

	return
}

func TerminalGetStringValueChoices(question string, choices []string, defs ...string) (result string) {
	def := ""
	if len(defs) > 0 {
		def = strings.TrimSpace(defs[0])
	}

outer:
	for {
		if def != "" {
			fmt.Printf("%v [choices: %v, default:%v]> ", question, choices, def)
		} else {
			fmt.Printf("%v [choices: %v]> ", question, choices)
		}

		reader := bufio.NewReader(os.Stdin)
		if response, err := reader.ReadString('\n'); err != nil {
			panic(err)
		} else {
			result = strings.TrimSpace(response)
		}

		if result == "" {
			result = def
		}

		for _, choice := range choices {
			if choice == result {
				break outer
			}
		}
	}

	return result
}

/*
GetTermIntValue waits for user input for number
*/
func TerminalGetIntValue(question string, def int) (result int) {

	result = def

	var (
		err error
	)

	for {
		value := TerminalGetStringValue(question, fmt.Sprintf("%v", def))
		if value == "" {
			return def
		}
		if result, err = strconv.Atoi(value); err == nil {
			return
		}
		println("error:", err.Error())
	}
	return def
}

func TerminalGetBoolValue(question string, def bool) (result bool) {
	result = def

	for {
		value := TerminalGetStringValue(question, fmt.Sprintf("%v", def))
		if value == "" {
			return def
		}

		if boolTrue.MatchString(value) {
			return true
		} else if boolFalse.MatchString(value) {
			return false
		} else {
			println("error:", "please answer")
		}
	}

	return
}
