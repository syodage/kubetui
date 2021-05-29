package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Kubetui struct {
	app      *tview.Application
	infoView *InfoView
	menuView *MenuView
	mainView *MainView
	logView  *LogView
	// let's start with tracking latest state, maybe need to keep an stack of states
	state   State
	context *KContext
}

// use to pass data to views
type KContext struct {
	queueUpdate        func(func())
	queueUpdateDraw    func(func())
	stateEvents        chan KEvent
	focusEvents        chan KFocusEvent
	logEvents          chan string
	lineSelectionStyle tcell.Style
	activeMenuTitle    string
}

func (ctx *KContext) LogCommand(cmds []string) {
	ctx.logEvents <- fmt.Sprintf(`Exec: %v`, strings.Join(cmds, " "))
}

func (ctx *KContext) LogFocusChange(from, to string) {
	ctx.logEvents <- fmt.Sprintf(`Focus: change %v -> %v`, from, to)
}

func (ctx *KContext) LogMsg(msg string) {
	ctx.logEvents <- fmt.Sprintf(`Msg: %v`, msg)
}

// use to pass state change events
type KEvent struct {
	State State
}

func NewKEvent(st State) KEvent {
	return KEvent{
		State: st,
	}
}

// use to change the focus
type KFocusEvent struct {
	kview    KView
	setFocus func(p tview.Primitive)
}

func NewKubetui(app *tview.Application) *Kubetui {
	ctx := &KContext{}
	ctx.stateEvents = make(chan KEvent, 1)
	ctx.focusEvents = make(chan KFocusEvent, 1)
	ctx.logEvents = make(chan string, 1)
	ctx.queueUpdate = func(invoke func()) {
		app.QueueUpdate(invoke)
	}
	ctx.queueUpdateDraw = func(invoke func()) {
		app.QueueUpdateDraw(invoke)
	}
	ctx.lineSelectionStyle = tcell.StyleDefault.
		Foreground(tcell.ColorGreenYellow).
		Bold(true)

	menuView := NewMenuView(ctx)

	infoView, err := NewInfoView(ctx)
	if err != nil {
		log.Panic(err)
	}
	mainView := NewMainView(ctx)
	mainView.SetViewTitle(ctx.activeMenuTitle)

	logView := NewLogView(ctx)

	kubetui := &Kubetui{
		app:      app,
		infoView: infoView,
		menuView: menuView,
		mainView: mainView,
		logView:  logView,
		state:    NOOP,
		context:  ctx,
	}

	// goroutine for kevents handling
	go func() {
		for {
			select {
			case kev := <-ctx.stateEvents:
				mainView.HandleStateChange(kev)
			case fev := <-ctx.focusEvents:
				switch fev.kview {
				case MAIN_VIEW:
					fev.setFocus(mainView)
				case MENU_VIEW:
					fev.setFocus(menuView)
				case INFO_VIEW:
					fev.setFocus(infoView)
				}
				app.Draw()
			case log := <-ctx.logEvents:
				logView.Log(log)
			}
		}
	}()

	return kubetui
}
