package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "alice@example.com", false},
		{"empty email", "", true},
		{"missing domain", "alice@", true},
		{"missing at", "alice.com", true},
		{"no tld", "alice@example", false},
		{"with plus", "alice+tag@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"exactly 8 chars", "12345678", false},
		{"exactly 32 chars", "12345678901234567890123456789012", false},
		{"too short", "1234567", true},
		{"too long", "123456789012345678901234567890123", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCreds(t *testing.T) {
	t.Run("valid creds", func(t *testing.T) {
		assert.NoError(t, validateCreds("a@b.com", "password1"))
	})

	t.Run("invalid email", func(t *testing.T) {
		assert.Error(t, validateCreds("not-an-email", "password1"))
	})

	t.Run("invalid password", func(t *testing.T) {
		assert.Error(t, validateCreds("a@b.com", "short"))
	})
}
