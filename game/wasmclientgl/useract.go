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

package wasmclientgl

import (
	"syscall/js"

	"github.com/kasworld/goguelike-single/enum/clientcontroltype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

func (app *WasmClient) actViewportView() {
	if gameOptions.GetByIDBase("ViewMode").State != 1 {
		// if not viewportview
		return
	}
	dx, dy := app.KeyboardPressedMap.SumMoveDxDy(Key2Dir)
	app.KeyDir = way9type.RemoteDxDy2Way9(dx, dy)

	if app.KeyDir != way9type.Center {
		app.ClientColtrolMode = clientcontroltype.Keyboard
		app.Path2dst = nil
		app.vp.ClearMovePath()
		dir := app.KeyDir
		app.floorVPPosX += dir.Dx()
		app.floorVPPosY += dir.Dy()
		return
	}

}

// from noti obj list
func (app *WasmClient) actPlayView() {
	if gameOptions.GetByIDBase("ViewMode").State != 0 {
		// if not playview
		return
	}

	if app.olNotiData.ActiveObj.AP > 0 {
		if app.olNotiData == nil || app.olNotiData.ActiveObj.HP <= 0 {
			return
		}
		dx, dy := app.KeyboardPressedMap.SumMoveDxDy(Key2Dir)
		app.KeyDir = way9type.RemoteDxDy2Way9(dx, dy)
		if app.KeyDir != way9type.Center {
			app.ClientColtrolMode = clientcontroltype.Keyboard
			app.Path2dst = nil
			app.vp.ClearMovePath()
			autoPlayButton := autoActs.GetByIDBase("AutoPlay")
			if autoPlayButton.State == 0 {
				autoPlayButton.JSFn(js.Null(), nil)
			}
		}
		if app.moveByUserInput() {
			return
		}
		if fo := app.onFieldObj; fo == nil ||
			(fo != nil && fieldobjacttype.ClientData[fo.ActType].ActOn) {
			for i, v := range autoActs.ButtonList {
				if v.State == 0 {
					if tryAutoActFn[i](app, v) {
						return
					}
				}
			}
		}
	}
}

func (app *WasmClient) moveByUserInput() bool {
	cf := app.CurrentFloor
	playerX, playerY := app.GetPlayerXY()
	if !cf.IsValidPos(playerX, playerY) {
		jslog.Errorf("ao out of floor %v %v", app.olNotiData.ActiveObj, cf)
		return false
	}
	w, h := cf.Tiles.GetXYLen()
	switch app.ClientColtrolMode {
	default:
		jslog.Errorf("invalid ClientColtrolMode %v", app.ClientColtrolMode)
	case clientcontroltype.Keyboard:
		if app.sendMovePacketByInput(app.KeyDir) {
			return true
		}
	case clientcontroltype.FollowMouse:
		if app.sendMovePacketByInput(app.MouseDir) {
			return true
		}
	case clientcontroltype.MoveToDest:
		playerPos := [2]int{playerX, playerY}
		if app.Path2dst == nil || len(app.Path2dst) == 0 {
			app.ClientColtrolMode = clientcontroltype.Keyboard
			return false
		}
		for i := len(app.Path2dst) - 1; i >= 0; i-- {
			nextPos := app.Path2dst[i]
			isContact, dir := way9type.CalcContactDirWrapped(playerPos, nextPos, w, h)
			if isContact {
				if dir == way9type.Center {
					// arrived
					app.Path2dst = nil
					app.vp.ClearMovePath()
					app.ClientColtrolMode = clientcontroltype.Keyboard
					return false
				} else {
					go app.sendPacket(c2t_idcmd.Move,
						&c2t_obj.ReqMove_data{Dir: dir},
					)
					return true
				}
			}
		}
		// fail to path2dest
		app.Path2dst = nil
		app.vp.ClearMovePath()
		app.ClientColtrolMode = clientcontroltype.Keyboard
	}
	return false
}
