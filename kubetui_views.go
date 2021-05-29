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

// ==========================Info View==============================================
type InfoView struct {
	*tview.Table
	ctx *KContext
}

// return a new HeaderView and not-nil value if some error occurs
func NewInfoView(ctx *KContext) (infoView *InfoView, err error) {
	infoView = &InfoView{
		Table: tview.NewTable(),
		ctx:   ctx,
	}

	t := infoView.Table
	b := infoView.Box
	b.SetTitle("Info").
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

type MenuView struct {
	*tview.Table
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

func NewMenuView(ctx *KContext) *MenuView {
	menu := &MenuView{
		Table: tview.NewTable(),
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
	menu.SetSelectedStyle(ctx.lineSelectionStyle)
	menu.updateView()
	return menu
}

func (m *MenuView) updateView() {
	for row, it := range m.menuItems {
		sel := "|"
		if row == m.selectIndex {
			sel = "| " // select 
		}

		if row == m.activeIndex {
			sel = "|❯ " // active 
		}

		resource := fmt.Sprintf(`%v %v`, sel, it.Name)
		m.SetCellSimple(row, 0, resource)	
	}

	m.Select(m.selectIndex, 0)
}

func (m *MenuView) Focus(delegate func(p tview.Primitive)) {
	m.Table.Focus(delegate)
	m.SetSelectable(true, false)
	m.ctx.LogMsg("[Menu] Focused Menu view")
}

func (m *MenuView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

	moveDown := func() {
		m.selectIndex++
		if m.selectIndex >= len(m.menuItems) {
			m.selectIndex = 0
		}
		m.updateView()
		m.ctx.LogMsg("[Menu] move down")
	}

	moveUp := func() {
		m.selectIndex--
		if m.selectIndex < 0 {
			m.selectIndex = len(m.menuItems) - 1
		}
		m.updateView()
		m.ctx.LogMsg("[Menu] move up")
	}

	enter := func() {
		m.activeIndex = m.selectIndex
		m.updateView()
		// FIXME: if order of following channel push change, UI get hang sometime
		m.ctx.LogMsg("[Menu] press enter")
		// send an event to channel which update the Main view as require
		m.ctx.stateEvents <- NewKEvent(m.menuItems[m.activeIndex].State)
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
				m.ctx.LogFocusChange("Menu", "Main")
				m.ctx.focusEvents <- KFocusEvent{
					kview:    MAIN_VIEW,
					setFocus: setFocus,
				}
				m.SetSelectable(false, false)
			}
		case tcell.KeyEnter:
			enter()
		}
	})
}

// ==========================Main View==============================================

type MainView struct {
	*tview.Table
	ctx       *KContext
	activeRow int
}

func NewMainView(ctx *KContext) *MainView {
	m := &MainView{
		Table: tview.NewTable(),
		ctx:   ctx,
	}
	// set Box properties
	m.SetTitle("Main").
		SetBorder(true)
	m.updateTable(KUBETUI_BANNER)
	m.SetSelectedStyle(ctx.lineSelectionStyle)
	return m
}

func (m *MainView) HandleStateChange(ev KEvent) {

	var cmd []string
	var update func(string)

	switch ev.State {
	case NOOP:
		cmd = []string{"echo", "NOOOOP"}
		update = m.updateSimple
	case NAMESPACES:
		cmd = NewKubectl().Get().Namespaces().Build()
		update = m.updateTable
	case CONTEXTS:
		cmd = NewKubectl().Configs("get-contexts").Build()
		update = m.updateTable
	case DEPLOYMENTS:
		cmd = NewKubectl().Get().Deployments().WithAllNamespaces().Build()
		update = m.updateTable
	case PODS:
		cmd = NewKubectl().Get().Pods().WithAllNamespaces().Build()
		update = m.updateTable
	case NODES:
		cmd = NewKubectl().Get().Nodes().WithAllNamespaces().Build()
		update = m.updateTable
	case SERVICES:
		cmd = NewKubectl().Get().Services().WithAllNamespaces().Build()
		update = m.updateTable
	case ENDPOINTS:
		cmd = NewKubectl().Get().Endpoints().WithAllNamespaces().Build()
		update = m.updateTable
	default:
		update = func(_ string) {
			m.Table.SetCellSimple(0, 0, "Not yet implemented")
		}
	}

	// FIXME: can we remove this new goroutine and just execute thc content as it is in current goroutine?
	go func(cmd []string, update func(data string)) {
		data := ""
		if cmd != nil {
			data = executeCmd(cmd)
		}
		m.ctx.LogCommand(cmd)
		m.ctx.queueUpdateDraw(func() {
			m.Clear()
			update(data)
		})
	}(cmd, update)

}

// Main view key-bindings
func (m *MainView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	moveUp := func() {
		m.activeRow--
		if m.activeRow < 0 {
			m.activeRow = 0
		}
		m.Select(m.activeRow, 0)
		m.ctx.LogMsg("[Main] move up")
	}

	moveDown := func() {
		m.activeRow++
		// TODO: why is this offset? should we strip table data?
		if m.activeRow >= m.GetRowCount()-1 {
			m.activeRow = m.GetRowCount() - 2
		}
		m.Select(m.activeRow, 0)
		m.ctx.LogMsg("[Main] move down")
	}

	enter := func() {
		m.ctx.LogMsg("[Main] press enter")
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
				m.ctx.LogMsg("[Main] Move focus to Menu view")
				m.SetSelectable(false, false)
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

func (m *MainView) Focus(delegate func(p tview.Primitive)) {
	m.Table.Focus(delegate)
	m.SetSelectable(true, false)
	m.ctx.LogMsg("[Main] Focused Main view")
}

func (m *MainView) updateSimple(data string) {
	m.Table.SetCellSimple(0, 0, data)
}

func (m *MainView) updateTable(data string) {
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
