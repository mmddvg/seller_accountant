package models

type Repository interface {
	CreateCustomer(name string) (Customer, error)
	GetAllCustomers() []Customer
	GetCustomerByID(id uint) (Customer, error)
	CreatePurchase(factors []Factor) (Purchase, error)
	GetPurchaseByID(id uint) (Purchase, error)
	GetAllPurchases() ([]Purchase, error)
	GetAllFactors() ([]Factor, error)
	CreateSale(customerID int, price uint64) (Sale, error)
	GetSales() ([]Sale, error)
	Charge(customerId int, charge uint) (Customer, error)
	GetNetProfit() (int, error)
	GetCustomerByName(string) (Customer, error)
}
