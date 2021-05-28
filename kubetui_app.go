package main

import (
	"log"
)

type Kubetui struct {
	contextView *ContextView
	menu        *Menu
	mainView    *Main
}

func NewKubetui() *Kubetui {
	menu := NewMenu()
	menu.Box.SetTitle("Menu").SetBorder(true)

	context, err := NewContextView()
	if err != nil {
		log.Panic(err)
	}
	main := NewMain()

	return &Kubetui{
		contextView: context,
		menu:        menu,
		mainView:    main,
	}
}
