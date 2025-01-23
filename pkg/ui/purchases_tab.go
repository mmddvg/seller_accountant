package ui

import (
	"database/sql"
	"fmt"
	"image"
	"inventory/pkg/models"
	"inventory/pkg/usecases"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func purchasesTab(app *usecases.Application, window fyne.Window, writer chan<- bool, reader <-chan bool) *fyne.Container {
	purchasesGrid, updatePurchases := gridPurchases(app, window)
	factorsContainer, createBtn := createPurchase(app, window, writer)

	go func() {
		for range reader {
			updatePurchases()
		}
	}()

	return container.NewVBox(
		widget.NewLabel("خرید ها"),
		purchasesGrid,
		widget.NewSeparator(),
		factorsContainer,
		createBtn,
	)
}

func gridPurchases(app *usecases.Application, window fyne.Window) (*container.Scroll, func()) {
	gridContainer := container.NewVBox()

	scroll := container.NewScroll(gridContainer)

	refreshGrid := func() {
		gridContainer.Objects = nil

		headers := container.NewGridWithColumns(3,
			widget.NewLabelWithStyle("آیدی", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("تاریخ", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("رسید ها", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)
		gridContainer.Add(headers)

		purchases, err := app.GetAllPurchases()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		for _, purchase := range purchases {
			row := container.NewGridWithColumns(3,
				widget.NewLabel(fmt.Sprintf("%d", purchase.ID)),
				widget.NewLabel(purchase.CreatedAt.Format(time.RFC1123)),
				container.NewVBox(createFactorsGrid(purchase.Factors, window)...),
			)
			gridContainer.Add(row)
		}

		gridContainer.Refresh()
	}

	refreshGrid()

	scroll.SetMinSize(fyne.NewSize(600, 300))
	return scroll, refreshGrid
}

func createFactorsGrid(factors []models.Factor, window fyne.Window) []fyne.CanvasObject {
	factorHeaders := container.NewHBox(
		widget.NewLabelWithStyle("نام فروشگاه", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("قیمت", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("عملیات", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	factorRows := []fyne.CanvasObject{factorHeaders}
	for _, factor := range factors {
		showImageBtn := widget.NewButton("نشان دادن عکس", func() {
			if factor.FileName.Valid {
				file, err := os.Open(factor.FileName.String)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to open image file: %w", err), window)
					return
				}
				defer file.Close()

				img, _, err := image.Decode(file)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to decode image: %w", err), window)
					return
				}

				imgCanvas := canvas.NewImageFromImage(img)
				imgCanvas.FillMode = canvas.ImageFillContain
				imgCanvas.SetMinSize(fyne.NewSize(300, 300))

				dialog.ShowCustom("Image", "Close", imgCanvas, window)
			} else {
				dialog.ShowInformation("No Image", "No image file provided for this factor.", window)
			}
		})

		row := container.NewHBox(
			widget.NewLabel(factor.StoreName),
			widget.NewLabel(fmt.Sprintf("%d", factor.Price)),
			showImageBtn,
		)
		factorRows = append(factorRows, row)
	}

	return factorRows
}

const message string = "فایلی انتخاب نشده است"

func createPurchase(app *usecases.Application, window fyne.Window, refresh chan<- bool) (*fyne.Container, *widget.Button) {
	factorsContainer := container.NewVBox()

	addFactorRow := func() {
		storeEntry := widget.NewEntry()
		storeEntry.PlaceHolder = "نام فروشگاه"

		priceEntry := widget.NewEntry()
		priceEntry.PlaceHolder = "قیمت"

		fileLabel := widget.NewLabel(message)
		filePicker := widget.NewButton("انتخاب فایل", func() {
			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				if reader != nil {
					fileLabel.SetText(reader.URI().Name())
					fileLabel.Text = reader.URI().Path()
				}
			}, window)
		})

		removeButton := widget.NewButton("پاک کردن", func() {
			factorsContainer.Remove(factorsContainer.Objects[len(factorsContainer.Objects)-1])
			factorsContainer.Refresh()
		})

		row := container.New(
			layout.NewGridLayoutWithColumns(5),
			storeEntry,
			priceEntry,
			fileLabel,
			filePicker,
			removeButton,
		)

		factorsContainer.Add(row)
		factorsContainer.Refresh()
	}

	addFactorRow()

	addFactorBtn := widget.NewButton("ایجاد رسید", addFactorRow)
	createBtn := widget.NewButton("ایجاد خرید", func() {
		factors := []models.Factor{}
		for _, obj := range factorsContainer.Objects {
			if row, ok := obj.(*fyne.Container); ok {
				storeEntry := row.Objects[0].(*widget.Entry)
				priceEntry := row.Objects[1].(*widget.Entry)
				fileLabel := row.Objects[2].(*widget.Label)

				storeName := strings.TrimSpace(storeEntry.Text)
				if storeName == "" {
					dialog.ShowError(fmt.Errorf("نام فروشگاه الزامی است"), window)
					return
				}

				price, err := strconv.Atoi(strings.TrimSpace(priceEntry.Text))
				if err != nil {
					dialog.ShowError(fmt.Errorf("قیمت نامعتبر است %s", storeName), window)
					return
				}

				factor := models.Factor{
					StoreName: storeName,
					Price:     price,
				}

				if fileLabel.Text != message {
					factor.FileName = sql.NullString{String: fileLabel.Text, Valid: true}
				}

				factors = append(factors, factor)
			}
		}

		_, err := app.CreatePurchase(factors)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		refresh <- true
		factorsContainer.Objects = nil
		addFactorRow()
	})

	return container.NewVBox(
		factorsContainer,
		addFactorBtn,
	), createBtn
}
