package auth

import (
	"os"
	"testing"
)

func TestHashPassword(t *testing.T) {
	if _, err := HashPassword("testpassword1!"); err != nil {
		t.Error(err)
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name           string
		unhashed       string
		hashed         string
		shouldValidate bool
	}{
		{
			"Should match",
			"testpassword1!",
			"$2a$15$S5TatUOm3iLafpVZr6syRucrE6cDP9XKkoehRHP5vMj..NYIdi5aS",
			true,
		},
		{
			"Should not match",
			"testpassword2!",
			"$2a$15$S5TatUOm3iLafpVZr6syRucrE6cDP9XKkoehRHP5vMj..NYIdi5aS",
			false,
		},
	}
	for _, tt := range tests {
		valid, _ := ValidatePassword(tt.unhashed, tt.hashed)
		if (tt.shouldValidate && !valid) || (!tt.shouldValidate && valid) {
			t.Fail()
		}
	}
}

func TestCreateAccessToken(t *testing.T) {
	if os.Getenv("SECRET_TOKEN") == "" {
		t.Skip("missing secret key")
	}

	if _, err := CreateAccessToken(); err != nil {
		t.Error(err)
	}
}

func TestValidateToken(t *testing.T) {
	if os.Getenv("SECRET_TOKEN") == "" {
		t.Skip("missing secret key")
	}

	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			"Valid token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMyNTAzNzAxNjAwLCJpYXQiOjE2OTY5NzYyNTAsImlzcyI6InR1bmVzLXNlcnZlciJ9.hGVJKbOiGUxVbpNuN9dbspyHhwSC7DU60u-wFM5obMo",
			true,
		},
		{
			"Expired token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTY5NzcxNTAsImlhdCI6MTY5Njk3NjI1MCwiaXNzIjoidHVuZXMtc2VydmVyIn0.llmj0EyWrFMOfqMZS8hihQojiEYNYMZGtW4DbFgvHdA",
			false,
		},
		{
			"Malformed token",
			"not.a.token",
			false,
		},
		{
			"Invalid signature",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTY5NzY5NjIsImlhdCI6MTY5Njk3NjA2MiwiaXNzIjoidHVuZXMtc2VydmVyIn0.lXxQksdtkbODp5H17LrGpOwWVLB8PACZw-aN6ny9SzQ",
			false,
		},
		{
			"Tampered payload",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMDI0NDQ4MDAsImlhdCI6MTY5Njk3NjI1MCwiaXNzIjoidHVuZXMtc2VydmVyIn0.llmj0EyWrFMOfqMZS8hihQojiEYNYMZGtW4DbFgvHdA",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := ValidateToken(tt.token)
			if actual != tt.expected {
				t.Fail()
			}
		})
	}
}

func TestShouldRefresh(t *testing.T) {
	if os.Getenv("SECRET_TOKEN") == "" {
		t.Skip("missing secret key")
	}

	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			"Valid expired token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTY5NzcxNTAsImlhdCI6MTY5Njk3NjI1MCwiaXNzIjoidHVuZXMtc2VydmVyIn0.llmj0EyWrFMOfqMZS8hihQojiEYNYMZGtW4DbFgvHdA",
			true,
		},
		{
			"Valid non-expired token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMyNTAzNzAxNjAwLCJpYXQiOjE2OTY5NzYyNTAsImlzcyI6InR1bmVzLXNlcnZlciJ9.hGVJKbOiGUxVbpNuN9dbspyHhwSC7DU60u-wFM5obMo",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ShouldRefresh(tt.token)
			if actual != tt.expected {
				t.Fail()
			}
		})
	}
}
