package main

// This is kubetui keys(KKey) use for keybindings
// Only used for runes for now.
// eventually mapping runes will be loaded by a user provided config file.
type KKey rune

const (
	KKeyUp     = rune('j')
	KKeyDown   = rune('k')
	KKeyLeft   = rune('l')
	KKeyRight  = rune('h')
	KKeySelect = rune(' ')
)

// This is the app main state, this will be primirally used to idetnify
// what data should be show in the Main view.
// Other views also can use this to update themselves
type State int16

const (
	NOOP State = iota
	// states related to kubernete contexts
	CtxMain
	// states related to kubernete namespaces
	NpMain
	// states related to kubernete deployments
	DpMain
)
