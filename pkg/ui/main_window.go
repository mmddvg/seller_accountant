package ui

import (
	"inventory/pkg/usecases"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/samber/lo"
)

func MainWindow(app *usecases.Application, fyneApp fyne.App) {
	mainWindow := fyneApp.NewWindow("Accounting Application")
	mainWindow.Resize(fyne.NewSize(800, 600))

	writer := make(chan bool)

	readers := lo.FanOut(3, 0, writer)

	tabs := container.NewAppTabs(
		container.NewTabItem("داشبورد", dashBoardTab(app, mainWindow, readers[0])),
		container.NewTabItem("مدیریت مشتری و فروش ", accountsTab(app, mainWindow, writer, readers[1])),
		container.NewTabItem("مدیریت خرید", purchasesTab(app, mainWindow, writer, readers[2])),
	)

	mainWindow.SetContent(tabs)
	mainWindow.Show()
}
