// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/kasworld/goguelike-single/lib/engine/appbase/appwindow"
	"github.com/kasworld/goguelike-single/lib/engine/dispatcher"
	"github.com/kasworld/goguelike-single/lib/engine/dispatcheri"
	"github.com/kasworld/goguelike-single/lib/engine/eventtype"
	"github.com/kasworld/goguelike-single/lib/engine/node"
	"github.com/kasworld/goguelike-single/lib/engine/timermanager"
)

// manager singleton
var gm *manager

// manager routes GUI events to the appropriate panels.
type manager struct {
	dispatcher.Dispatcher                             // Embedded Dispatcher
	timermanager.TimerManager                         // Embedded TimerManager
	win                       appwindow.AppWindowI    // The current AppWindowI
	scene                     node.NodeI              // NodeI containing IPanels to dispatch events to (can contain non-IPanels as well)
	modal                     PanelI                  // Panel which along its descendants will exclusively receive all events
	target                    PanelI                  // Panel immediately under the cursor
	keyFocus                  dispatcheri.DispatcherI // DispatcherI which will exclusively receive all key and char events
	cursorFocus               dispatcheri.DispatcherI // DispatcherI which will exclusively receive all OnCursor events
	cev                       *appwindow.CursorEvent  // DispatcherI which will exclusively receive all OnCursor events
}

// Manager returns the GUI manager singleton (creating it the first time)
func Manager() *manager {

	// Return singleton if already created
	if gm != nil {
		return gm
	}

	gm = new(manager)
	gm.Dispatcher.Initialize()
	gm.TimerManager.Initialize()

	// Subscribe to appwindow events
	gm.win = appwindow.Get()
	gm.win.Subscribe(eventtype.OnKeyUp, gm.onKeyboard)
	gm.win.Subscribe(eventtype.OnKeyDown, gm.onKeyboard)
	gm.win.Subscribe(eventtype.OnKeyRepeat, gm.onKeyboard)
	gm.win.Subscribe(eventtype.OnChar, gm.onKeyboard)
	gm.win.Subscribe(eventtype.OnCursor, gm.onCursor)
	gm.win.Subscribe(eventtype.OnMouseUp, gm.onMouse)
	gm.win.Subscribe(eventtype.OnMouseDown, gm.onMouse)
	gm.win.Subscribe(eventtype.OnScroll, gm.onScroll)

	return gm
}

// Set sets the NodeI to watch for events.
// It's usually a scene containing a hierarchy of INodes.
// The manager only cares about IPanels inside that hierarchy.
func (gm *manager) Set(scene node.NodeI) {

	gm.scene = scene
}

// SetModal sets the specified panel and its descendants to be the exclusive receivers of events.
func (gm *manager) SetModal(ipan PanelI) {

	gm.modal = ipan
	gm.SetKeyFocus(nil)
	gm.SetCursorFocus(nil)
}

// SetKeyFocus sets the key-focused DispatcherI, which will exclusively receive key and char events.
func (gm *manager) SetKeyFocus(disp dispatcheri.DispatcherI) {

	if gm.keyFocus == disp {
		return
	}
	if gm.keyFocus != nil {
		gm.keyFocus.Dispatch(eventtype.OnFocusLost, nil)
	}
	gm.keyFocus = disp
	if gm.keyFocus != nil {
		gm.keyFocus.Dispatch(eventtype.OnFocus, nil)
	}
}

// SetCursorFocus sets the cursor-focused DispatcherI, which will exclusively receive OnCursor events.
func (gm *manager) SetCursorFocus(disp dispatcheri.DispatcherI) {

	if gm.cursorFocus == disp {
		return
	}
	gm.cursorFocus = disp
	if gm.cursorFocus == nil {
		gm.onCursor(eventtype.OnCursor, gm.cev)
	}
}

// onKeyboard is called when char or key events are received.
// The events are dispatched to the focused DispatcherI or to non-GUI.
func (gm *manager) onKeyboard(evname eventtype.EventType, ev interface{}) {

	if gm.keyFocus != nil {
		if gm.modal == nil {
			gm.keyFocus.Dispatch(evname, ev)
		} else if ipan, ok := gm.keyFocus.(PanelI); ok && gm.modal.IsAncestorOf(ipan) {
			gm.keyFocus.Dispatch(evname, ev)
		}
	} else {
		gm.Dispatch(evname, ev)
	}
}

// onMouse is called when mouse events are received.
// OnMouseDown/OnMouseUp are dispatched to gm.target or to non-GUI, while
// OnMouseDownOut/OnMouseUpOut are dispatched to all non-target panels.
func (gm *manager) onMouse(evname eventtype.EventType, ev interface{}) {

	// Check if gm.scene is nil and if so then there are no IPanels to send events to
	if gm.scene == nil {
		gm.Dispatch(evname, ev) // Dispatch event to non-GUI since event was not filtered by any GUI component
		return
	}

	// Dispatch OnMouseDownOut/OnMouseUpOut to all panels except ancestors of target
	gm.forEachIPanel(func(ipan PanelI) {
		if gm.target == nil || !ipan.IsAncestorOf(gm.target) {
			switch evname {
			case eventtype.OnMouseDown:
				ipan.Dispatch(eventtype.OnMouseDownOut, ev)
			case eventtype.OnMouseUp:
				ipan.Dispatch(eventtype.OnMouseUpOut, ev)
			}
		}
	})

	// Appropriately dispatch the event to target panel's lowest subscribed ancestor or to non-GUI or not at all
	if gm.target != nil {
		if gm.modal == nil || gm.modal.IsAncestorOf(gm.target) {
			sendAncestry(gm.target, false, nil, gm.modal, evname, ev)
		}
	} else if gm.modal == nil {
		gm.Dispatch(evname, ev)
	}
}

// onScroll is called when scroll events are received.
// The events are dispatched to the target panel or to non-GUI.
func (gm *manager) onScroll(evname eventtype.EventType, ev interface{}) {

	// Check if gm.scene is nil and if so then there are no IPanels to send events to
	if gm.scene == nil {
		gm.Dispatch(evname, ev) // Dispatch event to non-GUI since event was not filtered by any GUI component
		return
	}

	// Appropriately dispatch the event to target panel's lowest subscribed ancestor or to non-GUI or not at all
	if gm.target != nil {
		if gm.modal == nil || gm.modal.IsAncestorOf(gm.target) {
			sendAncestry(gm.target, false, nil, gm.modal, evname, ev)
		}
	} else if gm.modal == nil {
		gm.Dispatch(evname, ev)
	}
}

// onCursor is called when (mouse) cursor events are received.
// Updates the target/click panels and dispatches OnCursor, OnCursorEnter, OnCursorLeave events.
func (gm *manager) onCursor(evname eventtype.EventType, ev interface{}) {

	// If an DispatcherI is capturing cursor events dispatch to it and return
	if gm.cursorFocus != nil {
		gm.cursorFocus.Dispatch(evname, ev)
		return
	}

	// If gm.scene is nil then there are no IPanels to send events to
	if gm.scene == nil {
		gm.Dispatch(evname, ev) // Dispatch event to non-GUI since event was not filtered by any GUI component
		return
	}

	// Get and store CursorEvent
	gm.cev = ev.(*appwindow.CursorEvent)

	// Temporarily store last target and clear current one
	oldTarget := gm.target
	gm.target = nil

	// Find PanelI immediately under the cursor and store it in gm.target
	gm.forEachIPanel(func(ipan PanelI) {
		if ipan.InsideBorders(gm.cev.Xpos, gm.cev.Ypos) && (gm.target == nil || ipan.Position().Z < gm.target.GetPanel().Position().Z) {
			gm.target = ipan
		}
	})

	// If the cursor is now over a different panel, dispatch OnCursorLeave/OnCursorEnter
	if gm.target != oldTarget {
		// We are only interested in sending events up to the lowest common ancestor of target and oldTarget
		var commonAnc PanelI
		if gm.target != nil && oldTarget != nil {
			commonAnc, _ = gm.target.LowestCommonAncestor(oldTarget).(PanelI)
		}
		// If just left a panel and the new panel is not a descendant of the old panel
		if oldTarget != nil && !oldTarget.IsAncestorOf(gm.target) && (gm.modal == nil || gm.modal.IsAncestorOf(oldTarget)) {
			sendAncestry(oldTarget, true, commonAnc, gm.modal, eventtype.OnCursorLeave, ev)
		}
		// If just entered a panel and it's not an ancestor of the old panel
		if gm.target != nil && !gm.target.IsAncestorOf(oldTarget) && (gm.modal == nil || gm.modal.IsAncestorOf(gm.target)) {
			sendAncestry(gm.target, true, commonAnc, gm.modal, eventtype.OnCursorEnter, ev)
		}
	}

	// Appropriately dispatch the event to target panel's lowest subscribed ancestor or to non-GUI or not at all
	if gm.target != nil {
		if gm.modal == nil || gm.modal.IsAncestorOf(gm.target) {
			sendAncestry(gm.target, false, nil, gm.modal, evname, ev)
		}
	} else if gm.modal == nil {
		gm.Dispatch(evname, ev)
	}
}

// sendAncestry sends the specified event (evname/ev) to the specified target panel and its ancestors.
// If all is false, then the event is only sent to the lowest subscribed ancestor.
// If uptoEx (i.e. excluding) is not nil then the event will not be dispatched to that ancestor nor any higher ancestors.
// If uptoIn (i.e. including) is not nil then the event will be dispatched to that ancestor but not to any higher ancestors.
// uptoEx and uptoIn can both be defined.
func sendAncestry(ipan PanelI, all bool, uptoEx PanelI, uptoIn PanelI, evname eventtype.EventType, ev interface{}) {

	var ok bool
	for ipan != nil {
		if uptoEx != nil && ipan == uptoEx {
			break
		}
		count := ipan.Dispatch(evname, ev)
		if (uptoIn != nil && ipan == uptoIn) || (!all && count > 0) {
			break
		}
		ipan, ok = ipan.Parent().(PanelI)
		if !ok {
			break
		}
	}
}

// traverseIPanel traverses the descendants of the provided PanelI,
// executing the specified function for each PanelI.
func traverseIPanel(ipan PanelI, f func(ipan PanelI)) {

	// If panel not visible, ignore entire hierarchy below this point
	if !ipan.Visible() {
		return
	}
	if ipan.Enabled() {
		f(ipan) // Call specified function
	}
	// Check descendants (can assume they are IPanels)
	for _, child := range ipan.Children() {
		traverseIPanel(child.(PanelI), f)
	}
}

// traverseINode traverses the descendants of the specified NodeI,
// executing the specified function for each PanelI.
func traverseINode(inode node.NodeI, f func(ipan PanelI)) {

	if ipan, ok := inode.(PanelI); ok {
		traverseIPanel(ipan, f)
	} else {
		for _, child := range inode.Children() {
			traverseINode(child, f)
		}
	}
}

// forEachIPanel executes the specified function for each enabled and visible PanelI in gm.scene.
func (gm *manager) forEachIPanel(f func(ipan PanelI)) {

	traverseINode(gm.scene, f)
}
