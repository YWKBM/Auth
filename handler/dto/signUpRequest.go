package dto

import (
	"errors"
	"net/mail"
)

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r *SignUpRequest) Validate() error {
	if r.Login == "" {
		return errors.New("invalid login")
	}

	if len(r.Password) < 6 {
		return errors.New("password should be at least 6 characters")
	}

	if !validEmail(r.Email) {
		return errors.New("invalid email")
	}

	return nil
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
