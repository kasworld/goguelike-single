// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Rect represents a rectangle.
type Rect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// Contains determines whether a 2D point is inside the Rect.
func (r *Rect) Contains(x, y float32) bool {

	if x < r.X || x > r.X+r.Width {
		return false
	}
	if y < r.Y || y > r.Y+r.Height {
		return false
	}
	return true
}
