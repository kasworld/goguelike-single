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

	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
)

func (tw *Tower) handle_c2tch() {
	for rpk := range tw.c2tCh {
		shd, sbody, serr := tw.handleRecvPacket(rpk)
		// process result
		if serr != nil {
			// handle error
		}
		spk := &c2t_packet.Packet{
			Header: shd,
			Body:   sbody,
		}
		tw.t2cCh <- spk

	}
}

func (tw *Tower) handleRecvPacket(rpk *c2t_packet.Packet) (c2t_packet.Header, interface{}, error) {
	switch body := rpk.Body.(type) {
	default:
		return c2t_packet.Header{}, nil, fmt.Errorf("invalid packet")
	case *c2t_obj.ReqInvalid_data:
		return tw.objAPIFn_ReqInvalid(rpk.Header, body)
	case *c2t_obj.ReqLogin_data:
		return tw.objAPIFn_ReqLogin(rpk.Header, body)
	case *c2t_obj.ReqHeartbeat_data:
		return tw.objAPIFn_ReqHeartbeat(rpk.Header, body)
	case *c2t_obj.ReqChat_data:
		return tw.objAPIFn_ReqChat(rpk.Header, body)
	case *c2t_obj.ReqAchieveInfo_data:
		return tw.objAPIFn_ReqAchieveInfo(rpk.Header, body)
	case *c2t_obj.ReqRebirth_data:
		return tw.objAPIFn_ReqRebirth(rpk.Header, body)
	case *c2t_obj.ReqMoveFloor_data:
		return tw.objAPIFn_ReqMoveFloor(rpk.Header, body)
	case *c2t_obj.ReqAIPlay_data:
		return tw.objAPIFn_ReqAIPlay(rpk.Header, body)
	case *c2t_obj.ReqVisitFloorList_data:
		return tw.objAPIFn_ReqVisitFloorList(rpk.Header, body)
	case *c2t_obj.ReqPassTurn_data:
		return tw.objAPIFn_ReqPassTurn(rpk.Header, body)
	case *c2t_obj.ReqMeditate_data:
		return tw.objAPIFn_ReqMeditate(rpk.Header, body)
	case *c2t_obj.ReqKillSelf_data:
		return tw.objAPIFn_ReqKillSelf(rpk.Header, body)
	case *c2t_obj.ReqMove_data:
		return tw.objAPIFn_ReqMove(rpk.Header, body)
	case *c2t_obj.ReqAttack_data:
		return tw.objAPIFn_ReqAttack(rpk.Header, body)
	case *c2t_obj.ReqAttackWide_data:
		return tw.objAPIFn_ReqAttackWide(rpk.Header, body)
	case *c2t_obj.ReqAttackLong_data:
		return tw.objAPIFn_ReqAttackLong(rpk.Header, body)
	case *c2t_obj.ReqPickup_data:
		return tw.objAPIFn_ReqPickup(rpk.Header, body)
	case *c2t_obj.ReqDrop_data:
		return tw.objAPIFn_ReqDrop(rpk.Header, body)
	case *c2t_obj.ReqEquip_data:
		return tw.objAPIFn_ReqEquip(rpk.Header, body)
	case *c2t_obj.ReqUnEquip_data:
		return tw.objAPIFn_ReqUnEquip(rpk.Header, body)
	case *c2t_obj.ReqDrinkPotion_data:
		return tw.objAPIFn_ReqDrinkPotion(rpk.Header, body)
	case *c2t_obj.ReqReadScroll_data:
		return tw.objAPIFn_ReqReadScroll(rpk.Header, body)
	case *c2t_obj.ReqRecycle_data:
		return tw.objAPIFn_ReqRecycle(rpk.Header, body)
	case *c2t_obj.ReqEnterPortal_data:
		return tw.objAPIFn_ReqEnterPortal(rpk.Header, body)
	case *c2t_obj.ReqActTeleport_data:
		return tw.objAPIFn_ReqActTeleport(rpk.Header, body)
	case *c2t_obj.ReqAdminTowerCmd_data:
		return tw.objAPIFn_ReqAdminTowerCmd(rpk.Header, body)
	case *c2t_obj.ReqAdminFloorCmd_data:
		return tw.objAPIFn_ReqAdminFloorCmd(rpk.Header, body)
	case *c2t_obj.ReqAdminActiveObjCmd_data:
		return tw.objAPIFn_ReqAdminActiveObjCmd(rpk.Header, body)
	case *c2t_obj.ReqAdminFloorMove_data:
		return tw.objAPIFn_ReqAdminFloorMove(rpk.Header, body)
	case *c2t_obj.ReqAdminTeleport_data:
		return tw.objAPIFn_ReqAdminTeleport(rpk.Header, body)
	case *c2t_obj.ReqAdminAddExp_data:
		return tw.objAPIFn_ReqAdminAddExp(rpk.Header, body)
	case *c2t_obj.ReqAdminPotionEffect_data:
		return tw.objAPIFn_ReqAdminPotionEffect(rpk.Header, body)
	case *c2t_obj.ReqAdminScrollEffect_data:
		return tw.objAPIFn_ReqAdminScrollEffect(rpk.Header, body)
	case *c2t_obj.ReqAdminCondition_data:
		return tw.objAPIFn_ReqAdminCondition(rpk.Header, body)
	case *c2t_obj.ReqAdminAddPotion_data:
		return tw.objAPIFn_ReqAdminAddPotion(rpk.Header, body)
	case *c2t_obj.ReqAdminAddScroll_data:
		return tw.objAPIFn_ReqAdminAddScroll(rpk.Header, body)
	case *c2t_obj.ReqAdminAddMoney_data:
		return tw.objAPIFn_ReqAdminAddMoney(rpk.Header, body)
	case *c2t_obj.ReqAdminAddEquip_data:
		return tw.objAPIFn_ReqAdminAddEquip(rpk.Header, body)
	case *c2t_obj.ReqAdminForgetFloor_data:
		return tw.objAPIFn_ReqAdminForgetFloor(rpk.Header, body)
	case *c2t_obj.ReqAdminFloorMap_data:
		return tw.objAPIFn_ReqAdminFloorMap(rpk.Header, body)
	}
}

// Invalid make empty packet error
func (tw *Tower) objAPIFn_ReqInvalid(hd c2t_packet.Header, robj *c2t_obj.ReqInvalid_data) (
	c2t_packet.Header, *c2t_obj.RspInvalid_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspInvalid_data{}
	return sendHeader, sendBody, nil
}

// Login
func (tw *Tower) objAPIFn_ReqLogin(hd c2t_packet.Header, robj *c2t_obj.ReqLogin_data) (
	c2t_packet.Header, *c2t_obj.RspLogin_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspLogin_data{}
	return sendHeader, sendBody, nil
}

// Heartbeat
func (tw *Tower) objAPIFn_ReqHeartbeat(hd c2t_packet.Header, robj *c2t_obj.ReqHeartbeat_data) (
	c2t_packet.Header, *c2t_obj.RspHeartbeat_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspHeartbeat_data{}
	return sendHeader, sendBody, nil
}

// Chat
func (tw *Tower) objAPIFn_ReqChat(hd c2t_packet.Header, robj *c2t_obj.ReqChat_data) (
	c2t_packet.Header, *c2t_obj.RspChat_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspChat_data{}
	return sendHeader, sendBody, nil
}

// AchieveInfo
func (tw *Tower) objAPIFn_ReqAchieveInfo(hd c2t_packet.Header, robj *c2t_obj.ReqAchieveInfo_data) (
	c2t_packet.Header, *c2t_obj.RspAchieveInfo_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAchieveInfo_data{}
	return sendHeader, sendBody, nil
}

// Rebirth
func (tw *Tower) objAPIFn_ReqRebirth(hd c2t_packet.Header, robj *c2t_obj.ReqRebirth_data) (
	c2t_packet.Header, *c2t_obj.RspRebirth_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspRebirth_data{}
	return sendHeader, sendBody, nil
}

// MoveFloor tower cmd
func (tw *Tower) objAPIFn_ReqMoveFloor(hd c2t_packet.Header, robj *c2t_obj.ReqMoveFloor_data) (
	c2t_packet.Header, *c2t_obj.RspMoveFloor_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspMoveFloor_data{}
	return sendHeader, sendBody, nil
}

// AIPlay
func (tw *Tower) objAPIFn_ReqAIPlay(hd c2t_packet.Header, robj *c2t_obj.ReqAIPlay_data) (
	c2t_packet.Header, *c2t_obj.RspAIPlay_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAIPlay_data{}
	return sendHeader, sendBody, nil
}

// VisitFloorList floor info of visited
func (tw *Tower) objAPIFn_ReqVisitFloorList(hd c2t_packet.Header, robj *c2t_obj.ReqVisitFloorList_data) (
	c2t_packet.Header, *c2t_obj.RspVisitFloorList_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspVisitFloorList_data{}
	return sendHeader, sendBody, nil
}

// PassTurn no action just trigger turn
func (tw *Tower) objAPIFn_ReqPassTurn(hd c2t_packet.Header, robj *c2t_obj.ReqPassTurn_data) (
	c2t_packet.Header, *c2t_obj.RspPassTurn_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspPassTurn_data{}
	return sendHeader, sendBody, nil
}

// Meditate rest and recover HP,SP
func (tw *Tower) objAPIFn_ReqMeditate(hd c2t_packet.Header, robj *c2t_obj.ReqMeditate_data) (
	c2t_packet.Header, *c2t_obj.RspMeditate_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspMeditate_data{}
	return sendHeader, sendBody, nil
}

// KillSelf
func (tw *Tower) objAPIFn_ReqKillSelf(hd c2t_packet.Header, robj *c2t_obj.ReqKillSelf_data) (
	c2t_packet.Header, *c2t_obj.RspKillSelf_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspKillSelf_data{}
	return sendHeader, sendBody, nil
}

// Move move 8way near tile
func (tw *Tower) objAPIFn_ReqMove(hd c2t_packet.Header, robj *c2t_obj.ReqMove_data) (
	c2t_packet.Header, *c2t_obj.RspMove_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspMove_data{}
	return sendHeader, sendBody, nil
}

// Attack attack near 1 tile
func (tw *Tower) objAPIFn_ReqAttack(hd c2t_packet.Header, robj *c2t_obj.ReqAttack_data) (
	c2t_packet.Header, *c2t_obj.RspAttack_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAttack_data{}
	return sendHeader, sendBody, nil
}

// AttackWide attack near 3 tile
func (tw *Tower) objAPIFn_ReqAttackWide(hd c2t_packet.Header, robj *c2t_obj.ReqAttackWide_data) (
	c2t_packet.Header, *c2t_obj.RspAttackWide_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAttackWide_data{}
	return sendHeader, sendBody, nil
}

// AttackLong attack 3 tile to direction
func (tw *Tower) objAPIFn_ReqAttackLong(hd c2t_packet.Header, robj *c2t_obj.ReqAttackLong_data) (
	c2t_packet.Header, *c2t_obj.RspAttackLong_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAttackLong_data{}
	return sendHeader, sendBody, nil
}

// Pickup pickup carryobj
func (tw *Tower) objAPIFn_ReqPickup(hd c2t_packet.Header, robj *c2t_obj.ReqPickup_data) (
	c2t_packet.Header, *c2t_obj.RspPickup_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspPickup_data{}
	return sendHeader, sendBody, nil
}

// Drop drop carryobj
func (tw *Tower) objAPIFn_ReqDrop(hd c2t_packet.Header, robj *c2t_obj.ReqDrop_data) (
	c2t_packet.Header, *c2t_obj.RspDrop_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspDrop_data{}
	return sendHeader, sendBody, nil
}

// Equip equip equipable carryobj
func (tw *Tower) objAPIFn_ReqEquip(hd c2t_packet.Header, robj *c2t_obj.ReqEquip_data) (
	c2t_packet.Header, *c2t_obj.RspEquip_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspEquip_data{}
	return sendHeader, sendBody, nil
}

// UnEquip unequip equipable carryobj
func (tw *Tower) objAPIFn_ReqUnEquip(hd c2t_packet.Header, robj *c2t_obj.ReqUnEquip_data) (
	c2t_packet.Header, *c2t_obj.RspUnEquip_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspUnEquip_data{}
	return sendHeader, sendBody, nil
}

// DrinkPotion
func (tw *Tower) objAPIFn_ReqDrinkPotion(hd c2t_packet.Header, robj *c2t_obj.ReqDrinkPotion_data) (
	c2t_packet.Header, *c2t_obj.RspDrinkPotion_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspDrinkPotion_data{}
	return sendHeader, sendBody, nil
}

// ReadScroll
func (tw *Tower) objAPIFn_ReqReadScroll(hd c2t_packet.Header, robj *c2t_obj.ReqReadScroll_data) (
	c2t_packet.Header, *c2t_obj.RspReadScroll_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspReadScroll_data{}
	return sendHeader, sendBody, nil
}

// Recycle sell carryobj
func (tw *Tower) objAPIFn_ReqRecycle(hd c2t_packet.Header, robj *c2t_obj.ReqRecycle_data) (
	c2t_packet.Header, *c2t_obj.RspRecycle_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspRecycle_data{}
	return sendHeader, sendBody, nil
}

// EnterPortal
func (tw *Tower) objAPIFn_ReqEnterPortal(hd c2t_packet.Header, robj *c2t_obj.ReqEnterPortal_data) (
	c2t_packet.Header, *c2t_obj.RspEnterPortal_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspEnterPortal_data{}
	return sendHeader, sendBody, nil
}

// ActTeleport
func (tw *Tower) objAPIFn_ReqActTeleport(hd c2t_packet.Header, robj *c2t_obj.ReqActTeleport_data) (
	c2t_packet.Header, *c2t_obj.RspActTeleport_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspActTeleport_data{}
	return sendHeader, sendBody, nil
}

// AdminTowerCmd generic cmd
func (tw *Tower) objAPIFn_ReqAdminTowerCmd(hd c2t_packet.Header, robj *c2t_obj.ReqAdminTowerCmd_data) (
	c2t_packet.Header, *c2t_obj.RspAdminTowerCmd_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminTowerCmd_data{}
	return sendHeader, sendBody, nil
}

// AdminFloorCmd generic cmd
func (tw *Tower) objAPIFn_ReqAdminFloorCmd(hd c2t_packet.Header, robj *c2t_obj.ReqAdminFloorCmd_data) (
	c2t_packet.Header, *c2t_obj.RspAdminFloorCmd_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminFloorCmd_data{}
	return sendHeader, sendBody, nil
}

// AdminActiveObjCmd generic cmd
func (tw *Tower) objAPIFn_ReqAdminActiveObjCmd(hd c2t_packet.Header, robj *c2t_obj.ReqAdminActiveObjCmd_data) (
	c2t_packet.Header, *c2t_obj.RspAdminActiveObjCmd_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminActiveObjCmd_data{}
	return sendHeader, sendBody, nil
}

// AdminFloorMove Next Before floorUUID
func (tw *Tower) objAPIFn_ReqAdminFloorMove(hd c2t_packet.Header, robj *c2t_obj.ReqAdminFloorMove_data) (
	c2t_packet.Header, *c2t_obj.RspAdminFloorMove_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminFloorMove_data{}
	return sendHeader, sendBody, nil
}

// AdminTeleport random pos in floor
func (tw *Tower) objAPIFn_ReqAdminTeleport(hd c2t_packet.Header, robj *c2t_obj.ReqAdminTeleport_data) (
	c2t_packet.Header, *c2t_obj.RspAdminTeleport_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminTeleport_data{}
	return sendHeader, sendBody, nil
}

// AdminAddExp  add arg to battle exp
func (tw *Tower) objAPIFn_ReqAdminAddExp(hd c2t_packet.Header, robj *c2t_obj.ReqAdminAddExp_data) (
	c2t_packet.Header, *c2t_obj.RspAdminAddExp_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddExp_data{}
	return sendHeader, sendBody, nil
}

// AdminPotionEffect buff by arg potion type
func (tw *Tower) objAPIFn_ReqAdminPotionEffect(hd c2t_packet.Header, robj *c2t_obj.ReqAdminPotionEffect_data) (
	c2t_packet.Header, *c2t_obj.RspAdminPotionEffect_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminPotionEffect_data{}
	return sendHeader, sendBody, nil
}

// AdminScrollEffect buff by arg Scroll type
func (tw *Tower) objAPIFn_ReqAdminScrollEffect(hd c2t_packet.Header, robj *c2t_obj.ReqAdminScrollEffect_data) (
	c2t_packet.Header, *c2t_obj.RspAdminScrollEffect_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminScrollEffect_data{}
	return sendHeader, sendBody, nil
}

// AdminCondition add arg condition for 100 turn
func (tw *Tower) objAPIFn_ReqAdminCondition(hd c2t_packet.Header, robj *c2t_obj.ReqAdminCondition_data) (
	c2t_packet.Header, *c2t_obj.RspAdminCondition_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminCondition_data{}
	return sendHeader, sendBody, nil
}

// AdminAddPotion add arg potion to inven
func (tw *Tower) objAPIFn_ReqAdminAddPotion(hd c2t_packet.Header, robj *c2t_obj.ReqAdminAddPotion_data) (
	c2t_packet.Header, *c2t_obj.RspAdminAddPotion_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddPotion_data{}
	return sendHeader, sendBody, nil
}

// AdminAddScroll add arg scroll to inven
func (tw *Tower) objAPIFn_ReqAdminAddScroll(hd c2t_packet.Header, robj *c2t_obj.ReqAdminAddScroll_data) (
	c2t_packet.Header, *c2t_obj.RspAdminAddScroll_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddScroll_data{}
	return sendHeader, sendBody, nil
}

// AdminAddMoney add arg money to inven
func (tw *Tower) objAPIFn_ReqAdminAddMoney(hd c2t_packet.Header, robj *c2t_obj.ReqAdminAddMoney_data) (
	c2t_packet.Header, *c2t_obj.RspAdminAddMoney_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddMoney_data{}
	return sendHeader, sendBody, nil
}

// AdminAddEquip add random equip to inven
func (tw *Tower) objAPIFn_ReqAdminAddEquip(hd c2t_packet.Header, robj *c2t_obj.ReqAdminAddEquip_data) (
	c2t_packet.Header, *c2t_obj.RspAdminAddEquip_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddEquip_data{}
	return sendHeader, sendBody, nil
}

// AdminForgetFloor forget current floor map
func (tw *Tower) objAPIFn_ReqAdminForgetFloor(hd c2t_packet.Header, robj *c2t_obj.ReqAdminForgetFloor_data) (
	c2t_packet.Header, *c2t_obj.RspAdminForgetFloor_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminForgetFloor_data{}
	return sendHeader, sendBody, nil
}

// AdminFloorMap complete current floor map
func (tw *Tower) objAPIFn_ReqAdminFloorMap(hd c2t_packet.Header, robj *c2t_obj.ReqAdminFloorMap_data) (
	c2t_packet.Header, *c2t_obj.RspAdminFloorMap_data, error) {
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminFloorMap_data{}
	return sendHeader, sendBody, nil
}
