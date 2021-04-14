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

package tower

import (
	"github.com/kasworld/goguelike-single/enum/achievetype"
	"github.com/kasworld/goguelike-single/enum/respawntype"
	"github.com/kasworld/goguelike-single/game/cmd2tower"
	"github.com/kasworld/goguelike-single/game/fieldobject"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
)

func (tw *Tower) processCmd(data interface{}) {
	tw.cmdActStat.Inc()
	switch pk := data.(type) {
	default:
		g2log.Fatal("unknown tower cmd %#v", data)

	case *cmd2tower.AdminFloorMove:
		pk.RspCh <- tw.Call_AdminFloorMove(pk.ActiveObj, pk.RecvPacket)

	case *cmd2tower.FloorMove:
		pk.RspCh <- tw.Call_FloorMove(pk.ActiveObj, pk.FloorName)

	case *cmd2tower.ActiveObjUsePortal:
		tw.Call_ActiveObjUsePortal(pk.ActiveObj, pk.SrcFloor, pk.P1, pk.P2)

	case *cmd2tower.ActiveObjTrapTeleport:
		tw.Call_ActiveObjTrapTeleport(pk.ActiveObj, pk.SrcFloor, pk.DstFloorName)

	case *cmd2tower.ActiveObjRebirth:
		tw.Call_ActiveObjRebirth(pk.ActiveObj)

	case *cmd2tower.Turn:
		tw.Turn(pk.Now)
	}
}

func (tw *Tower) Call_AdminFloorMove(
	ActiveObj gamei.ActiveObjectI,
	RecvPacket *c2t_obj.ReqAdminFloorMove_data) c2t_error.ErrorCode {

	switch cmd := RecvPacket.Floor; cmd {
	case "Before":
		aoFloor := tw.ao2Floor.GetFloorByActiveObjID(ActiveObj.GetUUID())
		aoFloorIndex, err := tw.floorMan.GetFloorIndexByName(aoFloor.GetName())
		if err != nil {
			g2log.Error("floor not found %v", aoFloor)
			return c2t_error.ObjectNotFound
		}
		dstFloor := tw.floorMan.GetFloorByIndexWrap(aoFloorIndex - 1)
		x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
		if err != nil {
			g2log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
		}
		if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
			g2log.Fatal("%v", err)
			return c2t_error.ActionProhibited
		}
		ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
		return c2t_error.None
	case "Next":
		aoFloor := tw.ao2Floor.GetFloorByActiveObjID(ActiveObj.GetUUID())
		aoFloorIndex, err := tw.floorMan.GetFloorIndexByName(aoFloor.GetName())
		if err != nil {
			g2log.Error("floor not found %v", aoFloor)
			return c2t_error.ObjectNotFound
		}
		dstFloor := tw.floorMan.GetFloorByIndexWrap(aoFloorIndex + 1)
		x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
		if err != nil {
			g2log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
		}
		if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
			g2log.Fatal("%v", err)
			return c2t_error.ActionProhibited
		}
		ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
		return c2t_error.None
	default:
		dstFloor := tw.floorMan.GetFloorByName(cmd)
		if dstFloor == nil {
			g2log.Error("floor not found %v", cmd)
			return c2t_error.ObjectNotFound
		}
		x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
		if err != nil {
			g2log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
		}
		if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
			g2log.Fatal("%v", err)
			return c2t_error.ActionProhibited
		}
		ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
		return c2t_error.None
	}
}

func (tw *Tower) Call_FloorMove(
	ActiveObj gamei.ActiveObjectI, FloorName string) c2t_error.ErrorCode {

	dstFloor := tw.floorMan.GetFloorByName(FloorName)
	if dstFloor == nil {
		g2log.Error("floor not found %v", FloorName)
		return c2t_error.ObjectNotFound
	}
	x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
	if err != nil {
		g2log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
	}
	if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
		g2log.Fatal("%v", err)
		return c2t_error.ActionProhibited
	}
	return c2t_error.None
}

func (tw *Tower) Call_ActiveObjUsePortal(
	ActiveObj gamei.ActiveObjectI,
	SrcFloor gamei.FloorI,
	P1, P2 *fieldobject.FieldObject,
) {
	srcPosMan := SrcFloor.GetActiveObjPosMan()
	if srcPosMan == nil {
		g2log.Warn("pos man nil %v", SrcFloor)
		return
	}
	if srcPosMan.GetByUUID(ActiveObj.GetUUID()) == nil {
		g2log.Warn("ActiveObj not in floor %v %v", ActiveObj, SrcFloor)
		return
	}
	dstFloor := tw.floorMan.GetFloorByName(P2.FloorName)
	if dstFloor == nil {
		g2log.Fatal("dstFloor not found %v", P2.FloorName)
		return
	}
	g2log.Debug("ActiveObjUsePortal %v %v to %v", ActiveObj, SrcFloor, dstFloor)

	x, y, exist := dstFloor.GetFieldObjPosMan().GetXYByUUID(P2.GetUUID())
	if !exist {
		g2log.Fatal("fieldobj not found %v", P2)
		return
	}

	if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
		g2log.Fatal("%v", err)
	}
}

func (tw *Tower) Call_ActiveObjTrapTeleport(
	ActiveObj gamei.ActiveObjectI,
	SrcFloor gamei.FloorI,
	DstFloorName string,
) {
	dstFloor := tw.floorMan.GetFloorByName(DstFloorName)
	if dstFloor == nil {
		g2log.Fatal("dstFloor not found %v", DstFloorName)
		return
	}
	g2log.Debug("ActiveObjTrapTeleport %v to %v", ActiveObj, dstFloor)
	x, y, err := dstFloor.SearchRandomActiveObjPos()
	if err != nil {
		g2log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
	}
	if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
		g2log.Fatal("%v", err)
	}
}

func (tw *Tower) Call_ActiveObjRebirth(ao gamei.ActiveObjectI) {
	if ao.IsAlive() {
		g2log.Fatal("ao is alive %v HP:%v/%v",
			ao, ao.GetHP(), ao.GetTurnData().HPMax)
	}
	var dstFloor gamei.FloorI

	switch ao.GetRespawnType() {
	default:
		g2log.Fatal("invalid respawntype %v %v", ao, ao.GetRespawnType())
	case respawntype.ToCurrentFloor:
		dstFloor = ao.GetCurrentFloor()
	case respawntype.ToHomeFloor:
		dstFloor = ao.GetHomeFloor()
	case respawntype.ToRandomFloor:
		dstFloor = tw.floorMan.GetFloorList()[tw.rnd.Intn(tw.floorMan.GetFloorCount())]
	}
	if err := tw.ao2Floor.ActiveObjRebirthToFloor(dstFloor, ao); err != nil {
		g2log.Fatal("%v", err)
		return
	}
	g2log.Debug("ActiveObjRebirth %v to %v", ao, dstFloor)
}
