package services

import "auth/repo"

type AuthorizationService interface {
	CreateUser(login, password, email string) error
	CreateProvider() (int, error)
	CreateTokenPair(login, password string) (string, string, error)
	// ParseAccessToken(accessToken string) (int, error)
	RenewToken(refreshToken string) (string, string, error)
}

type Services struct {
	AuthService AuthorizationService
}

func NewServicess(repos *repo.Repos) *Services {
	return &Services{
		AuthService: NewAuthService(repos),
	}
}
