package ui

import (
	"fmt"
	"inventory/pkg/models"
	"inventory/pkg/usecases"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/samber/lo"
)

func dashBoardTab(app *usecases.Application, window fyne.Window, reader <-chan bool) *fyne.Container {
	customersBinding := binding.NewString()
	spentBinding := binding.NewString()
	profitBinding := binding.NewString()
	netProfitBinding := binding.NewString()

	updateCustomers(app, customersBinding)
	updateSpendings(app, spentBinding, window)
	updateProfits(app, profitBinding, window)
	updateNetProfit(app, netProfitBinding, window)

	go func() {
		for range reader {
			updateCustomers(app, customersBinding)
			updateSpendings(app, spentBinding, window)
			updateProfits(app, profitBinding, window)
			updateNetProfit(app, netProfitBinding, window)
		}
	}()

	var (
		w1 = widget.NewLabelWithData(customersBinding)
		w2 = widget.NewLabelWithData(spentBinding)
		w3 = widget.NewLabelWithData(profitBinding)
		w4 = widget.NewLabelWithData(netProfitBinding)
	)
	w1.Alignment = fyne.TextAlignCenter
	w2.Alignment = fyne.TextAlignCenter
	w3.Alignment = fyne.TextAlignCenter
	w4.Alignment = fyne.TextAlignCenter

	return container.NewVBox(
		w1, w2, w3, w4,
	)
}

func updateCustomers(app *usecases.Application, customersBinding binding.String) {
	customers := app.ListAccounts()
	customersBinding.Set(fmt.Sprintf("تعداد مشتری ها : %d", len(customers)))
}

func updateSpendings(app *usecases.Application, spentBinding binding.String, window fyne.Window) {
	factors, err := app.GetAllFactors()
	if err != nil {
		dialog.ShowError(fmt.Errorf("error getting factors : %s", err.Error()), window)
		return
	}
	spent := lo.SumBy(factors, func(f models.Factor) int { return f.Price })
	spentBinding.Set(fmt.Sprintf("خرید کل : %d", spent))
}

func updateProfits(app *usecases.Application, profitBinding binding.String, window fyne.Window) {
	sales, err := app.GetSales()
	if err != nil {
		dialog.ShowError(fmt.Errorf("error listing sales : %s", err.Error()), window)
		return
	}
	profit := lo.SumBy(sales, func(s models.Sale) int { return s.Price })
	profitBinding.Set(fmt.Sprintf("فروش کل : %d", profit))
}

func updateNetProfit(app *usecases.Application, netProfitBinding binding.String, window fyne.Window) {
	p, err := app.GetNetProfit()
	if err != nil {
		dialog.ShowError(fmt.Errorf("error listing net profit : %s", err.Error()), window)
		return
	}
	netProfitBinding.Set(fmt.Sprintf(" سود خالص : %d", p))
}
