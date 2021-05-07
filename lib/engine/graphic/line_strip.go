// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/kasworld/goguelike-single/lib/engine/geometry"
	"github.com/kasworld/goguelike-single/lib/engine/gls"
	"github.com/kasworld/goguelike-single/lib/engine/material"
	"github.com/kasworld/goguelike-single/lib/engine/renderinfo"
)

// LineStrip is a Graphic which is rendered as a collection of connected lines.
type LineStrip struct {
	Graphic             // Embedded graphic object
	uniMVPm gls.Uniform // Model view projection matrix uniform location cache
}

// NewLineStrip creates and returns a pointer to a new LineStrip graphic
// with the specified geometry and material.
func NewLineStrip(igeom geometry.GeometryI, imat material.MaterialI) *LineStrip {

	l := new(LineStrip)
	l.Graphic.Init(l, igeom, gls.LINE_STRIP)
	l.AddMaterial(l, imat, 0, 0)
	l.uniMVPm.Init("MVP")
	return l
}

// RenderSetup is called by the engine before drawing this geometry.
func (l *LineStrip) RenderSetup(gs *gls.GLS, rinfo *renderinfo.RenderInfo) {

	// Transfer model view projection matrix uniform
	mvpm := l.ModelViewProjectionMatrix()
	location := l.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}
