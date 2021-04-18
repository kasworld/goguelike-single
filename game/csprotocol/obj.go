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
	"fmt"
	"sort"
	"time"

	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/enum/aiplan"
	"github.com/kasworld/goguelike-single/enum/carryingobjecttype"
	"github.com/kasworld/goguelike-single/enum/condition_flag"
	"github.com/kasworld/goguelike-single/enum/dangertype"
	"github.com/kasworld/goguelike-single/enum/equipslottype"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/returncode"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/enum/turnaction"
	"github.com/kasworld/goguelike-single/enum/turnresulttype"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/aoactreqrsp"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/prettystring"
)

type GameInfo struct {
	Version         string
	ProtocolVersion string
	DataVersion     string

	ActiveObjUUID string
	NickName      string

	TowerSeed     int64
	TowerName     string
	Factor        [3]int64 `prettystring:"simple"`
	TotalFloorNum int
	StartTime     time.Time `prettystring:"simple"`
}

func (info *GameInfo) StringForm() string {
	return prettystring.PrettyString(info, 4)
}

type FloorInfo struct {
	Name       string
	W          int
	H          int
	Tiles      int
	Bias       bias.Bias
	VisitCount int
}

func (fi FloorInfo) GetName() string {
	return fi.Name
}
func (fi FloorInfo) GetWidth() int {
	return fi.W
}
func (fi FloorInfo) GetHeight() int {
	return fi.H
}
func (fi FloorInfo) VisitableCount() int {
	return fi.Tiles
}

type FloorInfoByName []*FloorInfo

func (objList FloorInfoByName) Len() int { return len(objList) }
func (objList FloorInfoByName) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList FloorInfoByName) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.Name > po2.Name
}
func (objList FloorInfoByName) Sort() {
	sort.Stable(objList)
}

type CarryObjClientOnFloor struct {
	UUID               string
	X                  int
	Y                  int
	CarryingObjectType carryingobjecttype.CarryingObjectType

	// for equip
	EquipType equipslottype.EquipSlotType
	Faction   factiontype.FactionType

	// for potion
	PotionType potiontype.PotionType

	// for scroll
	ScrollType scrolltype.ScrollType

	// for money
	Value int
}

func (po CarryObjClientOnFloor) String() string {
	switch po.CarryingObjectType {
	default:
		return fmt.Sprintf("unknonw po %#v", po)
	case carryingobjecttype.Equip:
		return fmt.Sprintf("%v[%v]",
			po.EquipType.String(),
			po.Faction.String(),
		)
	case carryingobjecttype.Money:
		return fmt.Sprintf("$%v", po.Value)
	case carryingobjecttype.Potion:
		return fmt.Sprintf("Potion%v", po.PotionType.String())
	case carryingobjecttype.Scroll:
		return fmt.Sprintf("Scroll%v", po.ScrollType.String())
	}
}

type FieldObjClient struct {
	ID          string
	X           int
	Y           int
	ActType     fieldobjacttype.FieldObjActType
	DisplayType fieldobjdisplaytype.FieldObjDisplayType
	Message     string
}

func (p *FieldObjClient) GetUUID() string {
	return p.ID
}

type FieldObjByType []*FieldObjClient

func (objList FieldObjByType) Len() int { return len(objList) }
func (objList FieldObjByType) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList FieldObjByType) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	if po1.ActType == po2.ActType {
		return po1.ID < po2.ID
	}
	return po1.ActType < po2.ActType
}
func (objList FieldObjByType) Sort() {
	sort.Stable(objList)
}

type DangerObjClient struct {
	UUID       string
	OwnerID    string
	DangerType dangertype.DangerType
	X          int
	Y          int
	AffectRate float64
}

type ActiveObjClient struct {
	UUID       string
	NickName   string
	Faction    factiontype.FactionType
	EquippedPo []*EquipClient
	Conditions condition_flag.ConditionFlag // not all condition
	X          int
	Y          int
	Alive      bool
	Chat       string

	// turn result
	Act        turnaction.TurnAction
	Dir        way9type.Way9Type
	Result     returncode.ReturnCode
	DamageGive int
	DamageTake int
}

type TurnResultClient struct {
	ResultType turnresulttype.TurnResultType
	DstUUID    string
	Arg        float64
}

func (ifa TurnResultClient) String() string {
	return fmt.Sprintf("%v DstObj:%v Arg:%3.1f",
		ifa.ResultType.String(),
		ifa.DstUUID,
		ifa.Arg,
	)
}

type EquipClient struct {
	UUID      string
	Name      string
	EquipType equipslottype.EquipSlotType
	Faction   factiontype.FactionType
	BiasLen   float64
}

func (po EquipClient) String() string {
	return fmt.Sprintf("%v %v%.0f",
		po.Name,
		po.Faction.Rune(),
		po.BiasLen,
	)
}

func (po EquipClient) GetBias() bias.Bias {
	return bias.NewByFaction(po.Faction, po.BiasLen)
}

func (po EquipClient) Weight() float64 {
	return po.BiasLen * gameconst.EquipABSGram
}

type EquipClientByUUID []*EquipClient

func (objList EquipClientByUUID) Len() int { return len(objList) }
func (objList EquipClientByUUID) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList EquipClientByUUID) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	if po1.EquipType == po2.EquipType {
		return po1.UUID < po2.UUID
	} else {
		return po1.EquipType < po2.EquipType
	}
}
func (objList EquipClientByUUID) Sort() {
	sort.Sort(objList)
}

type CarryObjEqByLen []*EquipClient

func (objList CarryObjEqByLen) Len() int { return len(objList) }
func (objList CarryObjEqByLen) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList CarryObjEqByLen) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.BiasLen < po2.BiasLen
}
func (objList CarryObjEqByLen) Sort() {
	sort.Stable(objList)
}

type PotionClient struct {
	UUID       string
	PotionType potiontype.PotionType
}

func (po PotionClient) String() string {
	return fmt.Sprintf("Potion%v", po.PotionType.String())
}

func (po PotionClient) Weight() int {
	return gameconst.PotionGram
}

type PotionClientByUUID []*PotionClient

func (objList PotionClientByUUID) Len() int { return len(objList) }
func (objList PotionClientByUUID) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList PotionClientByUUID) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.UUID < po2.UUID
}
func (objList PotionClientByUUID) Sort() {
	sort.Sort(objList)
}

type ScrollClient struct {
	UUID       string
	ScrollType scrolltype.ScrollType
}

func (po ScrollClient) String() string {
	return fmt.Sprintf("Scroll%v", po.ScrollType.String())
}

func (po ScrollClient) Weight() int {
	return gameconst.ScrollGram
}

type ScrollClientByUUID []*ScrollClient

func (objList ScrollClientByUUID) Len() int { return len(objList) }
func (objList ScrollClientByUUID) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList ScrollClientByUUID) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.UUID < po2.UUID
}
func (objList ScrollClientByUUID) Sort() {
	sort.Sort(objList)
}

type ActiveObjBuff struct {
	Name        string
	RemainCount int
}

type PlayerActiveObjInfo struct {
	Bias       bias.Bias
	Conditions condition_flag.ConditionFlag
	Exp        int
	TotalAO    int // total activeobject count
	Ranking    int // ranking / totalao
	Death      int
	Kill       int
	Sight      float64
	HP         float64
	HPMax      float64
	SP         float64
	SPMax      float64
	AIPlan     aiplan.AIPlan
	EquippedPo []*EquipClient
	EquipBag   []*EquipClient
	PotionBag  []*PotionClient
	ScrollBag  []*ScrollClient
	Wallet     int
	Wealth     int
	ActiveBuff []*ActiveObjBuff
	AP         float64

	Act        *aoactreqrsp.ActReqRsp
	TurnResult []TurnResultClient
}

func (pao PlayerActiveObjInfo) CalcWeight() float64 {
	weight := 0.0
	for _, v := range pao.EquippedPo {
		weight += v.Weight()
	}
	for _, v := range pao.EquipBag {
		weight += v.Weight()
	}
	weight += float64(len(pao.PotionBag)) * gameconst.PotionGram
	weight += float64(len(pao.ScrollBag)) * gameconst.ScrollGram
	weight += float64(pao.Wallet) * gameconst.MoneyGram
	return weight
}

func (pao PlayerActiveObjInfo) CalcDamageGive() float64 {
	var DamageGive float64
	for _, v := range pao.TurnResult {
		switch v.ResultType {
		case turnresulttype.AttackTo:
			DamageGive += v.Arg
		}
	}
	return DamageGive
}
func (pao PlayerActiveObjInfo) CalcDamageTake() float64 {
	var DamageTake float64
	for _, v := range pao.TurnResult {
		switch v.ResultType {
		case turnresulttype.AttackedFrom:
			DamageTake += v.Arg
		case turnresulttype.DamagedByTile:
			DamageTake += v.Arg
		}
	}
	return DamageTake
}

func (pao PlayerActiveObjInfo) Exist(id string) bool {
	for _, v := range pao.EquippedPo {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	for _, v := range pao.EquipBag {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	for _, v := range pao.PotionBag {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	for _, v := range pao.ScrollBag {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	return false
}
