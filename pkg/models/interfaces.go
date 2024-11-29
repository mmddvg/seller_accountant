package models

type Database interface {
	CreateAccount(name string) (Account, error)
	ListAccounts() []Account
	GetAccount(uint) (Account, error)
	ChargeAccount(userId uint, amount uint) (Account, error)

	CreateProduct(NewProduct) (Product, error)
	UpdateProduct(prodId uint, price uint) (Product, error)
	CreateFactor(NewFactor) (Factor, error)
	ListFactors() []Factor
	ListProducts() []Product
	GetProducts([]uint) ([]Product, error)
}
