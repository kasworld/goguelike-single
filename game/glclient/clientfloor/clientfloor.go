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

package clientfloor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/kasworld/findnear"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike-single/config/goguelikeconfig"
	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/tilearea"
	"github.com/kasworld/goguelike-single/game/tilearea4pathfind"
	"github.com/kasworld/goguelike-single/game/visitarea"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/uuidposman_map"
	"github.com/kasworld/goguelike-single/lib/uuidposmani"
	"github.com/kasworld/wrapper"
)

type ClientFloor struct {
	FloorInfo *csprotocol.FloorInfo
	Tiles     tilearea.TileArea `prettystring:"simple"`

	Visited        *visitarea.VisitArea `prettystring:"simple"`
	XWrapper       *wrapper.Wrapper     `prettystring:"simple"`
	YWrapper       *wrapper.Wrapper     `prettystring:"simple"`
	XWrapSafe      func(i int) int
	YWrapSafe      func(i int) int
	Tiles4PathFind *tilearea4pathfind.TileArea4PathFind `prettystring:"simple"`
	visitTime      time.Time                            `prettystring:"simple"`

	FieldObjPosMan uuidposmani.UUIDPosManI `prettystring:"simple"`

	// for g3n
	Scene *core.Node

	// tile, x, y
	TerrainTiles [][][]*graphic.Mesh
}

func New(
	config *goguelikeconfig.GoguelikeConfig,
	FloorInfo *csprotocol.FloorInfo) *ClientFloor {
	cf := ClientFloor{
		Tiles:        tilearea.New(FloorInfo.W, FloorInfo.H),
		Visited:      visitarea.New(FloorInfo),
		FloorInfo:    FloorInfo,
		XWrapper:     wrapper.New(FloorInfo.W),
		YWrapper:     wrapper.New(FloorInfo.H),
		TerrainTiles: make([][][]*graphic.Mesh, tile.Tile_Count),
	}
	cf.XWrapSafe = cf.XWrapper.GetWrapSafeFn()
	cf.YWrapSafe = cf.YWrapper.GetWrapSafeFn()
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	cf.FieldObjPosMan = uuidposman_map.New(FloorInfo.W, FloorInfo.H)

	fw := float32(cf.FloorInfo.W)
	fh := float32(cf.FloorInfo.H)
	cf.Scene = core.NewNode()

	// make terrain layers
	rnd := g2rand.New()
	geo := geometry.NewPlane(1, 1)
	for i := range cf.TerrainTiles {
		texFilename := tile.Tile(i).String() + ".png"
		tex, err := texture.NewTexture2DFromImage(
			config.ClientDataFolder + "/tiles/" + texFilename)
		if err != nil {
			g2log.Fatal("Error loading texture: %s", err)
		}
		tex.SetWrapS(gls.REPEAT)
		tex.SetWrapT(gls.REPEAT)
		tex.SetRepeat(fw/128, fh/128)

		mat := material.NewStandard(math32.NewColor("White"))
		mat.SetOpacity(1)
		mat.SetTransparent(true)
		mat.AddTexture(tex)

		for x := 0; x < cf.FloorInfo.W; x++ {
			cf.TerrainTiles[i] = make([][]*graphic.Mesh, cf.FloorInfo.W)
			for y := 0; y < cf.FloorInfo.H; y++ {
				cf.TerrainTiles[i][x] = make([]*graphic.Mesh, cf.FloorInfo.H)
				mesh := graphic.NewMesh(geo, mat)
				mesh.SetPositionX(float32(x))
				mesh.SetPositionY(float32(y))
				mesh.SetPositionZ(float32(i - tile.Tile_Count))

				cf.Scene.Add(mesh)
				cf.TerrainTiles[i][x][y] = mesh
				mesh.SetVisible(rnd.Intn(10) == 0)
			}
		}
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
func (cf *ClientFloor) ReplaceFloorTiles(fta *csprotocol.NotiFloorTiles) {
	for x, xv := range fta.Tiles {
		xpos := fta.X + x
		for y, yv := range xv {
			ypos := fta.Y + y
			cf.Tiles[xpos][ypos] = yv
			if yv != 0 {
				cf.Visited.CheckAndSetNolock(xpos, ypos)
			}
		}
	}
}

func (cf *ClientFloor) String() string {
	return fmt.Sprintf("ClientFloor[%v %v %v %v]",
		cf.FloorInfo.Name,
		cf.Visited,
		cf.XWrapper.GetWidth(),
		cf.YWrapper.GetWidth(),
	)
}

func (cf *ClientFloor) UpdateFromViewportTile(
	vp *csprotocol.NotiVPTiles,
	vpXYLenList findnear.XYLenList,
) error {

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
		}
	}
	return nil
}

func (cf *ClientFloor) UpdateFieldObjList(folsit []*csprotocol.FieldObjClient) {
	for _, v := range folsit {
		cf.FieldObjPosMan.AddOrUpdateToXY(v, v.X, v.Y)
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

func (cf *ClientFloor) EnterFloor() {
	cf.visitTime = time.Now()
}

func (cf *ClientFloor) GetFieldObjAt(x, y int) *csprotocol.FieldObjClient {
	po, ok := cf.FieldObjPosMan.Get1stObjAt(x, y).(*csprotocol.FieldObjClient)
	if !ok {
		return nil
	}
	return po
}
