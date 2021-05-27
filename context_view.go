package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type ContextView struct {
	*tview.Table
}

// return a new ContextView and not-nil value if some error occurs
func NewContextView() (ctx *ContextView, err error) {
	ctx = &ContextView{
		Table: tview.NewTable(),
	}
	t := ctx.Table
	b := ctx.Box
	b.SetTitle("Context").
		SetBorder(true)

	versions := GetVersion()

	data := []*Kv{
		NewKv("cluster", "Minikube"),
		NewKv("context", "default-context"),
		NewKv("namespace", "default"),
		NewKv("kubernete", versions.Kubernetes),
		NewKv("kubectl", versions.Kubectl),
		NewKv("kubetui", "v0.0.1"),
	}

	for r, kv := range data {
		t.SetCellSimple(r, 0, fmt.Sprintf("%v: %v", kv.Key, kv.Value))
	}
	return
}
