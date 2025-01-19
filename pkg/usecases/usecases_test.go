package usecases_test

import (
	"inventory/pkg/repositories/sqlite"
	"inventory/pkg/usecases"
	"log"
	"os"
	"testing"
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
		log.Println("error deleting db :  ", err)
	}

	os.Exit(exitCode)

}
