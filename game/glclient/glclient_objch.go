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

	"github.com/kasworld/goguelike-single/config/leveldata"
	"github.com/kasworld/goguelike-single/enum/flowtype"
	"github.com/kasworld/goguelike-single/enum/turnaction"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/glclient/pid2rspfn"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

func (app *GLClient) reqAIPlay(onoff bool) error {
	return app.sendReqObjWithRspFn(
		&csprotocol.ReqAIPlay{On: onoff},
		func(pk *csprotocol.Packet) error {
			return nil
		},
	)
}

func (app *GLClient) sendReqObjWithRspFn(
	body interface{}, fn pid2rspfn.HandleRspFn) error {
	pid := app.pid2recv.NewPID(fn)
	spk := csprotocol.Packet{
		PacketID: pid,
		FlowType: flowtype.Request,
		Body:     body,
	}
	app.c2tCh <- &spk
	return nil
}

func (app *GLClient) handle_t2ch() {
	for len(app.t2cCh) > 0 {
		rpk := <-app.t2cCh
		g2log.TraceClient("recv %v", rpk)
		switch rpk.FlowType {
		default:
			g2log.Fatal("invalid packet type %v %v", rpk, rpk.Body)
			return
		case flowtype.Response:
			if err := app.pid2recv.HandleRsp(rpk); err != nil {
				g2log.Fatal("%v %v %v", app, rpk, err)
				return
			}
		case flowtype.Notification:
			err := app.handleRecvNotiObj(rpk)
			// process result
			if err != nil {
				g2log.Fatal("%v %v %v", app, rpk, err)
				return
			}
		}
	}
}

func (app *GLClient) handleRecvNotiObj(rpk *csprotocol.Packet) error {
	switch body := rpk.Body.(type) {
	default:
		return fmt.Errorf("invalid packet")

	case *csprotocol.NotiEnterFloor:
		return app.objRecvNotiFn_EnterFloor(body)
	case *csprotocol.NotiLeaveFloor:
		return app.objRecvNotiFn_LeaveFloor(body)
	case *csprotocol.NotiAgeing:
		return app.objRecvNotiFn_Ageing(body)
	case *csprotocol.NotiDeath:
		return app.objRecvNotiFn_Death(body)
	case *csprotocol.NotiReadyToRebirth:
		return app.objRecvNotiFn_ReadyToRebirth(body)
	case *csprotocol.NotiRebirthed:
		return app.objRecvNotiFn_Rebirthed(body)
	case *csprotocol.NotiVPObjList:
		return app.objRecvNotiFn_VPObjList(body)
	case *csprotocol.NotiVPTiles:
		return app.objRecvNotiFn_VPTiles(body)
	case *csprotocol.NotiFloorComplete:
		return app.objRecvNotiFn_FloorComplete(body)
	case *csprotocol.NotiFoundFieldObj:
		return app.objRecvNotiFn_FoundFieldObj(body)
	case *csprotocol.NotiForgetFloor:
		return app.objRecvNotiFn_ForgetFloor(body)
	case *csprotocol.NotiActivateTrap:
		return app.objRecvNotiFn_ActivateTrap(body)

	}
}

func (app *GLClient) objRecvNotiFn_EnterFloor(body *csprotocol.NotiEnterFloor) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FI.Name {
		newFl := NewClientFloor(app.meshMaker, body.FI)
		app.Name2Floor[body.FI.Name] = newFl
		app.CurrentFloor = newFl
	}
	app.scene.Add(app.CurrentFloor.Scene)
	app.CurrentFloor.EnterFloor()
	return nil
}
func (app *GLClient) objRecvNotiFn_LeaveFloor(body *csprotocol.NotiLeaveFloor) error {
	oldFl := app.Name2Floor[body.FI.Name]
	app.scene.Remove(oldFl.Scene)
	// do nothing
	return nil
}

func (app *GLClient) objRecvNotiFn_Ageing(body *csprotocol.NotiAgeing) error {
	// do nothing
	return nil
}
func (app *GLClient) objRecvNotiFn_Death(body *csprotocol.NotiDeath) error {
	// do nothing
	return nil
}
func (app *GLClient) objRecvNotiFn_ReadyToRebirth(body *csprotocol.NotiReadyToRebirth) error {
	app.sendReqObjWithRspFn(
		&csprotocol.ReqTurnAction{
			Act: turnaction.Rebirth,
		},
		func(pk *csprotocol.Packet) error {
			return nil
		})
	return nil
}
func (app *GLClient) objRecvNotiFn_Rebirthed(body *csprotocol.NotiRebirthed) error {
	// do nothing
	return nil
}

func (app *GLClient) objRecvNotiFn_VPObjList(body *csprotocol.NotiVPObjList) error {
	oldOLNotiData := app.OLNotiData
	app.OLNotiData = body
	newOLNotiData := body
	app.onFieldObj = nil

	csprotocol.EquipClientByUUID(body.ActiveObj.EquipBag).Sort()
	csprotocol.PotionClientByUUID(body.ActiveObj.PotionBag).Sort()
	csprotocol.ScrollClientByUUID(body.ActiveObj.ScrollBag).Sort()

	if oldOLNotiData != nil {
		app.HPdiff = newOLNotiData.ActiveObj.HP - oldOLNotiData.ActiveObj.HP
		app.SPdiff = newOLNotiData.ActiveObj.SP - oldOLNotiData.ActiveObj.SP
	}
	newLevel := int(leveldata.CalcLevelFromExp(float64(newOLNotiData.ActiveObj.Exp)))

	app.playerActiveObjClient = nil
	for _, v := range body.ActiveObjList {
		if v.UUID == app.GameInfo.ActiveObjUUID {
			app.playerActiveObjClient = v
			app.moveGLPos()
		}
	}

	app.IsOverLoad = newOLNotiData.ActiveObj.CalcWeight() >= leveldata.WeightLimit(newLevel)

	if app.CurrentFloor.FloorInfo == nil {
		g2log.Error("app.CurrentFloor.FloorInfo not set")
		return nil
	}
	if app.CurrentFloor.FloorInfo.Name != newOLNotiData.FloorName {
		g2log.Error("not current floor objlist data %v %v",
			app.CurrentFloor.FloorInfo.Name, newOLNotiData.FloorName,
		)
		return nil
	}

	for _, v := range newOLNotiData.FieldObjList {
		app.CurrentFloor.AddOrUpdateFieldObj(v)
	}

	playerX, playerY := app.GetPlayerXY()
	if app.playerActiveObjClient != nil && app.CurrentFloor.IsValidPos(playerX, playerY) {
		app.onFieldObj = app.CurrentFloor.GetFieldObjAt(playerX, playerY)
	}
	app.actByControlMode()
	return nil
}
func (app *GLClient) objRecvNotiFn_VPTiles(body *csprotocol.NotiVPTiles) error {
	if app.CurrentFloor.FloorInfo == nil {
		g2log.Warn("OrangeRed app.CurrentFloor.FloorInfo not set")
		return nil
	}
	if app.CurrentFloor.FloorInfo.Name != body.FloorName {
		g2log.Warn("not current floor vptile data %v %v",
			app.CurrentFloor.FloorInfo.Name, body.FloorName,
		)
		return nil
	}

	oldComplete := app.CurrentFloor.Visited.IsComplete()
	if err := app.CurrentFloor.UpdateFromViewportTile(body, app.ViewportXYLenList); err != nil {
		g2log.Warn("%v", err)
		return nil
	}
	if !oldComplete && app.CurrentFloor.Visited.IsComplete() {
		// just completed
	}
	return nil
}
func (app *GLClient) objRecvNotiFn_FloorComplete(body *csprotocol.NotiFloorComplete) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FI.Name {
		// new floor
		app.CurrentFloor = NewClientFloor(app.meshMaker, body.FI)
	}

	oldComplete := app.CurrentFloor.Visited.IsComplete()
	app.CurrentFloor.ReplaceFloorTiles(body.Tiles)
	if !oldComplete && app.CurrentFloor.Visited.IsComplete() {
		// floor complete
	}
	app.CurrentFloor.UpdateFieldObjList(body.FOList)
	return nil
}
func (app *GLClient) objRecvNotiFn_FoundFieldObj(body *csprotocol.NotiFoundFieldObj) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FloorName {
		g2log.Fatal("FoundFieldObj unknonw floor %v", body)
		return fmt.Errorf("FoundFieldObj unknonw floor %v", body)
	}
	app.CurrentFloor.AddOrUpdateFieldObj(body.FieldObj)
	return nil
}
func (app *GLClient) objRecvNotiFn_ForgetFloor(body *csprotocol.NotiForgetFloor) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FloorName {
	} else {
		app.CurrentFloor.Forget()
	}
	return nil
}
func (app *GLClient) objRecvNotiFn_ActivateTrap(body *csprotocol.NotiActivateTrap) error {
	// do nothing
	return nil
}
