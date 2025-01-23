package ui

import (
	"fmt"
	"inventory/pkg/models"
	"inventory/pkg/usecases"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func accountsTab(app *usecases.Application, window fyne.Window, writer chan<- bool, reader <-chan bool) *fyne.Container {
	table, updateTable := customersGrid(app, window)

	nameEntry, createBtn := createCustomer(window, app, updateTable)

	idEntry, chargeEntry, updateBtn := updateCustomer(window, app, updateTable)

	sellSection := sell(app, window, writer)

	go func() {
		for range reader {
			updateTable()
		}
	}()

	inputSection := container.NewVBox(
		container.NewVBox(
			widget.NewLabel("ایجاد مشتری"),
			container.NewGridWithColumns(2, nameEntry, createBtn),
		),
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabel("افزایش اعتبار مشتری"),
			container.NewGridWithColumns(3, idEntry, chargeEntry, updateBtn),
		),
		widget.NewSeparator(),
	)

	return container.NewVBox(
		widget.NewLabel("مشتری ها"),
		table,
		widget.NewSeparator(),
		inputSection,
		widget.NewSeparator(),
		widget.NewLabel("فروش"),
		sellSection,
	)
}

func createCustomer(window fyne.Window, app *usecases.Application, updateTable func()) (*widget.Entry, *widget.Button) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("نام مشتری")

	createBtn := widget.NewButton("ایجاد مشتری", func() {
		name := nameEntry.Text
		if name == "" {
			dialog.ShowError(fmt.Errorf("نام مشتری الزامی است"), window)
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
	idEntry.SetPlaceHolder("آیدی مشتری")

	chargeEntry := widget.NewEntry()
	chargeEntry.SetPlaceHolder("اعتبار")

	updateBtn := widget.NewButton("افزایش اعتبار", func() {
		id, err := strconv.Atoi(idEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("آیدی نامعتبر است"), window)
			return
		}
		balance, err := strconv.Atoi(chargeEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("ورودی نامعتبر است"), window)
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

func customersGrid(app *usecases.Application, window fyne.Window) (*container.Scroll, func()) {
	gridContainer := container.NewVBox()

	scroll := container.NewScroll(gridContainer)

	refreshGrid := func() {
		sales, err := app.GetSales()
		if err != nil {
			dialog.ShowError(err, window)
		}

		salesMap := make(map[uint][]models.Sale)
		for _, sale := range sales {
			salesMap[sale.CustomerId] = append(salesMap[sale.CustomerId], sale)
		}

		gridContainer.Objects = nil

		headers := container.NewGridWithColumns(4,
			widget.NewLabelWithStyle("آیدی", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("نام", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("اعتبار", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("فروش", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)
		gridContainer.Add(headers)

		accounts := app.ListAccounts()
		for _, account := range accounts {
			salesGrid := createSalesGrid(salesMap[uint(account.ID)])

			row := container.NewGridWithColumns(4,
				widget.NewLabel(fmt.Sprintf("%d", account.ID)),
				widget.NewLabel(account.Name),
				widget.NewLabel(fmt.Sprintf("%d", account.Charge)),
				salesGrid,
			)
			gridContainer.Add(row)
		}

		gridContainer.Refresh()
	}
	refreshGrid()

	scroll.SetMinSize(fyne.NewSize(600, 300))
	return scroll, refreshGrid
}

func createSalesGrid(sales []models.Sale) *fyne.Container {
	if len(sales) == 0 {
		return container.NewVBox(widget.NewLabel("No sales"))
	}

	var rows []fyne.CanvasObject
	for _, sale := range sales {
		row := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("آیدی فروش: %d", sale.Id)),
			widget.NewLabel(fmt.Sprintf("قیمت: %d", sale.Price)),
		)
		rows = append(rows, row)
	}

	return container.NewVBox(rows...)
}

func sell(app *usecases.Application, window fyne.Window, writer chan<- bool) *fyne.Container {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("نام مشتری")

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("قیمت")

	sellBtn := widget.NewButton("فروش", func() {
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

	return container.NewVBox(
		container.NewGridWithColumns(3, nameEntry, priceEntry, sellBtn),
	)
}
