package core

import "regexp"

const (
	VERSION = "0.5.3"
	LOGO    = `
                                                              o
                                                            _<|>_

   o__ __o/    o__ __o    \o_ __o     o      o   \o_ __o      o
  /v     |    /v     v\    |    v\   <|>    <|>   |    v\    <|>
 />     / \  />       <\  / \    <\  < >    < >  / \    <\   / \
 \      \o/  \         /  \o/     /   \o    o/   \o/     /   \o/
  o      |    o       o    |     o     v\  /v     |     o     |
  <\__  < >   <\__ __/>   / \ __/>      <\/>     / \ __/>    / \
         |                \o/            /       \o/
 o__     o                 |            o         |
 <\__ __/>                / \        __/>        / \

                                                     v ` + VERSION
)

var (
	AVAILABLE_DB_DRIVERS = []string{"postgres", "mysql"}
	DEFAULT_DB_DRIVER    = "postgres"
)

// token related constants
const (
	TOKEN_HEADER_NAME = "Authorization"
	TOKEN_ISSUER      = "gopypi"
	TOKEN_EXPIRATION  = 24 * 3600
)

// Context constants
const (
	CONTEXT_TOKEN_USER = iota + 1000
	CONTEXT_ROUTE_NAME
)

// confgen constants
const (
	CONFGEN_SECRET_KEY_LENGTH = 64
)

// WHen setup.py calls gopypi it can define following actions
const (
	POST_PACKAGE_ACTION_VERIFY = iota + 1
	POST_PACKAGE_ACTION_SUBMIT
	POST_PACKAGE_ACTION_DOC_UPLOAD
	POST_PACKAGE_ACTION_REMOVE_PKG
	POST_PACKAGE_ACTION_FILE_UPLOAD

	// find out how this works
	POST_PACKAGE_ACTION_USER
	POST_PACKAGE_ACTION_PASSWORD_RESET
)

// constants for pasword hashing
const (
	PASSWORD_SALT_BYTES = 32
	PASSWORD_HASH_BYTES = 128
	PASSWORD_ITERATIONS = 16384
)

var (
	// when package is posted following actions can apply
	POST_PACKAGE_ACTIONS = map[string]int{
		"verify":         POST_PACKAGE_ACTION_VERIFY,
		"submit":         POST_PACKAGE_ACTION_SUBMIT,
		"doc_upload":     POST_PACKAGE_ACTION_DOC_UPLOAD,
		"remove_pkg":     POST_PACKAGE_ACTION_REMOVE_PKG,
		"file_upload":    POST_PACKAGE_ACTION_FILE_UPLOAD,
		"user":           POST_PACKAGE_ACTION_USER,
		"password_reset": POST_PACKAGE_ACTION_PASSWORD_RESET,
	}
)

// Model related constants
const (
	PASSWORD_MIN_LENGTH = 5
	PASSWORD_MAX_LENGTH = 20

	USERNAME_MIN_LENGTH = 5
	USERNAME_MAX_LENGTH = 20
)

// Regular expressions for terminal input handling functions
var (
	boolTrue  *regexp.Regexp
	boolFalse *regexp.Regexp
)

func init() {
	boolTrue = regexp.MustCompile("(?i)^(yes|y|t|true)$")
	boolFalse = regexp.MustCompile("(?i)^(no|n|f|false)$")
}
