package ui

import (
	"fmt"
	"inventory/pkg/models"
	"inventory/pkg/usecases"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func purchasesTab(app *usecases.Application, window fyne.Window, writer chan<- bool, reader <-chan bool) *fyne.Container {
	purchasesGrid, updatePurchases := gridPurchases(app, window)
	factorsEntry, createBtn := createPurchase(app, window, writer)

	go func() {
		for range reader {
			updatePurchases()
		}
	}()

	return container.NewVBox(
		widget.NewLabel("Purchases"),
		purchasesGrid,
		widget.NewSeparator(),
		factorsEntry,
		createBtn,
	)
}

func gridPurchases(app *usecases.Application, window fyne.Window) (*container.Scroll, func()) {
	gridContainer := container.NewVBox()

	scroll := container.NewScroll(gridContainer)

	refreshGrid := func() {
		gridContainer.Objects = nil

		headers := container.NewHBox(
			widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Created At", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Factors", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)
		gridContainer.Add(headers)

		purchases, err := app.GetAllPurchases()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		for _, purchase := range purchases {
			row := container.NewHBox(
				widget.NewLabel(fmt.Sprintf("%d", purchase.ID)),
				widget.NewLabel(purchase.CreatedAt.Format(time.RFC1123)),
				container.NewVBox(createFactorsGrid(purchase.Factors)...),
			)
			gridContainer.Add(row)
		}

		gridContainer.Refresh()
	}

	refreshGrid()

	scroll.SetMinSize(fyne.NewSize(600, 300))

	return scroll, refreshGrid
}

func createFactorsGrid(factors []models.Factor) []fyne.CanvasObject {
	factorHeaders := container.NewHBox(
		widget.NewLabelWithStyle("Store Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Price", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	factorRows := []fyne.CanvasObject{factorHeaders}
	for _, factor := range factors {
		row := container.NewHBox(
			widget.NewLabel(factor.StoreName),
			widget.NewLabel(fmt.Sprintf("%d", factor.Price)),
		)
		factorRows = append(factorRows, row)
	}

	return factorRows
}

func createPurchase(app *usecases.Application, window fyne.Window, refresh chan<- bool) (*widget.Entry, *widget.Button) {
	factorsEntry := widget.NewEntry()
	factorsEntry.PlaceHolder = "store1 : price , store2 : price"

	createBtn := widget.NewButton("Create Purchase", func() {
		factors := splitAndConvert(factorsEntry.Text, window)
		defer factorsEntry.SetText("")
		_, err := app.CreatePurchase(factors)
		if err != nil {
			dialog.ShowError(err, window)
		}

		refresh <- true
	})

	return factorsEntry, createBtn
}

func splitAndConvert(arg string, window fyne.Window) []models.Factor {
	res := []models.Factor{}
	for _, v := range strings.Split(arg, ",") {
		d := strings.Split(v, ":")
		if len(d) < 2 {
			dialog.ShowError(fmt.Errorf("invalid format, expected 'store:price'"), window)
			return res
		}

		price, err := strconv.ParseUint(strings.TrimSpace(d[1]), 10, 0)
		if err != nil {
			dialog.ShowError(err, window)
			return res
		}

		res = append(res, models.Factor{StoreName: strings.TrimSpace(d[0]), Price: int(price)})
	}
	return res
}
