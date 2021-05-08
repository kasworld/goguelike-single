// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shape

// Shape is a collision shape.
// It can be an analytical geometry such as a sphere, plane, etc.. or it can be defined by a polygonal Geometry.
type Shape struct {

	// TODO
	//material

	// Collision filtering
	colFilterGroup int
	colFilterMask  int
}

func (s *Shape) initialize() {

	// Collision filtering
	s.colFilterGroup = 1
	s.colFilterMask = -1
}
