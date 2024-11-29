package usecases

import (
	"inventory/pkg/models"
)

func (app *Application) CreateProduct(prd models.NewProduct) (models.Product, error) {
	return app.DB.CreateProduct(prd)
}

func (app *Application) ListProducts() []models.Product {
	return app.DB.ListProducts()
}

func (app *Application) UpdateProduct(prodId uint, price uint) (models.Product, error) {
	return app.DB.UpdateProduct(prodId, price)
}
