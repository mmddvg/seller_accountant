package apperrors

import "fmt"

type NotFound struct {
	Entity string
	Id     string
}

func (err NotFound) Error() string {
	return fmt.Sprintf("%s : %s not found", err.Entity, err.Id)
}
