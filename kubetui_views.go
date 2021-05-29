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
	ctx *KContext
}

// return a new ContextView and not-nil value if some error occurs
func NewContextView(ctx *KContext) (ctxView *ContextView, err error) {
	ctxView = &ContextView{
		Table: tview.NewTable(),
		ctx:   ctx,
	}

	t := ctxView.Table
	b := ctxView.Box
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
	ctx         *KContext
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

func NewMenu(ctx *KContext) *Menu {
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
		ctx: ctx,
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
		m.ctx.logEvents <- "Menu: move down"
	}

	moveUp := func() {
		m.activeIndex--
		if m.activeIndex < 0 {
			m.activeIndex = len(m.menuItems) - 1
		}
		m.ctx.logEvents <- "Menu: move up"
	}

	enter := func() {
		m.selectIndex = m.activeIndex
		m.ctx.logEvents <- "Menu: press enter"
		// send an event to channel which update the Main view as require
		m.ctx.stateEvents <- NewKEvent(m.menuItems[m.selectIndex].State)
	}

	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			moveUp()
		case tcell.KeyDown:
			moveDown()
		case tcell.KeyRune:
			switch event.Rune() {
			case KKeyDown:
				moveDown()
			case KKeyUp:
				moveUp()
			case KKeySelect:
				enter()
			case KKeyLeft:
				m.ctx.logEvents <- "Menu: move focus to Main view"
				m.ctx.focusEvents <- KFocusEvent{
					kview:    MAIN_VIEW,
					setFocus: setFocus,
				}
			}
		case tcell.KeyEnter:
			enter()
		}
	})
}

// ==========================Main View==============================================

type Main struct {
	*tview.Table
	ctx *KContext
}

func NewMain(ctx *KContext) *Main {
	m := &Main{
		Table: tview.NewTable(),
		ctx:   ctx,
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

	// FIXME: menu shouldn't have full access to the app, just depend on functions
	m.ctx.queueUpdateDraw(func() {
		m.Table.Clear()
		update()
	})
}
func (m *Main) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	moveUp := func() {
		m.ctx.logEvents <- "Main: move up"
	}
	moveDown := func() {
		m.ctx.logEvents <- "Main: move down"
	}
	enter := func() {
		m.ctx.logEvents <- "Main: press enter"
	}
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			moveUp()
		case tcell.KeyDown:
			moveDown()
		case tcell.KeyRune:
			switch event.Rune() {
			case KKeyDown:
				moveDown()
			case KKeyUp:
				moveUp()
			case KKeySelect:
				enter()
			case KKeyRight:
				m.ctx.logEvents <- "Main: Move focus to Menu view"
				m.ctx.focusEvents <- KFocusEvent{
					kview:    MENU_VIEW,
					setFocus: setFocus,
				}
			}
		case tcell.KeyEnter:
			enter()
		}
	})
}
func (m *Main) Focus(delegate func(p tview.Primitive)) {
	m.Table.Focus(delegate)
	// m.Table.Clear()
	// m.Table.SetCellSimple(0, 0, "Focused")
	m.Table.ScrollToBeginning()
	m.ctx.logEvents <- "Main: Focused Main View"
	// m.app.Draw()
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

// ==========================Log View==============================================

type LogView struct {
	*tview.TextView
	ctx *KContext
	ln  int64
}

func NewLogView(ctx *KContext) *LogView {
	tv := tview.NewTextView()
	tv.SetTitle("Logs").SetBorder(true)
	lv := &LogView{
		TextView: tv,
		ctx:      ctx,
	}
	lv.Log("Waiting...")
	return lv
}

func (lv *LogView) Log(line string) {
	lv.ln++
	lv.SetText(fmt.Sprintf(`%d: %v`, lv.ln, line))
}
