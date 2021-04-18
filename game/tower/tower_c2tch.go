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

package tower

import (
	"fmt"
	"time"

	"github.com/kasworld/goguelike-single/enum/achievetype"
	"github.com/kasworld/goguelike-single/enum/flowtype"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/returncode"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/enum/statusoptype"
	"github.com/kasworld/goguelike-single/enum/turnaction"
	"github.com/kasworld/goguelike-single/game/aoactreqrsp"
	"github.com/kasworld/goguelike-single/game/carryingobject"
	"github.com/kasworld/goguelike-single/game/cmd2floor"
	"github.com/kasworld/goguelike-single/game/cmd2tower"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

func (tw *Tower) SendNoti(body interface{}) {
	spk := csprotocol.Packet{
		FlowType: flowtype.Notification,
		Body:     body,
	}
	tw.t2cCh <- &spk
}

func (tw *Tower) triggerTurnByCmd(cmd turnaction.TurnAction) {
	if cmd.TriggerTurn() {
		tw.GetTurnCh() <- time.Now()
	}
}

// handle recv req obj
func (tw *Tower) handle_c2tch() {
	for rpk := range tw.c2tCh {

		if rpk.FlowType != flowtype.Request {
			g2log.Error("Unexpected rpk.Header packet type: %v", rpk)
		}

		// call api
		rcode, sbody, apierr := tw.handleRecvReqObj(rpk)

		if apierr != nil {
			g2log.Fatal("%v", apierr)
		}
		if sbody == nil {
			g2log.Fatal("Response body nil")
		}

		spk := &csprotocol.Packet{
			FlowType:   flowtype.Response,
			PacketID:   rpk.PacketID,
			ReturnCode: rcode,
			Body:       sbody,
		}

		// send rsp
		tw.t2cCh <- spk
	}
}

func (tw *Tower) handleRecvReqObj(rpk *csprotocol.Packet) (returncode.ReturnCode, interface{}, error) {
	switch body := rpk.Body.(type) {
	default:
		return returncode.Success, nil, fmt.Errorf("invalid packet")

	case *csprotocol.ReqAchieveInfo:
		return tw.objAPIFn_ReqAchieveInfo(body)
	case *csprotocol.ReqAIPlay:
		return tw.objAPIFn_ReqAIPlay(body)
	case *csprotocol.ReqVisitFloorList:
		return tw.objAPIFn_ReqVisitFloorList(body)

	case *csprotocol.ReqTurnAction:
		return tw.objAPIFn_ReqTurnAction(body)

	case *csprotocol.ReqAdminFloorMove:
		return tw.objAPIFn_ReqAdminFloorMove(body)
	case *csprotocol.ReqAdminTeleport:
		return tw.objAPIFn_ReqAdminTeleport(body)
	case *csprotocol.ReqAdminAddExp:
		return tw.objAPIFn_ReqAdminAddExp(body)
	case *csprotocol.ReqAdminPotionEffect:
		return tw.objAPIFn_ReqAdminPotionEffect(body)
	case *csprotocol.ReqAdminScrollEffect:
		return tw.objAPIFn_ReqAdminScrollEffect(body)
	case *csprotocol.ReqAdminCondition:
		return tw.objAPIFn_ReqAdminCondition(body)
	case *csprotocol.ReqAdminAddPotion:
		return tw.objAPIFn_ReqAdminAddPotion(body)
	case *csprotocol.ReqAdminAddScroll:
		return tw.objAPIFn_ReqAdminAddScroll(body)
	case *csprotocol.ReqAdminAddMoney:
		return tw.objAPIFn_ReqAdminAddMoney(body)
	case *csprotocol.ReqAdminAddEquip:
		return tw.objAPIFn_ReqAdminAddEquip(body)
	case *csprotocol.ReqAdminForgetFloor:
		return tw.objAPIFn_ReqAdminForgetFloor(body)
	case *csprotocol.ReqAdminFloorMap:
		return tw.objAPIFn_ReqAdminFloorMap(body)
	}
}

// AchieveInfo
func (tw *Tower) objAPIFn_ReqAchieveInfo(robj *csprotocol.ReqAchieveInfo) (
	returncode.ReturnCode, *csprotocol.RspAchieveInfo, error) {

	sendBody := &csprotocol.RspAchieveInfo{
		AchieveStat:   *tw.playerAO.GetAchieveStat(),
		PotionStat:    *tw.playerAO.GetPotionStat(),
		ScrollStat:    *tw.playerAO.GetScrollStat(),
		FOActStat:     *tw.playerAO.GetFieldObjActStat(),
		ConditionStat: *tw.playerAO.GetConditionStat(),
	}
	return returncode.Success, sendBody, nil
}

// AIPlay
func (tw *Tower) objAPIFn_ReqAIPlay(robj *csprotocol.ReqAIPlay) (
	returncode.ReturnCode, *csprotocol.RspAIPlay, error) {

	sendBody := &csprotocol.RspAIPlay{}
	if err := tw.playerAO.DoAIOnOff(robj.On); err != nil {
		g2log.Error("fail to AIOn")
	}

	return returncode.Success, sendBody, nil
}

// VisitFloorList floor info of visited
func (tw *Tower) objAPIFn_ReqVisitFloorList(robj *csprotocol.ReqVisitFloorList) (
	returncode.ReturnCode, *csprotocol.RspVisitFloorList, error) {

	floorList := make([]*csprotocol.FloorInfo, 0)
	for _, f4c := range tw.playerAO.GetFloor4ClientList() {
		f := tw.floorMan.GetFloorByName(f4c.GetName())
		fi := f.ToPacket_FloorInfo()
		fi.VisitCount = f4c.Visit.GetDiscoveredTileCount()
		floorList = append(floorList, fi)
	}
	return returncode.Success, &csprotocol.RspVisitFloorList{
		FloorList: floorList,
	}, nil
}

// TurnAction no action just trigger turn
func (tw *Tower) objAPIFn_ReqTurnAction(robj *csprotocol.ReqTurnAction) (
	returncode.ReturnCode, *csprotocol.RspTurnAction, error) {

	sendBody := &csprotocol.RspTurnAction{}
	defer tw.triggerTurnByCmd(robj.Act)

	ec := returncode.Success
	switch robj.Act {
	default:
		tw.playerAO.SetReq2Handle(&aoactreqrsp.Act{
			Act:  robj.Act,
			Dir:  robj.Dir,
			UUID: robj.UUID,
		})
	case turnaction.Rebirth:
		if err := tw.playerAO.TryRebirth(); err != nil {
			ec = returncode.ActionProhibited
			g2log.Error("%v", err)
		}
	case turnaction.MoveFloor:
		tw.GetCmdCh() <- &cmd2tower.FloorMove{
			ActiveObj: tw.playerAO,
			FloorName: robj.UUID,
		}
	case turnaction.PassTurn:

	}

	return ec, sendBody, nil
}

// AdminFloorMove Next Before floorUUID
func (tw *Tower) objAPIFn_ReqAdminFloorMove(robj *csprotocol.ReqAdminFloorMove) (
	returncode.ReturnCode, *csprotocol.RspAdminFloorMove, error) {

	rspCh := make(chan returncode.ReturnCode, 1)
	tw.GetCmdCh() <- &cmd2tower.AdminFloorMove{
		ActiveObj:  tw.playerAO,
		RecvPacket: robj,
		RspCh:      rspCh,
	}
	ec := <-rspCh
	return ec, &csprotocol.RspAdminFloorMove{}, nil
}

// AdminTeleport random pos in floor
func (tw *Tower) objAPIFn_ReqAdminTeleport(robj *csprotocol.ReqAdminTeleport) (
	returncode.ReturnCode, *csprotocol.RspAdminTeleport, error) {

	f := tw.playerAO.GetCurrentFloor()
	if f == nil {
		return returncode.Success, nil, fmt.Errorf("user not in floor")
	}
	rspCh := make(chan returncode.ReturnCode, 1)
	f.GetCmdCh() <- &cmd2floor.APIAdminTeleport2Floor{
		ActiveObj: tw.playerAO,
		ReqPk:     robj,
		RspCh:     rspCh,
	}
	ec := <-rspCh
	return ec, &csprotocol.RspAdminTeleport{}, nil
}

// AdminAddExp  add arg to battle exp
func (tw *Tower) objAPIFn_ReqAdminAddExp(robj *csprotocol.ReqAdminAddExp) (
	returncode.ReturnCode, *csprotocol.RspAdminAddExp, error) {

	sendBody := &csprotocol.RspAdminAddExp{}
	tw.playerAO.AddBattleExp(float64(robj.Exp))
	return returncode.Success, sendBody, nil
}

// AdminPotionEffect buff by arg potion type
func (tw *Tower) objAPIFn_ReqAdminPotionEffect(robj *csprotocol.ReqAdminPotionEffect) (
	returncode.ReturnCode, *csprotocol.RspAdminPotionEffect, error) {

	sendBody := &csprotocol.RspAdminPotionEffect{}
	tw.playerAO.GetBuffManager().Add(robj.Potion.String(), false, false,
		potiontype.GetBuffByPotionType(robj.Potion),
	)
	return returncode.Success, sendBody, nil
}

// AdminScrollEffect buff by arg Scroll type
func (tw *Tower) objAPIFn_ReqAdminScrollEffect(robj *csprotocol.ReqAdminScrollEffect) (
	returncode.ReturnCode, *csprotocol.RspAdminScrollEffect, error) {

	sendBody := &csprotocol.RspAdminScrollEffect{}
	tw.playerAO.GetBuffManager().Add(robj.Scroll.String(), false, false,
		scrolltype.GetBuffByScrollType(robj.Scroll),
	)
	return returncode.Success, sendBody, nil
}

// AdminCondition add arg condition for 100 turn
func (tw *Tower) objAPIFn_ReqAdminCondition(robj *csprotocol.ReqAdminCondition) (
	returncode.ReturnCode, *csprotocol.RspAdminCondition, error) {

	sendBody := &csprotocol.RspAdminCondition{}
	buff2add := statusoptype.Repeat(100,
		statusoptype.OpArg{statusoptype.SetCondition, robj.Condition},
	)
	tw.playerAO.GetBuffManager().Add(
		"admin"+robj.Condition.String(),
		true, true, buff2add)

	return returncode.Success, sendBody, nil
}

// AdminAddPotion add arg potion to inven
func (tw *Tower) objAPIFn_ReqAdminAddPotion(robj *csprotocol.ReqAdminAddPotion) (
	returncode.ReturnCode, *csprotocol.RspAdminAddPotion, error) {

	sendBody := &csprotocol.RspAdminAddPotion{}
	pt := carryingobject.NewPotion(robj.Potion)
	tw.playerAO.GetInven().AddToBag(pt)
	return returncode.Success, sendBody, nil
}

// AdminAddScroll add arg scroll to inven
func (tw *Tower) objAPIFn_ReqAdminAddScroll(robj *csprotocol.ReqAdminAddScroll) (
	returncode.ReturnCode, *csprotocol.RspAdminAddScroll, error) {

	sendBody := &csprotocol.RspAdminAddScroll{}
	pt := carryingobject.NewScroll(robj.Scroll)
	tw.playerAO.GetInven().AddToBag(pt)
	return returncode.Success, sendBody, nil
}

// AdminAddMoney add arg money to inven
func (tw *Tower) objAPIFn_ReqAdminAddMoney(robj *csprotocol.ReqAdminAddMoney) (
	returncode.ReturnCode, *csprotocol.RspAdminAddMoney, error) {

	sendBody := &csprotocol.RspAdminAddMoney{}
	tw.playerAO.GetInven().AddToWallet(carryingobject.NewMoney(float64(robj.Money)))
	tw.playerAO.GetAchieveStat().Add(achievetype.MoneyGet, float64(robj.Money))
	return returncode.Success, sendBody, nil
}

// AdminAddEquip add random equip to inven
func (tw *Tower) objAPIFn_ReqAdminAddEquip(robj *csprotocol.ReqAdminAddEquip) (
	returncode.ReturnCode, *csprotocol.RspAdminAddEquip, error) {

	sendBody := &csprotocol.RspAdminAddEquip{}
	eq := carryingobject.NewEquipByFactionSlot("admin",
		robj.Faction, robj.Equip,
		tw.rnd,
	)
	tw.playerAO.GetInven().AddToBag(eq)
	return returncode.Success, sendBody, nil
}

// AdminForgetFloor forget current floor map
func (tw *Tower) objAPIFn_ReqAdminForgetFloor(robj *csprotocol.ReqAdminForgetFloor) (
	returncode.ReturnCode, *csprotocol.RspAdminForgetFloor, error) {

	sendBody := &csprotocol.RspAdminForgetFloor{}
	f := tw.playerAO.GetCurrentFloor()
	if f == nil {
		return returncode.Success, nil, fmt.Errorf("user not in floor")
	}
	if err := tw.playerAO.ForgetFloorByName(f.GetName()); err != nil {
		g2log.Error("%v", err)
	}
	return returncode.Success, sendBody, nil
}

// AdminFloorMap complete current floor map
func (tw *Tower) objAPIFn_ReqAdminFloorMap(robj *csprotocol.ReqAdminFloorMap) (
	returncode.ReturnCode, *csprotocol.RspAdminFloorMap, error) {

	sendBody := &csprotocol.RspAdminFloorMap{}
	f := tw.playerAO.GetCurrentFloor()
	if f == nil {
		return returncode.Success, nil, fmt.Errorf("user not in floor")
	}
	if err := tw.playerAO.MakeFloorComplete(f); err != nil {
		g2log.Error("%v", err)
	}
	return returncode.Success, sendBody, nil
}
