package customErrors

import "fmt"

type AlreadyExistsError struct {
}

type NotFoundError struct {
}

type ValidationError struct {
	Message string
}

func (u *AlreadyExistsError) Error() string {
	return fmt.Sprintln("alreday exists")
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintln("user not found")
}

func (v *ValidationError) Error() string {
	return fmt.Sprintln(v.Message)
}
