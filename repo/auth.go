package repo

import (
	"auth/entities"
	"errors"
	"time"

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

func (r *AuthRepo) CreateToken(jti string, userId int, expiry time.Time) error {
	token := &entities.UserToken{}

	err := r.db.Model(token).Where("token.UserId == ?", userId).First()
	if err != nil {
		return err
	}

	token.Jti = jti
	token.UserId = userId
	token.Expiry = expiry

	return nil
}

func (r *AuthRepo) GetUserByTokenId(jti string) (int, entities.Role, error) {
	token := &entities.UserToken{}

	err := r.db.Model(token).Where("token.Jti == ?", jti).First()
	if err != nil {
		return 0, "", err
	}

	if token.Id < 0 {
		return 0, "", errors.New("некорректный RefreshToken")
	}

	user := &entities.User{}

	err = r.db.Model(user).Where("user.TokenId = ?", token.Id).First()

	if user.Id < 0 {
		return 0, "", errors.New("пользователь не найден")
	}

	return user.Id, user.UserRole, nil

}

func (r *AuthRepo) GetUser(login, password string) (entities.User, error) {
	user := &entities.User{}

	err := r.db.Model(user).Where("user.Login == ?", login).Where("user.Password == ?", password).First()
	if err != nil {
		return *user, err
	}

	if user.Id < 0 {
		return *user, errors.New("Пользователь не найден")
	}

	return *user, nil
}
