package ui

import (
	"fmt"
	"inventory/pkg/usecases"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// accountsTab := func() *fyne.Container {
// 	accounts := app.ListAccounts()
// 	accountsList = widget.NewList(
// 		func() int { return len(accounts) },
// 		func() fyne.CanvasObject { return widget.NewLabel("") },
// 		func(i widget.ListItemID, item fyne.CanvasObject) {
// 			item.(*widget.Label).SetText(
// 				fmt.Sprintf("ID: %d, Name: %s, Balance: %d", accounts[i].Id, accounts[i].Name, accounts[i].Charge),
// 			)
// 		},
// 	)

// 	nameEntry := widget.NewEntry()
// 	nameEntry.SetPlaceHolder("Account Name")
// 	createBtn := widget.NewButton("Create Account", func() {
// 		name := nameEntry.Text
// 		if name == "" {
// 			dialog.ShowError(fmt.Errorf("account name cannot be empty"), mainWindow)
// 			return
// 		}
// 		_, err := app.CreateAccount(name)
// 		if err != nil {
// 			dialog.ShowError(err, mainWindow)
// 		}
// 		nameEntry.SetText("")
// 		refreshAccountsList()
// 	})

// 	idEntry := widget.NewEntry()
// 	idEntry.SetPlaceHolder("Account ID")
// 	chargeEntry := widget.NewEntry()
// 	chargeEntry.SetPlaceHolder("New Balance")
// 	updateBtn := widget.NewButton("Update Balance", func() {
// 		id, err := strconv.Atoi(idEntry.Text)
// 		if err != nil {
// 			dialog.ShowError(fmt.Errorf("invalid entry"), mainWindow)
// 			return
// 		}
// 		balance, err := strconv.Atoi(chargeEntry.Text)
// 		if err != nil {
// 			dialog.ShowError(fmt.Errorf("invalid entry"), mainWindow)
// 			return
// 		}
// 		_, err = app.ChargeAccount(uint(id), uint(balance))
// 		if err != nil {
// 			dialog.ShowError(err, mainWindow)
// 		}
// 		idEntry.SetText("")
// 		chargeEntry.SetText("")
// 		refreshAccountsList()
// 	})

// 	return container.NewVBox(
// 		widget.NewLabel("Accounts"),
// 		accountsList,
// 		widget.NewSeparator(),
// 		widget.NewLabel("Add Account"),
// 		nameEntry,
// 		createBtn,
// 		widget.NewSeparator(),
// 		widget.NewLabel("Update Account Balance"),
// 		idEntry,
// 		chargeEntry,
// 		updateBtn,
// 	)
// }

func accountsTab(app *usecases.Application, window fyne.Window, ref chan bool) *fyne.Container {

	accounts := app.ListAccounts()
	accountsList := widget.NewList(
		func() int {
			return len(accounts)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(
				fmt.Sprintf("ID: %d, Name: %s, Balance: %d", accounts[i].Id, accounts[i].Name, accounts[i].Charge),
			)
		},
	)

	refreshList := func() {
		accounts = app.ListAccounts()
		accountsList.Refresh()
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Account Name")
	createBtn := widget.NewButton("Create Account", func() {
		name := nameEntry.Text
		if name == "" {
			dialog.ShowError(fmt.Errorf("account name cannot be empty"), window)
			return
		}
		_, err := app.CreateAccount(name)
		if err != nil {
			dialog.ShowError(err, window)
		}
		nameEntry.SetText("")
		refreshList()
	})

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Account ID")
	chargeEntry := widget.NewEntry()
	chargeEntry.SetPlaceHolder("Balance")
	updateBtn := widget.NewButton("Update Balance", func() {
		id, err := strconv.Atoi(idEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid Account ID"), window)
			return
		}
		balance, err := strconv.Atoi(chargeEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid Balance"), window)
			return
		}
		_, err = app.ChargeAccount(uint(id), uint(balance))
		if err != nil {
			dialog.ShowError(err, window)
		}
		idEntry.SetText("")
		chargeEntry.SetText("")
		refreshList()
	})

	go func() {
		for range ref {
			refreshList()
		}
	}()

	return container.NewVBox(
		widget.NewLabel("Accounts"),
		accountsList,
		widget.NewSeparator(),
		widget.NewLabel("Add Account"),
		nameEntry,
		createBtn,
		widget.NewSeparator(),
		widget.NewLabel("Update Account Balance"),
		idEntry,
		chargeEntry,
		updateBtn,
	)
}
