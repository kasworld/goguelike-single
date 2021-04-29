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

package gamei

import (
	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/enum/achievetype_vector"
	"github.com/kasworld/goguelike-single/enum/aotype"
	"github.com/kasworld/goguelike-single/enum/condition_vector"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype_vector"
	"github.com/kasworld/goguelike-single/enum/potiontype_vector"
	"github.com/kasworld/goguelike-single/enum/respawntype"
	"github.com/kasworld/goguelike-single/enum/scrolltype_vector"
	"github.com/kasworld/goguelike-single/enum/turnaction_vector"
	"github.com/kasworld/goguelike-single/game/activeobject/activebuff"
	"github.com/kasworld/goguelike-single/game/activeobject/aoturndata"
	"github.com/kasworld/goguelike-single/game/activeobject/turnresult"
	"github.com/kasworld/goguelike-single/game/aoactreqrsp"
	"github.com/kasworld/goguelike-single/game/aoscore"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/floor4client"
)

type ActiveObjectI interface {
	Cleanup()
	GetUUID() string
	// GetNickName() string
	String() string

	GetInven() InventoryI

	GetHP() float64
	// GetSP() float64
	GetHPRate() float64
	GetSPRate() float64

	IsAlive() bool

	Rebirth()
	EnterFloor(f FloorI)
	Death(f FloorI)

	GetBias() bias.Bias
	// GetBornFaction() factiontype.FactionType

	AddBattleExp(v float64)

	// for sort
	GetExpCopy() float64
	UpdateExpCopy()

	NeedCharge(limit float64) bool
	Charged(limit float64) bool
	ReduceSP(apToReduce float64) (apReduced float64)
	ReduceHP(hpToReduce float64) (hpReduced float64)

	GetTurnResultList() []turnresult.TurnResult
	SetTurnActReqRsp(actrsp *aoactreqrsp.ActReqRsp)
	GetTurnActReqRsp() *aoactreqrsp.ActReqRsp
	SetNeedTANoti()
	GetAndClearNeedTANoti() bool
	GetAP() float64

	SetReq2Handle(req *aoactreqrsp.Act)
	GetClearReq2Handle() *aoactreqrsp.Act

	GetTurnData() *aoturndata.ActiveObjTurnData
	GetBuffManager() *activebuff.BuffManager

	GetActiveObjType() aotype.ActiveObjType
	GetRespawnType() respawntype.RespawnType

	IsAIUse() bool
	SetUseAI(b bool)
	RunAI(TurnCount int)

	GetChat() string
	SetChat(c string, TurnCount int)

	GetRemainTurn2Rebirth() int
	TryRebirth() error

	GetCurrentFloor() FloorI

	PrepareNewTurn(TurnCount int)
	ApplyTurnAct()
	AppendTurnResult(turnResult turnresult.TurnResult)

	ApplyDamageFromDangerObj() bool
	ApplyHPSPDecByActOnTile(hp, sp float64)
	Kill(dst ActiveObjectI)

	DoEquip(poid string) error
	DoUnEquip(poid string) error
	DoUseCarryObj(poid string) error
	DoRecycleCarryObj(poid string) error
	DoAIOnOff(onoff bool) error
	DoPickup(po CarryingObjectI) error

	// for tower
	GetHomeFloor() FloorI

	ToPacket_ActiveObjClient(x, y int) *csprotocol.ActiveObjClient
	ToPacket_PlayerActiveObjInfo() *csprotocol.PlayerActiveObjInfo
	To_ActiveObjScore(int) *aoscore.ActiveObjScore

	GetAchieveStat() *achievetype_vector.AchieveTypeVector
	GetFieldObjActStat() *fieldobjacttype_vector.FieldObjActTypeVector
	GetPotionStat() *potiontype_vector.PotionTypeVector
	GetScrollStat() *scrolltype_vector.ScrollTypeVector
	GetActStat() *turnaction_vector.TurnActionVector
	GetConditionStat() *condition_vector.ConditionVector

	UpdateVisitAreaBySightMat2(f FloorI, vpCenterX, vpCenterY int,
		sightMat *viewportdata.ViewportSight2, sight float32)

	GetFloor4ClientList() []*floor4client.Floor4Client
	GetFloor4Client(floorname string) *floor4client.Floor4Client
	ForgetFloorByName(floorname string) error
	MakeFloorComplete(f FloorI) error
}

const (
	ActiveObjectI_HTML_tableheader = `
	<tr>
	<td>ActiveObj</td>
	<td>bias</td>
	<td>level</td>
	<td>exp</td>
	<td>AI</td>
	</tr>	
`
	ActiveObjectI_HTML_row = `
	<tr>
		<td>
		<a href= "/ActiveObj?aoid={{$v.GetUUID}}" >
			{{$v}}
		</a>
		</td>
		<td>
		{{$v.GetBias}} 
		</td>
		<td>
		{{$v.GetTurnData.Level | printf "%4.2f" }} 
		</td>
		<td>
		{{$v.GetTurnData.TotalExp | printf "%4.2f" }} 
		</td>
		<td>
		{{if $v.GetAIObj }}
			{{$v.GetAIObj}}
		{{end}}
		</td>
	</tr>
`
)
