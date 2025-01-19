package ui

import (
	"database/sql"
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
	factorsEntry, filePickersContainer, createBtn := createPurchase(app, window, writer)

	go func() {
		for range reader {
			updatePurchases()
		}
	}()

	inputSection := container.NewVBox(
		widget.NewLabel("Add a Purchase"),
		container.NewVBox(
			widget.NewLabel("Factors (store:price, store:price)"),
			factorsEntry,
			filePickersContainer,
		),
		createBtn,
	)

	return container.NewVBox(
		widget.NewLabel("Purchases"),
		purchasesGrid,
		widget.NewSeparator(),
		inputSection,
	)
}

func gridPurchases(app *usecases.Application, window fyne.Window) (*container.Scroll, func()) {
	gridContainer := container.NewVBox()

	scroll := container.NewScroll(gridContainer)

	refreshGrid := func() {
		gridContainer.Objects = nil

		headers := container.NewGridWithColumns(4,
			widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Created At", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Factors", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Actions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)
		gridContainer.Add(headers)

		purchases, err := app.GetAllPurchases()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		for _, purchase := range purchases {
			factorRows := createFactorsGrid(purchase.Factors, window)
			row := container.NewGridWithColumns(4,
				widget.NewLabel(fmt.Sprintf("%d", purchase.ID)),
				widget.NewLabel(purchase.CreatedAt.Format(time.RFC1123)),
				container.NewVBox(factorRows...),
				widget.NewButton("View Details", func() {
					dialog.ShowInformation("Purchase Details", fmt.Sprintf("Purchase ID: %d\nCreated At: %s", purchase.ID, purchase.CreatedAt), window)
				}),
			)
			gridContainer.Add(row)
		}

		gridContainer.Refresh()
	}

	refreshGrid()

	scroll.SetMinSize(fyne.NewSize(700, 300))
	return scroll, refreshGrid
}

func createFactorsGrid(factors []models.Factor, window fyne.Window) []fyne.CanvasObject {
	var factorRows []fyne.CanvasObject

	for _, factor := range factors {
		viewImageBtn := widget.NewButton("View Image", func() {
			if factor.FileName.Valid {
				imagePath := factor.FileName.String
				image := widget.NewLabel(fmt.Sprintf("Displaying image: %s", imagePath)) // Replace with actual image widget later
				dialog.ShowCustom("Factor Image", "Close", image, window)
			} else {
				dialog.ShowInformation("No Image", "No image was provided for this factor.", window)
			}
		})

		row := container.NewGridWithColumns(3,
			widget.NewLabel(factor.StoreName),
			widget.NewLabel(fmt.Sprintf("%d", factor.Price)),
			viewImageBtn,
		)
		factorRows = append(factorRows, row)
	}

	return factorRows
}

func createPurchase(app *usecases.Application, window fyne.Window, refresh chan<- bool) (*widget.Entry, *fyne.Container, *widget.Button) {
	factorsEntry := widget.NewEntry()
	factorsEntry.SetPlaceHolder("store1:price,store2:price")

	filePickersContainer := container.NewVBox()
	imagePaths := map[string]string{}

	updateFilePickers := func() {
		filePickersContainer.Objects = nil

		for _, factor := range strings.Split(factorsEntry.Text, ",") {
			store := strings.Split(factor, ":")[0]
			store = strings.TrimSpace(store)

			if store != "" {
				label := widget.NewLabel(fmt.Sprintf("Image for %s:", store))
				filePickerBtn := widget.NewButton("Choose Image", func() {
					dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
						if err != nil || reader == nil {
							return
						}
						imagePaths[store] = reader.URI().Path()
					}, window)
				})
				filePickersContainer.Add(container.NewHBox(label, filePickerBtn))
			}
		}

		filePickersContainer.Refresh()
	}

	factorsEntry.OnChanged = func(_ string) {
		updateFilePickers()
	}

	createBtn := widget.NewButton("Create Purchase", func() {
		factors := splitAndConvertWithFiles(factorsEntry.Text, imagePaths, window)
		defer func() {
			factorsEntry.SetText("")
			imagePaths = map[string]string{}
			updateFilePickers()
		}()
		_, err := app.CreatePurchase(factors)
		if err != nil {
			dialog.ShowError(err, window)
		}

		refresh <- true
	})

	return factorsEntry, filePickersContainer, createBtn
}

func splitAndConvertWithFiles(factorsInput string, imagePaths map[string]string, window fyne.Window) []models.Factor {
	factors := []models.Factor{}

	for _, entry := range strings.Split(factorsInput, ",") {
		parts := strings.Split(entry, ":")
		if len(parts) < 2 {
			dialog.ShowError(fmt.Errorf("invalid format, expected 'store:price'"), window)
			return factors
		}

		price, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid price format for store: %s", strings.TrimSpace(parts[0])), window)
			return factors
		}

		storeName := strings.TrimSpace(parts[0])
		factor := models.Factor{
			StoreName: storeName,
			Price:     price,
		}

		if imagePath, ok := imagePaths[storeName]; ok {
			factor.FileName = sql.NullString{String: imagePath, Valid: true}
		}

		factors = append(factors, factor)
	}

	return factors
}
