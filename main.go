package main

import (
	"inventory/pkg/repositories/sqlite"
	"inventory/pkg/ui"
	"inventory/pkg/usecases"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	db, err := sqlite.InitializeDatabase("db.sql")
	if err != nil {
		log.Fatal(err)
	}

	appInstance := usecases.NewApp(sqlite.NewSqlxRepository(db))

	app := app.New()
	window := app.NewWindow("Login")
	window.Resize(fyne.NewSize(600, 400))

	ui.InitiateUI(&appInstance, app, window)
	window.ShowAndRun()
}
