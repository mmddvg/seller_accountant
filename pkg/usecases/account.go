package usecases

import (
	"inventory/pkg/models"
)

func (app *Application) CreateAccount(name string) (models.Account, error) {
	return app.DB.CreateAccount(name)
}

func (app *Application) ChargeAccount(accId uint, amount uint) (models.Account, error) {
	_, err := app.DB.GetAccount(accId)
	if err != nil {
		return models.Account{}, err
	}

	return app.DB.ChargeAccount(accId, amount)
}

func (app *Application) ListAccounts() []models.Account {
	return app.DB.ListAccounts()
}
