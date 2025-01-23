package usecases_test

import (
	"inventory/pkg/repositories/sqlite"
	"inventory/pkg/usecases"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
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

	app = usecases.NewApp(sqlite.NewSqlxRepository(db))

	exitCode := m.Run()

	err = os.Remove("test.db")
	if err != nil {
		log.Println("error deleting db: ", err)
	}

	os.Exit(exitCode)
}

func TestValidateImageFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		filePath    string
		setup       func() string
		expectedErr string
	}{
		{
			name:        "Non-existent file",
			filePath:    "non_existent_file.jpg",
			setup:       func() string { return "non_existent_file.jpg" },
			expectedErr: "file does not exist",
		},
		{
			name:        "File is a directory",
			filePath:    tempDir,
			setup:       func() string { return tempDir },
			expectedErr: "path is a directory",
		},
		{
			name:     "Invalid file type",
			filePath: filepath.Join(tempDir, "invalid_file.txt"),
			setup: func() string {
				file := filepath.Join(tempDir, "invalid_file.txt")
				os.WriteFile(file, []byte("dummy content"), 0644)
				return file
			},
			expectedErr: "invalid file type",
		},
		{
			name:     "Valid JPEG file",
			filePath: filepath.Join(tempDir, "image.jpg"),
			setup: func() string {
				file := filepath.Join(tempDir, "image.jpg")
				os.WriteFile(file, []byte("dummy content"), 0644)
				return file
			},
			expectedErr: "",
		},
		{
			name:     "Valid PNG file",
			filePath: filepath.Join(tempDir, "image.png"),
			setup: func() string {
				file := filepath.Join(tempDir, "image.png")
				os.WriteFile(file, []byte("dummy content"), 0644)
				return file
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup()
			err := usecases.ValidateImageFile(filePath)

			if tt.expectedErr == "" {
				assert.NoError(t, err, "expected no error, got: %v", err)
			} else {
				assert.Error(t, err, "expected an error but got none")
				assert.Contains(t, err.Error(), tt.expectedErr, "error message mismatch")
			}
		})
	}
}

func TestSellWithCommission(t *testing.T) {
	id, err := uuid.NewV7()
	require.NoError(t, err, "failed to generate UUID")

	db, err := sqlite.InitializeDatabase(id.String())
	require.NoError(t, err, "failed to initialize database")

	app := usecases.NewApp(sqlite.NewSqlxRepository(db))

	customer1, err := app.CreateAccount("John Doe")
	require.NoError(t, err, "failed to create customer John Doe")

	customer2, err := app.CreateAccount("Jane Smith")
	require.NoError(t, err, "failed to create customer Jane Smith")

	_, err = app.ChargeAccount(uint(customer1.ID), 1000)
	require.NoError(t, err, "failed to charge customer1 account")

	_, err = app.ChargeAccount(uint(customer2.ID), 2000)
	require.NoError(t, err, "failed to charge customer2 account")

	sales := []struct {
		name     string
		price    uint
		expected float64
	}{
		{"John Doe", 100, 110},
		{"Jane Smith", 200, 220},
	}

	for _, sale := range sales {
		err := app.Sell(sale.name, sale.price)
		require.NoError(t, err, "failed to create sale for %s", sale.name)
	}

	salesData, err := app.GetSales()
	require.NoError(t, err, "failed to fetch sales")

	var totalSales uint64
	for _, sale := range salesData {
		totalSales += uint64(sale.Price)
	}

	var expectedTotal float64
	for _, sale := range sales {
		expectedTotal += sale.expected
	}

	assert.Equal(t, expectedTotal, float64(totalSales), "total sales mismatch")

	netProfit, err := app.GetNetProfit()
	require.NoError(t, err, "failed to fetch net profit")

	expectedProfit := uint64(expectedTotal)
	assert.Equal(t, expectedProfit, uint64(netProfit), "net profit mismatch")
}
