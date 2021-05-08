// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/kasworld/goguelike-single/lib/engine/dispatcheri"
	"github.com/kasworld/goguelike-single/lib/engine/eventenum"
	"github.com/kasworld/goguelike-single/lib/engine/math32"
)

// Folder represents a folder GUI element.
type Folder struct {
	Panel               // Embedded panel
	label        Label  // Folder label
	icon         Label  // Folder icon
	contentPanel PanelI // Content panel
	styles       *FolderStyles
	cursorOver   bool
	alignRight   bool
}

// FolderStyle contains the styling of a Folder.
type FolderStyle struct {
	PanelStyle
	FgColor math32.Color4
	Icons   [2]string
}

// FolderStyles contains a FolderStyle for each valid GUI state.
type FolderStyles struct {
	Normal   FolderStyle
	Over     FolderStyle
	Focus    FolderStyle
	Disabled FolderStyle
}

// NewFolder creates and returns a pointer to a new folder widget
// with the specified text and initial width.
func NewFolder(text string, width float32, contentPanel PanelI) *Folder {

	f := new(Folder)
	f.Initialize(text, width, contentPanel)
	return f
}

// Initialize initializes the Folder with the specified text and initial width
// It is normally used when the folder is embedded in another object.
func (f *Folder) Initialize(text string, width float32, contentPanel PanelI) {

	f.Panel.Initialize(f, width, 0)
	f.styles = &StyleDefault().Folder

	// Initialize label
	f.label.initialize(text, StyleDefault().Font)
	f.Panel.Add(&f.label)

	// Create icon
	f.icon.initialize("", StyleDefault().FontIcon)
	f.icon.SetFontSize(StyleDefault().Label.PointSize * 1.3)
	f.Panel.Add(&f.icon)

	// Setup content panel
	f.contentPanel = contentPanel
	contentPanel.GetPanel().bounded = false
	contentPanel.GetPanel().zLayerDelta = 1
	contentPanel.GetPanel().SetVisible(false)
	f.Panel.Add(f.contentPanel)

	// Set event callbacks
	f.Panel.Subscribe(eventenum.OnMouseDown, f.onMouse)
	f.Panel.Subscribe(eventenum.OnCursorEnter, f.onCursor)
	f.Panel.Subscribe(eventenum.OnCursorLeave, f.onCursor)

	f.Subscribe(eventenum.OnMouseDownOut, func(s dispatcheri.EventName, i interface{}) {
		// Hide list when clicked out
		if f.contentPanel.Visible() {
			f.contentPanel.SetVisible(false)
		}
	})

	f.contentPanel.Subscribe(eventenum.OnCursorEnter, func(evname dispatcheri.EventName, ev interface{}) {
		f.Dispatch(eventenum.OnCursorLeave, ev)
	})
	f.contentPanel.Subscribe(eventenum.OnCursorLeave, func(evname dispatcheri.EventName, ev interface{}) {
		f.Dispatch(eventenum.OnCursorEnter, ev)
	})

	f.alignRight = true
	f.update()
	f.recalc()
}

// SetStyles set the folder styles overriding the default style.
func (f *Folder) SetStyles(fs *FolderStyles) {

	f.styles = fs
	f.update()
}

// SetAlignRight sets the side of the alignment of the content panel
// in relation to the folder.
func (f *Folder) SetAlignRight(state bool) {

	f.alignRight = state
	f.recalc()
}

// Height returns this folder total height
// considering the contents panel, if visible.
func (f *Folder) Height() float32 {

	height := f.Height()
	if f.contentPanel.GetPanel().Visible() {
		height += f.contentPanel.GetPanel().Height()
	}
	return height
}

// onMouse receives mouse button events over the folder panel.
func (f *Folder) onMouse(evname dispatcheri.EventName, ev interface{}) {

	switch evname {
	case eventenum.OnMouseDown:
		cont := f.contentPanel.GetPanel()
		if !cont.Visible() {
			cont.SetVisible(true)
		} else {
			cont.SetVisible(false)
		}
		f.update()
		f.recalc()
	default:
		return
	}
}

// onCursor receives cursor events over the folder panel
func (f *Folder) onCursor(evname dispatcheri.EventName, ev interface{}) {

	switch evname {
	case eventenum.OnCursorEnter:
		f.cursorOver = true
		f.update()
	case eventenum.OnCursorLeave:
		f.cursorOver = false
		f.update()
	default:
		return
	}
}

// update updates the folder visual state
func (f *Folder) update() {

	if f.cursorOver {
		f.applyStyle(&f.styles.Over)
		return
	}
	f.applyStyle(&f.styles.Normal)
}

// applyStyle applies the specified style
func (f *Folder) applyStyle(s *FolderStyle) {

	f.Panel.ApplyStyle(&s.PanelStyle)

	icode := 0
	if f.contentPanel.GetPanel().Visible() {
		icode = 1
	}
	f.icon.SetText(string(s.Icons[icode]))
	f.icon.SetColor4(&s.FgColor)
	f.label.SetBgColor4(&s.BgColor)
	f.label.SetColor4(&s.FgColor)
}

func (f *Folder) recalc() {

	// icon position
	f.icon.SetPosition(0, 0)

	// Label position and width
	f.label.SetPosition(f.icon.Width()+4, 0)
	f.Panel.SetContentHeight(f.label.Height())

	// Sets position of the base folder scroller panel
	cont := f.contentPanel.GetPanel()
	if f.alignRight {
		cont.SetPosition(0, f.Panel.Height())
	} else {
		dx := cont.Width() - f.Panel.Width()
		cont.SetPosition(-dx, f.Panel.Height())
	}
}
