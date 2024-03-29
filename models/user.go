package models

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

//User is a generated model from buffalo-auth, it serves as the base for username/password authentication.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`

	Password             string `json:"password" db:"-"`
	PasswordConfirmation string `json:"password_confirmation" db:"-"`
	Tokens               string `json:"tokens" db:"tokens"`
	// Phases        Phases     `has_many:"phases"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	fmt.Printf("Emaail: %s\n", u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	fmt.Printf("Password: %s\n", u.PasswordHash)
	return tx.ValidateAndCreate(u)
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Function to convert
func (u User) View() UserView {
	user := UserView{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt}
	return user
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// AddUser example
type AddUser struct {
	email                 string `json:"email" example:"example_user@example.com"`
	password              string `json:"password" example:"password"`
	password_confirmation string `json:"password_confirmation" example:"password"`
}

// AuthenticateUser example
type AuthenticateUser struct {
	email    string `json:"email" example:"example_user@example.com"`
	password string `json:"password" example:"password"`
}

// Authentication Header Example
type AuthenticationHeader struct {
	uid         string `json: "uid" example: "example_user@example.com"`
	AccessToken string `json: "access-token" example: "oUfHwqhMHQu-r28_n6Crh5SWV7vcCGJis3NuZOwOjBU="`
	client      string `json: "client" example: "_nF39TTNZrZmtFUlwhdSvVuH3ea42JCmiNquRYMqkjs="`
	expiry      string `json: "expiry" example: "1564117755395"`
}

// User View
type UserView struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Email     string    `json:"email" db:"email"`
}
