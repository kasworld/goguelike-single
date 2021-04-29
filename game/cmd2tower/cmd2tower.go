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

package cmd2tower

import (
	"fmt"

	"github.com/kasworld/goguelike-single/enum/returncode"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/fieldobject"
	"github.com/kasworld/goguelike-single/game/gamei"
)

type AdminFloorMove struct {
	ActiveObj  gamei.ActiveObjectI
	RecvPacket *csprotocol.ReqAdminFloorMove
	RspCh      chan<- returncode.ReturnCode
}

func (cet AdminFloorMove) String() string {
	return fmt.Sprintf("AdminFloorMove[%v %v]",
		cet.ActiveObj,
		cet.RecvPacket,
	)
}

type FloorMove struct {
	ActiveObj gamei.ActiveObjectI
	FloorName string
	RspCh     chan<- returncode.ReturnCode
}

func (cet FloorMove) String() string {
	return fmt.Sprintf("FloorMove[%v %v]",
		cet.ActiveObj,
		cet.FloorName,
	)
}

type ActiveObjTrapTeleport struct {
	ActiveObj    gamei.ActiveObjectI
	SrcFloor     gamei.FloorI
	DstFloorName string
}

type ActiveObjUsePortal struct {
	ActiveObj gamei.ActiveObjectI
	SrcFloor  gamei.FloorI
	P1, P2    *fieldobject.FieldObject
}

func (pk ActiveObjUsePortal) String() string {
	return fmt.Sprintf(
		"ActiveObjUsePortal[%v %v %v %v]",
		pk.SrcFloor,
		pk.ActiveObj,
		pk.P1,
		pk.P2,
	)
}

type ActiveObjRebirth struct {
	ActiveObj gamei.ActiveObjectI
}

func (pk ActiveObjRebirth) String() string {
	return fmt.Sprintf(
		"ActiveObjRebirth[%v]",
		pk.ActiveObj,
	)
}

type Turn struct {
	TurnCount int
}
