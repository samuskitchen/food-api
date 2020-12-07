package model

import (
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// Data of User
// swagger:model
type User struct {
	ID           string     `json:"id,omitempty"`
	// Required: true
	Names        string     `json:"names,omitempty"`
	// Required: true
	LastNames    string     `json:"last_names,omitempty"`
	// Required: true
	Email        string     `json:"email,omitempty"`
	// Required: true
	Password     string     `json:"password,omitempty"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `json:"-"`
}

// Information from user
// swagger:parameters getAllUser
type SwaggerAllUserRequest struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string
}

// Information from user
// swagger:parameters userRequest
type SwaggerUserRequest struct {
	// in: body
	Body User
}

// Information from user for update
// swagger:parameters userUpdateRequest
type SwaggerUserUpdateRequest struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string

	// in: path
	// Required: true
	ID string

	// in: body
	Body struct {
		// Required: true
		Names     string `json:"names,omitempty"`
		// Required: true
		LastNames string `json:"last_names,omitempty"`
		// Required: true
		Email     string `json:"email,omitempty"`
	}
}

// swagger:parameters idUserPath
type SwaggerUser struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string

	// in: path
	// Required: true
	ID string
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
			errorMessages["names_required"] = "names is required"
		}

		if u.LastNames == "" {
			errorMessages["lastnames_required"] = "last names is required"
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
