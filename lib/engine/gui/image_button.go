// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/kasworld/goguelike-single/lib/engine/appbase/appwindow"
	"github.com/kasworld/goguelike-single/lib/engine/eventtype"
	"github.com/kasworld/goguelike-single/lib/engine/texture"
)

// ImageButton represents an image button GUI element
type ImageButton struct {
	*Panel                                             // Embedded Panel
	label       *Label                                 // Label panel
	iconLabel   bool                                   // True if icon
	image       *Image                                 // pointer to button image (may be nil)
	styles      *ImageButtonStyles                     // pointer to current button styles
	mouseOver   bool                                   // true if mouse is over button
	pressed     bool                                   // true if button is pressed
	stateImages [ButtonDisabled + 1]*texture.Texture2D // array of images for each button state
}

// ButtonState specifies a button state.
type ButtonState int

// The possible button states.
const (
	ButtonNormal ButtonState = iota
	ButtonOver
	ButtonPressed
	ButtonDisabled
	// ButtonFocus
)

// ImageButtonStyle contains the styling of an ImageButton.
type ImageButtonStyle BasicStyle

// ImageButtonStyles contains one ImageButtonStyle for each possible ImageButton state.
type ImageButtonStyles struct {
	Normal   ImageButtonStyle
	Over     ImageButtonStyle
	Focus    ImageButtonStyle
	Pressed  ImageButtonStyle
	Disabled ImageButtonStyle
}

// NewImageButton creates and returns a pointer to a new ImageButton widget
// with the specified image.
func NewImageButton(normalImgPath string) (*ImageButton, error) {

	b := new(ImageButton)
	b.styles = &StyleDefault().ImageButton

	tex, err := texture.NewTexture2DFromImage(normalImgPath)
	if err != nil {
		return nil, err
	}
	b.stateImages[ButtonNormal] = tex
	b.image = NewImageFromTex(tex)

	// Initializes the button panel
	b.Panel = NewPanel(0, 0)
	b.Panel.SetContentSize(b.image.Width(), b.image.Height())
	b.Panel.SetBorders(5, 5, 5, 5)
	b.Panel.Add(b.image)

	// Subscribe to panel events
	b.Panel.Subscribe(eventtype.OnKeyDown, b.onKey)
	b.Panel.Subscribe(eventtype.OnKeyUp, b.onKey)
	b.Panel.Subscribe(eventtype.OnMouseUp, b.onMouse)
	b.Panel.Subscribe(eventtype.OnMouseDown, b.onMouse)
	b.Panel.Subscribe(eventtype.OnCursor, b.onCursor)
	b.Panel.Subscribe(eventtype.OnCursorEnter, b.onCursor)
	b.Panel.Subscribe(eventtype.OnCursorLeave, b.onCursor)
	b.Panel.Subscribe(eventtype.OnEnable, func(name eventtype.EventType, ev interface{}) { b.update() })
	b.Panel.Subscribe(eventtype.OnResize, func(name eventtype.EventType, ev interface{}) { b.recalc() })

	b.recalc()
	b.update()
	return b, nil
}

// SetText sets the text of the label
func (b *ImageButton) SetText(text string) {

	if b.iconLabel && b.label != nil {
		b.Panel.Remove(b.label)
		b.label.Dispose()
		b.label = nil
	}

	b.iconLabel = false
	if b.label == nil {
		// Create label
		b.label = NewLabel(text)
		b.Panel.Add(b.label)
	} else {
		b.label.SetText(text)
	}
	b.recalc()
}

// SetIcon sets the icon
func (b *ImageButton) SetIcon(icode string) {

	if b.iconLabel == false && b.label != nil {
		b.Panel.Remove(b.label)
		b.label.Dispose()
		b.label = nil
	}

	b.iconLabel = true
	if b.label == nil {
		// Create icon
		b.label = NewIcon(icode)
		b.Panel.Add(b.label)
	} else {
		b.label.SetText(icode)
	}
	b.recalc()
}

// SetFontSize sets the font size of the label/icon
func (b *ImageButton) SetFontSize(size float64) {

	if b.label != nil {
		b.label.SetFontSize(size)
		b.recalc()
	}
}

// SetImage sets the button left image from the specified filename
// If there is currently a selected icon, it is removed
func (b *ImageButton) SetImage(state ButtonState, imgfile string) error {

	tex, err := texture.NewTexture2DFromImage(imgfile)
	if err != nil {
		return err
	}

	if b.stateImages[state] != nil {
		b.stateImages[state].Dispose()
	}
	b.stateImages[state] = tex
	b.update()

	return nil
}

// Dispose releases resources used by this widget
func (b *ImageButton) Dispose() {
	b.Panel.Dispose()
	for _, tex := range b.stateImages {
		if tex != nil {
			tex.Dispose()
		}
	}
}

// SetStyles set the button styles overriding the default style
func (b *ImageButton) SetStyles(bs *ImageButtonStyles) {

	b.styles = bs
	b.update()
}

// onCursor process subscribed cursor events
func (b *ImageButton) onCursor(evname eventtype.EventType, ev interface{}) {

	switch evname {
	case eventtype.OnCursorEnter:
		b.mouseOver = true
		b.update()
	case eventtype.OnCursorLeave:
		b.pressed = false
		b.mouseOver = false
		b.update()
	}
}

// onMouseEvent process subscribed mouse events
func (b *ImageButton) onMouse(evname eventtype.EventType, ev interface{}) {

	switch evname {
	case eventtype.OnMouseDown:
		Manager().SetKeyFocus(b)
		b.pressed = true
		b.update()
		b.Dispatch(eventtype.OnClick, nil)
	case eventtype.OnMouseUp:
		b.pressed = false
		b.update()
	default:
		return
	}
}

// onKey processes subscribed key events
func (b *ImageButton) onKey(evname eventtype.EventType, ev interface{}) {

	kev := ev.(*appwindow.KeyEvent)
	if evname == eventtype.OnKeyDown && kev.Key == appwindow.KeyEnter {
		b.pressed = true
		b.update()
		b.Dispatch(eventtype.OnClick, nil)
		return
	}
	if evname == eventtype.OnKeyUp && kev.Key == appwindow.KeyEnter {
		b.pressed = false
		b.update()
		return
	}
	return
}

// update updates the button visual state
func (b *ImageButton) update() {

	if !b.Enabled() {
		if b.stateImages[ButtonDisabled] != nil {
			b.image.SetTexture(b.stateImages[ButtonDisabled])
		}
		b.applyStyle(&b.styles.Disabled)
		return
	}
	if b.pressed {
		if b.stateImages[ButtonPressed] != nil {
			b.image.SetTexture(b.stateImages[ButtonPressed])
		}
		b.applyStyle(&b.styles.Pressed)
		return
	}
	if b.mouseOver {
		if b.stateImages[ButtonOver] != nil {
			b.image.SetTexture(b.stateImages[ButtonOver])
		}
		b.applyStyle(&b.styles.Over)
		return
	}
	b.image.SetTexture(b.stateImages[ButtonNormal])
	b.applyStyle(&b.styles.Normal)
}

// applyStyle applies the specified button style
func (b *ImageButton) applyStyle(bs *ImageButtonStyle) {

	b.Panel.ApplyStyle(&bs.PanelStyle)
	if b.label != nil {
		b.label.SetColor4(&bs.FgColor)
	}
}

// recalc recalculates all dimensions and position from inside out
func (b *ImageButton) recalc() {

	// Only need to recal if there's a label preset
	if b.label != nil {
		width := b.Panel.ContentWidth()
		height := b.Panel.ContentHeight()

		x := (width - b.label.Width()) / 2
		y := (height - b.label.Height()) / 2

		b.label.SetPosition(x, y)
	}
}
