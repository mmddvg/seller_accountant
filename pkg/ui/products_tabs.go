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

func productsTab(app *usecases.Application, window fyne.Window) *fyne.Container {
	products := app.ListProducts()
	productsList := widget.NewList(
		func() int {
			return len(products)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(
				fmt.Sprintf("ID: %d, Name: %s, Price: %d", products[i].Id, products[i].Name, products[i].Price),
			)
		},
	)

	refreshList := func() {
		products = app.ListProducts()
		productsList.Refresh()
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Product Name")
	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Product Price")
	createBtn := widget.NewButton("Create Product", func() {
		name := nameEntry.Text
		price, err := strconv.Atoi(priceEntry.Text)
		if name == "" || err != nil || price <= 0 {
			dialog.ShowError(fmt.Errorf("invalid product details"), window)
			return
		}
		_, err = app.CreateProduct(models.NewProduct{Name: name, Price: uint(price)})
		if err != nil {
			dialog.ShowError(err, window)
		}
		nameEntry.SetText("")
		priceEntry.SetText("")
		refreshList()
	})

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Product ID")
	newPriceEntry := widget.NewEntry()
	newPriceEntry.SetPlaceHolder("New Price")
	updateBtn := widget.NewButton("Update Price", func() {
		id, err := strconv.Atoi(idEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid Product ID"), window)
			return
		}
		price, err := strconv.Atoi(newPriceEntry.Text)
		if err != nil || price <= 0 {
			dialog.ShowError(fmt.Errorf("invalid Price"), window)
			return
		}
		_, err = app.UpdateProduct(uint(id), uint(price))
		if err != nil {
			dialog.ShowError(err, window)
		}
		idEntry.SetText("")
		newPriceEntry.SetText("")
		refreshList()
	})

	return container.NewVBox(
		widget.NewLabel("Products"),
		productsList,
		widget.NewSeparator(),
		widget.NewLabel("Add Product"),
		nameEntry,
		priceEntry,
		createBtn,
		widget.NewSeparator(),
		widget.NewLabel("Update Product Price"),
		idEntry,
		newPriceEntry,
		updateBtn,
	)
}
