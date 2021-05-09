// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/kasworld/goguelike-single/lib/engine/appbase/appwindow"
	"github.com/kasworld/goguelike-single/lib/engine/eventtype"
)

/***************************************

 Button Panel
 +-------------------------------+
 |  Image/Icon      Label        |
 |  +----------+   +----------+  |
 |  |          |   |          |  |
 |  |          |   |          |  |
 |  +----------+   +----------+  |
 +-------------------------------+

****************************************/

// Button represents a button GUI element
type Button struct {
	Panel                   // Embedded Panel
	Label     *Label        // Label panel
	image     *Image        // pointer to button image (may be nil)
	icon      *Label        // pointer to button icon (may be nil
	styles    *ButtonStyles // pointer to current button styles
	mouseOver bool          // true if mouse is over button
	pressed   bool          // true if button is pressed
}

// ButtonStyle contains the styling of a Button
type ButtonStyle BasicStyle

// ButtonStyles contains one ButtonStyle for each possible button state
type ButtonStyles struct {
	Normal   ButtonStyle
	Over     ButtonStyle
	Focus    ButtonStyle
	Pressed  ButtonStyle
	Disabled ButtonStyle
}

// NewButton creates and returns a pointer to a new button widget
// with the specified text for the button label.
func NewButton(text string) *Button {

	b := new(Button)
	b.styles = &StyleDefault().Button

	// Initializes the button panel
	b.Panel.Initialize(b, 0, 0)

	// Subscribe to panel events
	b.Subscribe(eventtype.OnKeyDown, b.onKey)
	b.Subscribe(eventtype.OnKeyUp, b.onKey)
	b.Subscribe(eventtype.OnMouseUp, b.onMouse)
	b.Subscribe(eventtype.OnMouseDown, b.onMouse)
	b.Subscribe(eventtype.OnMouseUpOut, b.onMouse)
	b.Subscribe(eventtype.OnCursor, b.onCursor)
	b.Subscribe(eventtype.OnCursorEnter, b.onCursor)
	b.Subscribe(eventtype.OnCursorLeave, b.onCursor)
	b.Subscribe(eventtype.OnEnable, func(name eventtype.EventType, ev interface{}) { b.update() })
	b.Subscribe(eventtype.OnResize, func(name eventtype.EventType, ev interface{}) { b.recalc() })

	// Creates label
	b.Label = NewLabel(text)
	b.Label.Subscribe(eventtype.OnResize, func(name eventtype.EventType, ev interface{}) { b.recalc() })
	b.Panel.Add(b.Label)

	b.recalc() // recalc first then update!
	b.update()
	return b
}

// SetIcon sets the button icon from the default Icon font.
// If there is currently a selected image, it is removed
func (b *Button) SetIcon(icode string) {

	ico := NewIcon(icode)
	if b.image != nil {
		b.Panel.Remove(b.image)
		b.image = nil
	}
	if b.icon != nil {
		b.Panel.Remove(b.icon)
	}
	b.icon = ico
	b.icon.SetFontSize(b.Label.FontSize() * 1.4)
	b.Panel.Add(b.icon)

	b.recalc()
	b.update()
}

// SetImage sets the button left image from the specified filename
// If there is currently a selected icon, it is removed
func (b *Button) SetImage(imgfile string) error {

	img, err := NewImage(imgfile)
	if err != nil {
		return err
	}
	if b.image != nil {
		b.Panel.Remove(b.image)
	}
	b.image = img
	b.Panel.Add(b.image)
	b.recalc()
	return nil
}

// SetStyles set the button styles overriding the default style
func (b *Button) SetStyles(bs *ButtonStyles) {

	b.styles = bs
	b.update()
}

// onCursor process subscribed cursor events
func (b *Button) onCursor(evname eventtype.EventType, ev interface{}) {

	switch evname {
	case eventtype.OnCursorEnter:
		b.mouseOver = true
		b.update()
	case eventtype.OnCursorLeave:
		b.mouseOver = false
		b.update()
	}
}

// onMouseEvent process subscribed mouse events
func (b *Button) onMouse(evname eventtype.EventType, ev interface{}) {

	if !b.Enabled() {
		return
	}

	switch evname {
	case eventtype.OnMouseDown:
		Manager().SetKeyFocus(b)
		b.pressed = true
		b.update()
	case eventtype.OnMouseUpOut:
		fallthrough
	case eventtype.OnMouseUp:
		if b.pressed && b.mouseOver {
			b.Dispatch(eventtype.OnClick, nil)
		}
		b.pressed = false
		b.update()
	default:
		return
	}
}

// onKey processes subscribed key events
func (b *Button) onKey(evname eventtype.EventType, ev interface{}) {

	kev := ev.(*appwindow.KeyEvent)
	if kev.Key != appwindow.KeyEnter {
		return
	}
	switch evname {
	case eventtype.OnKeyDown:
		b.pressed = true
		b.update()
		b.Dispatch(eventtype.OnClick, nil)
	case eventtype.OnKeyUp:
		b.pressed = false
		b.update()
	}
}

// update updates the button visual state
func (b *Button) update() {

	if !b.Enabled() {
		b.applyStyle(&b.styles.Disabled)
		return
	}
	if b.pressed && b.mouseOver {
		b.applyStyle(&b.styles.Pressed)
		return
	}
	if b.mouseOver {
		b.applyStyle(&b.styles.Over)
		return
	}
	b.applyStyle(&b.styles.Normal)
}

// applyStyle applies the specified button style
func (b *Button) applyStyle(bs *ButtonStyle) {

	b.Panel.ApplyStyle(&bs.PanelStyle)
	if b.icon != nil {
		b.icon.SetColor4(&bs.FgColor)
	}
	b.Label.SetColor4(&bs.FgColor)
}

// recalc recalculates all dimensions and position from inside out
func (b *Button) recalc() {

	// Current width and height of button content area
	width := b.Panel.ContentWidth()
	height := b.Panel.ContentHeight()

	// Image or icon width
	imgWidth := float32(0)
	spacing := float32(4)
	if b.image != nil {
		imgWidth = b.image.Width()
	} else if b.icon != nil {
		imgWidth = b.icon.Width()
	}
	if imgWidth == 0 {
		spacing = 0
	}

	// If the label is empty and an icon of image was defined ignore the label widthh
	// to centralize the icon/image in the button
	labelWidth := spacing + b.Label.Width()
	if b.Label.Text() == "" && imgWidth > 0 {
		labelWidth = 0
	}

	// Sets new content width and height if necessary
	minWidth := imgWidth + labelWidth
	minHeight := b.Label.Height()
	resize := false
	if width < minWidth {
		width = minWidth
		resize = true
	}
	if height < minHeight {
		height = minHeight
		resize = true
	}
	if resize {
		b.SetContentSize(width, height)
	}

	// Centralize horizontally
	px := (width - minWidth) / 2

	// Set label position
	ly := (height - b.Label.Height()) / 2
	b.Label.SetPosition(px+imgWidth+spacing, ly)

	// Image/icon position
	if b.image != nil {
		iy := (height - b.image.height) / 2
		b.image.SetPosition(px, iy)
	} else if b.icon != nil {
		b.icon.SetPosition(px, ly)
	}
}
