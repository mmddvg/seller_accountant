package ui

import (
	"fmt"
	"inventory/pkg/models"
	"inventory/pkg/usecases"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/samber/lo"
)

func factorsTab(app *usecases.Application, window fyne.Window, refreshAccs chan bool) *fyne.Container {
	factors := app.ListFactors()

	factorsList := widget.NewList(
		func() int {
			return len(factors)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			factor := factors[i]

			var tmp string = ""
			for _, v := range factor.Products {
				tmp += v.ToStr() + " "
			}
			item.(*widget.Label).SetText(
				fmt.Sprintf("ID: %d, Account ID: %d, %v ", factor.Id, factor.AccountId, tmp),
			)
		},
	)

	refreshFactorsList := func() {
		factors = app.ListFactors()
		factorsList.Refresh()
	}

	accIDEntry := widget.NewEntry()
	accIDEntry.SetPlaceHolder("Account ID")
	productIDsEntry := widget.NewEntry()
	productIDsEntry.SetPlaceHolder("Product IDs (comma separated)")

	createBtn := widget.NewButton("Create Factor", func() {
		accID, err := strconv.Atoi(accIDEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid account id"), window)
			return
		}

		productIDs := splitAndConvert(productIDsEntry.Text)
		if len(productIDs) == 0 {
			dialog.ShowError(fmt.Errorf("invalid product ids"), window)
			return
		}

		_, err = app.CreateFactor(uint(accID), productIDs)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("Success", "Factor created successfully!", window)
		accIDEntry.SetText("")
		productIDsEntry.SetText("")

		refreshFactorsList()
		refreshAccs <- true
	})

	return container.NewVBox(
		widget.NewLabel("Create Factor"),
		accIDEntry,
		productIDsEntry,
		createBtn,
		widget.NewSeparator(),
		widget.NewLabel("List of Factors"),
		factorsList,
	)
}
func splitAndConvert(input string) []models.FactorProduct {
	res := []models.FactorProduct{}
	for _, part := range strings.Split(input, ",") {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil {
			_, index, ok := lo.FindIndexOf(res, func(fp models.FactorProduct) bool {
				return fp.ProductId == uint(id)
			})
			if ok {
				res[index].Count += 1
			} else {
				res = append(res, models.FactorProduct{ProductId: uint(id), Count: 1})
			}
		}
	}
	return res
}
