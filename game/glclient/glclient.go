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
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/goguelikeconfig"
	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/clientfloor"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_pid2rspfn"
)

type GLClient struct {
	GameInfo *c2t_obj.GameInfo

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

	// client to tower packet channel
	c2tCh chan *c2t_packet.Packet
	// tower to client packet channel
	t2cCh chan *c2t_packet.Packet

	// g3n field
	app      *app.Application
	scene    *core.Node
	cam      *camera.Camera
	pLight   *light.Point
	boundBox *graphic.Mesh
	playerAO *graphic.Mesh
}

func New(
	config *goguelikeconfig.GoguelikeConfig,
	gameInfo *c2t_obj.GameInfo,
	c2tch, t2cch chan *c2t_packet.Packet) *GLClient {
	app := &GLClient{
		pid2recv:          c2t_pid2rspfn.New(),
		ViewportXYLenList: viewportdata.ViewportXYLenList,
		c2tCh:             c2tch,
		t2cCh:             t2cch,
		GameInfo:          gameInfo,
	}
	return app
}

func (app *GLClient) TowerBias() bias.Bias {
	if app.OLNotiData == nil {
		return bias.Bias{}
	}
	ft := app.GameInfo.Factor
	dur := app.OLNotiData.Time.Sub(app.GameInfo.StartTime)
	return bias.MakeBiasByProgress(ft, dur.Seconds(), gameconst.TowerBaseBiasLen)
}

func (app *GLClient) GetPlayerXY() (int, int) {
	ao := app.playerActiveObjClient
	if ao != nil {
		return ao.X, ao.Y
	}
	return 0, 0
}
