package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ==========================Context View==============================================
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
		NewKv("cluster", GetClusterName()),
		NewKv("context", GetCurrentContext()),
		NewKv("namespace", GetCurrentNamespace()),
		NewKv("kubernete", versions.Kubernetes),
		NewKv("kubectl", versions.Kubectl),
		NewKv("kubetui", "v0.0.1"),
	}

	for r, kv := range data {
		t.SetCellSimple(r, 0, fmt.Sprintf("%v: %v", kv.Key, kv.Value))
	}
	return
}

// ==========================Menu View==============================================

type Menu struct {
	*tview.Box
	menuItems   []*MenuItem
	activeIndex int
	selectIndex int
}

type MenuItem struct {
	Name string
}

func newMenuItem(name string) *MenuItem {
	return &MenuItem{
		Name: name,
	}
}

func NewMenu() *Menu {
	menu := &Menu{
		Box: tview.NewBox(),
		menuItems: []*MenuItem{
			newMenuItem("contexts"),
			newMenuItem("deployment"),
			newMenuItem("namespace"),
			newMenuItem("pods"),
			newMenuItem("services"),
			newMenuItem("nodes"),
			newMenuItem("endpoints"),
		},
	}

	return menu
}

func (m *Menu) Draw(screen tcell.Screen) {
	m.Box.DrawForSubclass(screen, m)
	x, y, width, height := m.GetInnerRect()
	for index, it := range m.menuItems {
		if index >= height {
			break
		}

		sel := "|" // not active
		if index == m.activeIndex {
			sel = "|❯ " // active
		}

		if index == m.selectIndex {
			sel = "|❯❯ " // active
		}
		line := fmt.Sprintf(`%s[white] %s`, sel, it.Name)
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorRed)
	}
}

func (m *Menu) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

	moveDown := func() {
		m.activeIndex++
		if m.activeIndex >= len(m.menuItems) {
			m.activeIndex = 0
		}
	}

	moveUp := func() {
		m.activeIndex--
		if m.activeIndex < 0 {
			m.activeIndex = len(m.menuItems) - 1
		}
	}

	enter := func() {
		m.selectIndex = m.activeIndex
	}

	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			moveUp()
		case tcell.KeyDown:
			moveDown()
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				moveDown()
			case 'k':
				moveUp()
			// case ' ':
			// 	moveDown()
			}
			case tcell.KeyEnter:
				enter()
		}
	})
}

// ==========================Main View==============================================
