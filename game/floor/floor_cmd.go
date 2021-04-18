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

package floor

import (
	"time"

	"github.com/kasworld/goguelike-single/enum/achievetype"
	"github.com/kasworld/goguelike-single/enum/aotype"
	"github.com/kasworld/goguelike-single/enum/returncode"
	"github.com/kasworld/goguelike-single/game/cmd2floor"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/uuidposmani"
)

func (f *Floor) processCmd(data interface{}) {
	f.cmdActStat.Inc()
	switch pk := data.(type) {
	default:
		g2log.Fatal("unknown pk recv %v %#v", f, data)

	case *cmd2floor.ReqLeaveFloor:
		if err := f.aoPosMan.Del(pk.ActiveObj); err != nil {
			g2log.Fatal("%v %v", f, err)
		}
		if pk.ActiveObj.GetActiveObjType() == aotype.User {
			f.tower.SendNoti(
				&csprotocol.NotiLeaveFloor{
					FI: f.ToPacket_FloorInfo(),
				},
			)
		}

	case *cmd2floor.ReqEnterFloor:
		err := f.aoPosMan.AddOrUpdateToXY(pk.ActiveObj, pk.X, pk.Y)
		if err != nil {
			g2log.Fatal("%v %v", f, err)
		}
		pk.ActiveObj.EnterFloor(f)
		if pk.ActiveObj.GetActiveObjType() == aotype.User {
			f.tower.SendNoti(
				&csprotocol.NotiEnterFloor{
					FI: f.ToPacket_FloorInfo(),
				},
			)
			// send known tile area
			f4c := pk.ActiveObj.GetFloor4Client(f.GetName())
			fi := f.ToPacket_FloorInfo()
			ta := f.GetTerrain().GetTiles().DupWithFilter(f4c.Visit.GetXYNolock)
			f.tower.SendNoti(
				&csprotocol.NotiFloorTiles{
					FI:    fi,
					X:     0,
					Y:     0,
					Tiles: ta,
				},
			)

			// send fieldobj list
			fol := make([]*csprotocol.FieldObjClient, 0)
			f4c.FOPosMan.IterAll(func(o uuidposmani.UUIDPosI, foX, foY int) bool {
				fo := o.(*csprotocol.FieldObjClient)
				fol = append(fol, fo)
				return false
			})
			f.tower.SendNoti(
				&csprotocol.NotiFieldObjList{
					FI:     fi,
					FOList: fol,
				},
			)
			f.sendTANoti2Player(pk.ActiveObj)
			f.sendVPObj2Player(pk.ActiveObj, time.Now())
		}

	case *cmd2floor.ReqRebirth2Floor:
		err := f.aoPosMan.AddOrUpdateToXY(pk.ActiveObj, pk.X, pk.Y)
		if err != nil {
			g2log.Fatal("%v %v", f, err)
		}
		pk.ActiveObj.Rebirth()
		if pk.ActiveObj.GetActiveObjType() == aotype.User {
			f.tower.SendNoti(
				&csprotocol.NotiRebirthed{},
			)
		}

	case *cmd2floor.APIAdminTeleport2Floor:
		pk.RspCh <- f.Call_APIAdminTeleport2Floor(pk.ActiveObj, pk.ReqPk)
	}
}

func (f *Floor) Call_APIAdminTeleport2Floor(
	ActiveObj gamei.ActiveObjectI, ReqPk *csprotocol.ReqAdminTeleport) returncode.ReturnCode {

	if f.aoPosMan.GetByUUID(ActiveObj.GetUUID()) == nil {
		g2log.Warn("ActiveObj not in floor %v %v", f, ActiveObj)
		return returncode.ActionProhibited
	}
	x, y, err := f.SearchRandomActiveObjPos()
	if err != nil {
		g2log.Error("fail to teleport %v %v %v", f, ActiveObj, err)
		return returncode.ActionCanceled
	}
	x, y = f.terrain.WrapXY(x, y)
	if err := f.aoPosMan.UpdateToXY(ActiveObj, x, y); err != nil {
		g2log.Fatal("move ao fail %v %v %v", f, ActiveObj, err)
		return returncode.ActionCanceled
	}
	ActiveObj.SetNeedTANoti()
	ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
	return returncode.Success
}
