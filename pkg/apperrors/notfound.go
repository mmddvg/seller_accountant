package apperrors

import "fmt"

type NotFound struct {
	Entity string
	Id     uint
}

func (err NotFound) Error() string {
	return fmt.Sprintf("%s with id : %d not found", err.Entity, err.Id)
}
