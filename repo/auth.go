package repo

import (
	"auth/customErrors"
	"auth/entities"
	"database/sql"
	"time"
)

type AuthRepo struct {
	db *sql.DB
}

func newAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(login, password, email string) error {
	var userId int

	err := r.db.QueryRow("SELECT id FROM users WHERE login = $1 OR email = $2", login, email).Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if userId > 0 {
		return &customErrors.AlreadyExistsError{}
	}

	err = r.db.QueryRow("INSERT INTO users (login, password, email, userrole) values ($1, $2, $3, $4) RETURNING id", login, password, email, "USER").Scan(&userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) CreateToken(jti string, userId int, expiry time.Time) error {
	token := &entities.UserToken{}

	err := r.db.QueryRow("SELECT * FROM usertoken WHERE userid = $1", userId).Scan(&token.Id, &token.Jti, &token.UserId, &token.Expiry)

	if err == sql.ErrNoRows {
		_, err := r.db.Exec("INSERT INTO usertoken (jti, expiry, userid) VALUES ($1, $2, $3)", jti, expiry, userId)
		if err != nil {
			return err
		}
	} else {
		r.db.Exec("UPDATE usertoken SET jti = $1, expiry = $2 WHERE userid = $3", jti, expiry, userId)
	}

	return nil
}

func (r *AuthRepo) DeleteToken(userId int) error {
	_, err := r.db.Exec("DELETE FROM usertoken WHERE userid = $1", userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) GetUserByTokenId(jti string) (int, entities.Role, error) {
	user := &entities.User{}
	token := &entities.UserToken{}

	err := r.db.QueryRow("SELECT * FROM users JOIN usertoken ON users.id = usertoken.userid WHERE usertoken.jti = $1", jti).
		Scan(&user.Id, &user.Login, &user.Password, &user.UserRole, &user.Email, &token.Id, &token.Jti, &token.UserId, &token.Expiry)

	if err != nil {
		return 0, "", err
	}

	return user.Id, user.UserRole, nil
}

func (r *AuthRepo) GetUserById(userId int) (entities.User, error) {
	user := &entities.User{}

	err := r.db.QueryRow("SELECT * FROM users Where id = $1", userId).Scan(&user.Id, &user.Login, &user.Password, &user.UserRole, &user.Email)
	if err != nil {
		return *user, err
	}

	return *user, nil
}

func (r *AuthRepo) GetUser(login, password string) (entities.User, error) {
	user := &entities.User{}

	err := r.db.QueryRow("SELECT * FROM users WHERE login = $1 AND password = $2", login, password).Scan(&user.Id, &user.Login, &user.Password, &user.UserRole, &user.Email)
	if err == sql.ErrNoRows {
		return *user, &customErrors.NotFoundError{}
	}

	if err != nil {
		return *user, err
	}

	return *user, nil
}

func (r *AuthRepo) ChangePassword(userId int, newPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password = $1 WHERE id = $2", newPassword, userId)
	if err != nil {
		return err
	}

	return nil
}
