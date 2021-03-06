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
	"math/rand"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/enum/tile_flag"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/tilearea"
	"github.com/kasworld/goguelike-single/game/tilearea4pathfind"
	"github.com/kasworld/goguelike-single/game/visitarea"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/uuidposman_map"
	"github.com/kasworld/goguelike-single/lib/uuidposmani"
	"github.com/kasworld/h4o/graphic"
	"github.com/kasworld/h4o/node"
	"github.com/kasworld/wrapper"
)

func (cf *ClientFloor) String() string {
	return fmt.Sprintf("ClientFloor[%v %v %v %v]",
		cf.FloorInfo.Name,
		cf.Visited,
		cf.FloorInfo.W,
		cf.FloorInfo.H,
	)
}

type ClientFloor struct {
	FloorInfo *csprotocol.FloorInfo
	Tiles     tilearea.TileArea `prettystring:"simple"`

	Visited        *visitarea.VisitArea `prettystring:"simple"`
	XWrapper       *wrapper.Wrapper     `prettystring:"simple"`
	YWrapper       *wrapper.Wrapper     `prettystring:"simple"`
	XWrapSafe      func(i int) int
	YWrapSafe      func(i int) int
	Tiles4PathFind *tilearea4pathfind.TileArea4PathFind `prettystring:"simple"`
	visitTurnCount int

	FieldObjPosMan uuidposmani.UUIDPosManI `prettystring:"simple"`

	// for g3n
	Scene *node.Node

	meshMaker *MeshMaker
	// tiles at x,y
	TileMeshs [tile.Tile_Count][][]*graphic.Mesh
	// fieldobj at x,y
	FieldObjMeshs [][]*graphic.Mesh
}

func NewClientFloor(
	meshMaker *MeshMaker,
	FloorInfo *csprotocol.FloorInfo) *ClientFloor {
	cf := ClientFloor{
		meshMaker: meshMaker,
		Tiles:     tilearea.New(FloorInfo.W, FloorInfo.H),
		Visited:   visitarea.New(FloorInfo),
		FloorInfo: FloorInfo,
		XWrapper:  wrapper.New(FloorInfo.W),
		YWrapper:  wrapper.New(FloorInfo.H),
	}
	cf.XWrapSafe = cf.XWrapper.GetWrapSafeFn()
	cf.YWrapSafe = cf.YWrapper.GetWrapSafeFn()
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	cf.FieldObjPosMan = uuidposman_map.New(FloorInfo.W, FloorInfo.H)

	cf.Scene = node.NewNode()

	for i := range cf.TileMeshs {
		cf.TileMeshs[i] = make([][]*graphic.Mesh, cf.FloorInfo.W)
		for x := 0; x < cf.FloorInfo.W; x++ {
			cf.TileMeshs[i][x] = make([]*graphic.Mesh, cf.FloorInfo.H)
		}
	}
	cf.FieldObjMeshs = make([][]*graphic.Mesh, cf.FloorInfo.W)
	for x := 0; x < cf.FloorInfo.W; x++ {
		cf.FieldObjMeshs[x] = make([]*graphic.Mesh, cf.FloorInfo.H)
	}
	return &cf
}

func (cf *ClientFloor) Cleanup() {
	cf.Tiles = nil
	if v := cf.Visited; v != nil {
		v.Cleanup()
	}
	cf.Visited = nil
	if t := cf.Tiles4PathFind; t != nil {
		t.Cleanup()
	}
	cf.Tiles4PathFind = nil
	if i := cf.FieldObjPosMan; i != nil {
		i.Cleanup()
	}
	cf.FieldObjPosMan = nil
}

func (cf *ClientFloor) Forget() {
	FloorInfo := cf.FloorInfo
	cf.Tiles = tilearea.New(FloorInfo.W, FloorInfo.H)
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	cf.Visited = visitarea.New(FloorInfo)
}

// replace tile rect at x,y
func (cf *ClientFloor) ReplaceFloorTiles(tiles tilearea.TileArea) {
	for x, xv := range tiles {
		for y, yv := range xv {
			cf.Tiles[x][y] = yv
			if yv != 0 {
				cf.Visited.CheckAndSetNolock(x, y)
				cf.updateTileMeshAtByTileFlag(yv, x, y)
			}
		}
	}
}

func (cf *ClientFloor) updateTileMeshAtByTileFlag(tf tile_flag.TileFlag, x, y int) {
	for i := 0; i < tile.Tile_Count; i++ {
		if tf.TestByTile(tile.Tile(i)) {
			if cf.TileMeshs[i][x][y] == nil {
				// add new mesh
				mesh := cf.meshMaker.GetTile(tile.Tile(i), x, y)
				cf.Scene.Add(mesh)
				cf.TileMeshs[i][x][y] = mesh
			} else {
				// do nothing
			}
		} else {
			if cf.TileMeshs[i][x][y] == nil {
				// do nothing
			} else {
				// del exist mesh
				mesh := cf.TileMeshs[i][x][y]
				cf.TileMeshs[i][x][y] = nil
				cf.Scene.Remove(mesh)
			}
		}
	}
}
func (cf *ClientFloor) UpdateFromViewportTile(
	vp *csprotocol.NotiVPTiles,
	vpXYLenList findnear.XYLenList) error {

	if cf.FloorInfo.Name != vp.FloorName {
		return fmt.Errorf("vptile data floor not match %v %v",
			cf.FloorInfo.Name, vp.FloorName)

	}
	cf.Visited.UpdateByViewport2(vp.VPX, vp.VPY, vp.VPTiles, vpXYLenList)

	for i, v := range vpXYLenList {
		fx := cf.XWrapSafe(v.X + vp.VPX)
		fy := cf.YWrapSafe(v.Y + vp.VPY)
		if vp.VPTiles[i] != 0 {
			cf.Tiles[fx][fy] = vp.VPTiles[i]
			cf.updateTileMeshAtByTileFlag(vp.VPTiles[i], fx, fy)
		}
	}
	return nil
}

func (cf *ClientFloor) AddOrUpdateFieldObj(v *csprotocol.FieldObjClient) {
	old := cf.FieldObjPosMan.GetByUUID(v.ID)
	if old != nil {
		oldfo := old.(*csprotocol.FieldObjClient)
		if oldfo.X != v.X || oldfo.Y != v.Y {
			// moved
			err := cf.FieldObjPosMan.AddOrUpdateToXY(v, v.X, v.Y)
			if err != nil {
				g2log.Fatal("fail to AddOrUpdateToXY %v", v)
			}
			oldmesh := cf.FieldObjMeshs[v.X][v.Y]
			cf.FieldObjMeshs[v.X][v.Y] = nil
			cf.Scene.Remove(oldmesh)
			cf.meshMaker.PutFieldObj(oldmesh)

			mesh := cf.meshMaker.GetFieldObj(v.ActType, v.DisplayType, v.X, v.Y)
			cf.FieldObjMeshs[v.X][v.Y] = mesh
			cf.Scene.Add(mesh)
		} else {
			// do nothing
		}
	} else {
		// add new
		err := cf.FieldObjPosMan.AddOrUpdateToXY(v, v.X, v.Y)
		if err != nil {
			g2log.Fatal("fail to AddOrUpdateToXY %v", v)
		}
		mesh := cf.meshMaker.GetFieldObj(v.ActType, v.DisplayType, v.X, v.Y)
		cf.FieldObjMeshs[v.X][v.Y] = mesh
		cf.Scene.Add(mesh)
	}
}

func (cf *ClientFloor) UpdateFieldObjList(folsit []*csprotocol.FieldObjClient) {
	for _, v := range folsit {
		cf.AddOrUpdateFieldObj(v)
	}
}

func (cf *ClientFloor) PosAddDir(x, y int, dir way9type.Way9Type) (int, int) {
	nextX := x + dir.Dx()
	nextY := y + dir.Dy()
	nextX = cf.XWrapper.Wrap(nextX)
	nextY = cf.YWrapper.Wrap(nextY)
	return nextX, nextY
}

func (cf *ClientFloor) FindMovableDir(x, y int, dir way9type.Way9Type) way9type.Way9Type {
	dirList := []way9type.Way9Type{
		dir,
		dir.TurnDir(1),
		dir.TurnDir(-1),
		dir.TurnDir(2),
		dir.TurnDir(-2),
	}
	if rand.Float64() >= 0.5 {
		dirList = []way9type.Way9Type{
			dir,
			dir.TurnDir(-1),
			dir.TurnDir(1),
			dir.TurnDir(-2),
			dir.TurnDir(2),
		}
	}
	for _, dir := range dirList {
		nextX, nextY := cf.PosAddDir(x, y, dir)
		if cf.Tiles[nextX][nextY].CharPlaceable() {
			return dir
		}
	}
	return way9type.Center
}

func (cf *ClientFloor) IsValidPos(x, y int) bool {
	return cf.XWrapper.IsIn(x) && cf.YWrapper.IsIn(y)
}

func (cf *ClientFloor) GetBias() bias.Bias {
	if cf.FloorInfo != nil {
		return cf.FloorInfo.Bias
	} else {
		return bias.Bias{}
	}
}

func (cf *ClientFloor) LeaveFloor() {
}

func (cf *ClientFloor) EnterFloor(TurnCount int) {
	cf.visitTurnCount = TurnCount
}

func (cf *ClientFloor) GetFieldObjAt(x, y int) *csprotocol.FieldObjClient {
	po, ok := cf.FieldObjPosMan.Get1stObjAt(x, y).(*csprotocol.FieldObjClient)
	if !ok {
		return nil
	}
	return po
}
