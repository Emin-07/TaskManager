package service

import (
	"fmt"
	"net/mail"
)

const passwordMinLength = 8
const passwordMaxLength = 32

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email can't be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("expected email like example@example.com, got '%v'. err: %s", email, err)
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < passwordMinLength {
		return fmt.Errorf("length of the password should be greater or equal to %d, not %v", passwordMinLength, len(password))
	}
	if len(password) > passwordMaxLength {
		return fmt.Errorf("length of the password should be less than or equal to %d, not %v", passwordMaxLength, len(password))
	}
	return nil
}

func validateCreds(email, password string) error {
	err := validateEmail(email)
	if err != nil {
		return err
	}
	return validatePassword(password)
}
