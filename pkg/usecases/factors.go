package usecases

import (
	"inventory/pkg/apperrors"
	"inventory/pkg/models"

	"github.com/samber/lo"
)

func (app *Application) CreateFactor(accId uint, factorProducts []models.FactorProduct) (models.Factor, error) {
	var factor models.Factor

	account, err := app.DB.GetAccount(accId)
	if err != nil {

		return factor, err
	}

	var ids []uint
	for i := range factorProducts {
		ids = append(ids, factorProducts[i].ProductId)
	}

	products, err := app.DB.GetProducts(ids)
	if err != nil {
		return factor, err
	}

	if len(products) != len(ids) {
		existingProductIDs := make(map[uint]bool)
		for _, product := range products {
			existingProductIDs[product.Id] = true
		}

		for _, pid := range ids {
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
		p, _ := lo.Find(factorProducts, func(fp models.FactorProduct) bool {
			return product.Id == fp.ProductId
		})

		sumPrice += (product.Price * p.Count)
	}

	if account.Charge < sumPrice {
		return factor, apperrors.InvalidCredit{
			Have: account.Charge,
			Need: sumPrice,
		}
	}

	factor, err = app.DB.CreateFactor(models.NewFactor{Products: factorProducts, AccountId: accId})
	if err != nil {
		return factor, err
	}

	return factor, nil
}

func (app *Application) ListFactors() []models.Factor {
	return app.DB.ListFactors()
}
