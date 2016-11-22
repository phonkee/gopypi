package xmlrpc

import (
	"errors"
	"fmt"
)

var (
	ErrMethodNotFound = errors.New("Method not found")
)

type Error interface {
	Code() int
	Error() string
}

/*
Errorf creates new xmlrpc error with given code
*/
func Errorf(code int, msg string, args ...interface{}) Error {
	return err{
		code:    code,
		message: fmt.Sprintf(msg, args...),
	}
}

/*
err is implementation of Error
*/
type err struct {
	message string
	code    int
}

/*
Error satisfies error interface
*/
func (e err) Error() string {
	return e.message
}

/*
Code is additional param when returning from service methods. Otherwise xmlrpc returns code 500
*/
func (e err) Code() int {
	return e.code
}
