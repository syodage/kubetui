package main

import (
	"github.com/rivo/tview"
)

func main() {

	app := tview.NewApplication()

	k := NewKubetui(app)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(k.contextView, 0, 3, false).
		AddItem(tview.NewFlex().
			AddItem(k.menu, 0, 1, false).
			AddItem(k.mainView, 0, 7, false),
			0, 18, false).
			// set fixed size 3, which shows only one line, set propotion to 1 otherwise.
		AddItem(k.logView, 3, 0, false).
		// place holder for additional shortcut keys
		AddItem(nil, 1, 0, false)

	if err := app.SetRoot(flex, true).SetFocus(k.menu).Run(); err != nil {
		panic(err)
	}

}
