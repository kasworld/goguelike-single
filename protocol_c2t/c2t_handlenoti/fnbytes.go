// Code generated by "genprotocol.exe -ver=fc37e02b6858cffd9591410bf9ff4f28fcf1782014d44a7d0e102918f2b1f57d -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_handlenoti

import (
	"fmt"

	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_json"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
)

// bytes base demux fn map

var DemuxNoti2ByteFnMap = [...]func(me interface{}, hd c2t_packet.Header, rbody []byte) error{
	c2t_idnoti.Invalid:        bytesRecvNotiFn_Invalid,        // Invalid make empty packet error
	c2t_idnoti.EnterTower:     bytesRecvNotiFn_EnterTower,     // EnterTower
	c2t_idnoti.EnterFloor:     bytesRecvNotiFn_EnterFloor,     // EnterFloor
	c2t_idnoti.LeaveFloor:     bytesRecvNotiFn_LeaveFloor,     // LeaveFloor
	c2t_idnoti.LeaveTower:     bytesRecvNotiFn_LeaveTower,     // LeaveTower
	c2t_idnoti.Ageing:         bytesRecvNotiFn_Ageing,         // Ageing          // floor
	c2t_idnoti.Death:          bytesRecvNotiFn_Death,          // Death
	c2t_idnoti.ReadyToRebirth: bytesRecvNotiFn_ReadyToRebirth, // ReadyToRebirth
	c2t_idnoti.Rebirthed:      bytesRecvNotiFn_Rebirthed,      // Rebirthed
	c2t_idnoti.Broadcast:      bytesRecvNotiFn_Broadcast,      // Broadcast       // global chat broadcast from web admin
	c2t_idnoti.VPObjList:      bytesRecvNotiFn_VPObjList,      // VPObjList       // in viewport, every turn
	c2t_idnoti.VPTiles:        bytesRecvNotiFn_VPTiles,        // VPTiles         // in viewport, when viewport changed only
	c2t_idnoti.FloorTiles:     bytesRecvNotiFn_FloorTiles,     // FloorTiles      // for rebuild known floor
	c2t_idnoti.FieldObjList:   bytesRecvNotiFn_FieldObjList,   // FieldObjList    // for rebuild known floor
	c2t_idnoti.FoundFieldObj:  bytesRecvNotiFn_FoundFieldObj,  // FoundFieldObj   // hidden field obj
	c2t_idnoti.ForgetFloor:    bytesRecvNotiFn_ForgetFloor,    // ForgetFloor
	c2t_idnoti.ActivateTrap:   bytesRecvNotiFn_ActivateTrap,   // ActivateTrap

}

// Invalid make empty packet error
func bytesRecvNotiFn_Invalid(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// EnterTower
func bytesRecvNotiFn_EnterTower(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiEnterTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// EnterFloor
func bytesRecvNotiFn_EnterFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiEnterFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// LeaveFloor
func bytesRecvNotiFn_LeaveFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiLeaveFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// LeaveTower
func bytesRecvNotiFn_LeaveTower(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiLeaveTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Ageing          // floor
func bytesRecvNotiFn_Ageing(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiAgeing_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Death
func bytesRecvNotiFn_Death(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiDeath_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// ReadyToRebirth
func bytesRecvNotiFn_ReadyToRebirth(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiReadyToRebirth_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Rebirthed
func bytesRecvNotiFn_Rebirthed(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiRebirthed_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// Broadcast       // global chat broadcast from web admin
func bytesRecvNotiFn_Broadcast(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiBroadcast_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// VPObjList       // in viewport, every turn
func bytesRecvNotiFn_VPObjList(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiVPObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// VPTiles         // in viewport, when viewport changed only
func bytesRecvNotiFn_VPTiles(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiVPTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// FloorTiles      // for rebuild known floor
func bytesRecvNotiFn_FloorTiles(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiFloorTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// FieldObjList    // for rebuild known floor
func bytesRecvNotiFn_FieldObjList(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiFieldObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// FoundFieldObj   // hidden field obj
func bytesRecvNotiFn_FoundFieldObj(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiFoundFieldObj_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// ForgetFloor
func bytesRecvNotiFn_ForgetFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiForgetFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

// ActivateTrap
func bytesRecvNotiFn_ActivateTrap(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiActivateTrap_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}
