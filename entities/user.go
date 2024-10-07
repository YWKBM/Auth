package entities

type User struct {
	Id       int
	Login    string
	Password string
	UserRole Role
	Email    string
	TokenId  string
}

type Role string

const (
	PROVIDER Role = "PROVIDER"
	USER     Role = "USER"
)
