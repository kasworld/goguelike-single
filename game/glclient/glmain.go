// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package glclient

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

// runtime.LockOSThread
// must run in same thread

func (ga *GLClient) glInit() error {
	// Create application and scene
	ga.app = app.App()
	ga.app.IWindow.(*window.GlfwWindow).SetTitle("gogguelike-single")
	ga.app.IWindow.(*window.GlfwWindow).SetSize(1920, 1080)
	ga.scene = core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(ga.scene)

	// Create perspective camera
	ga.cam = camera.New(1)
	ga.cam.SetPosition(0, 0, 3)
	ga.scene.Add(ga.cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(ga.cam)

	ga.app.Subscribe(window.OnWindowSize, ga.onResize)
	ga.onResize("", nil)

	// Create and add lights to the scene
	ga.scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	ga.scene.Add(pointLight)

	// Create and add an axis helper to the scene
	ga.scene.Add(helper.NewAxes(0.5))
	return nil
}

func (ga *GLClient) Run() error {
	if err := ga.glInit(); err != nil {
		return err
	}

	// Create a blue torus and add it to the scene
	geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	ga.scene.Add(mesh)

	// Create and add a button to the scene
	btn := gui.NewButton("Make Red")
	btn.SetPosition(100, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		mat.SetColor(math32.NewColor("DarkRed"))
	})
	ga.scene.Add(btn)

	// Set background color to gray
	ga.app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// login and start
	if err := ga.reqLogin(); err != nil {
		return err
	}

	ga.app.Run(ga.updateGL)
	return nil
}

func (ga *GLClient) updateGL(renderer *renderer.Renderer, deltaTime time.Duration) {
	ga.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
	renderer.Render(ga.scene, ga.cam)
	ga.handle_t2ch()
}

// Set up callback to update viewport and camera aspect ratio when the window is resized
func (ga *GLClient) onResize(evname string, ev interface{}) {
	// Get framebuffer size and update viewport accordingly
	width, height := ga.app.GetSize()
	ga.app.Gls().Viewport(0, 0, int32(width), int32(height))
	// Update the camera's aspect ratio
	ga.cam.SetAspect(float32(width) / float32(height))
}
