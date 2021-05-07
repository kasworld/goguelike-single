package window

// Key corresponds to a keyboard key.
type Key int

// ModifierKey corresponds to a set of modifier keys (bitmask).
type ModifierKey int

// MouseButton corresponds to a mouse button.
type MouseButton int

// InputMode corresponds to an input mode.
type InputMode int

// InputMode corresponds to an input mode.
type CursorMode int

// Cursor corresponds to a g3n standard or user-created cursor icon.
type Cursor int

// Standard cursors for G3N.
const (
	ArrowCursor = Cursor(iota)
	IBeamCursor
	CrosshairCursor
	HandCursor
	HResizeCursor
	VResizeCursor
	DiagResize1Cursor
	DiagResize2Cursor
	CursorLast = DiagResize2Cursor
)

// Window event names. See availability per platform below ("x" indicates available).
const ( //                             Desktop | Browser |
	OnWindowPos  = "w.OnWindowPos"  //    x    |         |
	OnWindowSize = "w.OnWindowSize" //    x    |         |
	OnKeyUp      = "w.OnKeyUp"      //    x    |    x    |
	OnKeyDown    = "w.OnKeyDown"    //    x    |    x    |
	OnKeyRepeat  = "w.OnKeyRepeat"  //    x    |         |
	OnChar       = "w.OnChar"       //    x    |    x    |
	OnCursor     = "w.OnCursor"     //    x    |    x    |
	OnMouseUp    = "w.OnMouseUp"    //    x    |    x    |
	OnMouseDown  = "w.OnMouseDown"  //    x    |    x    |
	OnScroll     = "w.OnScroll"     //    x    |    x    |
)

// PosEvent describes a windows position changed event
type PosEvent struct {
	Xpos int
	Ypos int
}

// SizeEvent describers a window size changed event
type SizeEvent struct {
	Width  int
	Height int
}

// KeyEvent describes a window key event
type KeyEvent struct {
	Key  Key
	Mods ModifierKey
}

// CharEvent describes a window char event
type CharEvent struct {
	Char rune
	Mods ModifierKey
}

// MouseEvent describes a mouse event over the window
type MouseEvent struct {
	Xpos   float32
	Ypos   float32
	Button MouseButton
	Mods   ModifierKey
}

// CursorEvent describes a cursor position changed event
type CursorEvent struct {
	Xpos float32
	Ypos float32
	Mods ModifierKey
}

// ScrollEvent describes a scroll event
type ScrollEvent struct {
	Xoffset float32
	Yoffset float32
	Mods    ModifierKey
}
