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

func accountsTab(app *usecases.Application, window fyne.Window, writer chan<- bool, reader <-chan bool) *fyne.Container {
	table, updateTable := customersGrid(app)

	nameEntry, createBtn := createCustomer(window, app, updateTable)

	idEntry, chargeEntry, updateBtn := updateCustomer(window, app, updateTable)

	go func() {
		for range reader {
			updateTable()
		}
	}()

	return container.NewVBox(
		widget.NewLabel("Customers"),
		table,
		widget.NewSeparator(),
		widget.NewLabel("Create Customer"),
		nameEntry,
		createBtn,
		widget.NewSeparator(),
		widget.NewLabel("Increase Customer Balance"),
		idEntry,
		chargeEntry,
		updateBtn,
		sell(app, window, writer),
	)
}

func createCustomer(window fyne.Window, app *usecases.Application, updateTable func()) (*widget.Entry, *widget.Button) {
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
		updateTable()
	})
	return nameEntry, createBtn
}

func updateCustomer(window fyne.Window, app *usecases.Application, updateTable func()) (*widget.Entry, *widget.Entry, *widget.Button) {
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
		updateTable()
	})

	return idEntry, chargeEntry, updateBtn
}
func customersGrid(app *usecases.Application) (*container.Scroll, func()) {
	gridContainer := container.NewVBox()

	scroll := container.NewScroll(gridContainer)

	refreshGrid := func() {
		gridContainer.Objects = nil

		headers := container.NewHBox(
			widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Balance", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)
		gridContainer.Add(headers)

		accounts := app.ListAccounts()
		for _, account := range accounts {
			row := container.NewHBox(
				widget.NewLabel(fmt.Sprintf("%d", account.ID)),
				widget.NewLabel(account.Name),
				widget.NewLabel(fmt.Sprintf("%d", account.Charge)),
			)
			gridContainer.Add(row)
		}

		gridContainer.Refresh()
	}

	refreshGrid()

	scroll.SetMinSize(fyne.NewSize(200, 200))
	return scroll, refreshGrid
}

func sell(app *usecases.Application, window fyne.Window, writer chan<- bool) *fyne.Container {

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("customer name ")
	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("price")

	sellBtn := widget.NewButton("sell", func() {
		price, err := strconv.ParseUint(priceEntry.Text, 10, 0)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		err = app.Sell(nameEntry.Text, uint(price))
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		nameEntry.SetText("")
		priceEntry.SetText("")

		writer <- true

	})
	content := container.NewVBox(
		nameEntry,
		priceEntry,
		sellBtn,
	)

	return content
}
