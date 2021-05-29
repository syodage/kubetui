package main

import (
	"github.com/rivo/tview"
)

func main() {

	app := tview.NewApplication()

	k := NewKubetui(app)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(k.contextView, 0, 2, false).
		AddItem(tview.NewFlex().
			AddItem(k.menu, 0, 1, false).
			AddItem(k.mainView, 0, 7, false),
			0, 8, false).
		AddItem(k.logView, 0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(k.menu).Run(); err != nil {
		panic(err)
	}

}
