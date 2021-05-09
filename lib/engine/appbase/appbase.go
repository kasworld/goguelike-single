// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !wasm

package appbase

import (
	"fmt"
	"time"

	"github.com/kasworld/goguelike-single/lib/engine/appbase/appwindow"
	"github.com/kasworld/goguelike-single/lib/engine/audio/al"
	"github.com/kasworld/goguelike-single/lib/engine/audio/vorbis"
	"github.com/kasworld/goguelike-single/lib/engine/eventtype"
	"github.com/kasworld/goguelike-single/lib/engine/log"
	"github.com/kasworld/goguelike-single/lib/engine/renderer"
)

// AppBase
type AppBase struct {
	appwindow.AppWindowI                     // Embedded GlfwWindow
	keyState             *appwindow.KeyState // Keep track of keyboard state
	renderer             *renderer.Renderer  // Renderer object
	audioDev             *al.Device          // Default audio device
	startTime            time.Time           // AppBase start time
	frameStart           time.Time           // Frame start time
	frameDelta           time.Duration       // Duration of last frame
}

// New returns the AppBase object
func New(title string, width, height int) *AppBase {
	appBase := new(AppBase)
	// Initialize appwindow
	err := appwindow.Init(width, height, title)
	if err != nil {
		panic(err)
	}
	appBase.AppWindowI = appwindow.Get()
	appBase.openDefaultAudioDevice()                  // Set up audio
	appBase.keyState = appwindow.NewKeyState(appBase) // Create KeyState
	// Create renderer and add default shaders
	appBase.renderer = renderer.NewRenderer(appBase.Gls())
	err = appBase.renderer.AddDefaultShaders()
	if err != nil {
		panic(fmt.Errorf("AddDefaultShaders:%v", err))
	}
	return appBase
}

// Run starts the update loop.
// It calls the user-provided update function every frame.
func (a *AppBase) Run(update func(rend *renderer.Renderer, deltaTime time.Duration)) {

	// Initialize start and frame time
	a.startTime = time.Now()
	a.frameStart = time.Now()

	// Set up recurring calls to user's update function
	for true {
		// If Exit() was called or there was an attempt to close the window dispatch OnExit event for subscribers.
		// If no subscriber cancelled the event, terminate the application.
		if a.AppWindowI.(*appwindow.GlfwWindow).ShouldClose() {
			a.Dispatch(eventtype.OnExit, nil)
			// TODO allow for cancelling exit e.g. showing dialog asking the user if he/she wants to save changes
			// if exit was cancelled {
			//     a.AppWindowI.(*appwindow.GlfwWindow).SetShouldClose(false)
			// } else {
			break
			// }
		}
		// Update frame start and frame delta
		now := time.Now()
		a.frameDelta = now.Sub(a.frameStart)
		a.frameStart = now
		// Call user's update function
		update(a.renderer, a.frameDelta)
		// Swap buffers and poll events
		a.AppWindowI.(*appwindow.GlfwWindow).SwapBuffers()
		a.AppWindowI.(*appwindow.GlfwWindow).PollEvents()
	}

	// Close default audio device
	if a.audioDev != nil {
		al.CloseDevice(a.audioDev)
	}
	// Destroy appwindow
	a.Destroy()
}

// Exit requests to terminate the application
// AppBase will dispatch OnQuit events to registered subscribers which
// can cancel the process by calling CancelDispatch().
func (a *AppBase) Exit() {

	a.AppWindowI.(*appwindow.GlfwWindow).SetShouldClose(true)
}

// Renderer returns the application's renderer.
func (a *AppBase) Renderer() *renderer.Renderer {

	return a.renderer
}

// KeyState returns the application's KeyState.
func (a *AppBase) KeyState() *appwindow.KeyState {

	return a.keyState
}

// RunTime returns the elapsed duration since the call to Run().
func (a *AppBase) RunTime() time.Duration {

	return time.Now().Sub(a.startTime)
}

// openDefaultAudioDevice opens the default audio device setting it to the current context
func (a *AppBase) openDefaultAudioDevice() error {

	// Opens default audio device
	var err error
	a.audioDev, err = al.OpenDevice("")
	if err != nil {
		return fmt.Errorf("opening OpenAL default device: %s", err)
	}
	// Check for OpenAL effects extension support
	var attribs []int
	if al.IsExtensionPresent("ALC_EXT_EFX") {
		attribs = []int{al.MAX_AUXILIARY_SENDS, 4}
	}
	// Create audio context
	acx, err := al.CreateContext(a.audioDev, attribs)
	if err != nil {
		return fmt.Errorf("creating OpenAL context: %s", err)
	}
	// Makes the context the current one
	err = al.MakeContextCurrent(acx)
	if err != nil {
		return fmt.Errorf("setting OpenAL context current: %s", err)
	}
	// Logs audio library versions
	log.Info("%s version: %s", al.GetString(al.Vendor), al.GetString(al.Version))
	log.Info("%s", vorbis.VersionString())
	return nil
}
