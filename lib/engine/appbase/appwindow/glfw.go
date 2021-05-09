// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !wasm

package appwindow

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kasworld/goguelike-single/lib/engine/dispatcher"
	"github.com/kasworld/goguelike-single/lib/engine/eventtype"
	"github.com/kasworld/goguelike-single/lib/engine/gls"
	"github.com/kasworld/goguelike-single/lib/engine/gui/assets"
)

// GlfwWindow describes one glfw window
type GlfwWindow struct {
	*glfw.Window                   // Embedded GLFW window
	dispatcher.Dispatcher          // Embedded event dispatcher
	gls                   *gls.GLS // Associated OpenGL State
	fullscreen            bool
	lastX                 int
	lastY                 int
	lastWidth             int
	lastHeight            int
	scaleX                float64
	scaleY                float64

	// Events
	keyEv    KeyEvent
	charEv   CharEvent
	mouseEv  MouseEvent
	posEv    PosEvent
	sizeEv   SizeEvent
	cursorEv CursorEvent
	scrollEv ScrollEvent

	mods ModifierKey // Current modifier keys

	// Cursors
	cursors       map[Cursor]*glfw.Cursor
	lastCursorKey Cursor
}

// Init initializes the GlfwWindow singleton with the specified width, height, and title.
func Init(width, height int, title string) error {

	// Panic if already created
	if win != nil {
		panic(fmt.Errorf("can only call window.Init() once"))
	}

	// OpenGL functions must be executed in the same thread where
	// the context was created (by wmgr.CreateWindow())
	runtime.LockOSThread()

	// Create wrapper window with dispatcher
	w := new(GlfwWindow)
	w.Dispatcher.Initialize()
	var err error

	// Initialize GLFW
	err = glfw.Init()
	if err != nil {
		return err
	}

	// Set window hints
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Samples, 8)
	// Set OpenGL forward compatible context only for OSX because it is required for OSX.
	// When this is set, glLineWidth(width) only accepts width=1.0 and generates an error
	// for any other values although the spec says it should ignore unsupported widths
	// and generate an error only when width <= 0.
	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	// Create window and set it as the current context.
	// The window is created always as not full screen because if it is
	// created as full screen it not possible to revert it to windowed mode.
	// At the end of this function, the window will be set to full screen if requested.
	w.Window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return err
	}
	w.MakeContextCurrent()

	// Create OpenGL state
	w.gls, err = gls.New()
	if err != nil {
		return err
	}

	// Compute and store scale
	fbw, fbh := w.GetFramebufferSize()
	w.scaleX = float64(fbw) / float64(width)
	w.scaleY = float64(fbh) / float64(height)

	// Create map for cursors
	w.cursors = make(map[Cursor]*glfw.Cursor)
	w.lastCursorKey = CursorLast

	// Preallocate GLFW standard cursors
	w.cursors[ArrowCursor] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	w.cursors[IBeamCursor] = glfw.CreateStandardCursor(glfw.IBeamCursor)
	w.cursors[CrosshairCursor] = glfw.CreateStandardCursor(glfw.CrosshairCursor)
	w.cursors[HandCursor] = glfw.CreateStandardCursor(glfw.HandCursor)
	w.cursors[HResizeCursor] = glfw.CreateStandardCursor(glfw.HResizeCursor)
	w.cursors[VResizeCursor] = glfw.CreateStandardCursor(glfw.VResizeCursor)

	// Preallocate extra G3N standard cursors (diagonal resize cursors)
	cursorDiag1Png := assets.MustAsset("cursors/diag1.png") // [/]
	cursorDiag2Png := assets.MustAsset("cursors/diag2.png") // [\]
	diag1Img, _, err := image.Decode(bytes.NewReader(cursorDiag1Png))
	diag2Img, _, err := image.Decode(bytes.NewReader(cursorDiag2Png))
	if err != nil {
		return err
	}
	w.cursors[DiagResize1Cursor] = glfw.CreateCursor(diag1Img, 8, 8) // [/]
	w.cursors[DiagResize2Cursor] = glfw.CreateCursor(diag2Img, 8, 8) // [\]

	// Set up key callback to dispatch event
	w.SetKeyCallback(func(x *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		w.keyEv.Key = Key(key)
		w.keyEv.Mods = ModifierKey(mods)
		w.mods = w.keyEv.Mods
		if action == glfw.Press {
			w.Dispatch(eventtype.OnKeyDown, &w.keyEv)
		} else if action == glfw.Release {
			w.Dispatch(eventtype.OnKeyUp, &w.keyEv)
		} else if action == glfw.Repeat {
			w.Dispatch(eventtype.OnKeyRepeat, &w.keyEv)
		}
	})

	// Set up char callback to dispatch event
	w.SetCharModsCallback(func(x *glfw.Window, char rune, mods glfw.ModifierKey) {
		w.charEv.Char = char
		w.charEv.Mods = ModifierKey(mods)
		w.Dispatch(eventtype.OnChar, &w.charEv)
	})

	// Set up mouse button callback to dispatch event
	w.SetMouseButtonCallback(func(x *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		xpos, ypos := x.GetCursorPos()
		w.mouseEv.Button = MouseButton(button)
		w.mouseEv.Mods = ModifierKey(mods)
		w.mouseEv.Xpos = float32(xpos * w.scaleX)
		w.mouseEv.Ypos = float32(ypos * w.scaleY)
		if action == glfw.Press {
			w.Dispatch(eventtype.OnMouseDown, &w.mouseEv)
		} else if action == glfw.Release {
			w.Dispatch(eventtype.OnMouseUp, &w.mouseEv)
		}
	})

	// Set up window size callback to dispatch event
	w.SetSizeCallback(func(x *glfw.Window, width int, height int) {
		fbw, fbh := x.GetFramebufferSize()
		w.sizeEv.Width = width
		w.sizeEv.Height = height
		w.scaleX = float64(fbw) / float64(width)
		w.scaleY = float64(fbh) / float64(height)
		w.Dispatch(eventtype.OnWindowSize, &w.sizeEv)
	})

	// Set up window position callback to dispatch event
	w.SetPosCallback(func(x *glfw.Window, xpos int, ypos int) {
		w.posEv.Xpos = xpos
		w.posEv.Ypos = ypos
		w.Dispatch(eventtype.OnWindowPos, &w.posEv)
	})

	// Set up window cursor position callback to dispatch event
	w.SetCursorPosCallback(func(x *glfw.Window, xpos float64, ypos float64) {
		w.cursorEv.Xpos = float32(xpos)
		w.cursorEv.Ypos = float32(ypos)
		w.cursorEv.Mods = w.mods
		w.Dispatch(eventtype.OnCursor, &w.cursorEv)
	})

	// Set up mouse wheel scroll callback to dispatch event
	w.SetScrollCallback(func(x *glfw.Window, xoff float64, yoff float64) {
		w.scrollEv.Xoffset = float32(xoff)
		w.scrollEv.Yoffset = float32(yoff)
		w.scrollEv.Mods = w.mods
		w.Dispatch(eventtype.OnScroll, &w.scrollEv)
	})

	win = w // Set singleton
	return nil
}

// Gls returns the associated OpenGL state.
func (w *GlfwWindow) Gls() *gls.GLS {

	return w.gls
}

// Fullscreen returns whether this windows is currently fullscreen.
func (w *GlfwWindow) Fullscreen() bool {

	return w.fullscreen
}

// SetFullscreen sets this window as fullscreen on the primary monitor
// TODO allow for fullscreen with resolutions different than the monitor's
func (w *GlfwWindow) SetFullscreen(full bool) {

	// If already in the desired state, nothing to do
	if w.fullscreen == full {
		return
	}
	// Set window fullscreen on the primary monitor
	if full {
		// Get size of primary monitor
		mon := glfw.GetPrimaryMonitor()
		vmode := mon.GetVideoMode()
		width := vmode.Width
		height := vmode.Height
		// Set as fullscreen on the primary monitor
		w.SetMonitor(mon, 0, 0, width, height, vmode.RefreshRate)
		w.fullscreen = true
		// Save current position and size of the window
		w.lastX, w.lastY = w.GetPos()
		w.lastWidth, w.lastHeight = w.GetSize()
	} else {
		// Restore window to previous position and size
		w.SetMonitor(nil, w.lastX, w.lastY, w.lastWidth, w.lastHeight, glfw.DontCare)
		w.fullscreen = false
	}
}

// Destroy destroys this window and its context
func (w *GlfwWindow) Destroy() {

	w.Window.Destroy()
	glfw.Terminate()
	runtime.UnlockOSThread() // Important when using the execution tracer
}

// Scale returns this window's DPI scale factor (FramebufferSize / Size)
func (w *GlfwWindow) GetScale() (x float64, y float64) {

	return w.scaleX, w.scaleY
}

// ScreenResolution returns the screen resolution
func (w *GlfwWindow) ScreenResolution(p interface{}) (width, height int) {

	mon := glfw.GetPrimaryMonitor()
	vmode := mon.GetVideoMode()
	return vmode.Width, vmode.Height
}

// PollEvents process events in the event queue
func (w *GlfwWindow) PollEvents() {

	glfw.PollEvents()
}

// SetSwapInterval sets the number of screen updates to wait from the time SwapBuffer()
// is called before swapping the buffers and returning.
func (w *GlfwWindow) SetSwapInterval(interval int) {

	glfw.SwapInterval(interval)
}

// SetCursor sets the window's cursor.
func (w *GlfwWindow) SetCursor(cursor Cursor) {

	cur, ok := w.cursors[cursor]
	if !ok {
		panic("Invalid cursor")
	}
	w.Window.SetCursor(cur)
}

// CreateCursor creates a new custom cursor and returns an int handle.
func (w *GlfwWindow) CreateCursor(imgFile string, xhot, yhot int) (Cursor, error) {

	// Open image file
	file, err := os.Open(imgFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}
	// Create and store cursor
	w.lastCursorKey += 1
	w.cursors[Cursor(w.lastCursorKey)] = glfw.CreateCursor(img, xhot, yhot)

	return w.lastCursorKey, nil
}

// DisposeCursor deletes the existing custom cursor with the provided int handle.
func (w *GlfwWindow) DisposeCursor(cursor Cursor) {

	if cursor <= CursorLast {
		panic("Can't dispose standard cursor")
	}
	w.cursors[cursor].Destroy()
	delete(w.cursors, cursor)
}

// DisposeAllCursors deletes all existing custom cursors.
func (w *GlfwWindow) DisposeAllCustomCursors() {

	// Destroy and delete all custom cursors
	for key := range w.cursors {
		if key > CursorLast {
			w.cursors[key].Destroy()
			delete(w.cursors, key)
		}
	}
	// Set the next cursor key as the last standard cursor key + 1
	w.lastCursorKey = CursorLast
}

// Center centers the window on the screen.
//func (w *GlfwWindow) Center() {
//
//	// TODO
//}
