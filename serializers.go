package core

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
)

/*
FeatureSerializer enables/disables feature
*/
type FeatureSerializer struct {
	Value bool `json:"value"`
}

/*
LoginSerializer is login form for json
*/
type LoginSerializer struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *LoginSerializer) Validate() (err error) {

	l.Username = strings.TrimSpace(l.Username)
	l.Password = strings.TrimSpace(l.Password)

	if l.Username == "" {
		return ErrLoginUsername
	}

	if l.Password == "" {
		return ErrLoginPassword
	}

	return
}

/*
Serializer for updating user profile
*/
type UserProfileSerializer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

/*
Validate validates serializer
*/
func (u *UserProfileSerializer) Validate(cfg Config) (result ValidationResult) {
	result = NewValidationResult()

	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)

	if u.Email != "" && !govalidator.IsEmail(u.Email) {
		result.AddFieldError("email", errors.New("Invalid email"), "invalid_email")
	}

	return
}

/*
Serializer to create new user
*/
type UserAddSerializer struct {
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Password2   string `json:"password2"`
	IsActive    bool   `json:"is_active"`
	IsAdmin     bool   `json:"is_admin"`
	CanList     bool   `json:"can_list"`
	CanDownload bool   `json:"can_download"`
	CanCreate   bool   `json:"can_create"`
	CanUpdate   bool   `json:"can_update"`
}

/*
Validate validates information about new user

@TODO: add ValidateUsername and ValidateEmail
*/
func (u *UserAddSerializer) Validate(cfg Config) (result ValidationResult) {
	result = NewValidationResult()
	u.Username = strings.TrimSpace(u.Username)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)
	u.Password2 = strings.TrimSpace(u.Password2)

	// validate username
	if ValidateUsername("username", &(u.Username), result) {
		user := User{}
		if !cfg.DB().First(&user, "username = ?", u.Username).RecordNotFound() {
			result.AddFieldError("username", ErrUsernameAlreadyExists)
		}
	}

	// validate email
	if u.Email != "" {
		if !govalidator.IsEmail(u.Email) {
			result.AddFieldError("email", ErrInvalidEmail)
		} else {
			user := User{}
			if !cfg.DB().First(&user, "email = ?", u.Email).RecordNotFound() {
				result.AddFieldError("email", ErrEmailAlreadyExists)
			}
		}
	}

	// validate passwords
	ValidatePassword("password", &(u.Password), result)
	ValidatePassword("password2", &(u.Password2), result)

	if !result.HasFieldError("password") && !result.HasFieldError("password2") {
		if u.Password != u.Password2 {
			result.AddFieldError("password2", ErrPasswordsMustMatch)
		}
	}

	return
}

/*
return user initialized with serializer values
*/
func (u *UserAddSerializer) GetUser(cfg Config) (user User) {
	user = User{
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		IsActive:    u.IsActive,
		IsAdmin:     u.IsAdmin,
		CanCreate:   u.CanCreate,
		CanList:     u.CanList,
		CanDownload: u.CanDownload,
		CanUpdate:   u.CanUpdate,
	}

	// set password
	cfg.Manager().User().SetPassword(&user, u.Password)

	return
}

/*
Serializer to update existing user
*/
type UserUpdateSerializer struct {
	ID          uint   `json:"-"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Password2   string `json:"password2"`
	IsActive    bool   `json:"is_active"`
	IsAdmin     bool   `json:"is_admin"`
	CanList     bool   `json:"can_list"`
	CanDownload bool   `json:"can_download"`
	CanCreate   bool   `json:"can_create"`
	CanUpdate   bool   `json:"can_update"`

	// mark that user is changing password
	passwordChange bool
}

/*
Validate validates data in serializer
*/
func (u *UserUpdateSerializer) Validate(cfg Config) (result ValidationResult) {
	result = NewValidationResult()

	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)
	u.Password2 = strings.TrimSpace(u.Password2)

	// validate username
	if ValidateUsername("username", &(u.Username), result) {
		count := 0
		cfg.DB().Model(User{}).Where("id != ? AND username = ?", u.ID, u.Username).Count(&count)
		if count > 0 {
			result.AddFieldError("username", ErrUsernameAlreadyExists)
		}
	}

	// validate email
	if u.Email != "" {
		if !govalidator.IsEmail(u.Email) {
			result.AddFieldError("email", ErrInvalidEmail)
		} else {
			user := User{}
			if !cfg.DB().First(&user, "email = ? AND id != ?", u.Email, u.ID).RecordNotFound() {
				result.AddFieldError("email", ErrEmailAlreadyExists)
			}
		}
	}

	// if any password has been given perform validation
	if u.Password != "" || u.Password2 != "" {
		// mark that serializer changes password
		u.passwordChange = true

		ValidatePassword("password", &(u.Password), result)
		ValidatePassword("password2", &(u.Password2), result)

		if !result.HasFieldError("password") && !result.HasFieldError("password2") {
			if u.Password != u.Password2 {
				result.AddFieldError("password2", ErrPasswordsMustMatch)
			}
		}
	}

	return
}

/*
UpdateUser updates user with serializer data
*/
func (u *UserUpdateSerializer) UpdateUser(cfg Config, user *User) {
	user.Username = u.Username
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.Email = u.Email

	if u.passwordChange {
		cfg.Manager().User().SetPassword(user, u.Password)
	}
	user.IsActive = u.IsActive
	user.IsAdmin = u.IsAdmin
	user.CanList = u.CanList
	user.CanDownload = u.CanDownload
	user.CanCreate = u.CanCreate
	user.CanUpdate = u.CanUpdate
}

/*
LicenseUpdateSerializer handles license update
*/
type LicenseUpdateSerializer struct {
	ID       uint   `json:"-"`
	Approved bool   `json:"approved"`
	Name     string `json:"name"`
	Content  string `json:"content"`
}

/*
Validate runs validation on serializer fields
*/
func (l *LicenseUpdateSerializer) Validate(cfg Config, ID int) (result ValidationResult) {
	result = NewValidationResult()
	l.ID = uint(ID)
	l.Name = strings.TrimSpace(l.Name)
	l.Content = strings.TrimSpace(l.Content)
	return
}

/*
Update license with data form serializer
*/
func (l *LicenseUpdateSerializer) UpdateLicense(license *License) {
	license.Name = l.Name
	license.Content = l.Content
	license.Approved = l.Approved
}

/*
Serializer to change password for currently logged user
*/
type UserChangePasswordSerializer struct {
	Current   string `json:"current"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
	User      User   `json:"-"`
}

/*
Validate validates serializer data
*/
func (u *UserChangePasswordSerializer) Validate(cfg Config) (result ValidationResult) {
	result = NewValidationResult()

	// validate current password
	if !cfg.Manager().User().VerifyPassword(u.User, u.Current) {
		result.AddFieldError("current", ErrLoginPassword)
		return
	}

	ValidatePassword("password", &(u.Password), result)
	ValidatePassword("password2", &(u.Password2), result)

	if !result.HasFieldError("password") && !result.HasFieldError("password2") {
		if u.Password != u.Password2 {
			result.AddFieldError("password2", ErrPasswordsMustMatch)
		}
	}

	return
}

/*
ChangePassword changes password for given user
*/
func (u *UserChangePasswordSerializer) ChangePassword(config Config, user *User) (err error) {
	config.Manager().User().SetPassword(user, u.Password)
	return
}
