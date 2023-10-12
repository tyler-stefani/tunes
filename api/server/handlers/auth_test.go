package handlers

import (
	"fmt"
	"os"
	"testing"
	"time"

	"tunes-service/data"
)

type userDBMock struct {
	getUser    func(nameOrEmail string) (data.User, error)
	createUser func(name, email, password string) (data.User, error)
}

func (db *userDBMock) GetUser(nameOrEmail string) (data.User, error) {
	return db.getUser(nameOrEmail)
}

func (db *userDBMock) CreateUser(name, email, password string) (data.User, error) {
	return db.createUser(name, email, password)
}

type authDBMock struct {
	writeRefreshToken func(id string, expires time.Time) (bool, error)
	findRefreshToken  func(id string) (data.RefreshToken, error)
}

func (db *authDBMock) WriteRefreshToken(id string, expires time.Time) (bool, error) {
	return db.writeRefreshToken(id, expires)
}

func (db *authDBMock) FindRefreshToken(id string) (data.RefreshToken, error) {
	return db.findRefreshToken(id)
}

func TestHandleRegistration(t *testing.T) {
	tests := []struct {
		name     string
		username string
		email    string
		password string
		db       *userDBMock
		expected bool
	}{
		{
			"Should register successfully",
			"test",
			"test@test.com",
			"testpassword1!",
			&userDBMock{
				func(nameOrEmail string) (data.User, error) {
					t.Fatalf("should not call this function")
					return data.User{}, nil
				},
				func(name, email, password string) (data.User, error) {
					return data.User{
						ID:       1,
						Name:     name,
						Email:    email,
						Password: password,
					}, nil
				},
			},
			true,
		},
		{
			"Duplicate username",
			"test",
			"test@test.com",
			"testpassword1!",
			&userDBMock{
				func(nameOrEmail string) (data.User, error) {
					t.Fatalf("should not call this function")
					return data.User{}, nil
				},
				func(name, email, password string) (data.User, error) {
					return data.User{}, fmt.Errorf("duplicate username")
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := HandleRegistration(tt.username, tt.email, tt.password, tt.db)
			if actual != tt.expected {
				if tt.expected {
					t.Fatalf("expected ok but got error: %s", err.Error())
				} else {
					t.Fatalf("expected error but got none")
				}
			}
		})
	}
}

func TestHandleLogin(t *testing.T) {
	if os.Getenv("SECRET_KEY") == "" {
		t.Skip("missing secret key")
	}

	tests := []struct {
		name            string
		usernameOrEmail string
		password        string
		udb             *userDBMock
		adb             *authDBMock
		expected        bool
	}{
		{
			"Successful login",
			"test",
			"testpassword1!",
			&userDBMock{
				func(nameOrEmail string) (data.User, error) {
					return data.User{
						ID:       1,
						Name:     "test",
						Email:    "test@test.com",
						Password: "$2a$15$S5TatUOm3iLafpVZr6syRucrE6cDP9XKkoehRHP5vMj..NYIdi5aS",
					}, nil
				},
				func(name, email, password string) (data.User, error) {
					t.Fatalf("should not call this function")
					return data.User{}, nil
				},
			},
			&authDBMock{
				func(id string, expires time.Time) (bool, error) {
					return true, nil
				},
				func(id string) (data.RefreshToken, error) {
					t.Fatalf("should not call this function")
					return data.RefreshToken{}, nil
				},
			},
			true,
		},
		{
			"Wrong Password",
			"test",
			"testpassword2!",
			&userDBMock{
				func(nameOrEmail string) (data.User, error) {
					return data.User{
						ID:       1,
						Name:     "test",
						Email:    "test@test.com",
						Password: "$2a$15$S5TatUOm3iLafpVZr6syRucrE6cDP9XKkoehRHP5vMj..NYIdi5aS",
					}, nil
				},
				func(name, email, password string) (data.User, error) {
					t.Fatalf("should not call this function")
					return data.User{}, nil
				},
			},
			&authDBMock{
				func(id string, expires time.Time) (bool, error) {
					return true, nil
				},
				func(id string) (data.RefreshToken, error) {
					t.Fatalf("should not call this function")
					return data.RefreshToken{}, nil
				},
			},
			false,
		},
		{
			"User doesn't exist",
			"test",
			"testpassword1!",
			&userDBMock{
				func(nameOrEmail string) (data.User, error) {
					return data.User{}, fmt.Errorf("user not found")
				},
				func(name, email, password string) (data.User, error) {
					t.Fatalf("should not call this function")
					return data.User{}, nil
				},
			},
			&authDBMock{
				func(id string, expires time.Time) (bool, error) {
					return true, nil
				},
				func(id string) (data.RefreshToken, error) {
					t.Fatalf("should not call this function")
					return data.RefreshToken{}, nil
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := HandleLogin(tt.usernameOrEmail, tt.password, tt.udb, tt.adb)
			if tt.expected != (err == nil) {
				if tt.expected {
					t.Fatalf("expected ok but got error: %s", err.Error())
				} else {
					t.Fatalf("expected error but got none")
				}
			}
		})
	}
}

func TestHandleRefresh(t *testing.T) {
	if os.Getenv("SECRET_KEY") == "" {
		t.Skip("missing secret key")
	}

	tests := []struct {
		name           string
		oldAccessToken string
		refreshToken   string
		adb            *authDBMock
		expected       bool
	}{
		{
			"Should Refresh",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTY5NzcxNTAsImlhdCI6MTY5Njk3NjI1MCwiaXNzIjoidHVuZXMtc2VydmVyIn0.llmj0EyWrFMOfqMZS8hihQojiEYNYMZGtW4DbFgvHdA",
			"test",
			&authDBMock{
				func(id string, expires time.Time) (bool, error) {
					t.Fatalf("should not call this function")
					return true, nil
				},
				func(id string) (data.RefreshToken, error) {
					return data.RefreshToken{
						ID:         id,
						Expiration: time.Now().Add(time.Hour),
					}, nil
				},
			},
			true,
		},
		{
			"Token not expired",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMyNTAzNzAxNjAwLCJpYXQiOjE2OTY5NzYyNTAsImlzcyI6InR1bmVzLXNlcnZlciJ9.hGVJKbOiGUxVbpNuN9dbspyHhwSC7DU60u-wFM5obMo",
			"test",
			&authDBMock{
				func(id string, expires time.Time) (bool, error) {
					t.Fatalf("should not call this function")
					return true, nil
				},
				func(id string) (data.RefreshToken, error) {
					return data.RefreshToken{
						ID:         id,
						Expiration: time.Now().Add(time.Hour),
					}, nil
				},
			},
			false,
		},
		{
			"Refresh token not found",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTY5NzcxNTAsImlhdCI6MTY5Njk3NjI1MCwiaXNzIjoidHVuZXMtc2VydmVyIn0.llmj0EyWrFMOfqMZS8hihQojiEYNYMZGtW4DbFgvHdA",
			"test",
			&authDBMock{
				func(id string, expires time.Time) (bool, error) {
					t.Fatalf("should not call this function")
					return true, nil
				},
				func(id string) (data.RefreshToken, error) {
					return data.RefreshToken{}, fmt.Errorf("refresh token not found")
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := HandleRefresh(tt.oldAccessToken, tt.refreshToken, tt.adb); tt.expected != (err == nil) {
				if tt.expected {
					t.Fatalf("expected ok but got error: %s", err.Error())
				} else {
					t.Fatalf("expected error but got none")
				}
			}
		})
	}
}
