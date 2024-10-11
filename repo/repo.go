package repo

import (
	"auth/entities"
	"database/sql"
	"time"
)

type AuthorizationRepo interface {
	CreateUser(login, password, email string) error
	GetUser(login, password string) (entities.User, error)
	CreateToken(jti string, userId int, expiry time.Time) error
	DeleteToken(userId int) error
	GetUserByTokenId(jti string) (int, entities.Role, error)
}

type Repos struct {
	Authorization AuthorizationRepo
}

func NewRepos(db *sql.DB) *Repos {
	return &Repos{
		Authorization: NewAuthRepo(db),
	}
}
