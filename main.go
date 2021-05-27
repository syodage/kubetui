package main

import (
	//	"fmt"

	"github.com/rivo/tview"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	main := newPrimitive("Main")
	context := newPrimitive("Context")

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
		
	grid := tview.NewGrid().
		SetRows(4, 0).
		SetColumns(20, 0).
		SetBorders(true).
		AddItem(context, 0, 0, 1, 2, 0, 0, false).
		AddItem(sidebar, 1, 0, 1, 1, 0, 0, false).
		AddItem(main, 1, 1, 1, 1, 0, 0, false)

	if err := tview.NewApplication().SetRoot(grid, true).SetFocus(sidebar).Run(); err != nil {
		panic(err)
	}
}
