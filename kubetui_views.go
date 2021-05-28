package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const KUBETUI_BANNER = `
 _   __      ___    _____        
| | / /     / _ \  (_   _)       
| |/ /_   _| |_) )___| |_   _ _  
|   <| | | |  _ </ __) | | | | | 
| |\ \ |_| | |_) > _)| | |_| | | 
|_| \_\___/|  __/\___)_|\___/ \_)
           | |                   
           |_|                   
                           v0.0.1
`

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
	events      chan<- KEvent
}

type MenuItem struct {
	Name  string
	State State
}

func newMenuItem(name string, st State) *MenuItem {
	return &MenuItem{
		Name:  name,
		State: st,
	}
}

func NewMenu(kev chan<- KEvent) *Menu {
	menu := &Menu{
		Box: tview.NewBox(),
		menuItems: []*MenuItem{
			newMenuItem("contexts", CONTEXTS),
			newMenuItem("deployment", DEPLOYMENTS),
			newMenuItem("namespace", NAMESPACES),
			newMenuItem("pods", PODS),
			newMenuItem("services", SERVICES),
			newMenuItem("nodes", NODES),
			newMenuItem("endpoints", ENDPOINTS),
		},
		events: kev,
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
		// send an event to channel which update the Main view as require
		m.events <- NewKEvent(m.menuItems[m.selectIndex].State)
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
			case ' ':
				enter()
			}
		case tcell.KeyEnter:
			enter()
		}
	})
}

// ==========================Main View==============================================

type Main struct {
	*tview.Table
	app *tview.Application
}

func NewMain(app *tview.Application) *Main {
	m := &Main{
		Table: tview.NewTable(),
		app:   app,
	}

	// set Box properties
	m.Box.SetTitle("Main").
		SetBorder(true)
	updateTable(m, KUBETUI_BANNER)

	return m
}

func (m *Main) HandleStateChange(ev KEvent) {

	update := func() {
		m.Table.SetCellSimple(0, 0, "Default Values")
	}
	switch ev.State {
	case NOOP:
		update = func() {
			updateSimple(m, "NOOOOOOP")
		}
	case NAMESPACES:
		data := executeCmd([]string{
			"kubectl", "get", "namespaces", "-A"})
		update = func() {
			updateTable(m, data)
		}
	case CONTEXTS:
		data := executeCmd([]string{
			"kubectl", "config", "get-contexts"})
		update = func() {
			updateTable(m, data)
		}
	case DEPLOYMENTS:
		data := executeCmd([]string{
			"kubectl", "get", "deploy", "-A"})
		update = func() {
			updateTable(m, data)
		}
	case PODS:
		data := executeCmd([]string{
			"kubectl", "get", "pods", "-A"})
		update = func() {
			updateTable(m, data)
		}
	case NODES:
		data := executeCmd([]string{
			"kubectl", "get", "nodes", "-A"})
		update = func() {
			updateTable(m, data)
		}
	case SERVICES:
		data := executeCmd([]string{
			"kubectl", "get", "services", "-A"})
		update = func() {
			updateTable(m, data)
		}
	case ENDPOINTS:
		data := executeCmd([]string{
			"kubectl", "get", "endpoints", "-A"})
		update = func() {
			updateTable(m, data)
		}
	default:
		update = func() {
			m.Table.SetCellSimple(0, 0, "Not yet implemented")
		}
	}

	m.app.QueueUpdateDraw(func() {
		m.Table.Clear()
		update()
	})
}

func updateSimple(m *Main, data string) {
	m.Table.SetCellSimple(0, 0, data)
}

func updateTable(m *Main, data string) {
	lines := strings.Split(data, "\n")
	for i, ln := range lines {
		m.Table.SetCellSimple(i, 0, ln)
	}
}
