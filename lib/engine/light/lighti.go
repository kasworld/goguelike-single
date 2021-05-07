// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/kasworld/goguelike-single/lib/engine/g3ncore"
	"github.com/kasworld/goguelike-single/lib/engine/gls"
)

// LightI is the interface that must be implemented for all light types.
type LightI interface {
	RenderSetup(gs *gls.GLS, rinfo *g3ncore.RenderInfo, idx int)
}
