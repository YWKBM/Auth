package repo

type Authorization interface {
	CreateUser(login, password, email string) error
}
