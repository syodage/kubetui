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
	Title string
}

func newMenuItem(name string, st State) *MenuItem {
	return &MenuItem{
		Name:  name,
		State: st,
		Title: strings.Title(name),
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
	menu.SetTitle("Menu").SetBorder(true)
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
		m.ctx.activeMenuTitle = m.menuItems[m.activeIndex].Title
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

type ViewUpdate struct {
	SubTitle   string
	Command    []string
	HasHeader  bool
	ToCells    func(data string) [][]string
	Update     func(data [][]string)
	ActionFunc func(row, col int, text string)
}

type MainView struct {
	*tview.Table
	ctx        *KContext
	activeRow  int
	actionFunc func(row, col int, text string)
}

func NewMainView(ctx *KContext) *MainView {
	m := &MainView{
		Table: tview.NewTable(),
		ctx:   ctx,
	}
	// set Box properties
	m.SetTitle("Main").
		SetBorder(true)
	ctx.activeMenuTitle = "KubeTui"
	// m.updateView(KUBETUI_BANNER)
	m.SetSelectedStyle(ctx.lineSelectionStyle)
	return m
}

func (m *MainView) HandleStateChange(ev KEvent) {

	vu := &ViewUpdate{}

	vu.SubTitle = m.ctx.activeMenuTitle
	vu.HasHeader = true
	vu.ToCells = FormatData
	vu.Update = func(data [][]string) {
		for r, cells := range data {
			for c, cell := range cells {
				m.Table.SetCellSimple(r, c, cell)
			}
		}
	}

	vu.ActionFunc = func(r, c int, txt string) {
		m.ctx.LogMsg(fmt.Sprintf(`{row:%d, col:%d, text:%v}`, r, c, txt))
	}

	switch ev.State {
	case NAMESPACES:
		vu.Command = NewKubectl().Get().Namespaces().Build()
	case CONTEXTS:
		vu.Command = NewKubectl().Configs("get-contexts").Build()
	case DEPLOYMENTS:
		vu.Command = NewKubectl().Get().Deployments().WithAllNamespaces().Build()
	case PODS:
		vu.Command = NewKubectl().Get().Pods().WithAllNamespaces().Build()
	case NODES:
		vu.Command = NewKubectl().Get().Nodes().WithAllNamespaces().Build()
	case SERVICES:
		vu.Command = NewKubectl().Get().Services().WithAllNamespaces().Build()
	case ENDPOINTS:
		vu.Command = NewKubectl().Get().Endpoints().WithAllNamespaces().Build()
	default:
		vu.Command = []string{"echo", "Not yet implemented"}
	}

	// update the view data
	m.update(vu)
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
		if m.actionFunc != nil {
			data := m.GetCell(m.activeRow, 0).Text
			m.actionFunc(m.activeRow, 0, data)
		}
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

func (m *MainView) SetViewTitle(subTitle string) {
	m.SetTitle(fmt.Sprintf(`Main - %v`, subTitle))
}

// Here we update the data in the view and redraw UI
func (m *MainView) update(vu *ViewUpdate) {
	if vu.Command == nil {
		panic("Invalid command")
	}
	if vu.SubTitle != "" {
		m.SetViewTitle(vu.SubTitle)
	}

	m.actionFunc = vu.ActionFunc
	data := executeCmd(vu.Command)
	m.ctx.LogCommand(vu.Command)
	m.ctx.queueUpdateDraw(func() {
		m.Clear()
		vu.Update(vu.ToCells(data))
	})
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
