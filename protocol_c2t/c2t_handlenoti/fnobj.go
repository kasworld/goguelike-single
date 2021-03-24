// Code generated by "genprotocol.exe -ver=fc37e02b6858cffd9591410bf9ff4f28fcf1782014d44a7d0e102918f2b1f57d -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_handlenoti

import (
	"fmt"

	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
)

// obj base demux fn map

var DemuxNoti2ObjFnMap = [...]func(me interface{}, hd c2t_packet.Header, body interface{}) error{
	c2t_idnoti.Invalid:        objRecvNotiFn_Invalid,        // Invalid make empty packet error
	c2t_idnoti.EnterTower:     objRecvNotiFn_EnterTower,     // EnterTower
	c2t_idnoti.EnterFloor:     objRecvNotiFn_EnterFloor,     // EnterFloor
	c2t_idnoti.LeaveFloor:     objRecvNotiFn_LeaveFloor,     // LeaveFloor
	c2t_idnoti.LeaveTower:     objRecvNotiFn_LeaveTower,     // LeaveTower
	c2t_idnoti.Ageing:         objRecvNotiFn_Ageing,         // Ageing          // floor
	c2t_idnoti.Death:          objRecvNotiFn_Death,          // Death
	c2t_idnoti.ReadyToRebirth: objRecvNotiFn_ReadyToRebirth, // ReadyToRebirth
	c2t_idnoti.Rebirthed:      objRecvNotiFn_Rebirthed,      // Rebirthed
	c2t_idnoti.Broadcast:      objRecvNotiFn_Broadcast,      // Broadcast       // global chat broadcast from web admin
	c2t_idnoti.VPObjList:      objRecvNotiFn_VPObjList,      // VPObjList       // in viewport, every turn
	c2t_idnoti.VPTiles:        objRecvNotiFn_VPTiles,        // VPTiles         // in viewport, when viewport changed only
	c2t_idnoti.FloorTiles:     objRecvNotiFn_FloorTiles,     // FloorTiles      // for rebuild known floor
	c2t_idnoti.FieldObjList:   objRecvNotiFn_FieldObjList,   // FieldObjList    // for rebuild known floor
	c2t_idnoti.FoundFieldObj:  objRecvNotiFn_FoundFieldObj,  // FoundFieldObj   // hidden field obj
	c2t_idnoti.ForgetFloor:    objRecvNotiFn_ForgetFloor,    // ForgetFloor
	c2t_idnoti.ActivateTrap:   objRecvNotiFn_ActivateTrap,   // ActivateTrap

}

// Invalid make empty packet error
func objRecvNotiFn_Invalid(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// EnterTower
func objRecvNotiFn_EnterTower(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiEnterTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// EnterFloor
func objRecvNotiFn_EnterFloor(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiEnterFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// LeaveFloor
func objRecvNotiFn_LeaveFloor(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiLeaveFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// LeaveTower
func objRecvNotiFn_LeaveTower(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiLeaveTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Ageing          // floor
func objRecvNotiFn_Ageing(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiAgeing_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Death
func objRecvNotiFn_Death(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiDeath_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// ReadyToRebirth
func objRecvNotiFn_ReadyToRebirth(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiReadyToRebirth_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Rebirthed
func objRecvNotiFn_Rebirthed(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiRebirthed_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Broadcast       // global chat broadcast from web admin
func objRecvNotiFn_Broadcast(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiBroadcast_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// VPObjList       // in viewport, every turn
func objRecvNotiFn_VPObjList(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiVPObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// VPTiles         // in viewport, when viewport changed only
func objRecvNotiFn_VPTiles(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiVPTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// FloorTiles      // for rebuild known floor
func objRecvNotiFn_FloorTiles(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiFloorTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// FieldObjList    // for rebuild known floor
func objRecvNotiFn_FieldObjList(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiFieldObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// FoundFieldObj   // hidden field obj
func objRecvNotiFn_FoundFieldObj(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiFoundFieldObj_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// ForgetFloor
func objRecvNotiFn_ForgetFloor(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiForgetFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// ActivateTrap
func objRecvNotiFn_ActivateTrap(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiActivateTrap_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}
