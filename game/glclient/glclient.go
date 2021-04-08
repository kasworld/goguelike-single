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
	"context"
	"fmt"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/kasworld/actjitter"
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/goguelikeconfig"
	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/game/clientfloor"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_connwsgorilla"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_pid2rspfn"
)

type GLClient struct {
	sendRecvStop func() `prettystring:"hide"`

	config *goguelikeconfig.GoguelikeConfig

	towerConn         *c2t_connwsgorilla.Connection
	ServiceInfo       *c2t_obj.ServiceInfo
	AccountInfo       *c2t_obj.AccountInfo
	TowerInfo         *c2t_obj.TowerInfo
	ViewportXYLenList findnear.XYLenList
	FloorInfoList     []*c2t_obj.FloorInfo
	CurrentFloor      *clientfloor.ClientFloor

	pid2recv *c2t_pid2rspfn.PID2RspFn

	// turn data
	OLNotiData            *c2t_obj.NotiVPObjList_data
	playerActiveObjClient *c2t_obj.ActiveObjClient
	onFieldObj            *c2t_obj.FieldObjClient
	IsOverLoad            bool
	HPdiff                int
	SPdiff                int
	level                 int
	ServerClientTimeDiff  time.Duration
	ServerJitter          *actjitter.ActJitter

	// client to tower packet channel
	c2tCh chan *c2t_packet.Packet
	// tower to client packet channel
	t2cCh chan *c2t_packet.Packet

	// g3n field
	g3napp *app.Application
	scene  *core.Node
	cam    *camera.Camera
}

func New(config *goguelikeconfig.GoguelikeConfig,
	c2tch, t2cch chan *c2t_packet.Packet) *GLClient {
	fmt.Printf("%v\n", config.StringForm())
	app := &GLClient{
		config:            config,
		ServerJitter:      actjitter.New("Server"),
		pid2recv:          c2t_pid2rspfn.New(),
		ViewportXYLenList: viewportdata.ViewportXYLenList,
		c2tCh:             c2tch,
		t2cCh:             t2cch,
	}
	app.sendRecvStop = func() {
		g2log.Error("Too early sendRecvStop call %v", app)
	}
	app.towerConn = c2t_connwsgorilla.New(10)
	go newG3N()
	return app
}

func (app *GLClient) Cleanup() {
	if tc := app.towerConn; tc != nil {
		tc.Cleanup()
	}
	app.ServerJitter = nil
}

func (app *GLClient) Run(mainctx context.Context) error {
	defer app.Cleanup()

	ctx, closeCtx := context.WithCancel(mainctx)
	app.sendRecvStop = closeCtx
	defer app.sendRecvStop()

	if err := app.towerConn.ConnectTo(app.config.ConnectToTower()); err != nil {
		return err
	}
	var rtnerr error
	go func() {
		err := app.towerConn.Run(ctx,
			gameconst.ClientReadTimeoutSec*time.Second,
			gameconst.ClientWriteTimeoutSec*time.Second,
			c2t_gob.MarshalBodyFn,
			app.handleRecvPacket,
			app.handleSentPacket,
		)

		if err != nil {
			rtnerr = err
		}
		app.sendRecvStop()
	}()

	if err := app.reqLogin(); err != nil {
		return err
	}

	timerPingTk := time.NewTicker(time.Second * gameconst.ServerPacketReadTimeOutSec / 2)
	defer timerPingTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerPingTk.C:
			go func() {
				err := app.reqHeartbeat()
				if err != nil {
					rtnerr = err
					app.sendRecvStop()
				}
			}()
		}
	}
	return rtnerr
}

func (app *GLClient) handleSentPacket(pk *c2t_packet.Packet) error {
	g2log.TraceClient("sent %v", pk.Header)
	return nil
}

func (app *GLClient) handleRecvPacket(header c2t_packet.Header, body []byte) error {
	g2log.TraceClient("recv %v", header)
	switch header.FlowType {
	default:
		return fmt.Errorf("invalid packet type %v %v", header, body)
	case c2t_packet.Response:
		if err := app.pid2recv.HandleRsp(header, body); err != nil {
			app.sendRecvStop()
			g2log.Fatal("%v %v %v %v", app, header, body, err)
			return err
		}
	case c2t_packet.Notification:
		fn := DemuxNoti2ByteFnMap[header.Cmd]
		if err := fn(app, header, body); err != nil {
			app.sendRecvStop()
			g2log.Fatal("%v %v %v %v", app, header, body, err)
			return err
		}
	}
	return nil
}
