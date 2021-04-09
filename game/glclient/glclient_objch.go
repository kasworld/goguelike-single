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
	"fmt"

	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_pid2rspfn"
)

func (app *GLClient) sendReqObjWithRspFn(cmd c2t_idcmd.CommandID, body interface{},
	fn c2t_pid2rspfn.HandleRspFn) error {
	pid := app.pid2recv.NewPID(fn)
	spk := c2t_packet.Packet{
		Header: c2t_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: c2t_packet.Request,
		},
		Body: body,
	}
	app.c2tCh <- &spk
	return nil
}

func (app *GLClient) handle_t2ch() {
	for rpk := range app.t2cCh {
		g2log.TraceClient("recv %v", rpk.Header)
		switch rpk.Header.FlowType {
		default:
			g2log.Fatal("invalid packet type %v %v", rpk.Header, rpk.Body)
		case c2t_packet.Response:
			if err := app.pid2recv.HandleRsp(rpk.Header, rpk.Body); err != nil {
				g2log.Fatal("%v %v %v %v", app, rpk.Header, rpk.Body, err)
				return
			}
		case c2t_packet.Notification:
			err := app.handleRecvNotiObj(rpk)
			// process result
			if err != nil {
				g2log.Fatal("%v %v %v %v", app, rpk.Header, rpk.Body, err)
				return
			}
		}
	}
}

func (app *GLClient) handleRecvNotiObj(rpk *c2t_packet.Packet) error {
	switch body := rpk.Body.(type) {
	default:
		return fmt.Errorf("invalid packet")

	case *c2t_obj.NotiInvalid_data:
		return app.objRecvNotiFn_Invalid(rpk.Header, body)
	case *c2t_obj.NotiEnterTower_data:
		return app.objRecvNotiFn_EnterTower(rpk.Header, body)
	case *c2t_obj.NotiEnterFloor_data:
		return app.objRecvNotiFn_EnterFloor(rpk.Header, body)
	case *c2t_obj.NotiLeaveFloor_data:
		return app.objRecvNotiFn_LeaveFloor(rpk.Header, body)
	case *c2t_obj.NotiLeaveTower_data:
		return app.objRecvNotiFn_LeaveTower(rpk.Header, body)
	case *c2t_obj.NotiAgeing_data:
		return app.objRecvNotiFn_Ageing(rpk.Header, body)
	case *c2t_obj.NotiDeath_data:
		return app.objRecvNotiFn_Death(rpk.Header, body)
	case *c2t_obj.NotiReadyToRebirth_data:
		return app.objRecvNotiFn_ReadyToRebirth(rpk.Header, body)
	case *c2t_obj.NotiRebirthed_data:
		return app.objRecvNotiFn_Rebirthed(rpk.Header, body)
	case *c2t_obj.NotiBroadcast_data:
		return app.objRecvNotiFn_Broadcast(rpk.Header, body)
	case *c2t_obj.NotiVPObjList_data:
		return app.objRecvNotiFn_VPObjList(rpk.Header, body)
	case *c2t_obj.NotiVPTiles_data:
		return app.objRecvNotiFn_VPTiles(rpk.Header, body)
	case *c2t_obj.NotiFloorTiles_data:
		return app.objRecvNotiFn_FloorTiles(rpk.Header, body)
	case *c2t_obj.NotiFieldObjList_data:
		return app.objRecvNotiFn_FieldObjList(rpk.Header, body)
	case *c2t_obj.NotiFoundFieldObj_data:
		return app.objRecvNotiFn_FoundFieldObj(rpk.Header, body)
	case *c2t_obj.NotiForgetFloor_data:
		return app.objRecvNotiFn_ForgetFloor(rpk.Header, body)
	case *c2t_obj.NotiActivateTrap_data:
		return app.objRecvNotiFn_ActivateTrap(rpk.Header, body)

	}
}

func (app *GLClient) objRecvNotiFn_Invalid(hd c2t_packet.Header, body *c2t_obj.NotiInvalid_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_EnterTower(hd c2t_packet.Header, body *c2t_obj.NotiEnterTower_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_EnterFloor(hd c2t_packet.Header, body *c2t_obj.NotiEnterFloor_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_LeaveFloor(hd c2t_packet.Header, body *c2t_obj.NotiLeaveFloor_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_LeaveTower(hd c2t_packet.Header, body *c2t_obj.NotiLeaveTower_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_Ageing(hd c2t_packet.Header, body *c2t_obj.NotiAgeing_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_Death(hd c2t_packet.Header, body *c2t_obj.NotiDeath_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_ReadyToRebirth(hd c2t_packet.Header, body *c2t_obj.NotiReadyToRebirth_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_Rebirthed(hd c2t_packet.Header, body *c2t_obj.NotiRebirthed_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_Broadcast(hd c2t_packet.Header, body *c2t_obj.NotiBroadcast_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_VPObjList(hd c2t_packet.Header, body *c2t_obj.NotiVPObjList_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_VPTiles(hd c2t_packet.Header, body *c2t_obj.NotiVPTiles_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_FloorTiles(hd c2t_packet.Header, body *c2t_obj.NotiFloorTiles_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_FieldObjList(hd c2t_packet.Header, body *c2t_obj.NotiFieldObjList_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_FoundFieldObj(hd c2t_packet.Header, body *c2t_obj.NotiFoundFieldObj_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_ForgetFloor(hd c2t_packet.Header, body *c2t_obj.NotiForgetFloor_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
func (app *GLClient) objRecvNotiFn_ActivateTrap(hd c2t_packet.Header, body *c2t_obj.NotiActivateTrap_data) error {
	return fmt.Errorf("Not implemented %v", body)
}
