package entities

import "time"

type userToken struct {
	Jti    string
	UserId int
	Expiry time.Time
}
