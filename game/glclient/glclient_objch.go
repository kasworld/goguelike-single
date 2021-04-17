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
	"github.com/kasworld/goguelike-single/game/clientfloor"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_pid2rspfn"
)

func (app *GLClient) reqAIPlay(onoff bool) error {
	return app.sendReqObjWithRspFn(
		c2t_idcmd.AIPlay,
		&c2t_obj.ReqAIPlay_data{On: onoff},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (app *GLClient) sendReqObjWithRspFn(cmd c2t_idcmd.CommandID, body interface{},
	fn c2t_pid2rspfn.HandleRspFn) error {
	pid := app.pid2recv.NewPID(fn)
	spk := c2t_packet.Packet{
		Header: c2t_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: c2t_packet.Request,
		},
		Body: body,
	}
	app.c2tCh <- &spk
	return nil
}

func (app *GLClient) handle_t2ch() {
	rpk, ok := <-app.t2cCh
	if !ok {
		return
	}
	g2log.TraceClient("recv %v", rpk.Header)
	switch rpk.Header.FlowType {
	default:
		g2log.Fatal("invalid packet type %v %v", rpk.Header, rpk.Body)
		return
	case c2t_packet.Response:
		if err := app.pid2recv.HandleRsp(rpk.Header, rpk.Body); err != nil {
			g2log.Fatal("%v %v %v %v", app, rpk.Header, rpk.Body, err)
			return
		}
	case c2t_packet.Notification:
		err := app.handleRecvNotiObj(rpk)
		// process result
		if err != nil {
			g2log.Fatal("%v %v %v %v", app, rpk.Header, rpk.Body, err)
			return
		}
	}
}

func (app *GLClient) handleRecvNotiObj(rpk *c2t_packet.Packet) error {
	switch body := rpk.Body.(type) {
	default:
		return fmt.Errorf("invalid packet")

	case *c2t_obj.NotiEnterFloor_data:
		return app.objRecvNotiFn_EnterFloor(rpk.Header, body)
	case *c2t_obj.NotiLeaveFloor_data:
		return app.objRecvNotiFn_LeaveFloor(rpk.Header, body)
	case *c2t_obj.NotiAgeing_data:
		return app.objRecvNotiFn_Ageing(rpk.Header, body)
	case *c2t_obj.NotiDeath_data:
		return app.objRecvNotiFn_Death(rpk.Header, body)
	case *c2t_obj.NotiReadyToRebirth_data:
		return app.objRecvNotiFn_ReadyToRebirth(rpk.Header, body)
	case *c2t_obj.NotiRebirthed_data:
		return app.objRecvNotiFn_Rebirthed(rpk.Header, body)
	case *c2t_obj.NotiVPObjList_data:
		return app.objRecvNotiFn_VPObjList(rpk.Header, body)
	case *c2t_obj.NotiVPTiles_data:
		return app.objRecvNotiFn_VPTiles(rpk.Header, body)
	case *c2t_obj.NotiFloorTiles_data:
		return app.objRecvNotiFn_FloorTiles(rpk.Header, body)
	case *c2t_obj.NotiFieldObjList_data:
		return app.objRecvNotiFn_FieldObjList(rpk.Header, body)
	case *c2t_obj.NotiFoundFieldObj_data:
		return app.objRecvNotiFn_FoundFieldObj(rpk.Header, body)
	case *c2t_obj.NotiForgetFloor_data:
		return app.objRecvNotiFn_ForgetFloor(rpk.Header, body)
	case *c2t_obj.NotiActivateTrap_data:
		return app.objRecvNotiFn_ActivateTrap(rpk.Header, body)

	}
}

func (app *GLClient) objRecvNotiFn_EnterFloor(hd c2t_packet.Header, body *c2t_obj.NotiEnterFloor_data) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FI.Name {
		app.CurrentFloor = clientfloor.New(body.FI)
	}
	app.CurrentFloor.EnterFloor()
	app.resizeGLFloor()
	return nil
}
func (app *GLClient) objRecvNotiFn_LeaveFloor(hd c2t_packet.Header, body *c2t_obj.NotiLeaveFloor_data) error {
	// do nothing
	return nil
}

func (app *GLClient) objRecvNotiFn_Ageing(hd c2t_packet.Header, body *c2t_obj.NotiAgeing_data) error {
	// do nothing
	return nil
}
func (app *GLClient) objRecvNotiFn_Death(hd c2t_packet.Header, body *c2t_obj.NotiDeath_data) error {
	// do nothing
	return nil
}
func (app *GLClient) objRecvNotiFn_ReadyToRebirth(hd c2t_packet.Header, body *c2t_obj.NotiReadyToRebirth_data) error {
	go app.sendReqObjWithRspFn(c2t_idcmd.Rebirth,
		&c2t_obj.ReqRebirth_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		})
	return nil
}
func (app *GLClient) objRecvNotiFn_Rebirthed(hd c2t_packet.Header, body *c2t_obj.NotiRebirthed_data) error {
	// do nothing
	return nil
}

func (app *GLClient) objRecvNotiFn_VPObjList(hd c2t_packet.Header, body *c2t_obj.NotiVPObjList_data) error {
	app.OLNotiData = body
	oldOLNotiData := app.OLNotiData
	app.OLNotiData = body
	newOLNotiData := body
	app.onFieldObj = nil

	c2t_obj.EquipClientByUUID(body.ActiveObj.EquipBag).Sort()
	c2t_obj.PotionClientByUUID(body.ActiveObj.PotionBag).Sort()
	c2t_obj.ScrollClientByUUID(body.ActiveObj.ScrollBag).Sort()

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
		app.CurrentFloor.FieldObjPosMan.AddOrUpdateToXY(v, v.X, v.Y)
	}

	playerX, playerY := app.GetPlayerXY()
	if app.playerActiveObjClient != nil && app.CurrentFloor.IsValidPos(playerX, playerY) {
		app.onFieldObj = app.CurrentFloor.GetFieldObjAt(playerX, playerY)
	}
	app.actByControlMode()
	return nil
}
func (app *GLClient) objRecvNotiFn_VPTiles(hd c2t_packet.Header, body *c2t_obj.NotiVPTiles_data) error {
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
func (app *GLClient) objRecvNotiFn_FloorTiles(hd c2t_packet.Header, body *c2t_obj.NotiFloorTiles_data) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FI.Name {
		// new floor
		app.CurrentFloor = clientfloor.New(body.FI)
	}

	oldComplete := app.CurrentFloor.Visited.IsComplete()
	app.CurrentFloor.ReplaceFloorTiles(body)
	if !oldComplete && app.CurrentFloor.Visited.IsComplete() {
		// floor complete
	}
	return nil
}
func (app *GLClient) objRecvNotiFn_FieldObjList(hd c2t_packet.Header, body *c2t_obj.NotiFieldObjList_data) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FI.Name {
		// new floor
		app.CurrentFloor = clientfloor.New(body.FI)
	}
	app.CurrentFloor.UpdateFieldObjList(body.FOList)
	return nil
}
func (app *GLClient) objRecvNotiFn_FoundFieldObj(hd c2t_packet.Header, body *c2t_obj.NotiFoundFieldObj_data) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FloorName {
		g2log.Fatal("FoundFieldObj unknonw floor %v", body)
		return fmt.Errorf("FoundFieldObj unknonw floor %v", body)
	}
	if app.CurrentFloor.FieldObjPosMan.Get1stObjAt(body.FieldObj.X, body.FieldObj.Y) == nil {
		app.CurrentFloor.FieldObjPosMan.AddOrUpdateToXY(body.FieldObj, body.FieldObj.X, body.FieldObj.Y)
	}
	return nil
}
func (app *GLClient) objRecvNotiFn_ForgetFloor(hd c2t_packet.Header, body *c2t_obj.NotiForgetFloor_data) error {
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != body.FloorName {
	} else {
		app.CurrentFloor.Forget()
	}
	return nil
}
func (app *GLClient) objRecvNotiFn_ActivateTrap(hd c2t_packet.Header, body *c2t_obj.NotiActivateTrap_data) error {
	// do nothing
	return nil
}
