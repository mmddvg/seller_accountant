package sqlite

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createAccountsTable = `
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		charge INTEGER NOT NULL
	);`

	createFactorsTable = `
	CREATE TABLE IF NOT EXISTS factors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_id INTEGER NOT NULL,
		FOREIGN KEY (account_id) REFERENCES accounts(id)
	);`

	createProductsTable = `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		price INTEGER NOT NULL
	);`

	createFactorProductsTable = `
	CREATE TABLE IF NOT EXISTS factor_products (
		factor_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		count INTEGER NOT NULL DEFAULT 1,
		FOREIGN KEY (factor_id) REFERENCES factors(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	);`
)

// Initialize the database
func InitializeDatabase(filePath string) (*sqlx.DB, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		file.Close()
		log.Printf("Created database file at %s\n", filePath)
	} else {
		log.Printf("Database file already exists at %s\n", filePath)
		db, err := sqlx.Open("sqlite3", filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return db, nil

	}

	db, err := sqlx.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec(createAccountsTable); err != nil {
		return nil, fmt.Errorf("failed to create accounts table: %w", err)
	}
	if _, err := db.Exec(createFactorsTable); err != nil {
		return nil, fmt.Errorf("failed to create factors table: %w", err)
	}
	if _, err := db.Exec(createProductsTable); err != nil {
		return nil, fmt.Errorf("failed to create products table: %w", err)
	}
	if _, err := db.Exec(createFactorProductsTable); err != nil {
		return nil, fmt.Errorf("failed to create factor_products table: %w", err)
	}

	log.Println("Database initialized and tables created")
	return db, nil
}
