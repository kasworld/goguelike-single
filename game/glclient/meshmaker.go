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
	"bytes"
	"fmt"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/enum/tile_vector"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

var tileAttrib = [tile.Tile_Count]struct {
	tranparent bool
}{
	tile.Swamp:  {false},
	tile.Soil:   {false},
	tile.Stone:  {false},
	tile.Sand:   {false},
	tile.Sea:    {false},
	tile.Magma:  {false},
	tile.Ice:    {false},
	tile.Grass:  {false},
	tile.Tree:   {false},
	tile.Road:   {false},
	tile.Room:   {false},
	tile.Wall:   {false},
	tile.Window: {true},
	tile.Door:   {false},
	tile.Fog:    {true},
	tile.Smoke:  {true},
}

func (mm *MeshMaker) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "MeshMaker[")
	for i, v := range mm.inUse {
		fmt.Fprintf(&buf, "%v:%v ", tile.Tile(i), v)
	}
	fmt.Fprintf(&buf, "]")
	return buf.String()
}

type MeshMaker struct {
	inUse tile_vector.TileVector
	tex   [tile.Tile_Count]*texture.Texture2D
	mat   [tile.Tile_Count]*material.Standard
	geo   [tile.Tile_Count]*geometry.Geometry
	// tile , free list
	tiles [tile.Tile_Count][]*graphic.Mesh
}

func NewMeshMaker(dataFolder string, initSize int) *MeshMaker {
	mm := MeshMaker{}
	for i := range mm.tex {
		texFilename := tile.Tile(i).String() + ".png"
		tex, err := texture.NewTexture2DFromImage(
			dataFolder + "/tiles/" + texFilename)
		if err != nil {
			g2log.Fatal("Error loading texture: %s", err)
			return nil
		}
		// tex.SetWrapS(gls.REPEAT)
		// tex.SetWrapT(gls.REPEAT)
		// tex.SetRepeat(fw/128, fh/128)
		mm.tex[i] = tex

		mat := material.NewStandard(math32.NewColor("White"))
		mat.AddTexture(tex)
		// mat.SetOpacity(1)
		mat.SetTransparent(tileAttrib[i].tranparent)

		mm.mat[i] = mat

		mm.geo[i] = geometry.NewPlane(1, 1)

		mm.tiles[i] = make([]*graphic.Mesh, 0, initSize)
	}
	return &mm
}

func (mm *MeshMaker) newTile(tl tile.Tile) *graphic.Mesh {
	mesh := graphic.NewMesh(mm.geo[tl], mm.mat[tl])
	return mesh
}

func (mm *MeshMaker) GetTile(tl tile.Tile, x, y int) *graphic.Mesh {
	mm.inUse.Inc(tl)
	var mesh *graphic.Mesh
	freeSize := len(mm.tiles[tl])
	if freeSize > 0 {
		mesh = mm.tiles[tl][freeSize-1]
		mm.tiles[tl] = mm.tiles[tl][:freeSize-1]
	} else {
		mesh = mm.newTile(tl)
	}
	mesh.SetPositionX(float32(x))
	mesh.SetPositionY(float32(y))
	mesh.SetPositionZ(float32(int(tl)-tile.Tile_Count) / float32(tile.Tile_Count))
	return mesh
}

func (mm *MeshMaker) PutTile(tl tile.Tile, mesh *graphic.Mesh) {
	mm.inUse.Dec(tl)
	mm.tiles[tl] = append(mm.tiles[tl], mesh)
}
