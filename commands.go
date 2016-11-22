/*
Set of commands that can be run from command line.
*/
package core

import (
	"errors"
	"fmt"
	"os/user"

	"github.com/asaskevich/govalidator"
)

/*
Command interface that all commands use
*/
type Command interface {

	// Run runs command
	Run() error
}

/*
CreateAdminCommand creates new admin in database
*/
type CreateAdminCommand struct {
	Config Config
}

/*
Run is method that runs the command
*/
func (c *CreateAdminCommand) Run() error {

	u, err := user.Current()
	if err != nil {
		return err
	}
	user := User{}
	user.Username = TerminalGetStringValue("Please enter username", u.Username)
	if c.Config.Manager().User().ExistsUsername(user.Username) {
		return fmt.Errorf("User with username %s already exists", user.Username)
	}

	user.Email = TerminalGetStringValue("Please enter email")
	if user.Email != "" {
		if !govalidator.IsEmail(user.Email) {
			return fmt.Errorf("Invalid email %s", user.Email)
		}
		if c.Config.Manager().User().ExistsEmail(user.Email) {
			return fmt.Errorf("User with email %s already exists", user.Email)
		}
	}

	// Get password
	password1 := TerminalGetPasswordValue("Please enter password")
	if password1 == "" {
		return errors.New("No password supplied")
	}

	password2 := TerminalGetPasswordValue("Please retype password")

	if password1 != password2 {
		return errors.New("Passwords don't match")
	}

	// set password
	c.Config.Manager().User().SetPassword(&user, password1)

	user.IsActive = true
	user.IsAdmin = true

	if errSave := c.Config.DB().Save(&user).Error; errSave != nil {
		return errSave
	}

	println("Admin has been succesfully created.")

	return nil
}

/*
ChangePasswordCommand command line command to change password
*/
type ChangePasswordCommand struct {
	Config Config
}

/*
Run
*/
func (c *ChangePasswordCommand) Run() (err error) {

	u, err := user.Current()
	if err != nil {
		return err
	}
	user := User{}
	user.Username = TerminalGetStringValue("Please enter username", u.Username)
	if !c.Config.Manager().User().ExistsUsername(user.Username) {
		return fmt.Errorf("User with username %s doesn't exist.", user.Username)
	}

	if c.Config.Manager().User().Get(&user).RecordNotFound() {
		return fmt.Errorf("User with username %s doesn't exist.", user.Username)
	}

	// Get password
	password1 := TerminalGetPasswordValue("Please enter password")
	if password1 == "" {
		return errors.New("No password supplied")
	}

	password2 := TerminalGetPasswordValue("Please retype password")

	if password1 != password2 {
		return errors.New("Passwords don't match")
	}

	// set password
	c.Config.Manager().User().SetPassword(&user, password1)

	if err = c.Config.DB().Save(&user).Error; err != nil {
		return fmt.Errorf("Error when saving user to database: %s", err.Error())
	}

	println("Password successfully changed for user", user.Username, ".")

	return
}
