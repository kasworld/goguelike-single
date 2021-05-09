// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// RectBounds specifies the size of the boundaries of a rectangle.
// It can represent the thickness of the borders, the margins, or the padding of a rectangle.
type RectBounds struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// Set sets the values of the RectBounds.
func (bs *RectBounds) Set(top, right, bottom, left float32) {

	if top >= 0 {
		bs.Top = top
	}
	if right >= 0 {
		bs.Right = right
	}
	if bottom >= 0 {
		bs.Bottom = bottom
	}
	if left >= 0 {
		bs.Left = left
	}
}
