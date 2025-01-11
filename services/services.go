package services

import (
	"auth/config"
	"auth/repo"
)

type AuthorizationService interface {
	CreateUser(login, password, email string) error
	ChangePassword(userId int, oldPassword, newPassword string) error
	CreateTokenPair(login, password string) (string, string, error)
	DeleteTokenPair(userId int) error
	ParseAccessToken(accessToken string) (int, error)
	RenewToken(refreshToken string) (string, string, error)
	ResolveAccess(accessToken string, expectedRole string) error
}

type ProviderIngterface interface {
	RequestCreateProvider(first_name, middle_name, second_name, email, phone string) error
	CreateProvider(login, password string) error
}

type Services struct {
	AuthService AuthorizationService
}

func NewServices(repos *repo.Repos, config config.Config) *Services {
	return &Services{
		AuthService: NewAuthService(repos, config.SECRET_KEY),
	}
}
