package usecases

import (
	"inventory/pkg/models"
)

func (app *Application) CreateAccount(name string) (models.Customer, error) {
	return app.db.CreateCustomer(name)
}

func (app *Application) ChargeAccount(accId uint, amount uint) (models.Customer, error) {
	_, err := app.db.GetCustomerByID(accId)
	if err != nil {
		return models.Customer{}, err
	}

	return app.db.Charge(int(accId), amount)
}

func (app *Application) ListAccounts() []models.Customer {
	return app.db.GetAllCustomers()
}

func (app *Application) GetSales() ([]models.Sale, error) {
	return app.db.GetSales()
}

func (app *Application) GetNetProfit() (int, error) {
	return app.db.GetNetProfit()
}

func (app *Application) Sell(name string, price uint) error {
	customer, err := app.db.GetCustomerByName(name)
	if err != nil {
		return err
	}

	finalPrice := float64(price) * 1.1

	_, err = app.db.CreateSale(customer.ID, uint64(finalPrice))
	return err
}
