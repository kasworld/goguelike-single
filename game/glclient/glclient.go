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
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_pid2rspfn"
)

type GLClient struct {
	sendRecvStop func() `prettystring:"hide"`

	config *goguelikeconfig.GoguelikeConfig

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
	app   *app.Application
	scene *core.Node
	cam   *camera.Camera
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
	return app
}

func (app *GLClient) Cleanup() {
	app.ServerJitter = nil
}

func (app *GLClient) Run(mainctx context.Context) error {
	defer app.Cleanup()

	ctx, closeCtx := context.WithCancel(mainctx)
	app.sendRecvStop = closeCtx
	defer app.sendRecvStop()

	go app.handle_t2ch()

	var rtnerr error
	if err := app.reqLogin(); err != nil {
		return err
	}

	go app.runG3N()

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
