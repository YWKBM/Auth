package services

import (
	"auth/repo"
	"crypto/sha1"
	"fmt"
)

const (
	salt = "qweqweasddfasdfasdfqwerqwetasdg"
)

type Role string

const (
	provider Role = "PROVIDER"
	client   Role = "CLIENT"
)

type tokenClaims struct {
	UserId   int
	UserRole Role
}

type AuthService struct {
	repo *repo.AuthRepo
}

func (a *AuthService) CreateUser(login, password, email string) error {
	pass := generateHashPassword(password)
	err := a.repo.CreateUser(login, pass, email)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) CreateProvider() (int, error) {
	return 0, nil
}

func (a *AuthService) CreateToken(login, password string) (string, error) {
	return "", nil
}

func (a *AuthService) ParseToken(accessToken string) (int, error) {
	return 0, nil
}

func generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
