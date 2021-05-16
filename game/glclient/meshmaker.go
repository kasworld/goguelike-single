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

	"github.com/kasworld/goguelike-single/enum/dangertype"
	"github.com/kasworld/goguelike-single/enum/dangertype_vector_int"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/factiontype_vector_int"
	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/enum/tile_vector_int"
	"github.com/kasworld/h4o/geometry"
	"github.com/kasworld/h4o/graphic"
	"github.com/kasworld/h4o/material"
	"github.com/kasworld/h4o/texture"
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
	tileInUse tile_vector_int.TileVector_int
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
	foMeshFreeList map[FOKey][]*graphic.Mesh

	// active object
	aoInUse factiontype_vector_int.FactionTypeVector_int
	aoMat   [factiontype.FactionType_Count]*material.Standard
	aoGeo   [factiontype.FactionType_Count]*geometry.Geometry
	// free list
	aoMeshFreeList [factiontype.FactionType_Count][]*graphic.Mesh

	// carry object
	coInUse map[COKey]int
	coMat   map[COKey]*material.Standard
	coGeo   map[COKey]*geometry.Geometry
	// free list
	coMeshFreeList map[COKey][]*graphic.Mesh

	// danger object
	doInUse dangertype_vector_int.DangerTypeVector_int
	doMat   [dangertype.DangerType_Count]*material.Standard
	doGeo   [dangertype.DangerType_Count]*geometry.Geometry
	// free list
	doMeshFreeList [dangertype.DangerType_Count][]*graphic.Mesh
}

func NewMeshMaker(dataFolder string) *MeshMaker {
	mm := MeshMaker{
		foInUse:        make(map[FOKey]int),
		foMat:          make(map[FOKey]*material.Standard),
		foGeo:          make(map[FOKey]*geometry.Geometry),
		foMeshFreeList: make(map[FOKey][]*graphic.Mesh),

		coInUse:        make(map[COKey]int),
		coMat:          make(map[COKey]*material.Standard),
		coGeo:          make(map[COKey]*geometry.Geometry),
		coMeshFreeList: make(map[COKey][]*graphic.Mesh),
	}
	mm.initTile(dataFolder)
	mm.initFieldObj(dataFolder)
	mm.initActiveObj(dataFolder)
	mm.initCarryObj(dataFolder)
	mm.initDangerObj(dataFolder)
	return &mm
}
