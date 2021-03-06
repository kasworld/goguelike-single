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

package terrain

import (
	"fmt"

	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/lib/scriptparse"

	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike-single/game/fieldobject"
	"github.com/kasworld/goguelike-single/game/terrain/roomsort"
)

func cmdAddPortal(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var x, y int
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var PortalID, DstPortalID string
	var acttype fieldobjacttype.FieldObjActType
	var message string
	if err := ca.GetArgs(&x, &y, &dispType, &acttype, &PortalID, &DstPortalID, &message); err != nil {
		return err
	}
	return tr.addPortal(PortalID, DstPortalID, x, y, dispType, acttype, message)
}

func cmdAddPortalRand(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var PortalID, DstPortalID string
	var acttype fieldobjacttype.FieldObjActType
	var message string
	if err := ca.GetArgs(&dispType, &acttype, &PortalID, &DstPortalID, &message); err != nil {
		return err
	}
	return tr.addPortalRand(PortalID, DstPortalID, dispType, acttype, message)
}

func cmdAddPortalRandInRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var PortalID, DstPortalID string
	var acttype fieldobjacttype.FieldObjActType
	var message string
	if err := ca.GetArgs(&dispType, &acttype, &PortalID, &DstPortalID, &message); err != nil {
		return err
	}
	return tr.addPortalRandInRoom(PortalID, DstPortalID, dispType, acttype, message)
}

func (tr *Terrain) isPlaceableWithVt(x, y int, vx, vy int) bool {
	tx, ty := tr.WrapXY(x+vx, y+vy)
	return tr.serviceTileArea[tx][ty].CharPlaceable()
}

func (tr *Terrain) isBlockWay(x, y int) bool {
	vt := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}
	contack := [4]bool{}
	contactSum := 0
	contackDiag := [2]int{}
	for i, v := range vt {
		if !tr.isPlaceableWithVt(x, y, v[0], v[1]) {
			contactSum++
			contack[i] = true
			contackDiag[0] += v[0]
			contackDiag[1] += v[1]
		}
	}
	if contactSum == 2 {
		if contack[0] && contack[1] {
			return true
		}
		if contack[2] && contack[3] {
			return true
		}
		return tr.isPlaceableWithVt(x, y, contackDiag[0], contackDiag[1])
	}
	return contactSum > 2
}

// forbid fieldobj contact
func (tr *Terrain) canPlaceFieldObjAt(x, y int) bool {
	for _, v := range [][2]int{
		{0, 0},
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if tr.foPosMan.Get1stObjAt(x+v[0], y+v[1]) != nil {
			return false
		}
	}
	tl := tr.serviceTileArea[x][y]
	return tl.CharPlaceable() && !tl.TestByTile(tile.Door)
}

func (tr *Terrain) addPortal(portalID string, dstPortalID string,
	x, y int,
	dispType fieldobjdisplaytype.FieldObjDisplayType,
	acttype fieldobjacttype.FieldObjActType, message string) error {

	x, y = x%tr.Xlen, y%tr.Ylen

	if !tr.canPlaceFieldObjAt(x, y) {
		return fmt.Errorf("can not add portal at NonCharPlaceable tile %v %v", x, y)
	}
	po := fieldobject.NewPortal(tr.Name,
		dispType, message, acttype,
		portalID, dstPortalID,
	)
	tr.foPosMan.AddToXY(po, x, y)
	if r := tr.roomManager.GetRoomByPos(x, y); r != nil {
		r.PortalCount++
	}

	return nil
}

func (tr *Terrain) addPortalRand(portalID string, dstPortalID string,
	dispType fieldobjdisplaytype.FieldObjDisplayType,
	acttype fieldobjacttype.FieldObjActType, message string) error {

	for try := 100; try > 0; try-- {
		x, y := tr.rnd.Intn(tr.Xlen), tr.rnd.Intn(tr.Ylen)
		if !tr.canPlaceFieldObjAt(x, y) {
			continue
		}
		if acttype == fieldobjacttype.PortalAutoIn && tr.isBlockWay(x, y) {
			continue
		}
		return tr.addPortal(portalID, dstPortalID, x, y, dispType, acttype, message)
	}
	return fmt.Errorf("fail to addPortalRand")
}

func (tr *Terrain) addPortalRandInRoom(portalID string, dstPortalID string,
	dispType fieldobjdisplaytype.FieldObjDisplayType,
	acttype fieldobjacttype.FieldObjActType, message string) error {

	if tr.roomManager.GetCount() == 0 {
		return fmt.Errorf("no room to add portal")
	}
	roomList := tr.roomManager.GetRoomList()
	for try := 100; try > 0; try-- {
		tr.rnd.Shuffle(len(roomList), func(i, j int) {
			roomList[i], roomList[j] = roomList[j], roomList[i]
		})
		rList := roomsort.ByPortalCount(roomList)
		rList.Sort()
		r := rList[0]
		x := tr.rnd.IntRange(r.Area.X, r.Area.X+r.Area.W)
		y := tr.rnd.IntRange(r.Area.Y, r.Area.Y+r.Area.H)
		if !tr.canPlaceFieldObjAt(x, y) {
			continue
		}
		if acttype == fieldobjacttype.PortalAutoIn && tr.isBlockWay(x, y) {
			continue
		}
		return tr.addPortal(portalID, dstPortalID, x, y, dispType, acttype, message)
	}
	return fmt.Errorf("fail to addPortalRandInRoom")
}
