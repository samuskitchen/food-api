package model

import (
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID           string     `json:"id,omitempty"`
	Names        string     `json:"names,omitempty"`
	LastNames    string     `json:"last_names,omitempty"`
	Email        string     `json:"email,omitempty"`
	Password     string     `json:"password,omitempty"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// HashPassword generates a hash of the password and places the result in PasswordHash.
func (u *User) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)

	return nil
}

// PasswordMatch compares HashPassword with the password and returns true if they match.
func (u User) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}

		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "email email"
			}
		}

	case "login":
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}

		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}

		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}

	case "forgot_password":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}

		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}

	default:
		if u.Names == "" {
			errorMessages["firstname_required"] = "names is required"
		}

		if u.LastNames == "" {
			errorMessages["lastname_required"] = "last names is required"
		}

		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}

		if u.Password != "" && len(u.Password) < 6 {
			errorMessages["invalid_password"] = "password should be at least 6 characters"
		}

		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}

		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	}

	return errorMessages
}
