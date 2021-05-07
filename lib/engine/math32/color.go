// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import (
	"strings"
)

// Color describes an RGB color
type Color struct {
	R float32
	G float32
	B float32
}

// NewColor creates and returns a pointer to a new Color
// with the specified web standard color name (case insensitive).
// Returns nil if the color name not found
func NewColor(name string) *Color {

	c, ok := mapColorNames[strings.ToLower(name)]
	if !ok {
		return nil
	}
	return &c
}

// ColorName returns a Color with the specified standard web color name (case insensitive).
// Returns black color if the specified color name not found
func ColorName(name string) Color {

	return mapColorNames[strings.ToLower(name)]
}

// NewColorHex creates and returns a pointer to a new color
// with its RGB components from the specified hex value
func NewColorHex(color uint) *Color {

	return (&Color{}).SetHex(color)
}

// Set sets this color individual R,G,B components
func (c *Color) Set(r, g, b float32) *Color {

	c.R = r
	c.G = g
	c.B = b
	return c
}

// SetHex sets the color RGB components from the
// specified integer interpreted as a color hex number
func (c *Color) SetHex(value uint) *Color {

	c.R = float32((value >> 16 & 255)) / 255
	c.G = float32((value >> 8 & 255)) / 255
	c.B = float32((value & 255)) / 255
	return c
}

// SetName sets the color RGB components from the
// specified standard web color name
func (c *Color) SetName(name string) *Color {

	color, ok := mapColorNames[strings.ToLower(name)]
	if ok {
		*c = color
	}
	return c
}

// Add adds to each RGB component of this color the correspondent component of other color
// Returns pointer to this updated color
func (c *Color) Add(other *Color) *Color {

	c.R += other.R
	c.G += other.G
	c.B += other.B
	return c
}

// AddColors adds to each RGB component of this color the correspondent component of color1 and color2
// Returns pointer to this updated color
func (c *Color) AddColors(color1, color2 *Color) *Color {

	c.R = color1.R + color2.R
	c.G = color1.G + color2.G
	c.B = color1.B + color2.B
	return c
}

// AddScalar adds the specified scalar value to each RGB component of this color
// Returns pointer to this updated color
func (c *Color) AddScalar(s float32) *Color {

	c.R += s
	c.G += s
	c.B += s
	return c
}

// Multiply multiplies each RGB component of this color by other
// Returns pointer to this updated color
func (c *Color) Multiply(other *Color) *Color {

	c.R *= other.R
	c.G *= other.G
	c.B *= other.B
	return c
}

// MultiplyScalar multiplies each RGB component of this color by the specified scalar.
// Returns pointer to this updated color
func (c *Color) MultiplyScalar(v float32) *Color {

	c.R *= v
	c.G *= v
	c.B *= v
	return c
}

// Lerp linear sets this color as the linear interpolation of itself
// with the specified color for the specified alpha.
// Returns pointer to this updated color
func (c *Color) Lerp(color *Color, alpha float32) *Color {

	c.R += (color.R - c.R) * alpha
	c.G += (color.G - c.G) * alpha
	c.B += (color.B - c.B) * alpha
	return c
}

// Equals returns if this color is equal to other
func (c *Color) Equals(other *Color) bool {

	return (c.R == other.R) && (c.G == other.G) && (c.B == other.B)
}

// IsColorName returns if the specified name is valid color name
func IsColorName(name string) (Color, bool) {

	c, ok := mapColorNames[strings.ToLower(name)]
	return c, ok
}
