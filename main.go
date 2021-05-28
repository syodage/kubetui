package main

import (
	"log"

	"github.com/rivo/tview"
)

func main() {
	menu := NewMenu()
	menu.Box.SetTitle("Menu").SetBorder(true)

	context, err := NewContextView()
	if err != nil {
		log.Panic(err)
	}
	main := tview.NewBox().SetBorder(true).SetTitle("Main")
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(context, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(menu, 0, 1, false).
			AddItem(main, 0, 9, false),
			0, 3, false)

	app := tview.NewApplication()

	if err := app.SetRoot(flex, true).SetFocus(menu).Run(); err != nil {
		panic(err)
	}

}
