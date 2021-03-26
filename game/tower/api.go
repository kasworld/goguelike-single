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
	"strings"
	"time"

	"github.com/kasworld/goguelike-single/config/authdata"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/game/activeobject"
	"github.com/kasworld/goguelike-single/game/cmd2tower"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/conndata"
	"github.com/kasworld/goguelike-single/lib/session"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_serveconnbyte"
	"github.com/kasworld/uuidstr"
	"github.com/kasworld/version"
)

func (tw *Tower) bytesAPIFn_ReqInvalid(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	rhd := c2t_packet.Header{}
	spacket := &c2t_obj.RspInvalid_data{}
	return rhd, spacket, fmt.Errorf("invalid packet")
}

func (tw *Tower) bytesAPIFn_ReqLogin(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	c2sc, ok := me.(*c2t_serveconnbyte.ServeConnByte)
	if !ok {
		panic(fmt.Sprintf("invalid me not c2t_serveconnbyte.ServeConnByte %#v", me))
	}

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqLogin_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}

	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}

	if err := authdata.UpdateByAuthKey(c2sc.GetAuthorCmdList(), robj.AuthKey); err != nil {
		return rhd, nil, err
	}

	connData := c2sc.GetConnData().(*conndata.ConnData)

	now := time.Now()

	// if reconnect
	if tw.playerSession != nil {
		// update exist session
		tw.playerSession.ConnUUID = connData.UUID
		tw.playerSession.RemoteAddr = connData.RemoteAddr
		tw.playerSession.LastUse = now
		oldAO, exist := tw.id2aoSuspend.GetByUUID(tw.playerSession.ActiveObjUUID)
		if !exist {
			panic("no exist ao")
		}
		oldAO.Resume(c2sc)
		rspCh := make(chan error, 1)
		tw.GetReqCh() <- &cmd2tower.ActiveObjResumeTower{
			ActiveObj: oldAO,
			RspCh:     rspCh,
		}
		err = <-rspCh
	} else {
		// new session
		tw.playerSession = &session.Session{
			SessionUUID:   uuidstr.New(),
			ConnUUID:      connData.UUID,
			RemoteAddr:    connData.RemoteAddr,
			Create:        now,
			LastUse:       now,
			NickName:      tw.Config().NickName,
			ActiveObjUUID: "",
		}
		// new ao
		var homeFloor gamei.FloorI
		homeFloor = tw.GetFloorManager().GetStartFloor()
		newAO := activeobject.NewUserActiveObj(
			tw.rnd.Int63(),
			homeFloor,
			tw.Config().NickName,
			tw.log,
			tw.towerAchieveStat,
			c2sc)
		tw.playerSession.ActiveObjUUID = newAO.GetUUID()
		rspCh := make(chan error, 1)
		tw.GetReqCh() <- &cmd2tower.ActiveObjEnterTower{
			ActiveObj: newAO,
			RspCh:     rspCh,
		}
		err = <-rspCh
	}

	connData.Session = tw.playerSession

	if err != nil {
		return rhd, nil, err
	} else {
		acinfo := &c2t_obj.AccountInfo{
			SessionUUID:   tw.playerSession.GetUUID(),
			ActiveObjUUID: tw.playerSession.ActiveObjUUID,
			NickName:      tw.playerSession.NickName,
			CmdList:       *c2sc.GetAuthorCmdList(),
		}
		return rhd, &c2t_obj.RspLogin_data{
			ServiceInfo: tw.serviceInfo,
			AccountInfo: acinfo,
		}, nil
	}
}

func (tw *Tower) bytesAPIFn_ReqHeartbeat(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqHeartbeat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}

	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	spacket := &c2t_obj.RspHeartbeat_data{
		Time: robj.Time,
	}

	return rhd, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqChat(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqChat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	robj.Chat = strings.TrimSpace(robj.Chat)
	if !version.IsRelease() && len(robj.Chat) > 1 && robj.Chat[0] == '/' {
		return c2t_packet.Header{
			ErrorCode: AdminCmd(ao, robj.Chat[1:]),
		}, &c2t_obj.RspChat_data{}, nil
	} else {
		if len(robj.Chat) > gameconst.MaxChatLen {
			robj.Chat = robj.Chat[:gameconst.MaxChatLen]
		}
		ao.SetChat(robj.Chat)
		return rhd, &c2t_obj.RspChat_data{}, nil
	}
}

func (tw *Tower) bytesAPIFn_ReqAchieveInfo(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	spacket := &c2t_obj.RspAchieveInfo_data{
		AchieveStat:   *ao.GetAchieveStat(),
		PotionStat:    *ao.GetPotionStat(),
		ScrollStat:    *ao.GetScrollStat(),
		FOActStat:     *ao.GetFieldObjActStat(),
		AOActionStat:  *ao.GetActStat(),
		ConditionStat: *ao.GetConditionStat(),
	}

	// for i, v := range ao.GetAchieveStat() {
	// 	spacket.Achieve[i] = int64(v)
	// }
	return rhd, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqRebirth(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	spacket := &c2t_obj.RspRebirth_data{}

	err = ao.TryRebirth()
	if err != nil {
		rhd.ErrorCode = c2t_error.ActionProhibited
		tw.log.Error("%v", err)
	}

	return rhd, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqMoveFloor(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqMoveFloor_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspMoveFloor_data{}

	tw.GetReqCh() <- &cmd2tower.FloorMove{
		ActiveObj: ao,
		FloorName: robj.UUID,
	}
	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqAIPlay(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAIPlay_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	if err := ao.DoAIOnOff(robj.On); err != nil {
		tw.log.Error("fail to AIOn %v %v", me)
	}
	return rhd, &c2t_obj.RspAIPlay_data{}, nil
}

func (tw *Tower) bytesAPIFn_ReqVisitFloorList(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	// r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// robj, ok := r.(*c2t_obj.ReqVisitFloorList_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	// }
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	floorList := make([]*c2t_obj.FloorInfo, 0)
	for _, f4c := range ao.GetFloor4ClientList() {
		f := tw.floorMan.GetFloorByName(f4c.GetName())
		fi := f.ToPacket_FloorInfo()
		fi.VisitCount = f4c.Visit.GetDiscoveredTileCount()
		floorList = append(floorList, fi)
	}
	return rhd, &c2t_obj.RspVisitFloorList_data{
		FloorList: floorList,
	}, nil
}
