package repo

import (
	"auth/entities"
	"errors"

	"github.com/go-pg/pg/v10/orm"
)

type AuthRepo struct {
	db orm.DB
}

func NewAuthRepo(db orm.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(login, password, email string) error {
	user := &entities.User{}

	err := r.db.Model(user).Where("user.Login == ?", login).WhereOr("user.Email == ?", email).First()

	if user.Id > 0 {
		return errors.New("Пользователь уже зарегистрирован")
	}

	user.Login = login
	user.Email = email
	user.Password = password

	_, err = r.db.Model(user).Insert()
	if err != nil {
		return err
	}

	return nil
}
