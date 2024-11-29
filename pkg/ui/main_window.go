package ui

import (
	"inventory/pkg/usecases"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func MainWindow(app *usecases.Application, fyneApp fyne.App) {
	mainWindow := fyneApp.NewWindow("Accounting Application")
	mainWindow.Resize(fyne.NewSize(800, 600))

	refreshAccs := make(chan bool)

	tabs := container.NewAppTabs(
		container.NewTabItem("Accounts", accountsTab(app, mainWindow, refreshAccs)),
		container.NewTabItem("Products", productsTab(app, mainWindow)),
		container.NewTabItem("Factors", factorsTab(app, mainWindow, refreshAccs)),
	)

	mainWindow.SetContent(tabs)
	mainWindow.Show()
}
