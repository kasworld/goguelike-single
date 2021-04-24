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
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

type MeshMaker struct {
	inUse []int
	tex   []*texture.Texture2D
	mat   []*material.Standard
	geo   []*geometry.Geometry
	// tile , free list
	tiles [][]*graphic.Mesh
}

func NewMeshMaker(dataFolder string, initSize int) *MeshMaker {
	mm := MeshMaker{
		inUse: make([]int, tile.Tile_Count),
		tex:   make([]*texture.Texture2D, tile.Tile_Count),
		mat:   make([]*material.Standard, tile.Tile_Count),
		geo:   make([]*geometry.Geometry, tile.Tile_Count),
		tiles: make([][]*graphic.Mesh, tile.Tile_Count),
	}
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
		// mat.SetTransparent(true)

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
	mm.inUse[tl]++
	var mesh *graphic.Mesh
	freeSize := len(mm.tiles[tl])
	if freeSize > 0 {
		mesh = mm.tiles[tl][freeSize-1]
		mm.tiles = mm.tiles[:freeSize-1]
	} else {
		mesh = mm.newTile(tl)
	}
	mesh.SetPositionX(float32(x))
	mesh.SetPositionY(float32(y))
	mesh.SetPositionZ(float32(int(tl) - tile.Tile_Count))
	return mesh
}

func (mm *MeshMaker) PutTile(tl tile.Tile, mesh *graphic.Mesh) {
	mm.inUse[tl]--
	mm.tiles[tl] = append(mm.tiles[tl], mesh)
}
