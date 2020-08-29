package payloads

import (
	"fmt"
	"go-cicd/app/validator"
)

// LoginEmailPayload payload for login email API
type LoginEmailPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validate login email payload
func (payload *LoginEmailPayload) Validate() (bool, error) {
	if payload.Email == "" || payload.Password == "" {
		return false, fmt.Errorf("Missing email or password field")
	}

	if validator.ValidateEmail(payload.Email) == false {
		return false, fmt.Errorf("Email is invalid")
	}

	return true, nil
}
