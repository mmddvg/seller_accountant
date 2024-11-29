package apperrors

import "fmt"

type InvalidCredit struct {
	Have uint
	Need uint
}

func (err InvalidCredit) Error() string {
	return fmt.Sprintf("not enough credit , have : %d , need : %d ", err.Have, err.Need)
}
