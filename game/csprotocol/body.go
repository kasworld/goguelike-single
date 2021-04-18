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

package csprotocol

import (
	"time"

	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/enum/achievetype_vector"
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/condition_vector"
	"github.com/kasworld/goguelike-single/enum/equipslottype"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype_vector"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/potiontype_vector"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/enum/scrolltype_vector"
	"github.com/kasworld/goguelike-single/enum/turnaction"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/tilearea"
)

////////////////////////////////////////////////////////////////////
// commnad protocol

// AchieveInfo
type ReqAchieveInfo struct {
}
type RspAchieveInfo struct {
	AchieveStat   achievetype_vector.AchieveTypeVector         `prettystring:"simple"`
	PotionStat    potiontype_vector.PotionTypeVector           `prettystring:"simple"`
	ScrollStat    scrolltype_vector.ScrollTypeVector           `prettystring:"simple"`
	FOActStat     fieldobjacttype_vector.FieldObjActTypeVector `prettystring:"simple"`
	ConditionStat condition_vector.ConditionVector             `prettystring:"simple"`
}

// AIPlay
type ReqAIPlay struct {
	On bool
}
type RspAIPlay struct {
}

// VisitFloorList floor info of visited
type ReqVisitFloorList struct {
}
type RspVisitFloorList struct {
	FloorList []*FloorInfo
}

////////////////////////////////////////////////////////////////////////////
// ao turn action

type ReqTurnAction struct {
	Act  turnaction.TurnAction
	Dir  way9type.Way9Type
	UUID string
}

////////////////////////////////////////////////////////////////////
// admin

// AdminFloorMove Next Before floorUUID
type ReqAdminFloorMove struct {
	Floor string
}
type RspAdminFloorMove struct {
}

// AdminTeleport random pos in floor
type ReqAdminTeleport struct {
	X int
	Y int
}
type RspAdminTeleport struct {
}

// AdminAddExp  add arg to battle exp
type ReqAdminAddExp struct {
	Exp int
}
type RspAdminAddExp struct {
}

// AdminPotionEffect buff by arg potion type
type ReqAdminPotionEffect struct {
	Potion potiontype.PotionType
}
type RspAdminPotionEffect struct {
}

// AdminScrollEffect buff by arg Scroll type
type ReqAdminScrollEffect struct {
	Scroll scrolltype.ScrollType
}
type RspAdminScrollEffect struct {
}

// AdminCondition add arg condition for 100 turn
type ReqAdminCondition struct {
	Condition condition.Condition
}
type RspAdminCondition struct {
}

// AdminAddPotion add arg potion to inven
type ReqAdminAddPotion struct {
	Potion potiontype.PotionType
}
type RspAdminAddPotion struct {
}

// AdminAddScroll add arg scroll to inven
type ReqAdminAddScroll struct {
	Scroll scrolltype.ScrollType
}
type RspAdminAddScroll struct {
}

// AdminAddMoney add arg money to inven
type ReqAdminAddMoney struct {
	Money int
}
type RspAdminAddMoney struct {
}

// AdminAddEquip add random equip to inven
type ReqAdminAddEquip struct {
	Faction factiontype.FactionType
	Equip   equipslottype.EquipSlotType
}
type RspAdminAddEquip struct {
}

// AdminForgetFloor forget current floor map
type ReqAdminForgetFloor struct {
}
type RspAdminForgetFloor struct {
}

// AdminFloorMap complete current floor map
type ReqAdminFloorMap struct {
}
type RspAdminFloorMap struct {
}

/////////////////////////////////////////////////////////////////
// noti

type NotiEnterFloor struct {
	FI *FloorInfo
}
type NotiLeaveFloor struct {
	FI *FloorInfo
}

type NotiAgeing struct {
	FloorName string
}

type NotiDeath struct {
}

type NotiReadyToRebirth struct {
}
type NotiRebirthed struct {
}

type NotiVPObjList struct {
	Time          time.Time `prettystring:"simple"`
	FloorName     string
	ActiveObj     *PlayerActiveObjInfo
	ActiveObjList []*ActiveObjClient
	CarryObjList  []*CarryObjClientOnFloor
	FieldObjList  []*FieldObjClient
	DangerObjList []*DangerObjClient
}

// NotiVPTiles contains tile info center from pos
type NotiVPTiles struct {
	FloorName string
	VPX       int // viewport center X
	VPY       int // viewport center Y
	VPTiles   *viewportdata.ViewportTileArea2
}

// NotiFloorTiles used for floor map, reconnect client
type NotiFloorTiles struct {
	FI    *FloorInfo
	X     int // X start position, not center
	Y     int // Y start position, not center
	Tiles tilearea.TileArea
}

// FieldObjList    // for rebuild known floor
type NotiFieldObjList struct {
	FI     *FloorInfo
	FOList []*FieldObjClient
}

type NotiFoundFieldObj struct {
	FloorName string
	FieldObj  *FieldObjClient
}

type NotiForgetFloor struct {
	FloorName string
}

type NotiActivateTrap struct {
	FieldObjAct fieldobjacttype.FieldObjActType
	Triggered   bool
}
