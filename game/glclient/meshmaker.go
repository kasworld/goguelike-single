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
	"github.com/kasworld/goguelike-single/enum/dangertype"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/enum/tile_vector"
)

func (mm *MeshMaker) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "MeshMaker[")
	fmt.Fprintf(&buf, "TileInUse:")
	for i, v := range mm.tileInUse {
		fmt.Fprintf(&buf, "%v:%v ", tile.Tile(i), v)
	}
	fmt.Fprintf(&buf, "FieldObjInUse:")
	for i, v := range mm.foInUse {
		fmt.Fprintf(&buf, "%v:%v ", i, v)
	}
	fmt.Fprintf(&buf, "]")
	return buf.String()
}

type MeshMaker struct {
	// tile
	tileInUse tile_vector.TileVector
	tileTex   [tile.Tile_Count]*texture.Texture2D
	tileMat   [tile.Tile_Count]*material.Standard
	tileGeo   [tile.Tile_Count]*geometry.Geometry
	// free list
	tileMeshFreeList [tile.Tile_Count][]*graphic.Mesh

	// field object
	foInUse map[FOKey]int
	foMat   map[FOKey]*material.Standard
	foGeo   map[FOKey]*geometry.Geometry
	// free list
	foMeshFreeLIst map[FOKey][]*graphic.Mesh

	// active object
	aoInUse [factiontype.FactionType_Count]int
	aoMat   [factiontype.FactionType_Count]*material.Standard
	aoGeo   [factiontype.FactionType_Count]*geometry.Geometry
	// free list
	aoMeshFreeLIst [factiontype.FactionType_Count][]*graphic.Mesh

	// carry object
	coInUse map[COKey]int
	coMat   map[COKey]*material.Standard
	coGeo   map[COKey]*geometry.Geometry
	// free list
	coMeshFreeLIst map[COKey][]*graphic.Mesh

	// danger object
	doInUse [dangertype.DangerType_Count]int
	doMat   [dangertype.DangerType_Count]*material.Standard
	doGeo   [dangertype.DangerType_Count]*geometry.Geometry
	// free list
	doMeshFreeLIst [dangertype.DangerType_Count][]*graphic.Mesh
}

func NewMeshMaker(dataFolder string, initSize int) *MeshMaker {
	mm := MeshMaker{
		foInUse:        make(map[FOKey]int),
		foMat:          make(map[FOKey]*material.Standard),
		foGeo:          make(map[FOKey]*geometry.Geometry),
		foMeshFreeLIst: make(map[FOKey][]*graphic.Mesh),

		coInUse:        make(map[COKey]int),
		coMat:          make(map[COKey]*material.Standard),
		coGeo:          make(map[COKey]*geometry.Geometry),
		coMeshFreeLIst: make(map[COKey][]*graphic.Mesh),
	}
	for i := range mm.tileTex {
		tex := loadTileTexture(dataFolder + "/tiles/" + tile.Tile(i).String() + ".png")
		mm.tileTex[i] = tex

		mat := material.NewStandard(math32.NewColor("White"))
		mat.AddTexture(tex)
		// mat.SetOpacity(1)
		mat.SetTransparent(tileAttrib[i].tranparent)

		mm.tileMat[i] = mat

		// mm.tileGeo[i] = geometry.NewPlane(1, 1)
		mm.tileGeo[i] = geometry.NewBox(1, 1, tileAttrib[i].height)

		mm.tileMeshFreeList[i] = make([]*graphic.Mesh, 0, initSize)
	}
	return &mm
}
