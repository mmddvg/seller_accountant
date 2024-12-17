package models

import "fmt"

type FactorProduct struct {
	ProductId uint `db:"product_id"`
	Count     uint `db:"count"`
}

func (fp FactorProduct) ToStr() string {
	return fmt.Sprintf(" (product : %d , count : %d )", fp.ProductId, fp.Count)
}

type NewFactor struct {
	Products  []FactorProduct
	AccountId uint
}

type Factor struct {
	Id        uint `db:"id"`
	Products  []FactorProduct
	AccountId uint `db:"account_id"`
}
