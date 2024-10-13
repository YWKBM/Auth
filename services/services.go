package services

import "auth/repo"

type AuthorizationService interface {
	CreateUser(login, password, email string) error
	CreateProvider() (int, error)
	ChangePassword(userId int, oldPassword, newPassword string) error
	CreateTokenPair(login, password string) (string, string, error)
	DeleteTokenPair(userId int) error
	ParseAccessToken(accessToken string) (int, error)
	RenewToken(refreshToken string) (string, string, error)
}

type Services struct {
	AuthService AuthorizationService
}

func NewServices(repos *repo.Repos) *Services {
	return &Services{
		AuthService: NewAuthService(repos),
	}
}
