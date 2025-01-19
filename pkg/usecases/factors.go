package usecases

import (
	"inventory/pkg/models"
)

func (app *Application) CreatePurchase(factors []models.Factor) (models.Purchase, error) {
	return app.db.CreatePurchase(factors)
}

func (app *Application) GetAllPurchases() ([]models.Purchase, error) {
	return app.db.GetAllPurchases()
}

func (app *Application) GetAllFactors() ([]models.Factor, error) {
	return app.db.GetAllFactors()
}
