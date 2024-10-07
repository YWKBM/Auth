package repo

import (
	"auth/entities"
	"time"

	"github.com/go-pg/pg/v10/orm"
)

type Authorization interface {
	CreateUser(login, password, email string) error
	GetUser(login, password string) (entities.User, error)
	CreateToken(jti string, userId int, expiry time.Time) error
	GetUserByTokenId(jti string) (int, entities.Role, error)
}

type Repos struct {
	Authorization AuthRepo
}

func NewRepos(db orm.DB) *Repos {
	return &Repos{
		Authorization: *NewAuthRepo(db),
	}
}
