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

package carryingobject

import (
	"fmt"

	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/enum/carryingobjecttype"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/idu64str"
)

var ScrollIDMaker = idu64str.New("ScrollID")

type Scroll struct {
	uuid              string
	remainTurnInFloor int

	scrollType scrolltype.ScrollType
}

func (p Scroll) String() string {
	return fmt.Sprintf("Scroll[%v %v]",
		p.uuid, p.scrollType,
	)
}

func NewScroll(st scrolltype.ScrollType) gamei.ScrollI {
	rtn := &Scroll{
		uuid:       ScrollIDMaker.New(),
		scrollType: st,
	}
	return rtn
}

func NewScrollByMakeRate(n int) gamei.ScrollI {
	for i := 0; i < scrolltype.ScrollType_Count; i++ {
		sct := scrolltype.ScrollType(i)
		n -= sct.MakeRate()
		if n < 0 {
			return NewScroll(sct)
		}
	}
	return NewScroll(scrolltype.Empty)
}

func (po *Scroll) ToPacket_CarryObjClientOnFloor(x, y int) *csprotocol.CarryObjClientOnFloor {
	poc := &csprotocol.CarryObjClientOnFloor{
		UUID:               po.uuid,
		CarryingObjectType: po.GetCarryingObjectType(),
		X:                  x,
		Y:                  y,

		ScrollType: po.scrollType,
	}
	return poc
}

func (po *Scroll) ToPacket_ScrollClient() *csprotocol.ScrollClient {
	poc := &csprotocol.ScrollClient{
		UUID:       po.uuid,
		ScrollType: po.scrollType,
	}
	return poc
}

func (po *Scroll) GetScrollType() scrolltype.ScrollType {
	return po.scrollType
}

// IDPosI interface
func (po *Scroll) GetUUID() string {
	return po.uuid
}

func (po *Scroll) GetName() string {
	return po.scrollType.String()
}

func (po *Scroll) GetCarryingObjectType() carryingobjecttype.CarryingObjectType {
	return carryingobjecttype.Scroll
}

func (po *Scroll) GetValue() float64 {
	return gameconst.ScrollValue
}

func (po *Scroll) GetWeight() float64 {
	return gameconst.ScrollGram
}

// life in floor handle

func (po *Scroll) GetRemainTurnInFloor() int {
	return po.remainTurnInFloor
}
func (po *Scroll) DecRemainTurnInFloor() int {
	if po.remainTurnInFloor > 0 {
		po.remainTurnInFloor--
	}
	return po.remainTurnInFloor
}
func (po *Scroll) SetRemainTurnInFloor() {
	po.remainTurnInFloor = gameconst.CarryingObjectLifeTurnInFloor
}
