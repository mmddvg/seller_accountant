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
	usernameEntry.SetPlaceHolder("نام کاربری")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("رمز عبور")

	loginBtn := widget.NewButton("ورود", func() {
		if app.Login(usernameEntry.Text, passwordEntry.Text) {
			window.Hide()
			MainWindow(app, fyneApp)
		} else {
			dialog.ShowError(fmt.Errorf("نام کاربری یا رمز عبور نامعتبر است"), window)
		}
	})

	form := container.NewVBox(
		widget.NewLabel("ورود"),
		usernameEntry,
		passwordEntry,
		loginBtn,
	)

	window.SetContent(container.NewVBox(form))
}
