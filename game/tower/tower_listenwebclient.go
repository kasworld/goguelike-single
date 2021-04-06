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
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/game/aoexpsort"
	"github.com/kasworld/goguelike-single/game/aoscore"
	"github.com/kasworld/goguelike-single/game/cmd2tower"
	"github.com/kasworld/goguelike-single/lib/conndata"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_authorize"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_serveconnbyte"
	"github.com/kasworld/weblib"
)

func (tw *Tower) initServiceWeb(ctx context.Context) {
	webMux := http.NewServeMux()
	webMux.Handle("/",
		http.FileServer(http.Dir(tw.Config().ClientDataFolder)),
	)

	webMux.HandleFunc("/highscore.json", tw.json_HighScore)

	webMux.HandleFunc("/TowerInfo", tw.json_TowerInfo)
	webMux.HandleFunc("/ServiceInfo", tw.json_ServiceInfo)
	webMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		tw.serveWebSocketClient(ctx, w, r)
	})
	g2log.TraceService("%v", webMux)
	tw.clientWeb = &http.Server{
		Handler: webMux,
		Addr:    fmt.Sprintf(":%v", tw.config.ServicePort),
	}
	tw.demuxReq2BytesAPIFnMap = [c2t_idcmd.CommandID_Count]func(
		me interface{}, hd c2t_packet.Header, rbody []byte) (
		c2t_packet.Header, interface{}, error){
		c2t_idcmd.Invalid:           tw.bytesAPIFn_ReqInvalid,           // Invalid make empty packet error
		c2t_idcmd.Login:             tw.bytesAPIFn_ReqLogin,             // Login
		c2t_idcmd.Heartbeat:         tw.bytesAPIFn_ReqHeartbeat,         // Heartbeat
		c2t_idcmd.Chat:              tw.bytesAPIFn_ReqChat,              // Chat
		c2t_idcmd.AchieveInfo:       tw.bytesAPIFn_ReqAchieveInfo,       // AchieveInfo
		c2t_idcmd.Rebirth:           tw.bytesAPIFn_ReqRebirth,           // Rebirth
		c2t_idcmd.MoveFloor:         tw.bytesAPIFn_ReqMoveFloor,         // MoveFloor tower cmd
		c2t_idcmd.AIPlay:            tw.bytesAPIFn_ReqAIPlay,            // AIPlay
		c2t_idcmd.VisitFloorList:    tw.bytesAPIFn_ReqVisitFloorList,    // VisitFloorList floor info of visited
		c2t_idcmd.PassTurn:          tw.bytesAPIFn_ReqPassTurn,          // PassTurn no action just trigger turn
		c2t_idcmd.Meditate:          tw.bytesAPIFn_ReqMeditate,          // Meditate turn act
		c2t_idcmd.KillSelf:          tw.bytesAPIFn_ReqKillSelf,          // KillSelf turn act
		c2t_idcmd.Move:              tw.bytesAPIFn_ReqMove,              // Move turn act
		c2t_idcmd.Attack:            tw.bytesAPIFn_ReqAttack,            // Attack turn act
		c2t_idcmd.AttackWide:        tw.bytesAPIFn_ReqAttackWide,        // Attack turn act
		c2t_idcmd.AttackLong:        tw.bytesAPIFn_ReqAttackLong,        // Attack turn act
		c2t_idcmd.Pickup:            tw.bytesAPIFn_ReqPickup,            // Pickup turn act
		c2t_idcmd.Drop:              tw.bytesAPIFn_ReqDrop,              // Drop turn act
		c2t_idcmd.Equip:             tw.bytesAPIFn_ReqEquip,             // Equip turn act
		c2t_idcmd.UnEquip:           tw.bytesAPIFn_ReqUnEquip,           // UnEquip turn act
		c2t_idcmd.DrinkPotion:       tw.bytesAPIFn_ReqDrinkPotion,       // DrinkPotion turn act
		c2t_idcmd.ReadScroll:        tw.bytesAPIFn_ReqReadScroll,        // ReadScroll turn act
		c2t_idcmd.Recycle:           tw.bytesAPIFn_ReqRecycle,           // Recycle turn act
		c2t_idcmd.EnterPortal:       tw.bytesAPIFn_ReqEnterPortal,       // EnterPortal turn act
		c2t_idcmd.ActTeleport:       tw.bytesAPIFn_ReqActTeleport,       // ActTeleport turn act
		c2t_idcmd.AdminTowerCmd:     tw.bytesAPIFn_ReqAdminTowerCmd,     // AdminTowerCmd generic cmd
		c2t_idcmd.AdminFloorCmd:     tw.bytesAPIFn_ReqAdminFloorCmd,     // AdminFloorCmd generic cmd
		c2t_idcmd.AdminActiveObjCmd: tw.bytesAPIFn_ReqAdminActiveObjCmd, // AdminActiveObjCmd generic cmd
		c2t_idcmd.AdminFloorMove:    tw.bytesAPIFn_ReqAdminFloorMove,    // AdminFloorMove Next Before floorUUID
		c2t_idcmd.AdminTeleport:     tw.bytesAPIFn_ReqAdminTeleport,     // AdminTeleport random pos in floor
		c2t_idcmd.AdminAddExp:       tw.bytesAPIFn_ReqAdminAddExp,       // AdminAddExp  add arg to battle exp
		c2t_idcmd.AdminPotionEffect: tw.bytesAPIFn_ReqAdminPotionEffect, // AdminPotionEffect buff by arg potion type
		c2t_idcmd.AdminScrollEffect: tw.bytesAPIFn_ReqAdminScrollEffect, // AdminScrollEffect buff by arg Scroll type
		c2t_idcmd.AdminCondition:    tw.bytesAPIFn_ReqAdminCondition,    // AdminCondition add arg condition for 100 turn
		c2t_idcmd.AdminAddPotion:    tw.bytesAPIFn_ReqAdminAddPotion,    // AdminAddPotion add arg potion to inven
		c2t_idcmd.AdminAddScroll:    tw.bytesAPIFn_ReqAdminAddScroll,    // AdminAddScroll add arg scroll to inven
		c2t_idcmd.AdminAddMoney:     tw.bytesAPIFn_ReqAdminAddMoney,     // AdminAddMoney add arg money to inven
		c2t_idcmd.AdminAddEquip:     tw.bytesAPIFn_ReqAdminAddEquip,     // AdminAddEquip add random equip to inven
		c2t_idcmd.AdminForgetFloor:  tw.bytesAPIFn_ReqAdminForgetFloor,  // AdminForgetFloor forget current floor map
		c2t_idcmd.AdminFloorMap:     tw.bytesAPIFn_ReqAdminFloorMap,     // AdminFloorMap complete current floor map
	}
}

func CheckOrigin(r *http.Request) bool {
	return true
}

func (tw *Tower) serveWebSocketClient(ctx context.Context,
	w http.ResponseWriter, r *http.Request) {

	if tw.IsListenClientPaused() {
		g2log.Warn("ListenClientPaused %v %v", w, r)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: CheckOrigin,
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		g2log.Error("upgrade %v", err)
		return
	}

	g2log.TraceClient("Start serveWebSocketClient %v", r.RemoteAddr)
	defer func() {
		g2log.TraceClient("End serveWebSocketClient %v", r.RemoteAddr)
	}()

	connData := &conndata.ConnData{
		RemoteAddr: r.RemoteAddr,
	}
	c2sc := c2t_serveconnbyte.NewWithStats(
		connData,
		tw.floorMan.GetSendBufferSize(),
		c2t_authorize.NewAllSet(),
		tw.sendStat, tw.recvStat,
		tw.protocolStat,
		tw.notiStat,
		tw.errorStat,
		tw.demuxReq2BytesAPIFnMap,
	)

	c2sc.StartServeWS(ctx, wsConn,
		gameconst.ServerPacketReadTimeOutSec*time.Second,
		gameconst.ServerPacketWriteTimeoutSec*time.Second,
		c2t_gob.MarshalBodyFn,
	)

	// connected user play

	// end play

	// TODO not only suspend ao but also pause tower and floor

	if connData.Logined {
		tw.playerAO.Suspend()
		rspCh := make(chan error, 1)
		tw.GetCmdCh() <- &cmd2tower.PlayerAOSuspendFromTower{
			RspCh: rspCh,
		}
		<-rspCh
		wsConn.Close()
		tw.playerConnection = nil
	}
}

func (tw *Tower) json_HighScore(w http.ResponseWriter, r *http.Request) {
	allActiveObj := aoexpsort.ByExp(tw.aoExpRanking)
	aoLen := len(allActiveObj)
	if aoLen >= 10 {
		aoLen = 10
	}
	aol := make([]*aoscore.ActiveObjScore, aoLen)
	for i := 0; i < aoLen; i++ {
		aol[i] = allActiveObj[i].To_ActiveObjScore()
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	weblib.ServeJSON2HTTP(aol, w)
}
