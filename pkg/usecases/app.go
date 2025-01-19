package usecases

import "inventory/pkg/models"

type Application struct {
	db models.Repository
}

func NewApp(db models.Repository) Application {
	return Application{db: db}
}

func (app *Application) Login(userName string, password string) bool {
	if userName == "admin" && password == "password" {
		return true
	}

	return false
}
