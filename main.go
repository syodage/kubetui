package main

import (
	"github.com/rivo/tview"
)

func main() {

	k := NewKubetui()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(k.contextView, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(k.menu, 0, 1, false).
			AddItem(k.mainView, 0, 9, false),
			0, 3, false)

	app := tview.NewApplication()

	if err := app.SetRoot(flex, true).SetFocus(k.menu).Run(); err != nil {
		panic(err)
	}

}
