package main

import (
	"log"

	"github.com/rivo/tview"
)

type Kubetui struct {
	app         *tview.Application
	contextView *ContextView
	menu        *Menu
	mainView    *Main
	logView     *LogView
	// let's start with tracking latest state, maybe need to keep an stack of states
	state   State
	context *KContext
}

// use to pass data to views
type KContext struct {
	queueUpdate     func(func())
	queueUpdateDraw func(func())
	stateEvents     chan KEvent
	focusEvents     chan KFocusEvent
	logEvents       chan string
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

	menu := NewMenu(ctx)
	menu.Box.SetTitle("Menu").SetBorder(true)

	context, err := NewContextView(ctx)
	if err != nil {
		log.Panic(err)
	}
	main := NewMain(ctx)
	logView := NewLogView(ctx)

	kubetui := &Kubetui{
		app:         app,
		contextView: context,
		menu:        menu,
		mainView:    main,
		logView:     logView,
		state:       NOOP,
		context:     ctx,
	}

	// goroutine for kevents handling
	go func() {
		for {
			select {
			case kev := <-ctx.stateEvents:
				main.HandleStateChange(kev)
			case fev := <-ctx.focusEvents:
				switch fev.kview {
				case MAIN_VIEW:
					fev.setFocus(main)
				case MENU_VIEW:
					fev.setFocus(menu)
				case CONTEXT_VIEW:
					fev.setFocus(context)
				}
			case log := <-ctx.logEvents:
				logView.Log(log)
			}
		}
	}()

	return kubetui
}
