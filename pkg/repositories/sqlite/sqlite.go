package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"inventory/pkg/apperrors"
	"inventory/pkg/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

type SqlxRepository struct {
	DB *sqlx.DB
}

func NewSqlxRepository(db *sqlx.DB) *SqlxRepository {
	return &SqlxRepository{DB: db}
}

func (repo *SqlxRepository) CreateCustomer(name string) (models.Customer, error) {
	var account models.Customer

	query := `INSERT INTO customers(name, charge) VALUES (?, 0) RETURNING *;`
	err := repo.DB.Get(&account, query, name, 0)
	if err != nil {
		if isDuplicateErr(err) {
			return account, apperrors.Duplicate{
				Entity: "customer",
				Id:     uint(account.ID),
			}
		}
		return account, err
	}

	return account, nil
}

func (repo *SqlxRepository) GetAllCustomers() []models.Customer {
	var accounts []models.Customer
	query := `SELECT * FROM customers;`
	err := repo.DB.Select(&accounts, query)
	if err != nil {
		return nil
	}
	return accounts
}

func (repo *SqlxRepository) GetCustomerByID(id uint) (models.Customer, error) {
	var account models.Customer
	query := `SELECT * FROM customers WHERE id = ?`
	err := repo.DB.Get(&account, query, id)
	if err != nil {
		if isNotFoundErr(err) {
			return account, apperrors.NotFound{
				Entity: "Customer",
				Id:     fmt.Sprint(id),
			}
		}
		return account, err
	}
	return account, nil
}

func (repo *SqlxRepository) CreatePurchase(factors []models.Factor) (models.Purchase, error) {

	purchase := models.Purchase{Factors: make([]models.Factor, len(factors))}
	query := `INSERT INTO purchases(created_at) VALUES (?) RETURNING *;`
	err := repo.DB.Get(&purchase, query, time.Now())
	if err != nil {
		return models.Purchase{}, err
	}

	for i, v := range factors {
		err := repo.DB.Get(&purchase.Factors[i], `INSERT INTO factors(purchase_id,store_name,price) VALUES(?,?,?) RETURNING *;`, purchase.ID, v.StoreName, v.Price)
		if err != nil {

			return models.Purchase{}, err
		}

	}

	return purchase, nil
}

func (repo *SqlxRepository) GetPurchaseByID(id uint) (models.Purchase, error) {
	var purchase models.Purchase
	query := `SELECT * FROM purchases WHERE id = ?;`
	err := repo.DB.Get(&purchase, query, id)
	if err != nil {
		if isNotFoundErr(err) {
			return purchase, apperrors.NotFound{
				Entity: "Customer",
				Id:     fmt.Sprint(id),
			}
		}
		return purchase, err
	}

	err = repo.DB.Select(&purchase.Factors, `SELECT * FROM factors WHERE purchase_id = ?;`, purchase.ID)
	return purchase, err
}

func (repo *SqlxRepository) GetAllPurchases() ([]models.Purchase, error) {
	var purchases []models.Purchase
	query := `SELECT * FROM purchases;`
	err := repo.DB.Select(&purchases, query)
	if err != nil {

		return purchases, err
	}

	for i := range purchases {
		err = repo.DB.Select(&purchases[i].Factors, "SELECT * FROM factors WHERE purchase_id = ?;", purchases[i].ID)
	}
	return purchases, err
}

func (repo *SqlxRepository) GetAllFactors() ([]models.Factor, error) {
	factors := []models.Factor{}
	err := repo.DB.Select(&factors, "SELECT * FROM factors;")
	return factors, err
}

func (repo *SqlxRepository) CreateSale(customerID int, price uint64) (models.Sale, error) {
	sale := models.Sale{}

	err := repo.DB.Get(&sale, "INSERT INTO sales(customer_id,price) VALUES(?,?) RETURNING *;", customerID, price)
	if err != nil {
		return sale, err
	}

	_, err = repo.DB.Exec("UPDATE customers SET charge = (charge - ?) WHERE id = ?;", price, customerID)

	return sale, err
}

func (repo *SqlxRepository) GetSales() ([]models.Sale, error) {
	sales := []models.Sale{}
	err := repo.DB.Select(&sales, "SELECT * FROM sales;")
	return sales, err
}

func (repo *SqlxRepository) Charge(customerId int, charge uint) (models.Customer, error) {
	customer := models.Customer{}
	err := repo.DB.Get(&customer, "UPDATE customers SET charge = (charge + ?) WHERE id = ? RETURNING *;", charge, customerId)
	return customer, err
}

func (repo *SqlxRepository) GetNetProfit() (int, error) {
	sales, err := repo.GetSales()
	if err != nil {
		return 0, err
	}

	var overallSale int = lo.SumBy(sales, func(sale models.Sale) int { return sale.Price })

	factors, err := repo.GetAllFactors()
	if err != nil {
		return 0, nil
	}

	overallSpent := lo.SumBy(factors, func(factor models.Factor) int { return factor.Price })
	return overallSale - overallSpent, nil
}

func (repo *SqlxRepository) GetCustomerByName(name string) (models.Customer, error) {
	var customer models.Customer

	err := repo.DB.Get(&customer, "SELECT * FROM customers WHERE name = ?;", name)
	if err != nil {
		if isNotFoundErr(err) {
			return customer, apperrors.NotFound{Entity: "Customer", Id: name}
		}
		return customer, err
	}

	if customer.Name != name {
		return customer, apperrors.NotFound{Entity: "Customer", Id: name}

	}
	return customer, nil
}

func isNotFoundErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func isDuplicateErr(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
