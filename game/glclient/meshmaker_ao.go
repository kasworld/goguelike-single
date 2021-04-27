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
	"math"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/kasworld/goguelike-single/enum/factiontype"
)

// manage active object

var aoAttrib = [factiontype.FactionType_Count]struct {
	Co string
}{
	factiontype.Black:           {"black"},
	factiontype.Maroon:          {"maroon"},
	factiontype.Red:             {"red"},
	factiontype.Green:           {"green"},
	factiontype.Olive:           {"olive"},
	factiontype.DarkOrange:      {"darkorange"},
	factiontype.Lime:            {"lime"},
	factiontype.Chartreuse:      {"chartreuse"},
	factiontype.Yellow:          {"yellow"},
	factiontype.Navy:            {"navy"},
	factiontype.Purple:          {"purple"},
	factiontype.DeepPink:        {"deeppink"},
	factiontype.Teal:            {"teal"},
	factiontype.Salmon:          {"salmon"},
	factiontype.SpringGreen:     {"springgreen"},
	factiontype.LightGreen:      {"lightgreen"},
	factiontype.Khaki:           {"khaki"},
	factiontype.Blue:            {"blue"},
	factiontype.DarkViolet:      {"darkviolet"},
	factiontype.Magenta:         {"magenta"},
	factiontype.DodgerBlue:      {"dodgerblue"},
	factiontype.MediumSlateBlue: {"mediumslateblue"},
	factiontype.Violet:          {"violet"},
	factiontype.Cyan:            {"cyan"},
	factiontype.Aquamarine:      {"aquamarine"},
	factiontype.White:           {"white"},
}

func newActiveObjMat(ft factiontype.FactionType) *material.Standard {
	return material.NewStandard(math32.NewColor(aoAttrib[ft].Co))
}

func newActiveObjGeo(ft factiontype.FactionType) *geometry.Geometry {
	return geometry.NewCylinder(0.4, 1, 16, 8, true, true)
}

func (mm *MeshMaker) initActiveObj(dataFolder string, initSize int) {
}

func (mm *MeshMaker) newActiveObj(ft factiontype.FactionType) *graphic.Mesh {
	var mat *material.Standard
	if mat = mm.aoMat[ft]; mat == nil {
		mat = newActiveObjMat(ft)
		mm.aoMat[ft] = mat
	}
	var geo *geometry.Geometry
	if geo = mm.aoGeo[ft]; geo == nil {
		geo = newActiveObjGeo(ft)
		mm.aoGeo[ft] = geo
	}
	return graphic.NewMesh(geo, mat)
}

func (mm *MeshMaker) GetActiveObj(ft factiontype.FactionType, x, y int) *graphic.Mesh {
	mm.aoInUse.Inc(ft)
	var mesh *graphic.Mesh
	freeSize := len(mm.aoMeshFreeList[ft])
	if freeSize > 0 {
		mesh = mm.aoMeshFreeList[ft][freeSize-1]
		mm.aoMeshFreeList[ft] = mm.aoMeshFreeList[ft][:freeSize-1]
	} else {
		mesh = mm.newActiveObj(ft)
		mesh.RotateX(math.Pi / 2)
	}
	mesh.SetPositionX(float32(x))
	mesh.SetPositionY(float32(y))
	mesh.SetPositionZ(0.5)
	mesh.SetUserData(ft)
	return mesh
}

func (mm *MeshMaker) PutActiveObj(mesh *graphic.Mesh) {
	ft := mesh.UserData().(factiontype.FactionType)
	mm.aoInUse.Dec(ft)
	mm.aoMeshFreeList[ft] = append(mm.aoMeshFreeList[ft], mesh)
}
