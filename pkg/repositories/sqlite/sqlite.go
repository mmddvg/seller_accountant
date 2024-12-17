package sqlite

import (
	"database/sql"
	"errors"
	"inventory/pkg/apperrors"
	"inventory/pkg/models"
	"log"
	"strings"

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

func (repo *SqlxRepository) CreateAccount(name string) (models.Account, error) {
	var account models.Account

	query := `INSERT INTO accounts (name, charge) VALUES (?, ?) RETURNING id, name, charge;`
	err := repo.DB.QueryRow(query, name, 0).Scan(&account.Id, &account.Name, &account.Charge)
	if err != nil {
		if isDuplicateErr(err) {
			return account, apperrors.Duplicate{
				Entity: "Account",
				Id:     account.Id,
			}
		}
		return account, err
	}

	return account, nil
}

func (repo *SqlxRepository) ListAccounts() []models.Account {
	var accounts []models.Account
	query := `SELECT * FROM accounts;`
	err := repo.DB.Select(&accounts, query)
	if err != nil {
		return nil
	}
	return accounts
}

func (repo *SqlxRepository) GetAccount(id uint) (models.Account, error) {
	var account models.Account
	query := `SELECT id, name, charge FROM accounts WHERE id = ?`
	err := repo.DB.Get(&account, query, id)
	if err != nil {
		if isNotFoundErr(err) {
			return account, apperrors.NotFound{
				Entity: "Account",
				Id:     id,
			}
		}
		return account, err
	}
	return account, nil
}

func (repo *SqlxRepository) ChargeAccount(userId uint, amount uint) (models.Account, error) {
	account, err := repo.GetAccount(userId)
	if err != nil {
		return account, err
	}

	newCharge := account.Charge + amount
	query := `UPDATE accounts SET charge = ? WHERE id = ?`
	_, err = repo.DB.Exec(query, newCharge, userId)
	if err != nil {
		return account, err
	}

	account.Charge = newCharge
	return account, nil
}

func (repo *SqlxRepository) CreateProduct(newProduct models.NewProduct) (models.Product, error) {
	var product models.Product

	query := `INSERT INTO products (name, price) VALUES (?, ?) RETURNING id, name, price`
	err := repo.DB.QueryRow(query, newProduct.Name, newProduct.Price).Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		if isDuplicateErr(err) {
			return product, apperrors.Duplicate{
				Entity: "Product",
				Id:     product.Id,
			}
		}
		return product, err
	}

	return product, nil
}
func (repo *SqlxRepository) ListProducts() []models.Product {
	var products []models.Product
	query := `SELECT id, name, price FROM products`
	err := repo.DB.Select(&products, query)
	if err != nil {
		return nil
	}
	return products
}

func (repo *SqlxRepository) GetProducts(ids []uint) ([]models.Product, error) {
	var products []models.Product

	query := `SELECT id, name, price FROM products WHERE id IN (?)`
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	query = repo.DB.Rebind(query)

	err = repo.DB.Select(&products, query, args...)
	if err != nil {
		return nil, err
	}

	if len(products) != len(ids) {
		existingIDs := make(map[uint]bool)
		for _, p := range products {
			existingIDs[p.Id] = true
		}

		for _, id := range ids {
			if !existingIDs[id] {
				return nil, apperrors.NotFound{
					Entity: "Product",
					Id:     id,
				}
			}
		}
	}

	return products, nil
}

func (repo *SqlxRepository) CreateFactor(newFactor models.NewFactor) (models.Factor, error) {
	var factor models.Factor

	var ids []uint
	for _, v := range newFactor.Products {
		ids = append(ids, v.ProductId)
	}

	products, err := repo.GetProducts(ids)
	if err != nil {
		return factor, err
	}

	totalPrice := uint(0)
	for _, product := range products {
		tmp, _ := lo.Find(newFactor.Products, func(arg models.FactorProduct) bool {
			return arg.ProductId == product.Id
		})
		totalPrice += (product.Price * tmp.Count)
	}

	account, err := repo.GetAccount(newFactor.AccountId)
	if err != nil {
		return factor, err
	}

	if account.Charge < totalPrice {
		return factor, apperrors.InvalidCredit{
			Have: account.Charge,
			Need: totalPrice,
		}
	}

	tx, err := repo.DB.Beginx()
	if err != nil {
		return factor, err
	}

	query := `INSERT INTO factors (account_id) VALUES (?) RETURNING id, account_id`
	err = tx.QueryRow(query, newFactor.AccountId).Scan(&factor.Id, &factor.AccountId)
	if err != nil {
		tx.Rollback()
		return factor, err
	}

	query = `INSERT INTO factor_products (factor_id, product_id,count) VALUES (?, ?,?)`
	for _, product := range products {
		fp, _ := lo.Find(newFactor.Products, func(arg models.FactorProduct) bool {
			return arg.ProductId == product.Id
		})
		_, err = tx.Exec(query, factor.Id, product.Id, fp.Count)
		if err != nil {
			tx.Rollback()
			return factor, err
		}
	}

	query = `UPDATE accounts SET charge = charge - ? WHERE id = ?`
	_, err = tx.Exec(query, totalPrice, newFactor.AccountId)
	if err != nil {
		tx.Rollback()
		return factor, err
	}

	tx.Commit()
	factor.Products = newFactor.Products
	return factor, nil
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

func (r *SqlxRepository) UpdateProduct(productId uint, newPrice uint) (models.Product, error) {
	var product models.Product

	querySelect := "SELECT id, name, price FROM products WHERE id = ?"
	err := r.DB.Get(&product, querySelect, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, apperrors.NotFound{
				Entity: "Product",
				Id:     productId,
			}
		}
		return product, err
	}

	queryUpdate := "UPDATE products SET price = ? WHERE id = ?"
	_, err = r.DB.Exec(queryUpdate, newPrice, productId)
	if err != nil {
		return product, err
	}

	product.Price = newPrice
	return product, nil
}

func (repo *SqlxRepository) ListFactors() []models.Factor {
	var factors []models.Factor
	query := `SELECT * FROM factors;`

	err := repo.DB.Select(&factors, query)
	if err != nil {
		log.Println("err : ", err)
		return nil
	}

	for i := range factors {
		var products []models.FactorProduct
		productQuery := `SELECT product_id ,count FROM factor_products WHERE factor_id = ?;`
		err := repo.DB.Select(&products, productQuery, factors[i].Id)
		if err != nil {
			log.Println("err: ", err)
			return nil
		}
		factors[i].Products = products
	}

	return factors
}
