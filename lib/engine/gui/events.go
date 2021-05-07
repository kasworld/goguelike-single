// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/kasworld/goguelike-single/lib/engine/dispatcheri"
	"github.com/kasworld/goguelike-single/lib/engine/window"
)

// Core events sent by the GUI manager.
// The target panel is the panel immediately under the mouse cursor.
const (
	// Events sent to target panel's lowest subscribed ancestor
	OnMouseDown = window.OnMouseDown // Any mouse button is pressed
	OnMouseUp   = window.OnMouseUp   // Any mouse button is released
	OnScroll    = window.OnScroll    // Scrolling mouse wheel

	// Events sent to all panels except the ancestors of the target panel
	OnMouseDownOut = dispatcheri.EventName("gui.OnMouseDownOut") // Any mouse button is pressed
	OnMouseUpOut   = dispatcheri.EventName("gui.OnMouseUpOut")   // Any mouse button is released

	// Event sent to new target panel and all of its ancestors up to
	// (not including) the common ancestor of the new and old targets
	// Cursor entered the panel or a descendant
	OnCursorEnter = dispatcheri.EventName("gui.OnCursorEnter")
	// Event sent to old target panel and all of its ancestors up to
	// (not including) the common ancestor of the new and old targets
	// Cursor left the panel or a descendant
	OnCursorLeave = dispatcheri.EventName("gui.OnCursorLeave")
	// Event sent to the cursor-focused DispatcherI if any, else sent to
	// target panel's lowest subscribed ancestor
	// Cursor is over the panel
	OnCursor = window.OnCursor

	// Event sent to the new key-focused DispatcherI,
	// specified on a call to gui.Manager().SetKeyFocus(dispatcheri.DispatcherI)
	// All keyboard events will be exclusively sent to the receiving DispatcherI
	OnFocus = dispatcheri.EventName("gui.OnFocus")
	// Event sent to the previous key-focused DispatcherI when another panel is key-focused
	// Keyboard events will stop being sent to the receiving DispatcherI
	OnFocusLost = dispatcheri.EventName("gui.OnFocusLost")

	// Events sent to the key-focused DispatcherI
	OnKeyDown   = window.OnKeyDown   // A key is pressed
	OnKeyUp     = window.OnKeyUp     // A key is released
	OnKeyRepeat = window.OnKeyRepeat // A key was pressed and is now automatically repeating
	OnChar      = window.OnChar      // A unicode key is pressed
)

const (
	// Panel size changed (no parameters)
	OnResize = dispatcheri.EventName("gui.OnResize")
	// Panel enabled/disabled (no parameters)
	OnEnable = dispatcheri.EventName("gui.OnEnable")
	// Widget clicked by mouse left button or via key press
	OnClick = dispatcheri.EventName("gui.OnClick")
	// Value was changed. Emitted by List, DropDownList, CheckBox and Edit
	OnChange = dispatcheri.EventName("gui.OnChange")
	// Radio button within a group changed state
	OnRadioGroup = dispatcheri.EventName("gui.OnRadioGroup")
)
