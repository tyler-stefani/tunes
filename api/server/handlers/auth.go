package handlers

import (
	"fmt"
	"time"

	"tunes-service/auth"
	"tunes-service/data"
)

func HandleRegistration(username string, email string, password string, db data.UserDB) (bool, error) {
	hashedPassword, _ := auth.HashPassword(password)

	_, err := db.CreateUser(username, email, hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to create user: %v", err)
	}

	return true, nil
}

func HandleLogin(usernameOrEmail string, password string, udb data.UserDB, adb data.AuthDB) (accessToken string, refreshToken string, err error) {
	u, err := udb.GetUser(usernameOrEmail)
	if err != nil {
		return accessToken, refreshToken, fmt.Errorf("user not found")
	}

	if ok, _ := auth.ValidatePassword(password, u.Password); !ok {
		return accessToken, refreshToken, fmt.Errorf("incorrect password")
	}

	accessToken, _ = auth.CreateAccessToken()
	refreshToken, refreshTokenExpiration := auth.CreateRefreshToken()

	adb.WriteRefreshToken(refreshToken, refreshTokenExpiration)

	return accessToken, refreshToken, nil
}

func HandleRefresh(oldAccessToken, refreshToken string, adb data.AuthDB) (newAccessToken string, err error) {
	if !auth.ShouldRefresh(oldAccessToken) {
		return "", fmt.Errorf("access token is invalid or does not need to be refreshed")
	}

	rt, err := adb.FindRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token was not found in the db")
	}

	if time.Now().After(rt.Expiration) {
		return "", fmt.Errorf("refresh token is expired")
	}

	return auth.CreateAccessToken()
}
