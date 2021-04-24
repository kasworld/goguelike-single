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

package activeobject

import (
	"fmt"

	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/enum/aotype"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/fieldobject"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/uuidposmani"
)

func (ao *ActiveObject) EnterFloor(f gamei.FloorI) {
	ao.currentFloor = f
	ao.SetNeedTANoti()
	if _, exist := ao.floor4ClientMan.GetByName(f.GetName()); !exist {
		ao.floor4ClientMan.Add(f)
	}
	if aio := ao.ai; aio != nil {
		ao.ResetPlan(aio)
	}
}

func (ao *ActiveObject) UpdateVisitAreaBySightMat2(
	f gamei.FloorI, vpCenterX, vpCenterY int,
	sightMat *viewportdata.ViewportSight2, sight float32) {
	f4c, exist := ao.floor4ClientMan.GetByName(f.GetName())
	if !exist {
		g2log.Fatal("floor not visited %v %v", ao, f.GetName())
		f4c = ao.floor4ClientMan.Add(f)
	}
	f4c.Visit.UpdateBySightMat2(
		f.GetTerrain().GetTiles(),
		vpCenterX, vpCenterY,
		sightMat,
		sight)
}

func (ao *ActiveObject) forgetAnyFloor() error {
	for _, floorName := range ao.floor4ClientMan.GetNameList() {
		return ao.ForgetFloorByName(floorName)
	}
	return fmt.Errorf("no visit floor")
}

func (ao *ActiveObject) ForgetFloorByName(floorName string) error {
	if err := ao.floor4ClientMan.Forget(floorName); err != nil {
		return err
	}
	if ao.aoType == aotype.User {
		ao.homefloor.GetTower().SendNoti(
			&csprotocol.NotiForgetFloor{
				FloorName: floorName,
			},
		)
	}
	return nil
}

func (ao *ActiveObject) MakeFloorComplete(f gamei.FloorI) error {
	f4c, exist := ao.floor4ClientMan.GetByName(f.GetName())
	if !exist {
		f4c = ao.floor4ClientMan.Add(f)
	}

	f4c.Visit.MakeComplete()
	f.GetFieldObjPosMan().IterAll(func(o uuidposmani.UUIDPosI, foX, foY int) bool {
		fo := o.(*fieldobject.FieldObject)
		f4c.FOPosMan.AddOrUpdateToXY(fo.ToPacket_FieldObjClient(foX, foY), foX, foY)
		return false
	})

	fi := f.ToPacket_FloorInfo()
	if ao.aoType == aotype.User {
		// send tile area
		ta := f.GetTerrain().GetTiles()
		fol := make([]*csprotocol.FieldObjClient, 0)
		f4c.FOPosMan.IterAll(func(o uuidposmani.UUIDPosI, foX, foY int) bool {
			fo := o.(*csprotocol.FieldObjClient)
			fol = append(fol, fo)
			return false
		})
		ao.homefloor.GetTower().SendNoti(
			&csprotocol.NotiFloorComplete{
				FI:     fi,
				Tiles:  ta,
				FOList: fol,
			},
		)
		// send fieldobj list
	}
	return nil
}
