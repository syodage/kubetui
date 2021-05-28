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
	// let's start with tracking latest state, maybe need to keep an stack of states
	state  State
	events chan KEvent
}

type KEvent struct {
	State State
}

func NewKEvent(st State) KEvent {
	return KEvent{
		State: st,
	}
}

func NewKubetui(app *tview.Application) *Kubetui {
	kevents := make(chan KEvent, 1)
	menu := NewMenu(kevents)
	menu.Box.SetTitle("Menu").SetBorder(true)

	context, err := NewContextView()
	if err != nil {
		log.Panic(err)
	}
	main := NewMain(app)

	kubetui := &Kubetui{
		app:         app,
		contextView: context,
		menu:        menu,
		mainView:    main,
		state:       NOOP,
		events:      kevents,
	}

	// goroutine for kevents handling
	go func() {
		for kev := range kevents {
			main.HandleStateChange(kev)
		}
	}()

	return kubetui
}
