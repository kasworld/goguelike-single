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
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/goguelikeconfig"
	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/glclient/pid2rspfn"
	"github.com/kasworld/goguelike-single/lib/engine/appbase"
	"github.com/kasworld/goguelike-single/lib/engine/camera"
	"github.com/kasworld/goguelike-single/lib/engine/gui"
	"github.com/kasworld/goguelike-single/lib/engine/light"
	"github.com/kasworld/goguelike-single/lib/engine/node"
	"github.com/kasworld/goguelike-single/lib/engine/util/framerater"
)

type GLClient struct {
	GameInfo *csprotocol.GameInfo
	config   *goguelikeconfig.GoguelikeConfig

	ViewportXYLenList findnear.XYLenList
	Name2Floor        map[string]*ClientFloor
	CurrentFloor      *ClientFloor

	pid2recv *pid2rspfn.PID2RspFn

	// turn data
	OLNotiData            *csprotocol.NotiVPObjList
	playerActiveObjClient *csprotocol.ActiveObjClient
	onFieldObj            *csprotocol.FieldObjClient
	IsOverLoad            bool
	HPdiff                float64
	SPdiff                float64
	level                 int

	// client to tower packet channel
	c2tCh chan *csprotocol.Packet
	// tower to client packet channel
	t2cCh chan *csprotocol.Packet

	// g3n field
	meshMaker  *MeshMaker
	app        *appbase.AppBase
	scene      *node.Node
	cam        *camera.Camera
	camZpos    float32
	pLight     *light.Point
	frameRater *framerater.FrameRater // Render loop frame rater
	labelFPS   *gui.Label             // header FPS label

	sceneAO *node.Node
	sceneCO *node.Node
	sceneDO *node.Node
}

func New(
	config *goguelikeconfig.GoguelikeConfig,
	gameInfo *csprotocol.GameInfo,
	c2tch, t2cch chan *csprotocol.Packet) *GLClient {
	app := &GLClient{
		Name2Floor:        make(map[string]*ClientFloor),
		pid2recv:          pid2rspfn.New(),
		ViewportXYLenList: viewportdata.ViewportXYLenList,
		c2tCh:             c2tch,
		t2cCh:             t2cch,
		GameInfo:          gameInfo,
		config:            config,
	}
	app.meshMaker = NewMeshMaker(config.ClientDataFolder)
	if app.meshMaker == nil {
		return nil
	}
	return app
}

func (app *GLClient) GetTurnCount() int {
	return app.OLNotiData.TurnCount
}

func (app *GLClient) TowerBias() bias.Bias {
	if app.OLNotiData == nil {
		return bias.Bias{}
	}
	ft := app.GameInfo.Factor
	return bias.MakeBiasByProgress(ft, float64(app.GetTurnCount()), gameconst.TowerBaseBiasLen)
}

func (app *GLClient) GetPlayerXY() (int, int) {
	ao := app.playerActiveObjClient
	if ao != nil {
		return ao.X, ao.Y
	}
	return 0, 0
}
