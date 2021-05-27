package main

import (
	//	"fmt"

	"log"

	"github.com/rivo/tview"
)

func main() {
	sidebar := tview.NewList().
		AddItem("contexts", "", 0, nil).
		AddItem("deployment", "", 0, nil).
		AddItem("namespace", "", 0, nil).
		AddItem("pods", "", 0, nil).
		AddItem("services", "", 0, nil).
		AddItem("nodes", "", 0, nil).
		AddItem("endpoints", "", 0, nil).
		// SetCurrentItem(0).
		SetSelectedFocusOnly(true).
		ShowSecondaryText(false)
	sidebar.Box.SetTitle("Menu").SetBorder(true)

	context := tview.NewBox().SetBorder(true).SetTitle("Context")
	main := tview.NewBox().SetBorder(true).SetTitle("Main")
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(context, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(sidebar, 0, 1, false).
			AddItem(main, 0, 9, false),
			0, 3, false)

	if err := tview.NewApplication().SetRoot(flex, true).SetFocus(sidebar).Run(); err != nil {
		panic(err)
	}

}
