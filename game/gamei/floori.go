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

package gamei

import (
	"context"
	"net/http"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/fieldobject"
	"github.com/kasworld/goguelike-single/game/terraini"
	"github.com/kasworld/goguelike-single/lib/uuidposmani"
)

type FloorI interface {
	Initialized() bool
	Cleanup()
	GetName() string
	GetWidth() int
	GetHeight() int
	VisitableCount() int // for visitarea

	GetTower() TowerI

	GetBias() bias.Bias
	GetEnvBias() bias.Bias

	GetTerrain() terraini.TerrainI

	GetCmdFloorActStat() *actpersec.ActPerSec

	GetActiveObjPosMan() uuidposmani.UUIDPosManI
	GetCarryObjPosMan() uuidposmani.UUIDPosManI
	GetFieldObjPosMan() uuidposmani.UUIDPosManI

	GetCmdCh() chan<- interface{}
	Turn(TurnCount int)
	ProcessAllCmds()
	Run(ctx context.Context, queuesize int)

	TotalActiveObjCount() int
	TotalCarryObjCount() int
	SearchRandomActiveObjPos() (int, int, error)
	SearchRandomActiveObjPosInRoomOrRandPos() (int, int, error)

	FindPath(dstx, dsty, srcx, srcy int, limit int) [][2]int

	Web_FloorInfo(w http.ResponseWriter, r *http.Request)
	Web_FloorImageZoom(w http.ResponseWriter, r *http.Request)
	Web_FloorImageAutoZoom(w http.ResponseWriter, r *http.Request)
	Web_TileInfo(w http.ResponseWriter, r *http.Request)

	ToPacket_FloorInfo() *csprotocol.FloorInfo

	FindUsablePortalPairAt(x, y int) (*fieldobject.FieldObject, *fieldobject.FieldObject, error)
}
