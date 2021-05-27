package main

import "github.com/rivo/tview"

type ContextView struct {
  *tview.Table	

}

// return a new ContextView and not-nil value if some error occurs
func NewContextView() (ctx *ContextView, err error) {
	ctx = &ContextView {
		Table: tview.NewTable(),
	}	
	ctx.Table.SetBorder(true)
	return
}