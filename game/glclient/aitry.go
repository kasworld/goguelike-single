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

package glclient

import (
	"github.com/kasworld/goguelike-single/config/leveldata"
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/equipslottype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
)

var tryAutoActFn = []func(app *GLClient) bool{
	tryAutoBattle,
	tryAutoPickup,
	tryAutoEquip,
	tryAutoUsePotion,
	tryAutoRecyclePotionScroll,
	tryAutoRecycleEquip,
}

func (app *GLClient) actByControlMode() {
	if app.OLNotiData != nil &&
		app.OLNotiData.ActiveObj.AP > 0 &&
		app.OLNotiData.ActiveObj.HP > 0 {
		for _, v := range tryAutoActFn {
			if v(app) {
				return
			}
		}
	}
	app.sendPacket(c2t_idcmd.PassTurn,
		&c2t_obj.ReqPassTurn_data{},
	)
}

func tryAutoBattle(app *GLClient) bool {
	cf := app.CurrentFloor
	if app.OLNotiData == nil {
		return false
	}
	playerX, playerY := app.GetPlayerXY()
	if !cf.IsValidPos(playerX, playerY) {
		return false
	}
	if !cf.Tiles[playerX][playerY].CanBattle() {
		return false
	}
	w, h := cf.Tiles.GetXYLen()
	for _, ao := range app.OLNotiData.ActiveObjList {
		if !ao.Alive {
			continue
		}
		if ao.UUID == app.AccountInfo.ActiveObjUUID {
			continue
		}
		if !cf.IsValidPos(ao.X, ao.Y) {
			continue
		}
		if !cf.Tiles[ao.X][ao.Y].CanBattle() {
			continue
		}
		isContact, dir := way9type.CalcContactDirWrappedXY(
			playerX, playerY, ao.X, ao.Y, w, h)
		if isContact && dir != way9type.Center {
			app.ReqWithRspFnWithAuth(c2t_idcmd.Attack,
				&c2t_obj.ReqAttack_data{Dir: dir},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func tryAutoPickup(app *GLClient) bool {
	cf := app.CurrentFloor
	if app.OLNotiData == nil {
		return false
	}
	if app.OLNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}

	playerX, playerY := app.GetPlayerXY()
	w, h := cf.Tiles.GetXYLen()
	for _, po := range app.OLNotiData.CarryObjList {
		isContact, dir := way9type.CalcContactDirWrappedXY(
			playerX, playerY, po.X, po.Y, w, h)
		if !isContact {
			continue
		}
		if dir == way9type.Center {
			app.ReqWithRspFnWithAuth(c2t_idcmd.Pickup,
				&c2t_obj.ReqPickup_data{UUID: po.UUID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		} else {
			app.ReqWithRspFnWithAuth(c2t_idcmd.Move,
				&c2t_obj.ReqMove_data{Dir: dir},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func tryAutoEquip(app *GLClient) bool {
	if app.OLNotiData == nil {
		return false
	}
	for _, po := range app.OLNotiData.ActiveObj.EquippedPo {
		if app.needUnEquipCarryObj(po.GetBias()) {
			app.ReqWithRspFnWithAuth(c2t_idcmd.UnEquip,
				&c2t_obj.ReqUnEquip_data{UUID: po.UUID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	for _, po := range app.OLNotiData.ActiveObj.EquipBag {
		if app.isBetterCarryObj(po.EquipType, po.GetBias()) {
			app.ReqWithRspFnWithAuth(c2t_idcmd.Equip,
				&c2t_obj.ReqEquip_data{UUID: po.UUID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func tryAutoUsePotion(app *GLClient) bool {
	if app.OLNotiData == nil {
		return false
	}
	for _, po := range app.OLNotiData.ActiveObj.PotionBag {
		if app.needUsePotion(po) {
			app.ReqWithRspFnWithAuth(c2t_idcmd.DrinkPotion,
				&c2t_obj.ReqDrinkPotion_data{UUID: po.UUID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}

	for _, po := range app.OLNotiData.ActiveObj.ScrollBag {
		if app.needUseScroll(po) {
			app.ReqWithRspFnWithAuth(c2t_idcmd.ReadScroll,
				&c2t_obj.ReqReadScroll_data{UUID: po.UUID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}

	return false
}

func tryAutoRecycleEquip(app *GLClient) bool {
	if app.OLNotiData == nil {
		return false
	}
	if app.OLNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}
	if app.onFieldObj == nil {
		return false
	}
	if app.onFieldObj.ActType != fieldobjacttype.RecycleCarryObj {
		return false
	}
	return app.recycleEqbag()
}

func tryAutoRecyclePotionScroll(app *GLClient) bool {
	if app.OLNotiData == nil {
		return false
	}
	if app.OLNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}
	if app.onFieldObj == nil {
		return false
	}
	if app.onFieldObj.ActType != fieldobjacttype.RecycleCarryObj {
		return false
	}
	if app.recycleUselessPotion() {
		return true
	}
	if app.recycleUselessScroll() {
		return true
	}
	return false
}

/////////

func (app *GLClient) isBetterCarryObj(EquipType equipslottype.EquipSlotType, PoBias bias.Bias) bool {
	aoEnvBias := app.TowerBias().Add(app.CurrentFloor.GetBias()).Add(app.OLNotiData.ActiveObj.Bias)
	newBiasAbs := aoEnvBias.Add(PoBias).AbsSum()
	for _, v := range app.OLNotiData.ActiveObj.EquippedPo {
		if v.EquipType == EquipType {
			return newBiasAbs > aoEnvBias.Add(v.GetBias()).AbsSum()+1
		}
	}
	return newBiasAbs > aoEnvBias.AbsSum()+1
}

func (app *GLClient) needUnEquipCarryObj(PoBias bias.Bias) bool {
	aoEnvBias := app.TowerBias().Add(app.CurrentFloor.GetBias()).Add(app.OLNotiData.ActiveObj.Bias)

	currentBias := aoEnvBias.Add(PoBias)
	newBias := aoEnvBias
	return newBias.AbsSum() > currentBias.AbsSum()+1
}

func (app *GLClient) needUseScroll(po *c2t_obj.ScrollClient) bool {
	cf := app.CurrentFloor
	switch po.ScrollType {
	case scrolltype.FloorMap:
		if cf.Visited.CalcCompleteRate() < 1.0 {
			return true
		}
	}
	return false
}

func (app *GLClient) needUsePotion(po *c2t_obj.PotionClient) bool {
	pao := app.OLNotiData.ActiveObj
	switch po.PotionType {
	case potiontype.RecoverHP10:
		return pao.HPMax-pao.HP > 10
	case potiontype.RecoverHP50:
		return pao.HPMax-pao.HP > 50
	case potiontype.RecoverHP100:
		return pao.HPMax-pao.HP > 100

	case potiontype.RecoverSP10:
		return pao.SPMax-pao.SP > 10
	case potiontype.RecoverSP50:
		return pao.SPMax-pao.SP > 50
	case potiontype.RecoverSP100:
		return pao.SPMax-pao.SP > 100

	case potiontype.RecoverHPRate10:
		return pao.HPMax-pao.HP > pao.HPMax/10
	case potiontype.RecoverHPRate50:
		return pao.HPMax/2 > pao.HP
	case potiontype.RecoverHPFull:
		return pao.HPMax/10 > pao.HP

	case potiontype.RecoverSPRate10:
		return pao.SPMax-pao.SP > pao.SPMax/10
	case potiontype.RecoverSPRate50:
		return pao.SPMax/2 > pao.SP
	case potiontype.RecoverSPFull:
		return pao.SPMax/10 > pao.SP

	case potiontype.BuffRecoverHP1:
		return pao.HPMax/2 > pao.HP
	case potiontype.BuffRecoverSP1:
		return pao.SPMax/2 > pao.SP

	case potiontype.BuffSight1:
		return pao.Sight <= leveldata.Sight(app.level)
	case potiontype.BuffSight5:
		return pao.Sight <= leveldata.Sight(app.level)
	case potiontype.BuffSightMax:
		return pao.Sight <= leveldata.Sight(app.level)
	}
	return false
}

func (app *GLClient) recycleEqbag() bool {
	var poList c2t_obj.CarryObjEqByLen
	poList = append(poList, app.OLNotiData.ActiveObj.EquipBag...)
	if len(poList) == 0 {
		return false
	}
	poList.Sort()
	app.ReqWithRspFnWithAuth(c2t_idcmd.Recycle,
		&c2t_obj.ReqRecycle_data{UUID: poList[0].UUID},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		})
	return true
}

func (app *GLClient) recycleUselessPotion() bool {
	for _, po := range app.OLNotiData.ActiveObj.PotionBag {
		if potiontype.AIRecycleMap[po.PotionType] {
			app.ReqWithRspFnWithAuth(c2t_idcmd.Recycle,
				&c2t_obj.ReqRecycle_data{UUID: po.UUID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func (app *GLClient) recycleUselessScroll() bool {
	for _, po := range app.OLNotiData.ActiveObj.ScrollBag {
		if scrolltype.AIRecycleMap[po.ScrollType] {
			app.ReqWithRspFnWithAuth(c2t_idcmd.Recycle,
				&c2t_obj.ReqRecycle_data{UUID: po.UUID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}
