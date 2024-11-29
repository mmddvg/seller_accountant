package sqlite_test

import (
	"inventory/pkg/apperrors"
	"inventory/pkg/models"
	"inventory/pkg/repositories/sqlite"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestAccountCRUD(t *testing.T) {
	db, err := sqlite.InitializeDatabase("./test.db")
	require.NoError(t, err)
	defer db.Close()

	repo := sqlite.NewSqlxRepository(db)

	accountName := randSeq(10)
	acc1, err := repo.CreateAccount(accountName)
	require.NoError(t, err)
	assert.NotZero(t, acc1.Id)
	assert.Equal(t, accountName, acc1.Name)
	assert.Equal(t, uint(0), acc1.Charge)

	retrievedAcc, err := repo.GetAccount(acc1.Id)
	require.NoError(t, err)
	assert.Equal(t, acc1, retrievedAcc)

	_, err = repo.CreateAccount(accountName)
	assert.Error(t, err)
	assert.IsType(t, apperrors.Duplicate{}, err)

	_, err = repo.GetAccount(99999)
	assert.Error(t, err)
	assert.IsType(t, apperrors.NotFound{}, err)
}

func TestProductCRUD(t *testing.T) {
	db, err := sqlite.InitializeDatabase("./test.db")
	require.NoError(t, err)
	defer db.Close()

	repo := sqlite.NewSqlxRepository(db)

	productName := randSeq(10)
	productPrice := uint(100)
	newProduct := models.NewProduct{Name: productName, Price: productPrice}
	prod1, err := repo.CreateProduct(newProduct)
	require.NoError(t, err)
	assert.NotZero(t, prod1.Id)
	assert.Equal(t, productName, prod1.Name)
	assert.Equal(t, productPrice, prod1.Price)

	products := repo.ListProducts()
	assert.NotEmpty(t, products)
	assert.Contains(t, products, prod1)

	_, err = repo.CreateProduct(newProduct)
	assert.Error(t, err)
	assert.IsType(t, apperrors.Duplicate{}, err)

	_, err = repo.GetProducts([]uint{99999})
	assert.Error(t, err)
	assert.IsType(t, apperrors.NotFound{}, err)
}

func TestCreateFactor(t *testing.T) {
	db, err := sqlite.InitializeDatabase("./test.db")
	require.NoError(t, err)
	defer db.Close()

	repo := sqlite.NewSqlxRepository(db)

	account, err := repo.CreateAccount(randSeq(10))
	require.NoError(t, err)
	product1, err := repo.CreateProduct(models.NewProduct{Name: randSeq(10), Price: 50})
	require.NoError(t, err)
	product2, err := repo.CreateProduct(models.NewProduct{Name: randSeq(10), Price: 30})
	require.NoError(t, err)

	account, err = repo.ChargeAccount(account.Id, 100)
	require.NoError(t, err)
	assert.Equal(t, uint(100), account.Charge)

	newFactor := models.NewFactor{
		Products:  []uint{product1.Id, product2.Id},
		AccountId: account.Id,
	}
	factor, err := repo.CreateFactor(newFactor)
	require.NoError(t, err)
	assert.NotZero(t, factor.Id)
	assert.Equal(t, account.Id, factor.AccountId)
	assert.ElementsMatch(t, []uint{product1.Id, product2.Id}, factor.Products)

	account, err = repo.GetAccount(account.Id)
	require.NoError(t, err)
	assert.Equal(t, uint(20), account.Charge)

	_, err = repo.CreateFactor(models.NewFactor{
		Products:  []uint{product1.Id, product2.Id},
		AccountId: account.Id,
	})
	assert.Error(t, err)
	assert.IsType(t, apperrors.InvalidCredit{}, err)
}

func TestListAccounts(t *testing.T) {
	db, err := sqlite.InitializeDatabase("./test.db")
	require.NoError(t, err)
	defer db.Close()

	repo := sqlite.NewSqlxRepository(db)

	account1, err := repo.CreateAccount(randSeq(10))
	require.NoError(t, err)
	account2, err := repo.CreateAccount(randSeq(10))
	require.NoError(t, err)

	res, err := db.DB.Query("SELECT COUNT(*) FROM accounts;")
	require.NoError(t, err)
	var count int
	res.Next()
	res.Scan(&count)
	res.Close()

	accounts := repo.ListAccounts()
	assert.Len(t, accounts, count)
	assert.Contains(t, accounts, account1)
	assert.Contains(t, accounts, account2)
}
