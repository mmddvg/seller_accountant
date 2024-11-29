package usecases

import "inventory/pkg/models"

type Application struct {
	DB models.Database
}

func (app *Application) Login(userName string, password string) bool {
	if userName == "admin" && password == "password" {
		return true
	}

	return false
}
