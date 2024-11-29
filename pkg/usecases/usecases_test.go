package usecases_test

// import (
// 	"errors"
// 	"inventory/pkg/apperrors"
// 	"inventory/pkg/models"
// 	"inventory/pkg/usecases"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestChargeAccount(t *testing.T) {
// 	mockRepo := new(mocks.SqlxRepositoryMock)
// 	app := usecases.Application{DB: mockRepo}

// 	account := models.Account{Id: 1, Name: "Test Account", Charge: 100}

// 	mockRepo.On("GetAccount", uint(1)).Return(account, nil).Once()
// 	mockRepo.On("ChargeAccount", uint(1), uint(50)).Return(models.Account{Id: 1, Name: "Test Account", Charge: 150}, nil).Once()

// 	updatedAccount, err := app.ChargeAccount(1, 50)
// 	assert.NoError(t, err)
// 	assert.Equal(t, uint(150), updatedAccount.Charge)

// 	mockRepo.AssertExpectations(t)
// }

// func TestChargeAccount_InvalidAccount(t *testing.T) {
// 	mockRepo := new(mocks.SqlxRepositoryMock)
// 	app := usecases.Application{DB: mockRepo}

// 	mockRepo.On("GetAccount", uint(2)).Return(models.Account{}, errors.New("account not found")).Once()

// 	_, err := app.ChargeAccount(2, 50)
// 	assert.Error(t, err)
// 	assert.EqualError(t, err, "account not found")

// 	mockRepo.AssertExpectations(t)
// }

// func TestLogin(t *testing.T) {
// 	app := usecases.Application{}

// 	assert.True(t, app.Login("admin", "password"))
// 	assert.False(t, app.Login("user", "wrongpassword"))
// }

// func TestCreateFactor_Success(t *testing.T) {
// 	mockRepo := new(mocks.SqlxRepositoryMock)
// 	app := usecases.Application{DB: mockRepo}

// 	account := models.Account{Id: 1, Name: "Test Account", Charge: 200}
// 	products := []models.Product{
// 		{Id: 1, Name: "Product 1", Price: 50},
// 		{Id: 2, Name: "Product 2", Price: 100},
// 	}

// 	mockRepo.On("GetAccount", uint(1)).Return(account, nil).Once()
// 	mockRepo.On("GetProducts", []uint{1, 2}).Return(products, nil).Once()
// 	mockRepo.On("CreateFactor", mock.Anything).Return(models.Factor{Id: 1}, nil).Once()

// 	factor, err := app.CreateFactor(1, []uint{1, 2})
// 	assert.NoError(t, err)
// 	assert.Equal(t, uint(1), factor.Id)

// 	mockRepo.AssertExpectations(t)
// }

// func TestCreateFactor_InsufficientFunds(t *testing.T) {
// 	mockRepo := new(mocks.SqlxRepositoryMock)
// 	app := usecases.Application{DB: mockRepo}

// 	account := models.Account{Id: 1, Name: "Test Account", Charge: 100}
// 	products := []models.Product{
// 		{Id: 1, Name: "Product 1", Price: 50},
// 		{Id: 2, Name: "Product 2", Price: 100},
// 	}

// 	mockRepo.On("GetAccount", uint(1)).Return(account, nil).Once()
// 	mockRepo.On("GetProducts", []uint{1, 2}).Return(products, nil).Once()

// 	_, err := app.CreateFactor(1, []uint{1, 2})
// 	assert.Error(t, err)
// 	assert.IsType(t, apperrors.InvalidCredit{}, err)

// 	mockRepo.AssertExpectations(t)
// }

// func TestCreateFactor_ProductNotFound(t *testing.T) {
// 	mockRepo := new(mocks.SqlxRepositoryMock)
// 	app := usecases.Application{DB: mockRepo}

// 	account := models.Account{Id: 1, Name: "Test Account", Charge: 200}
// 	products := []models.Product{
// 		{Id: 1, Name: "Product 1", Price: 50},
// 	}

// 	mockRepo.On("GetAccount", uint(1)).Return(account, nil).Once()
// 	mockRepo.On("GetProducts", []uint{1, 2}).Return(products, nil).Once()

// 	_, err := app.CreateFactor(1, []uint{1, 2})
// 	assert.Error(t, err)
// 	assert.IsType(t, apperrors.NotFound{}, err)

// 	mockRepo.AssertExpectations(t)
// }
