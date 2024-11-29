package ui

import (
	"fmt"
	"inventory/pkg/usecases"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func InitiateUI(app *usecases.Application, fyneApp fyne.App, window fyne.Window) {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	loginBtn := widget.NewButton("Login", func() {
		if app.Login(usernameEntry.Text, passwordEntry.Text) {
			window.Hide()
			MainWindow(app, fyneApp)
		} else {
			dialog.ShowError(fmt.Errorf("invalid username or password"), window)
		}
	})

	form := container.NewVBox(
		widget.NewLabel("Login"),
		usernameEntry,
		passwordEntry,
		loginBtn,
	)

	window.SetContent(container.NewVBox(form))
}
