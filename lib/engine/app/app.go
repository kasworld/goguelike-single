// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package app implements a cross-platform G3N app.
package app

import "github.com/kasworld/goguelike-single/lib/engine/util/logger"

// Package logger
var log = logger.New("APP", logger.Default)

// Application singleton
var a *Application

// OnExit is the event generated by Application when the user
// tries to close the window (desktop) or the Exit() method is called.
const OnExit = "app.OnExit"
