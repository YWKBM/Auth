package entities

import "time"

type UserToken struct {
	Id     int
	Jti    string
	UserId int
	Expiry time.Time
}
