package entities

import "fmt"

type User struct {
	Id       int
	Login    string
	Password string
	UserRole Role
	Email    string
}

type Role string

const (
	MANAGER  Role = "MANAGER"
	PROVIDER Role = "PROVIDER"
	USER     Role = "USER"
)

func (r Role) String() string {
	return string(r)
}

func ParseRole(s string) (r Role, err error) {
	capabilities := map[Role]struct{}{
		MANAGER:  {},
		PROVIDER: {},
		USER:     {},
	}

	cap := Role(s)
	_, ok := capabilities[cap]
	if !ok {
		return r, fmt.Errorf(`cannot parse:[%s] as roles`, s)
	}
	return cap, nil
}
