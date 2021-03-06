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

	"github.com/kasworld/goguelike-single/config/leveldata"
	"github.com/kasworld/goguelike-single/enum/achievetype"
	"github.com/kasworld/goguelike-single/enum/aotype"
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/game/activeobject/aoturndata"
	"github.com/kasworld/goguelike-single/game/activeobject/turnresult"
	"github.com/kasworld/goguelike-single/game/aoactreqrsp"
	"github.com/kasworld/goguelike-single/game/cmd2tower"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

// PrepareNewTurn before turn start
func (ao *ActiveObject) PrepareNewTurn(TurnCount int) {
	ao.turnResultList = ao.turnResultList[:0] // for ao to ao damage record
	ao.turnActReqRsp = nil
	ao.aoClientCache = nil
	if TurnCount-ao.chatOldTurnCount > 10 {
		ao.chat = ""
	}
	ao.achieveStat.Inc(achievetype.Turn)
}

// SetTurnActReqRsp set turn act result
func (ao *ActiveObject) SetTurnActReqRsp(actrsp *aoactreqrsp.ActReqRsp) {
	if needAP := actrsp.Req.CalcAPByActAndCondition(ao.AOTurnData.Condition); needAP > 0 {
		ao.ap += -needAP
	}
	ao.turnActReqRsp = actrsp
	if actrsp.IsSuccess() {
		ao.aoActionStat.Inc(actrsp.Done.Act)
	}
}

func (ao *ActiveObject) updateActiveObjTurnData() {
	old := ao.AOTurnData
	ao.AOTurnData = &aoturndata.ActiveObjTurnData{}

	if ao.IsAlive() {
		oparglist := ao.buffManager.GetOpArgListToApply()
		for _, v := range oparglist {
			ao.applyOpArg(old, ao.AOTurnData, v)
		}
	}

	nonBattleExp := 0
	for _, v := range ao.floor4ClientMan.GetList() {
		nonBattleExp += v.Visit.GetDiscoverExp()
	}
	ao.AOTurnData.NonBattleExp = float64(nonBattleExp)
	ao.AOTurnData.TotalExp = ao.AOTurnData.NonBattleExp + ao.battleExp
	ao.achieveStat.SetIfGt(achievetype.MaxExp, ao.AOTurnData.TotalExp)
	ao.AOTurnData.Level = leveldata.CalcLevelFromExp(ao.AOTurnData.TotalExp)

	// sight buff applied
	ao.AOTurnData.Sight += leveldata.Sight(int(ao.AOTurnData.Level))
	if ao.AOTurnData.Condition.TestByCondition(condition.Blind) {
		ao.AOTurnData.Sight = 0
	}
	if ao.AOTurnData.Sight <= 0 { // if sight debuf make <=0
		ao.AOTurnData.Condition.SetByCondition(condition.Blind)
	}
	EquipAttackBias := ao.inven.SumEquipAttackBias()
	ao.AOTurnData.AttackBias = ao.currentBias.Add(EquipAttackBias)
	ao.AOTurnData.SPMax = leveldata.MaxSP(int(ao.AOTurnData.Level)) +
		ao.AOTurnData.AttackBias.AbsSum()

	EquipDefenceBias := ao.inven.SumEquipDefenceBias()
	ao.AOTurnData.DefenceBias = ao.currentBias.Add(EquipDefenceBias)
	ao.AOTurnData.HPMax = leveldata.MaxHP(int(ao.AOTurnData.Level)) +
		ao.AOTurnData.DefenceBias.AbsSum()

	totalWeight := ao.inven.GetTotalWeight()
	if ao.AOTurnData.Condition.TestByCondition(condition.Burden) {
		totalWeight *= 2
	}
	if ao.AOTurnData.Condition.TestByCondition(condition.Float) {
		totalWeight /= 2
	}
	ao.AOTurnData.LoadRate = float64(totalWeight) / leveldata.WeightLimit(int(ao.AOTurnData.Level))
	if ao.AOTurnData.Sight != old.Sight {
		ao.SetNeedTANoti()
	}
}

// ApplyTurnAct at all act end in turn
// apply turn result and prepare next turn info to send
// can die ao
func (ao *ActiveObject) ApplyTurnAct() {
	ao.updateActiveObjTurnData()
	intLv := int(ao.AOTurnData.Level)
	if ao.IsAlive() {
		hpLvMax := leveldata.MaxHP(intLv)
		apLvMax := leveldata.MaxSP(intLv)
		if ao.AOTurnData.LoadRate > 1 {
			// overload penalty
			ao.hp += -ao.AOTurnData.LoadRate
			ao.sp += -ao.AOTurnData.LoadRate
		} else {
			// add no act/ no interaction hp/sp recover bonus
			if len(ao.turnResultList) == 0 && ao.ap > 0 {
				ao.hp += hpLvMax / 100
				ao.sp += apLvMax / 100
			}
		}
		ao.ap++
	}
	if ao.hp > ao.AOTurnData.HPMax {
		ao.hp = ao.AOTurnData.HPMax
	}
	if ao.sp > ao.AOTurnData.SPMax {
		ao.sp = ao.AOTurnData.SPMax
	}
	if ao.ap > leveldata.MaxAP(intLv) { // limit max ap
		ao.ap = leveldata.MaxAP(intLv)
	}

	if ao.remainTurn2Rebirth > 0 {
		ao.remainTurn2Rebirth--
		if ao.remainTurn2Rebirth == 0 {
			if ao.aoType == aotype.User {
				ao.homefloor.GetTower().SendNoti(
					&csprotocol.NotiReadyToRebirth{},
				)
			} else {
				err := ao.TryRebirth()
				if err != nil {
					g2log.Error("%v", err)
				}
			}
			// send noti
			ao.ResetPlan(ao.ai)
		}
	}
}

func (ao *ActiveObject) TryRebirth() error {
	if ao.remainTurn2Rebirth == 0 && !ao.IsAlive() {
		ao.buffManager.ClearOnRebirth()
		ao.homefloor.GetTower().GetCmdCh() <- &cmd2tower.ActiveObjRebirth{
			ActiveObj: ao,
		}
		return nil
	}
	return fmt.Errorf("no need rebirth")
}

func (ao *ActiveObject) AppendTurnResult(turnResult turnresult.TurnResult) {
	ao.turnResultList = append(ao.turnResultList, turnResult)
}

////////////////////////////////////////////////////////////////////////////////
// process packet

func (ao *ActiveObject) DoPickup(co gamei.CarryingObjectI) error {
	var err error
	switch po := co.(type) {
	default:
		err = fmt.Errorf("unknown type %T %v", po, po)
	case gamei.MoneyI:
		err = ao.GetInven().AddToWallet(po.(gamei.MoneyI))
		if err == nil {
			ao.GetAchieveStat().Add(achievetype.MoneyGet, float64(po.GetValue()))
		}
	case gamei.EquipObjI:
		err = ao.GetInven().AddToBag(po)
	case gamei.PotionI:
		err = ao.GetInven().AddToBag(po)
	case gamei.ScrollI:
		err = ao.GetInven().AddToBag(po)
	}
	if err == nil {
		ao.GetAchieveStat().Inc(achievetype.PickupCarryObj)
	}
	return err
}

func (ao *ActiveObject) DoEquip(poid string) error {
	err := ao.GetInven().EquipFromBagByUUID(poid)
	if err == nil {
		ao.achieveStat.Inc(achievetype.EquipCarryObj)
		return nil
	} else {
		return fmt.Errorf("not in inventory %v %v", ao, err)
	}
}

func (ao *ActiveObject) DoUnEquip(poid string) error {
	_, err := ao.GetInven().UnEquipToBagByUUID(poid)
	if err == nil {
		ao.achieveStat.Inc(achievetype.UnEquipCarryObj)
		return nil
	} else {
		return fmt.Errorf("not in inventory %v", ao)
	}
}

func (ao *ActiveObject) DoUseCarryObj(poid string) error {
	po := ao.GetInven().RemoveByUUID(poid)
	if po == nil {
		return fmt.Errorf("not in inventory %v %v", ao, poid)
	}
	switch o := po.(type) {
	default:
		return fmt.Errorf("invalid objtype %v", po)
	case gamei.PotionI:
		ao.achieveStat.Inc(achievetype.UseCarryObj)
		ao.potionStat.Inc(o.GetPotionType())
		tb := potiontype.GetBuffByPotionType(o.GetPotionType())
		if tb != nil { // potion data exist
			ao.buffManager.Add(o.GetPotionType().String(), false, false, tb)
			return nil
		}
		return fmt.Errorf("no potiondata %v %v", ao, o)
	case gamei.ScrollI:
		ao.achieveStat.Inc(achievetype.UseCarryObj)
		ao.scrollStat.Inc(o.GetScrollType())
		tb := scrolltype.GetBuffByScrollType(o.GetScrollType())
		if tb != nil { // scroll data exist
			ao.buffManager.Add(o.GetScrollType().String(), false, false, tb)
			return nil
		}
		switch o.GetScrollType() {
		default:
			g2log.Fatal("unknown scrolltype %v", po)
		case scrolltype.Empty:
		case scrolltype.FloorMap:
			return ao.MakeFloorComplete(ao.currentFloor)
		case scrolltype.Teleport:
			g2log.Fatal("Scroll_Teleport must processed in floor %v", ao)
		}
		return nil
	}
}

func (ao *ActiveObject) DoRecycleCarryObj(poid string) error {
	v, err := ao.GetInven().RecycleCarryObjByID(poid)
	if err == nil {
		ao.achieveStat.Add(achievetype.MoneyGet, float64(v))
		ao.foActStat.Inc(fieldobjacttype.RecycleCarryObj)
		return nil
	} else {
		return fmt.Errorf("not in inventory %v %v", ao, err)
	}
}

func (ao *ActiveObject) DoAIOnOff(onoff bool) error {
	if ao.aoType == aotype.User {
		ao.SetUseAI(onoff)
	} else {
		return fmt.Errorf("not user ao %v", ao)
	}
	return nil
}
