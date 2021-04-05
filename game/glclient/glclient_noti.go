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

	"github.com/kasworld/goguelike-single/config/leveldata"
	"github.com/kasworld/goguelike-single/game/clientfloor"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
)

var DemuxNoti2ByteFnMap = [...]func(me interface{}, hd c2t_packet.Header, rbody []byte) error{
	c2t_idnoti.Invalid:        bytesRecvNotiFn_Invalid,        // Invalid make empty packet error
	c2t_idnoti.EnterTower:     bytesRecvNotiFn_EnterTower,     // EnterTower
	c2t_idnoti.EnterFloor:     bytesRecvNotiFn_EnterFloor,     // EnterFloor
	c2t_idnoti.LeaveFloor:     bytesRecvNotiFn_LeaveFloor,     // LeaveFloor
	c2t_idnoti.LeaveTower:     bytesRecvNotiFn_LeaveTower,     // LeaveTower
	c2t_idnoti.Ageing:         bytesRecvNotiFn_Ageing,         // Ageing          // floor
	c2t_idnoti.Death:          bytesRecvNotiFn_Death,          // Death
	c2t_idnoti.ReadyToRebirth: bytesRecvNotiFn_ReadyToRebirth, // ReadyToRebirth
	c2t_idnoti.Rebirthed:      bytesRecvNotiFn_Rebirthed,      // Rebirthed
	c2t_idnoti.Broadcast:      bytesRecvNotiFn_Broadcast,      // Broadcast       // global chat broadcast from web admin
	c2t_idnoti.VPObjList:      bytesRecvNotiFn_VPObjList,      // VPObjList       // in viewport, every turn
	c2t_idnoti.VPTiles:        bytesRecvNotiFn_VPTiles,        // VPTiles         // in viewport, when viewport changed only
	c2t_idnoti.FloorTiles:     bytesRecvNotiFn_FloorTiles,     // FloorTiles      // for rebuild known floor
	c2t_idnoti.FieldObjList:   bytesRecvNotiFn_FieldObjList,   // FieldObjList    // for rebuild known floor
	c2t_idnoti.FoundFieldObj:  bytesRecvNotiFn_FoundFieldObj,  // FoundFieldObj   // hidden field obj
	c2t_idnoti.ForgetFloor:    bytesRecvNotiFn_ForgetFloor,    // ForgetFloor
	c2t_idnoti.ActivateTrap:   bytesRecvNotiFn_ActivateTrap,   // ActivateTrap
}

func bytesRecvNotiFn_Invalid(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return fmt.Errorf("not implemented %v", hd)
}

func bytesRecvNotiFn_EnterTower(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiEnterTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	app.TowerInfo = pkbody.TowerInfo

	return nil
}
func bytesRecvNotiFn_LeaveTower(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiLeaveTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	_ = pkbody
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = app
	return nil
}

func bytesRecvNotiFn_EnterFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiEnterFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != pkbody.FI.Name {
		app.CurrentFloor = clientfloor.New(pkbody.FI)
	}

	app.CurrentFloor.EnterFloor()

	return nil
}
func bytesRecvNotiFn_LeaveFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiLeaveFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	_ = pkbody
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = app
	return nil
}

func bytesRecvNotiFn_Ageing(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return nil
}

func bytesRecvNotiFn_Death(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = app
	if app.config.DisconnectOnDeath {
		app.sendRecvStop()
	}
	return nil
}

func bytesRecvNotiFn_ReadyToRebirth(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	go app.ReqWithRspFnWithAuth(c2t_idcmd.Rebirth,
		&c2t_obj.ReqRebirth_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		})
	return nil
}

func bytesRecvNotiFn_Rebirthed(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return nil
}

func bytesRecvNotiFn_Broadcast(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return nil
}

func bytesRecvNotiFn_VPObjList(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiVPObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	app.OLNotiData = pkbody
	app.ServerClientTimeDiff = pkbody.Time.Sub(time.Now())
	oldOLNotiData := app.OLNotiData
	app.OLNotiData = pkbody
	newOLNotiData := pkbody
	app.onFieldObj = nil

	app.ServerJitter.ActByValue(pkbody.Time)

	c2t_obj.EquipClientByUUID(pkbody.ActiveObj.EquipBag).Sort()
	c2t_obj.PotionClientByUUID(pkbody.ActiveObj.PotionBag).Sort()
	c2t_obj.ScrollClientByUUID(pkbody.ActiveObj.ScrollBag).Sort()

	if oldOLNotiData != nil {
		app.HPdiff = newOLNotiData.ActiveObj.HP - oldOLNotiData.ActiveObj.HP
		app.SPdiff = newOLNotiData.ActiveObj.SP - oldOLNotiData.ActiveObj.SP
	}
	newLevel := int(leveldata.CalcLevelFromExp(float64(newOLNotiData.ActiveObj.Exp)))

	app.playerActiveObjClient = nil
	if ainfo := app.AccountInfo; ainfo != nil {
		for _, v := range pkbody.ActiveObjList {
			if v.UUID == app.AccountInfo.ActiveObjUUID {
				app.playerActiveObjClient = v
			}
		}
	}

	app.IsOverLoad = newOLNotiData.ActiveObj.CalcWeight() >= leveldata.WeightLimit(newLevel)

	if app.CurrentFloor.FloorInfo == nil {
		app.log.Error("app.CurrentFloor.FloorInfo not set")
		return nil
	}
	if app.CurrentFloor.FloorInfo.Name != newOLNotiData.FloorName {
		app.log.Error("not current floor objlist data %v %v",
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

func bytesRecvNotiFn_VPTiles(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiVPTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}

	if app.CurrentFloor.FloorInfo == nil {
		app.log.Warn("OrangeRed app.CurrentFloor.FloorInfo not set")
		return nil
	}
	if app.CurrentFloor.FloorInfo.Name != pkbody.FloorName {
		app.log.Warn("not current floor vptile data %v %v",
			app.CurrentFloor.FloorInfo.Name, pkbody.FloorName,
		)
		return nil
	}

	oldComplete := app.CurrentFloor.Visited.IsComplete()
	if err := app.CurrentFloor.UpdateFromViewportTile(pkbody, app.ViewportXYLenList); err != nil {
		app.log.Warn("%v", err)
		return nil
	}
	if !oldComplete && app.CurrentFloor.Visited.IsComplete() {
		// just completed
	}

	return nil
}

func bytesRecvNotiFn_FloorTiles(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiFloorTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != pkbody.FI.Name {
		// new floor
		app.CurrentFloor = clientfloor.New(pkbody.FI)
	}

	oldComplete := app.CurrentFloor.Visited.IsComplete()
	app.CurrentFloor.ReplaceFloorTiles(pkbody)
	if !oldComplete && app.CurrentFloor.Visited.IsComplete() {
		// floor complete
	}
	return nil
}

// FieldObjList    // for rebuild known floor
func bytesRecvNotiFn_FieldObjList(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiFieldObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != pkbody.FI.Name {
		// new floor
		app.CurrentFloor = clientfloor.New(pkbody.FI)
	}
	app.CurrentFloor.UpdateFieldObjList(pkbody.FOList)
	return nil
}

func bytesRecvNotiFn_FoundFieldObj(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiFoundFieldObj_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != pkbody.FloorName {
		app.log.Fatal("FoundFieldObj unknonw floor %v", pkbody)
		return fmt.Errorf("FoundFieldObj unknonw floor %v", pkbody)
	}
	if app.CurrentFloor.FieldObjPosMan.Get1stObjAt(pkbody.FieldObj.X, pkbody.FieldObj.Y) == nil {
		app.CurrentFloor.FieldObjPosMan.AddOrUpdateToXY(pkbody.FieldObj, pkbody.FieldObj.X, pkbody.FieldObj.Y)
	}
	return nil
}

func bytesRecvNotiFn_ForgetFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiForgetFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}

	if app.CurrentFloor == nil || app.CurrentFloor.FloorInfo.Name != pkbody.FloorName {
	} else {
		app.CurrentFloor.Forget()
	}
	return nil
}

func bytesRecvNotiFn_ActivateTrap(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiActivateTrap_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	app, ok := me.(*GLClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = app
	_ = pkbody
	// g2log.Debug("%v", robj)
	return nil
}
