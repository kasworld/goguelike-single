// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/kasworld/goguelike-single/lib/engine/appbase/appwindow"
	"github.com/kasworld/goguelike-single/lib/engine/eventtype"
	"github.com/kasworld/goguelike-single/lib/engine/gui/assets/icon"
)

const (
	checkON  = string(icon.CheckBox)
	checkOFF = string(icon.CheckBoxOutlineBlank)
	radioON  = string(icon.RadioButtonChecked)
	radioOFF = string(icon.RadioButtonUnchecked)
)

// CheckRadio is a GUI element that can be either a checkbox or a radio button
type CheckRadio struct {
	Panel             // Embedded panel
	Label      *Label // Text label
	icon       *Label
	styles     *CheckRadioStyles
	check      bool
	group      string // current group name
	cursorOver bool
	state      bool
	codeON     string
	codeOFF    string
	subroot    bool // indicates root subcription
}

// CheckRadioStyle contains the styling of a CheckRadio
type CheckRadioStyle BasicStyle

// CheckRadioStyles contains an CheckRadioStyle for each valid GUI state
type CheckRadioStyles struct {
	Normal   CheckRadioStyle
	Over     CheckRadioStyle
	Focus    CheckRadioStyle
	Disabled CheckRadioStyle
}

// NewCheckBox creates and returns a pointer to a new CheckBox widget
// with the specified text
func NewCheckBox(text string) *CheckRadio {

	return newCheckRadio(true, text)
}

// NewRadioButton creates and returns a pointer to a new RadioButton widget
// with the specified text
func NewRadioButton(text string) *CheckRadio {

	return newCheckRadio(false, text)
}

// newCheckRadio creates and returns a pointer to a new CheckRadio widget
// with the specified type and text
func newCheckRadio(check bool, text string) *CheckRadio {

	cb := new(CheckRadio)
	cb.styles = &StyleDefault().CheckRadio

	// Adapts to specified type: CheckBox or RadioButton
	cb.check = check
	cb.state = false
	if cb.check {
		cb.codeON = checkON
		cb.codeOFF = checkOFF
	} else {
		cb.codeON = radioON
		cb.codeOFF = radioOFF
	}

	// Initialize panel
	cb.Panel.Initialize(cb, 0, 0)

	// Subscribe to events
	cb.Panel.Subscribe(eventtype.OnKeyDown, cb.onKey)
	cb.Panel.Subscribe(eventtype.OnCursorEnter, cb.onCursor)
	cb.Panel.Subscribe(eventtype.OnCursorLeave, cb.onCursor)
	cb.Panel.Subscribe(eventtype.OnMouseDown, cb.onMouse)
	cb.Panel.Subscribe(eventtype.OnEnable, func(evname eventtype.EventType, ev interface{}) { cb.update() })

	// Creates label
	cb.Label = NewLabel(text)
	cb.Label.Subscribe(eventtype.OnResize, func(evname eventtype.EventType, ev interface{}) { cb.recalc() })
	cb.Panel.Add(cb.Label)

	// Creates icon label
	cb.icon = NewIcon(" ")
	cb.Panel.Add(cb.icon)

	cb.recalc()
	cb.update()
	return cb
}

// Value returns the current state of the checkbox
func (cb *CheckRadio) Value() bool {

	return cb.state
}

// SetValue sets the current state of the checkbox
func (cb *CheckRadio) SetValue(state bool) *CheckRadio {

	if state == cb.state {
		return cb
	}
	cb.state = state
	cb.update()
	cb.Dispatch(eventtype.OnChange, nil)
	return cb
}

// Group returns the name of the radio group
func (cb *CheckRadio) Group() string {

	return cb.group
}

// SetGroup sets the name of the radio group
func (cb *CheckRadio) SetGroup(group string) *CheckRadio {

	cb.group = group
	return cb
}

// SetStyles set the button styles overriding the default style
func (cb *CheckRadio) SetStyles(bs *CheckRadioStyles) {

	cb.styles = bs
	cb.update()
}

// toggleState toggles the current state of the checkbox/radiobutton
func (cb *CheckRadio) toggleState() {

	// Subscribes once to the root panel for OnRadioGroup events
	// The root panel is used to dispatch events to all checkradios
	if !cb.subroot {
		Manager().Subscribe(eventtype.OnRadioGroup, func(name eventtype.EventType, ev interface{}) {
			cb.onRadioGroup(ev.(*CheckRadio))
		})
		cb.subroot = true
	}

	if cb.check {
		cb.state = !cb.state
	} else {
		if len(cb.group) == 0 {
			cb.state = !cb.state
		} else {
			if cb.state {
				return
			}
			cb.state = !cb.state
		}
	}
	cb.update()
	cb.Dispatch(eventtype.OnChange, nil)
	if !cb.check && len(cb.group) > 0 {
		Manager().Dispatch(eventtype.OnRadioGroup, cb)
	}
}

// onMouse process OnMouseDown events
func (cb *CheckRadio) onMouse(evname eventtype.EventType, ev interface{}) {

	// Dispatch OnClick for left mouse button down
	if evname == eventtype.OnMouseDown {
		mev := ev.(*appwindow.MouseEvent)
		if mev.Button == appwindow.MouseButtonLeft && cb.Enabled() {
			Manager().SetKeyFocus(cb)
			cb.toggleState()
			cb.Dispatch(eventtype.OnClick, nil)
		}
	}
}

// onCursor process OnCursor* events
func (cb *CheckRadio) onCursor(evname eventtype.EventType, ev interface{}) {

	if evname == eventtype.OnCursorEnter {
		cb.cursorOver = true
	} else {
		cb.cursorOver = false
	}
	cb.update()
}

// onKey receives subscribed key events
func (cb *CheckRadio) onKey(evname eventtype.EventType, ev interface{}) {

	kev := ev.(*appwindow.KeyEvent)
	if evname == eventtype.OnKeyDown && kev.Key == appwindow.KeyEnter {
		cb.toggleState()
		cb.update()
		cb.Dispatch(eventtype.OnClick, nil)
		return
	}
	return
}

// onRadioGroup receives subscribed OnRadioGroup events
func (cb *CheckRadio) onRadioGroup(other *CheckRadio) {

	// If event is for this button, ignore
	if cb == other {
		return
	}
	// If other radio group is not the group of this button, ignore
	if cb.group != other.group {
		return
	}
	// Toggle this button state
	cb.SetValue(!other.Value())
}

// update updates the visual appearance of the checkbox
func (cb *CheckRadio) update() {

	if cb.state {
		cb.icon.SetText(cb.codeON)
	} else {
		cb.icon.SetText(cb.codeOFF)
	}

	if !cb.Enabled() {
		cb.applyStyle(&cb.styles.Disabled)
		return
	}
	if cb.cursorOver {
		cb.applyStyle(&cb.styles.Over)
		return
	}
	cb.applyStyle(&cb.styles.Normal)
}

// setStyle sets the specified checkradio style
func (cb *CheckRadio) applyStyle(s *CheckRadioStyle) {

	cb.Panel.ApplyStyle(&s.PanelStyle)
	cb.icon.SetColor4(&s.FgColor)
	cb.Label.SetColor4(&s.FgColor)
}

// recalc recalculates dimensions and position from inside out
func (cb *CheckRadio) recalc() {

	// Sets icon position
	cb.icon.SetFontSize(cb.Label.FontSize() * 1.3)
	cb.icon.SetPosition(0, 0)

	// Label position
	spacing := float32(4)
	cb.Label.SetPosition(cb.icon.Width()+spacing, 0)

	// Content width
	width := cb.icon.Width() + spacing + cb.Label.Width()
	cb.SetContentSize(width, cb.Label.Height())
}
