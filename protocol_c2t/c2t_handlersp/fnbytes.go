// Code generated by "genprotocol.exe -ver=dbfb5de2348191665a2535cbf8ffbfe5d10d4a9b3eea92f9b2ced3abe883c572 -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_handlersp

import (
	"fmt"

	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_json"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
)

// bytes base demux fn map

var DemuxRsp2BytesFnMap = [...]func(me interface{}, hd c2t_packet.Header, rbody []byte) error{
	c2t_idcmd.Invalid:           bytesRecvRspFn_Invalid,           // Invalid make empty packet error
	c2t_idcmd.Login:             bytesRecvRspFn_Login,             // Login
	c2t_idcmd.Heartbeat:         bytesRecvRspFn_Heartbeat,         // Heartbeat
	c2t_idcmd.Chat:              bytesRecvRspFn_Chat,              // Chat
	c2t_idcmd.AchieveInfo:       bytesRecvRspFn_AchieveInfo,       // AchieveInfo
	c2t_idcmd.Rebirth:           bytesRecvRspFn_Rebirth,           // Rebirth
	c2t_idcmd.MoveFloor:         bytesRecvRspFn_MoveFloor,         // MoveFloor tower cmd
	c2t_idcmd.AIPlay:            bytesRecvRspFn_AIPlay,            // AIPlay
	c2t_idcmd.VisitFloorList:    bytesRecvRspFn_VisitFloorList,    // VisitFloorList floor info of visited
	c2t_idcmd.PassTurn:          bytesRecvRspFn_PassTurn,          // PassTurn no action just trigger turn
	c2t_idcmd.Meditate:          bytesRecvRspFn_Meditate,          // Meditate rest and recover HP,SP
	c2t_idcmd.KillSelf:          bytesRecvRspFn_KillSelf,          // KillSelf
	c2t_idcmd.Move:              bytesRecvRspFn_Move,              // Move move 8way near tile
	c2t_idcmd.Attack:            bytesRecvRspFn_Attack,            // Attack attack near 1 tile
	c2t_idcmd.AttackWide:        bytesRecvRspFn_AttackWide,        // AttackWide attack near 3 tile
	c2t_idcmd.AttackLong:        bytesRecvRspFn_AttackLong,        // AttackLong attack 3 tile to direction
	c2t_idcmd.Pickup:            bytesRecvRspFn_Pickup,            // Pickup pickup carryobj
	c2t_idcmd.Drop:              bytesRecvRspFn_Drop,              // Drop drop carryobj
	c2t_idcmd.Equip:             bytesRecvRspFn_Equip,             // Equip equip equipable carryobj
	c2t_idcmd.UnEquip:           bytesRecvRspFn_UnEquip,           // UnEquip unequip equipable carryobj
	c2t_idcmd.DrinkPotion:       bytesRecvRspFn_DrinkPotion,       // DrinkPotion
	c2t_idcmd.ReadScroll:        bytesRecvRspFn_ReadScroll,        // ReadScroll
	c2t_idcmd.Recycle:           bytesRecvRspFn_Recycle,           // Recycle sell carryobj
	c2t_idcmd.EnterPortal:       bytesRecvRspFn_EnterPortal,       // EnterPortal
	c2t_idcmd.ActTeleport:       bytesRecvRspFn_ActTeleport,       // ActTeleport
	c2t_idcmd.AdminTowerCmd:     bytesRecvRspFn_AdminTowerCmd,     // AdminTowerCmd generic cmd
	c2t_idcmd.AdminFloorCmd:     bytesRecvRspFn_AdminFloorCmd,     // AdminFloorCmd generic cmd
	c2t_idcmd.AdminActiveObjCmd: bytesRecvRspFn_AdminActiveObjCmd, // AdminActiveObjCmd generic cmd
	c2t_idcmd.AdminFloorMove:    bytesRecvRspFn_AdminFloorMove,    // AdminFloorMove Next Before floorUUID
	c2t_idcmd.AdminTeleport:     bytesRecvRspFn_AdminTeleport,     // AdminTeleport random pos in floor
	c2t_idcmd.AdminAddExp:       bytesRecvRspFn_AdminAddExp,       // AdminAddExp  add arg to battle exp
	c2t_idcmd.AdminPotionEffect: bytesRecvRspFn_AdminPotionEffect, // AdminPotionEffect buff by arg potion type
	c2t_idcmd.AdminScrollEffect: bytesRecvRspFn_AdminScrollEffect, // AdminScrollEffect buff by arg Scroll type
	c2t_idcmd.AdminCondition:    bytesRecvRspFn_AdminCondition,    // AdminCondition add arg condition for 100 turn
	c2t_idcmd.AdminAddPotion:    bytesRecvRspFn_AdminAddPotion,    // AdminAddPotion add arg potion to inven
	c2t_idcmd.AdminAddScroll:    bytesRecvRspFn_AdminAddScroll,    // AdminAddScroll add arg scroll to inven
	c2t_idcmd.AdminAddMoney:     bytesRecvRspFn_AdminAddMoney,     // AdminAddMoney add arg money to inven
	c2t_idcmd.AdminAddEquip:     bytesRecvRspFn_AdminAddEquip,     // AdminAddEquip add random equip to inven
	c2t_idcmd.AdminForgetFloor:  bytesRecvRspFn_AdminForgetFloor,  // AdminForgetFloor forget current floor map
	c2t_idcmd.AdminFloorMap:     bytesRecvRspFn_AdminFloorMap,     // AdminFloorMap complete current floor map

}

// Invalid make empty packet error
func bytesRecvRspFn_Invalid(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Login
func bytesRecvRspFn_Login(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspLogin_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Heartbeat
func bytesRecvRspFn_Heartbeat(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspHeartbeat_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Chat
func bytesRecvRspFn_Chat(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspChat_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AchieveInfo
func bytesRecvRspFn_AchieveInfo(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAchieveInfo_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Rebirth
func bytesRecvRspFn_Rebirth(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspRebirth_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// MoveFloor tower cmd
func bytesRecvRspFn_MoveFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspMoveFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AIPlay
func bytesRecvRspFn_AIPlay(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAIPlay_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// VisitFloorList floor info of visited
func bytesRecvRspFn_VisitFloorList(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspVisitFloorList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// PassTurn no action just trigger turn
func bytesRecvRspFn_PassTurn(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspPassTurn_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Meditate rest and recover HP,SP
func bytesRecvRspFn_Meditate(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspMeditate_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// KillSelf
func bytesRecvRspFn_KillSelf(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspKillSelf_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Move move 8way near tile
func bytesRecvRspFn_Move(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspMove_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Attack attack near 1 tile
func bytesRecvRspFn_Attack(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAttack_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AttackWide attack near 3 tile
func bytesRecvRspFn_AttackWide(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAttackWide_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AttackLong attack 3 tile to direction
func bytesRecvRspFn_AttackLong(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAttackLong_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Pickup pickup carryobj
func bytesRecvRspFn_Pickup(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspPickup_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Drop drop carryobj
func bytesRecvRspFn_Drop(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspDrop_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Equip equip equipable carryobj
func bytesRecvRspFn_Equip(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspEquip_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// UnEquip unequip equipable carryobj
func bytesRecvRspFn_UnEquip(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspUnEquip_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// DrinkPotion
func bytesRecvRspFn_DrinkPotion(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspDrinkPotion_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// ReadScroll
func bytesRecvRspFn_ReadScroll(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspReadScroll_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Recycle sell carryobj
func bytesRecvRspFn_Recycle(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspRecycle_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// EnterPortal
func bytesRecvRspFn_EnterPortal(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspEnterPortal_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// ActTeleport
func bytesRecvRspFn_ActTeleport(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspActTeleport_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminTowerCmd generic cmd
func bytesRecvRspFn_AdminTowerCmd(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminTowerCmd_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminFloorCmd generic cmd
func bytesRecvRspFn_AdminFloorCmd(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminFloorCmd_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminActiveObjCmd generic cmd
func bytesRecvRspFn_AdminActiveObjCmd(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminActiveObjCmd_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminFloorMove Next Before floorUUID
func bytesRecvRspFn_AdminFloorMove(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminFloorMove_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminTeleport random pos in floor
func bytesRecvRspFn_AdminTeleport(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminTeleport_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminAddExp  add arg to battle exp
func bytesRecvRspFn_AdminAddExp(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminAddExp_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminPotionEffect buff by arg potion type
func bytesRecvRspFn_AdminPotionEffect(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminPotionEffect_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminScrollEffect buff by arg Scroll type
func bytesRecvRspFn_AdminScrollEffect(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminScrollEffect_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminCondition add arg condition for 100 turn
func bytesRecvRspFn_AdminCondition(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminCondition_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminAddPotion add arg potion to inven
func bytesRecvRspFn_AdminAddPotion(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminAddPotion_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminAddScroll add arg scroll to inven
func bytesRecvRspFn_AdminAddScroll(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminAddScroll_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminAddMoney add arg money to inven
func bytesRecvRspFn_AdminAddMoney(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminAddMoney_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminAddEquip add random equip to inven
func bytesRecvRspFn_AdminAddEquip(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminAddEquip_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminForgetFloor forget current floor map
func bytesRecvRspFn_AdminForgetFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminForgetFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// AdminFloorMap complete current floor map
func bytesRecvRspFn_AdminFloorMap(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.RspAdminFloorMap_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}
