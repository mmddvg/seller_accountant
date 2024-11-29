package apperrors

import "fmt"

type Duplicate struct {
	Entity string
	Id     uint
}

func (err Duplicate) Error() string {
	return fmt.Sprintf("entry is conficlted with %s id %d ", err.Entity, err.Id)
}
