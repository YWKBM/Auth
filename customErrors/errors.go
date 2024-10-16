package customErrors

import "fmt"

type AlreadyExistsError struct {
}

type NotFoundError struct {
}

func (u *AlreadyExistsError) Error() string {
	return fmt.Sprintln("alreday exists")
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintln("user not found")
}
