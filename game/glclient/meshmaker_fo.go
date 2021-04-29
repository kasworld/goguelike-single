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
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/fieldobjdisplaytype"
)

// mange field object

type FOKey struct {
	AT fieldobjacttype.FieldObjActType
	DT fieldobjdisplaytype.FieldObjDisplayType
}

var foAttrib = []struct {
	Co string
}{
	fieldobjacttype.None:             {"black"},
	fieldobjacttype.PortalInOut:      {"mediumvioletred"},
	fieldobjacttype.PortalIn:         {"mediumvioletred"},
	fieldobjacttype.PortalOut:        {"mediumvioletred"},
	fieldobjacttype.PortalAutoIn:     {"mediumvioletred"},
	fieldobjacttype.RecycleCarryObj:  {"green"},
	fieldobjacttype.Teleport:         {"red"},
	fieldobjacttype.ForgetFloor:      {"orangered"},
	fieldobjacttype.ForgetOneFloor:   {"orangered"},
	fieldobjacttype.AlterFaction:     {"red"},
	fieldobjacttype.AllFaction:       {"red"},
	fieldobjacttype.Bleeding:         {"crimson"},
	fieldobjacttype.Chilly:           {"darkturquoise"},
	fieldobjacttype.Blind:            {"darkred"},
	fieldobjacttype.Invisible:        {"lemonchiffon"},
	fieldobjacttype.Burden:           {"deeppink"},
	fieldobjacttype.Float:            {"wheat"},
	fieldobjacttype.Greasy:           {"papayawhip"},
	fieldobjacttype.Drunken:          {"plum"},
	fieldobjacttype.Sleepy:           {"lightcoral"},
	fieldobjacttype.Contagion:        {"darkgreen"},
	fieldobjacttype.Slow:             {"darkblue"},
	fieldobjacttype.Haste:            {"lightblue"},
	fieldobjacttype.RotateLineAttack: {"lavender"},
	fieldobjacttype.Mine:             {"orange"},
}

func newFieldObjMat(fokey FOKey) *material.Standard {
	return material.NewStandard(math32.NewColor(foAttrib[fokey.AT].Co))
}

func newFieldObjGeo(fokey FOKey) *geometry.Geometry {
	return geometry.NewCone(0.5, 1, int(fokey.DT)+2, 8, true)
}

func (mm *MeshMaker) initFieldObj(dataFolder string) {
	// do nothing
}

func (mm *MeshMaker) newFieldObj(fokey FOKey) *graphic.Mesh {
	var mat *material.Standard
	var exist bool
	if mat, exist = mm.foMat[fokey]; !exist {
		mat = newFieldObjMat(fokey)
		mm.foMat[fokey] = mat
	}
	var geo *geometry.Geometry
	if geo, exist = mm.foGeo[fokey]; !exist {
		geo = newFieldObjGeo(fokey)
		mm.foGeo[fokey] = geo
	}
	return graphic.NewMesh(geo, mat)
}

func (mm *MeshMaker) GetFieldObj(
	at fieldobjacttype.FieldObjActType,
	dt fieldobjdisplaytype.FieldObjDisplayType,
	x, y int) *graphic.Mesh {
	fokey := FOKey{at, dt}
	mm.foInUse[fokey]++
	var mesh *graphic.Mesh
	freeSize := len(mm.foMeshFreeList[fokey])
	if freeSize > 0 {
		mesh = mm.foMeshFreeList[fokey][freeSize-1]
		mm.foMeshFreeList[fokey] = mm.foMeshFreeList[fokey][:freeSize-1]
	} else {
		mesh = mm.newFieldObj(fokey)
		mesh.RotateX(math.Pi / 2)
	}
	mesh.SetPositionX(float32(x))
	mesh.SetPositionY(float32(y))
	mesh.SetPositionZ(0.5)
	mesh.SetUserData(fokey)
	return mesh
}

func (mm *MeshMaker) PutFieldObj(mesh *graphic.Mesh) {
	fokey := mesh.UserData().(FOKey)
	mm.foInUse[fokey]--
	mm.foMeshFreeList[fokey] = append(mm.foMeshFreeList[fokey], mesh)
}
