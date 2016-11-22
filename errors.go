package core

import "errors"

var (

	// Config
	ErrUnknownDBDriver = errors.New("Unknown database driver")

	// Auth
	ErrInvalidAuthHeader = errors.New("Invalid authorization header")
	ErrLoginUsername     = errors.New("invalid username")
	ErrLoginPassword     = errors.New("invalid password")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenInvalid      = errors.New("invalid token")
	ErrTokenUserInvalid  = errors.New("invalid user in token")

	// Model errors
	ErrUsernameAlreadyExists = errors.New("user with this username already exists")
	ErrPasswordsMustMatch    = errors.New("passwords must match")
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrEmailAlreadyExists    = errors.New("user with this email already exists")

	ErrLicenseNotFound  = errors.New("license not found")
	ErrPlatformNotFound = errors.New("license not found")

	ErrPackageNotFound = errors.New("package not found")

	// http errors
	ErrUserCannotRetrievePackages = errors.New("user cannot retrieve packages")
	ErrUserCannotDownloadPackages = errors.New("user cannot download packages")

	// Post package errors
	ErrPostPackageInvalidAction  = errors.New("action not recognized")
	ErrPostPackageInvalidName    = errors.New("invalid name")
	ErrPostPackageInvalidVersion = errors.New("invalid version")

	// generic error for all methods that return single object
	ErrObjectNotFound = errors.New("object not found")
)
