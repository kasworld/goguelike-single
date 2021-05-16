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
	"fmt"
	"time"

	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/h4o/appbase"
	"github.com/kasworld/h4o/appbase/appwindow"
	"github.com/kasworld/h4o/camera"
	"github.com/kasworld/h4o/eventtype"
	"github.com/kasworld/h4o/gls"
	"github.com/kasworld/h4o/graphic"
	"github.com/kasworld/h4o/gui"
	"github.com/kasworld/h4o/light"
	"github.com/kasworld/h4o/math32"
	"github.com/kasworld/h4o/node"
	"github.com/kasworld/h4o/renderer"
	"github.com/kasworld/h4o/util/framerater"
	"github.com/kasworld/h4o/util/helper"
)

// runtime.LockOSThread
// must run in same thread

func (ga *GLClient) glInit() error {
	// Create application and scene
	ga.app = appbase.New("goguelike-single", 1920, 1080)
	ga.scene = node.NewNode()

	ga.sceneAO = node.NewNode()
	ga.scene.Add(ga.sceneAO)
	ga.sceneCO = node.NewNode()
	ga.scene.Add(ga.sceneCO)
	ga.sceneDO = node.NewNode()
	ga.scene.Add(ga.sceneDO)

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(ga.scene)

	// Create perspective camera
	ga.cam = camera.New(1)
	ga.cam.SetFar(1400)
	ga.camZpos = 100
	ga.cam.SetPosition(0, 0, ga.camZpos)
	ga.scene.Add(ga.cam)

	// Set up orbit control for the camera
	// camera.NewOrbitControl(ga.cam)

	ga.app.Subscribe(eventtype.OnWindowSize, ga.onResize)
	ga.onResize(eventtype.OnResize, nil)

	// Create and add lights to the scene
	ga.scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	ga.pLight = light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	ga.pLight.SetPosition(1, 0, 2)
	ga.scene.Add(ga.pLight)

	// Create and add an axis helper to the scene
	ga.scene.Add(helper.NewAxes(100))

	ga.frameRater = framerater.NewFrameRater(60)
	ga.labelFPS = gui.NewLabel(" ")
	ga.labelFPS.SetFontSize(20)
	ga.labelFPS.SetLayoutParams(&gui.HBoxLayoutParams{AlignV: gui.AlignCenter})
	lightTextColor := math32.Color4{0.8, 0.8, 0.8, 1}
	ga.labelFPS.SetColor4(&lightTextColor)
	ga.scene.Add(ga.labelFPS)

	gui.Manager().SubscribeID(eventtype.OnMouseUp, ga, ga.onMouse)
	gui.Manager().SubscribeID(eventtype.OnMouseDown, ga, ga.onMouse)
	gui.Manager().SubscribeID(eventtype.OnScroll, &ga, ga.onScroll)

	return nil
}

// onMouse is called when an OnMouseDown/OnMouseUp event is received.
func (ga *GLClient) onMouse(evname eventtype.EventType, ev interface{}) {

	switch evname {
	case eventtype.OnMouseDown:
		// gui.Manager().SetCursorFocus(ga)
		mev := ev.(*appwindow.MouseEvent)
		switch mev.Button {
		case appwindow.MouseButtonLeft: // Rotate
		case appwindow.MouseButtonMiddle: // Zoom
		case appwindow.MouseButtonRight: // Pan
		}
	case eventtype.OnMouseUp:
		// gui.Manager().SetCursorFocus(nil)
	}
}

// onScroll is called when an OnScroll event is received.
func (ga *GLClient) onScroll(evname eventtype.EventType, ev interface{}) {
	zF := float32(1.5)
	sev := ev.(*appwindow.ScrollEvent)
	if sev.Yoffset > 0 {
		ga.camZpos *= zF
		if ga.camZpos > 1000 {
			ga.camZpos = 1000
		}
	} else if sev.Yoffset < 0 {
		ga.camZpos /= zF
		if ga.camZpos < 10 {
			ga.camZpos = 10
		}
	}
	ga.moveGLPos()
}

func (ga *GLClient) Run() error {
	if err := ga.glInit(); err != nil {
		return err
	}

	// Create and add a button to the scene
	// btn := gui.NewButton("Make Red")
	// btn.SetPosition(100, 40)
	// btn.SetSize(40, 40)
	// btn.Subscribe(eventtype.gui.OnClick, func(name string, ev interface{}) {
	// 	mat.SetColor(math32.NewColor("DarkRed"))
	// })
	// ga.scene.Add(btn)

	// Set background color to gray
	ga.app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	ga.reqAIPlay(true)

	ga.app.Run(ga.updateGL)
	return nil
}

func (ga *GLClient) updateGL(renderer *renderer.Renderer, deltaTime time.Duration) {
	// Start measuring this frame
	ga.frameRater.Start()

	ga.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
	renderer.Render(ga.scene, ga.cam)
	ga.handle_t2ch()

	// Control and update FPS
	ga.frameRater.Wait()
	ga.updateFPS()

}

func (ga *GLClient) moveGLPos() {
	aox, aoy := ga.GetPlayerXY()
	ga.cam.SetPosition(float32(aox), float32(aoy), ga.camZpos)
	ga.pLight.SetPosition(float32(aox), float32(aoy), ga.camZpos)
}

// Set up callback to update viewport and camera aspect ratio when the window is resized
func (ga *GLClient) onResize(evname eventtype.EventType, ev interface{}) {
	// Get framebuffer size and update viewport accordingly
	width, height := ga.app.GetSize()
	ga.app.Gls().Viewport(0, 0, int32(width), int32(height))
	// Update the camera's aspect ratio
	ga.cam.SetAspect(float32(width) / float32(height))
}

// UpdateFPS updates the fps value in the window title or header label
func (ga *GLClient) updateFPS() {

	// Get the FPS and potential FPS from the frameRater
	fps, pfps, ok := ga.frameRater.FPS(time.Duration(60) * time.Millisecond)
	if !ok {
		return
	}

	// Show the FPS in the header label
	ga.labelFPS.SetText(fmt.Sprintf("%3.1f / %3.1f", fps, pfps))
}

func (ga *GLClient) updateVPObjList(body *csprotocol.NotiVPObjList) {
	// update active object
	for _, v := range ga.sceneAO.Children() {
		aoMesh := v.(*graphic.Mesh)
		ga.meshMaker.PutActiveObj(aoMesh)
	}
	ga.sceneAO.RemoveAll(true)
	for _, v := range body.ActiveObjList {
		aoMesh := ga.meshMaker.GetActiveObj(v.Faction, v.X, v.Y)
		ga.sceneAO.Add(aoMesh)
	}

	// update danger object
	for _, v := range ga.sceneDO.Children() {
		doMesh := v.(*graphic.Mesh)
		ga.meshMaker.PutDangerObj(doMesh)
	}
	ga.sceneDO.RemoveAll(true)
	for _, v := range body.DangerObjList {
		doMesh := ga.meshMaker.GetDangerObj(v.DangerType, v.X, v.Y)
		ga.sceneDO.Add(doMesh)
	}

	// update carry object
	for _, v := range ga.sceneCO.Children() {
		coMesh := v.(*graphic.Mesh)
		ga.meshMaker.PutCarryObj(coMesh)
	}
	ga.sceneCO.RemoveAll(true)
	for _, v := range body.CarryObjList {
		cokey := NewCOKeyFromCarryObjClientOnFloor(v)
		coMesh := ga.meshMaker.GetCarryObj(cokey, v.X, v.Y)
		ga.sceneCO.Add(coMesh)
	}

}
