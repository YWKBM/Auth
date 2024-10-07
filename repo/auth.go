package repo

import (
	"auth/entities"
	"database/sql"
	"errors"
	"time"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(login, password, email string) error {
	var userId int

	err := r.db.QueryRow("SELECT Id FROM User WHERE Login == $1 OR Email == $2", login, email).Scan(userId)
	if err != sql.ErrNoRows {
		return errors.New("Пользователь уже зарегистрирован")
	}

	err = r.db.QueryRow("INSERT INTO User (email, password, login) values ($1, $2, $3) RETURNING Id", login, password, email).Scan(userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) CreateToken(jti string, userId int, expiry time.Time) error {
	_, err := r.db.Exec("UPDATE UserToken SET Jti = $1, Expiry = $2 WHERE UserId = $3", jti, expiry, userId)
	_, err = r.db.Exec("INSERT INTO UserToken SET (Jti, Expiry, UserId) VALUES ($1, $2, $3) WHERE NOT EXISTS (SELECT 1 FROM UserToken WHERE UserId = $3)", jti, expiry, userId)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) GetUserByTokenId(jti string) (int, entities.Role, error) {
	user := &entities.User{}
	token := &entities.UserToken{}

	err := r.db.QueryRow("SELECT * FROM User JOIN UserToken ON User.Id = UserToken.UserId WHERE UserToken.Jti == $1", jti).
		Scan(user.Id, user.Email, user.UserRole, user.Password, token.Id, token.Jti, token.Expiry)

	if err == sql.ErrNoRows {
		return 0, "", errors.New("некорректный RefreshToken")
	}

	return user.Id, user.UserRole, nil
}

func (r *AuthRepo) GetUser(login, password string) (entities.User, error) {
	user := &entities.User{}

	err := r.db.QueryRow("SELECT * FROM User WHERE Login == $1 AND Password == $2", login, password).Scan(user.Id, user.Login, user.Password, user.Email, user.UserRole)
	if err != nil {
		return *user, err
	}

	return *user, nil
}
