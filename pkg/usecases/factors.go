package usecases

import (
	"inventory/pkg/apperrors"
	"inventory/pkg/models"
)

func (app *Application) CreateFactor(accId uint, productIds []uint) (models.Factor, error) {
	var factor models.Factor

	account, err := app.DB.GetAccount(accId)
	if err != nil {

		return factor, err
	}

	products, err := app.DB.GetProducts(productIds)
	if err != nil {
		return factor, err
	}

	if len(products) != len(productIds) {
		existingProductIDs := make(map[uint]bool)
		for _, product := range products {
			existingProductIDs[product.Id] = true
		}

		for _, pid := range productIds {
			if !existingProductIDs[pid] {
				return factor, apperrors.NotFound{
					Entity: "Product",
					Id:     pid,
				}
			}
		}
	}

	sumPrice := uint(0)
	for _, product := range products {
		sumPrice += product.Price
	}

	if account.Charge < sumPrice {
		return factor, apperrors.InvalidCredit{
			Have: account.Charge,
			Need: sumPrice,
		}
	}

	factor, err = app.DB.CreateFactor(models.NewFactor{Products: productIDsFromProducts(products), AccountId: accId})
	if err != nil {
		return factor, err
	}

	return factor, nil
}

func (app *Application) ListFactors() []models.Factor {
	return app.DB.ListFactors()
}

func productIDsFromProducts(products []models.Product) []uint {
	var ids []uint
	for _, p := range products {
		ids = append(ids, p.Id)
	}
	return ids
}
