package sqlite_test

import (
	"inventory/pkg/apperrors"
	"inventory/pkg/models"
	"inventory/pkg/repositories/sqlite"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	id, err := uuid.NewV7()
	assert.NoError(t, err)

	db, err := sqlite.InitializeDatabase(id.String() + ".tmp.db")
	assert.NoError(t, err)
	return db
}

func TestCreateCustomer(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	name := "John Doe"
	customer, err := repo.CreateCustomer(name)

	assert.NoError(t, err)
	assert.Equal(t, name, customer.Name)
	assert.Equal(t, 0, customer.Charge)
}

func TestGetAllCustomers(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec("INSERT INTO customers (name, charge) VALUES (?, ?), (?, ?);", "Alice", 100, "Bob", 200)

	customers := repo.GetAllCustomers()
	assert.Len(t, customers, 2)
	assert.Equal(t, "Alice", customers[0].Name)
	assert.Equal(t, 100, customers[0].Charge)
	assert.Equal(t, "Bob", customers[1].Name)
	assert.Equal(t, 200, customers[1].Charge)
}

func TestGetCustomerByID(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec("INSERT INTO customers (id, name, charge) VALUES (?, ?, ?);", 1, "Charlie", 300)

	customer, err := repo.GetCustomerByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Charlie", customer.Name)
	assert.Equal(t, 300, customer.Charge)
}

func TestCreatePurchase(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	factors := []models.Factor{
		{StoreName: "Store1", Price: 50},
		{StoreName: "Store2", Price: 100},
	}

	purchase, err := repo.CreatePurchase(factors)

	assert.NoError(t, err)
	assert.Len(t, purchase.Factors, 2)
	assert.Equal(t, "Store1", purchase.Factors[0].StoreName)
	assert.Equal(t, 50, purchase.Factors[0].Price)
	assert.Equal(t, "Store2", purchase.Factors[1].StoreName)
	assert.Equal(t, 100, purchase.Factors[1].Price)
}

func TestCreateSale(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec("INSERT INTO customers (id, name, charge) VALUES (?, ?, ?);", 1, "Dave", 500)

	sale, err := repo.CreateSale(1, 300)

	assert.NoError(t, err)
	assert.Equal(t, sale.CustomerId, uint(1))

	var remainingCharge int
	db.Get(&remainingCharge, "SELECT charge FROM customers WHERE id = ?;", 1)
	assert.Equal(t, 200, remainingCharge)
}

func TestGetNetProfit(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec(`
		INSERT INTO sales (id,  customer_id, price) VALUES (1, 1, 500), (2, 2, 700);
		INSERT INTO factors (id, purchase_id, store_name, price) VALUES (1, 1, 'Store1', 300), (2, 1, 'Store2', 200);
	`)

	netProfit, err := repo.GetNetProfit()

	assert.NoError(t, err)
	assert.Equal(t, 700, netProfit)
}

func TestGetPurchaseByID(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec(`
		INSERT INTO purchases (id, created_at) VALUES (1, ?);
		INSERT INTO factors (purchase_id, store_name, price) VALUES (1, 'Store1', 100), (1, 'Store2', 200);
	`, time.Now())

	purchase, err := repo.GetPurchaseByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(purchase.Factors))
	assert.Equal(t, "Store1", purchase.Factors[0].StoreName)
	assert.Equal(t, 100, purchase.Factors[0].Price)
}

func TestGetAllPurchases(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec(`
		INSERT INTO purchases (id, created_at) VALUES (1, ?), (2, ?);
		INSERT INTO factors (purchase_id, store_name, price) VALUES (1, 'Store1', 100), (2, 'Store2', 200);
	`, time.Now(), time.Now())

	purchases, err := repo.GetAllPurchases()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(purchases))
	assert.Equal(t, "Store1", purchases[0].Factors[0].StoreName)
	assert.Equal(t, 100, purchases[0].Factors[0].Price)
}

func TestGetAllFactors(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec(`INSERT INTO factors (id, purchase_id, store_name, price) VALUES (1, 1, 'Store1', 100), (2, 1, 'Store2', 200);`)

	factors, err := repo.GetAllFactors()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(factors))
	assert.Equal(t, "Store1", factors[0].StoreName)
	assert.Equal(t, 100, factors[0].Price)
}

func TestCharge(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec("INSERT INTO customers (id, name, charge) VALUES (?, ?, ?);", 1, "Eve", 200)

	customer, err := repo.Charge(1, 300)

	assert.NoError(t, err)
	assert.Equal(t, 500, customer.Charge)
}

func TestGetSales(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec(`INSERT INTO sales (id,customer_id, price) VALUES (1, 1, 300), (2,2, 400);`)

	sales, err := repo.GetSales()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(sales))
	assert.Equal(t, 300, sales[0].Price)
}

func TestCreateCustomer_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	name := "John Doe"
	_, err := repo.CreateCustomer(name)
	assert.NoError(t, err)

	_, err = repo.CreateCustomer(name)
	assert.Error(t, err)
	assert.IsType(t, apperrors.Duplicate{}, err)
}

func TestGetAllCustomers_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	customers := repo.GetAllCustomers()
	assert.Len(t, customers, 0)
}
func TestCreateSale_Concurrent(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewSqlxRepository(db)

	db.MustExec("INSERT INTO customers (id, name, charge) VALUES (?, ?, ?);", 1, "Concurrent Customer", 1000)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := repo.CreateSale(1, 100)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	var remainingCharge int
	db.Get(&remainingCharge, "SELECT charge FROM customers WHERE id = ?;", 1)
	assert.Equal(t, 0, remainingCharge) // Assuming total charge doesn't exceed initial amount
}
