package main

// This is kubetui keys(KKey) use for keybindings
// Only used for runes for now.
// eventually mapping runes will be loaded by a user provided config file.
type KKey rune

const (
	KKeyUp     = rune('k')
	KKeyDown   = rune('j')
	KKeyLeft   = rune('l')
	KKeyRight  = rune('h')
	KKeySelect = rune(' ')
)

type KView int16

const (
	MAIN_VIEW KView = iota
	MENU_VIEW
	INFO_VIEW
)

// This is the app main state, this will be primirally used to idetnify
// what data should be show in the Main view.
// Other views also can use this to update themselves
type State int16

const (
	NOOP State = iota
	// states related to kubernete contexts
	CONTEXTS
	// states related to kubernete namespaces
	NAMESPACES
	// states related to kubernete deployments
	DEPLOYMENTS
	// pods
	PODS
	// nodes
	NODES
	// services
	SERVICES
	// endpoints
	ENDPOINTS
)
