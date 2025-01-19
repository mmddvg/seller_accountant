package sqlite

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createCustomersTable = `
	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		charge INTEGER NOT NULL
	);`

	createPurchasesTable = `
	CREATE TABLE IF NOT EXISTS purchases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP NOT NULL
	);`

	createFactorsTable = `
	CREATE TABLE IF NOT EXISTS factors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		purchase_id INTEGER NOT NULL,
		store_name TEXT NOT NULL,
		price INTEGER NOT NULL,
		file_name TEXT,
		FOREIGN KEY(purchase_id) REFERENCES purchases(id)
	);`

	createSalesTable = `
	CREATE TABLE IF NOT EXISTS sales (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER NOT NULL,
		price INTEGER NOT NULL, -- Sale price with 10% markup
		FOREIGN KEY(customer_id) REFERENCES customers(id)
	);`
)

func InitializeDatabase(filePath string) (*sqlx.DB, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
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

	if _, err := db.Exec(createCustomersTable); err != nil {
		return nil, fmt.Errorf("failed to create customers table: %w", err)
	}

	if _, err := db.Exec(createPurchasesTable); err != nil {
		return nil, fmt.Errorf("failed to create purchases table: %w", err)
	}

	if _, err := db.Exec(createFactorsTable); err != nil {
		return nil, fmt.Errorf("failed to create factors table: %w", err)
	}

	if _, err := db.Exec(createSalesTable); err != nil {
		return nil, fmt.Errorf("failed to create sales table: %w", err)
	}

	log.Println("Database initialized and tables created")
	return db, nil
}
