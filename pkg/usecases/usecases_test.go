package usecases_test

import (
	"inventory/pkg/models"
	"inventory/pkg/repositories/sqlite"
	"inventory/pkg/usecases"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	app usecases.Application
)

func TestMain(m *testing.M) {
	db, err := sqlite.InitializeDatabase("test.db")
	if err != nil {
		log.Fatal(err)
	}

	app = usecases.Application{DB: sqlite.NewSqlxRepository(db)}

	exitCode := m.Run()

	err = os.Remove("test.db")
	if err != nil {
		log.Println("error deleting db :  ", err)
	}

	os.Exit(exitCode)

}

func TestChargeAccount_InvalidAccount(t *testing.T) {
	_, err := app.ChargeAccount(1532, 50)
	assert.Error(t, err)

	assert.EqualError(t, err, "Account with id : 1532 not found")
}

func TestLogin(t *testing.T) {
	assert.True(t, app.Login("admin", "password"))
	assert.False(t, app.Login("user", "wrongpassword"))
}

func TestChargeAccount(t *testing.T) {
	acc, err := app.CreateAccount("test account1")
	require.NoError(t, err)

	acc, err = app.ChargeAccount(acc.Id, 100)
	require.NoError(t, err)

	assert.Equal(t, uint(100), acc.Charge)
}

func TestSale(t *testing.T) {
	// prepare
	acc, _ := app.CreateAccount("test account2")
	acc, _ = app.ChargeAccount(acc.Id, 1000)

	prod1, err := app.CreateProduct(models.NewProduct{Name: "product 1", Price: 30})
	require.NoError(t, err)

	prod2, err := app.CreateProduct(models.NewProduct{Name: "product 2", Price: 40})
	require.NoError(t, err)

	// test

	_, err = app.CreateFactor(acc.Id, []models.FactorProduct{{ProductId: 1315, Count: 12}})
	require.Error(t, err, "Product with id : 1315 not found")

	factor, err := app.CreateFactor(acc.Id, []models.FactorProduct{{ProductId: prod1.Id, Count: 5}, {ProductId: prod2.Id, Count: 7}})
	require.NoError(t, err)

	require.Equal(t, uint(5), factor.Products[0].Count)

	_, err = app.CreateFactor(acc.Id, []models.FactorProduct{{ProductId: prod1.Id, Count: 10000}})
	require.Error(t, err)
}
