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

// Package activeobject manage actor object
package activeobject

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/gamedata"
	"github.com/kasworld/goguelike-single/enum/achievetype"
	"github.com/kasworld/goguelike-single/enum/achievetype_vector"
	"github.com/kasworld/goguelike-single/enum/aotype"
	"github.com/kasworld/goguelike-single/enum/condition_vector"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype_vector"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/potiontype_vector"
	"github.com/kasworld/goguelike-single/enum/respawntype"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/enum/scrolltype_vector"
	"github.com/kasworld/goguelike-single/enum/towerachieve_vector"
	"github.com/kasworld/goguelike-single/enum/turnaction_vector"
	"github.com/kasworld/goguelike-single/game/activeobject/activebuff"
	"github.com/kasworld/goguelike-single/game/activeobject/aoturndata"
	"github.com/kasworld/goguelike-single/game/activeobject/turnresult"
	"github.com/kasworld/goguelike-single/game/aoactreqrsp"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/carryingobject"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/floor4clientman"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/game/inventory"
	"github.com/kasworld/goguelike-single/lib/idu64str"
	"github.com/kasworld/uuidstr"
)

var SysAOIDMaker = idu64str.New("SysAOID")

// var UserAOIDMaker = idu64str.New("UserAOID")

var _ gamei.ActiveObjectI = &ActiveObject{}

func (ao *ActiveObject) String() string {
	return fmt.Sprintf("ActiveObject[%v %v]",
		ao.nickName,
		ao.uuid,
	)
}

// ActiveObject is act by ai or user action
type ActiveObject struct {
	// mutex sync.RWMutex   `prettystring:"hide"`
	rnd *g2rand.G2Rand `prettystring:"hide"`

	uuid        string // aouuid
	nickName    string
	homefloor   gamei.FloorI
	aoType      aotype.ActiveObjType
	respawnType respawntype.RespawnType
	bornFaction factiontype.FactionType `prettystring:"simple"`
	ai          *ServerAIState          // for server side ai
	isAIInUse   bool

	towerAchieveStat *towerachieve_vector.TowerAchieveVector      `prettystring:"simple"`
	achieveStat      achievetype_vector.AchieveTypeVector         `prettystring:"simple"`
	potionStat       potiontype_vector.PotionTypeVector           `prettystring:"simple"`
	scrollStat       scrolltype_vector.ScrollTypeVector           `prettystring:"simple"`
	foActStat        fieldobjacttype_vector.FieldObjActTypeVector `prettystring:"simple"`
	aoActionStat     turnaction_vector.TurnActionVector           `prettystring:"simple"`
	conditionStat    condition_vector.ConditionVector             `prettystring:"simple"`

	// info known to activeobject
	floor4ClientMan *floor4clientman.Floor4ClientMan `prettystring:"simple"`

	currentFloor       gamei.FloorI
	remainTurn2Rebirth int

	chat     string
	chatTime time.Time `prettystring:"simple"`

	ap float64 // action point to use,  -inf ~ 1
	// battle relate
	battleExp    float64
	currentBias  bias.Bias `prettystring:"simple"`
	hp           float64
	sp           float64
	inven        *inventory.Inventory
	buffManager  *activebuff.BuffManager
	expCopy4Sort float64

	// valid in a Turn
	needTANoti     bool // ao move, sight change , floor change, floor aged .. etc
	turnResultList []turnresult.TurnResult
	req2Handle     unsafe.Pointer         // aoactreqrsp.Act
	turnActReqRsp  *aoactreqrsp.ActReqRsp `prettystring:"simple"`

	aoClientCache *csprotocol.ActiveObjClient `prettystring:"simple"`
	AOTurnData    *aoturndata.ActiveObjTurnData
}

func newActiveObj(
	seed int64,
	homefloor gamei.FloorI,
	towerAchieveStat *towerachieve_vector.TowerAchieveVector,
) *ActiveObject {

	ao := &ActiveObject{
		rnd:              g2rand.NewWithSeed(seed),
		towerAchieveStat: towerAchieveStat,
		homefloor:        homefloor,

		// battle
		hp:              100,
		sp:              100,
		inven:           inventory.New(towerAchieveStat),
		buffManager:     activebuff.New(),
		floor4ClientMan: floor4clientman.New(),
		AOTurnData:      &aoturndata.ActiveObjTurnData{},
	}
	ao.bornFaction = factiontype.FactionType(ao.rnd.Intn(factiontype.FactionType_Count))
	ao.currentBias = bias.Bias(ao.bornFaction.FactorBase()).MakeAbsSumTo(gameconst.ActiveObjBaseBiasLen)
	return ao
}

func NewUserActiveObj(seed int64, homefloor gamei.FloorI, nickname string,
	towerAchieveStat *towerachieve_vector.TowerAchieveVector,
) *ActiveObject {

	ao := newActiveObj(seed, homefloor, towerAchieveStat)
	ao.uuid = uuidstr.New()
	// ao.uuid = UserAOIDMaker.New()
	ao.nickName = nickname
	ao.isAIInUse = false
	ao.aoType = aotype.User
	ao.respawnType = respawntype.ToCurrentFloor
	ao.ai = ao.NewServerAI()
	ao.addRandFactionCarryObjEquip(ao.nickName, ao.currentBias.NearFaction(), gameconst.InitCarryObjEquipCount*2)
	ao.addRandPotion(gameconst.InitPotionCount * 2)
	ao.addRandScroll(gameconst.InitScrollCount * 2)
	ao.addInitGold()
	ao.updateActiveObjTurnData()
	return ao
}

func NewSystemActiveObj(seed int64, homefloor gamei.FloorI,
	towerAchieveStat *towerachieve_vector.TowerAchieveVector,
) *ActiveObject {
	ao := newActiveObj(seed, homefloor, towerAchieveStat)
	ao.uuid = SysAOIDMaker.New()
	ao.nickName = gamedata.ActiveObjNameList[ao.rnd.Intn(len(gamedata.ActiveObjNameList))]
	ao.isAIInUse = true
	ao.aoType = aotype.System
	ao.respawnType = respawntype.ToHomeFloor
	ao.ai = ao.NewServerAI()
	ao.addRandFactionCarryObjEquip(ao.nickName, ao.currentBias.NearFaction(), gameconst.InitCarryObjEquipCount)
	ao.addRandPotion(gameconst.InitPotionCount)
	ao.addRandScroll(gameconst.InitScrollCount)
	ao.addInitGold()
	ao.updateActiveObjTurnData()
	return ao
}

func (ao *ActiveObject) Cleanup() {
	ao.isAIInUse = false
}

func (ao *ActiveObject) Rebirth() {
	ao.hp = ao.AOTurnData.HPMax * gameconst.RebirthHPRate
	ao.sp = ao.AOTurnData.SPMax * gameconst.RebirthSPRate
	ao.SetNeedTANoti()
	if ao.aoType == aotype.System {
		eqCount, potionCount, scrollCount := ao.inven.GetTypeCount()
		if eqCount < gameconst.InitCarryObjEquipCount {
			ao.addRandFactionCarryObjEquip(ao.nickName, ao.currentBias.NearFaction(), gameconst.InitCarryObjEquipCount-eqCount)
		}
		if potionCount < gameconst.InitPotionCount {
			ao.addRandPotion(gameconst.InitPotionCount - potionCount)
		}
		if scrollCount < gameconst.InitScrollCount {
			ao.addRandScroll(gameconst.InitScrollCount - scrollCount)
		}
		ao.addInitGold()
	}
	ao.ResetPlan(ao.ai)
}

func (ao *ActiveObject) addInitGold() {
	if ao.inven.GetWalletValue() < gameconst.InitGoldMean {
		v := ao.rnd.NormFloat64Range(gameconst.InitGoldMean, gameconst.InitGoldMean/2)
		if v < 1 {
			v = 1
		}
		ao.inven.AddToWallet(carryingobject.NewMoney(v))
		ao.achieveStat.Add(achievetype.MoneyGet, float64(v))
	}
}

func (ao *ActiveObject) addRandCarryObjEquip(basename string, n int) {
	for ; n > 0; n-- {
		po := carryingobject.NewRandFactionEquipObj(
			basename,
			factiontype.FactionType(ao.rnd.Intn(factiontype.FactionType_Count)),
			ao.rnd)
		ao.inven.AddToBag(po)
	}
}

func (ao *ActiveObject) addRandFactionCarryObjEquip(basename string, ft factiontype.FactionType, n int) {
	for ; n > 0; n-- {
		po := carryingobject.NewRandFactionEquipObj(basename, ft, ao.rnd)
		ao.inven.AddToBag(po)
	}
}
func (ao *ActiveObject) addRandPotion(n int) {
	for ; n > 0; n-- {
		n := ao.rnd.Intn(potiontype.TotalPotionMakeRate)
		po := carryingobject.NewPotionByMakeRate(n)
		ao.inven.AddToBag(po)
	}
}

func (ao *ActiveObject) addRandScroll(n int) {
	for ; n > 0; n-- {
		n := ao.rnd.Intn(scrolltype.TotalScrollMakeRate)
		po := carryingobject.NewScrollByMakeRate(n)
		ao.inven.AddToBag(po)
	}
}
