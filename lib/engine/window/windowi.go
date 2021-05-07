// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package window abstracts a platform-specific window.
// Depending on the build tags it can be a GLFW desktop window or a browser WebGlCanvas.
package window

import (
	"fmt"

	"github.com/kasworld/goguelike-single/lib/engine/dispatcheri"
	"github.com/kasworld/goguelike-single/lib/engine/gls"
	"github.com/kasworld/goguelike-single/lib/engine/util/logger"
)

// Package logger
var log = logger.New("WIN", logger.Default)

// WindowI singleton
var win WindowI

// Get returns the WindowI singleton.
func Get() WindowI {
	// Return singleton if already created
	if win != nil {
		return win
	}
	panic(fmt.Errorf("need to call window.Init() first"))
}

// WindowI is the interface for all windows
type WindowI interface {
	dispatcheri.DispatcherI
	Gls() *gls.GLS
	GetFramebufferSize() (width int, height int)
	GetSize() (width int, height int)
	GetScale() (x float64, y float64)
	CreateCursor(imgFile string, xhot, yhot int) (Cursor, error)
	SetCursor(cursor Cursor)
	DisposeAllCustomCursors()
	Destroy()
}
