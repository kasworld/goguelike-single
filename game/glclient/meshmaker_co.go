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
	"github.com/kasworld/goguelike-single/enum/carryingobjecttype"
	"github.com/kasworld/goguelike-single/enum/equipslottype"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/game/csprotocol"
)

// manage carry object

type COKey struct {
	CT  carryingobjecttype.CarryingObjectType
	ET  equipslottype.EquipSlotType
	FT  factiontype.FactionType
	PT  potiontype.PotionType
	ST  scrolltype.ScrollType
	Val int // log10 value
}

func NewCOKeyFromCarryObjClientOnFloor(
	src *csprotocol.CarryObjClientOnFloor) COKey {
	return COKey{
		CT:  src.CarryingObjectType,
		ET:  src.EquipType,
		FT:  src.Faction,
		PT:  src.PotionType,
		ST:  src.ScrollType,
		Val: int(math.Log10(float64(src.Value))),
	}
}

// TODO update
func newCarryObjMat(cokey COKey) *material.Standard {
	return material.NewStandard(math32.NewColor("yellow"))
}

// TODO update
func newCarryObjGeo(cokey COKey) *geometry.Geometry {
	return geometry.NewBox(0.1, 0.1, 0.1)
}

func (mm *MeshMaker) initCarryObj(dataFolder string, initSize int) {

}

func (mm *MeshMaker) newCarryObj(cokey COKey) *graphic.Mesh {
	var mat *material.Standard
	var exist bool
	if mat, exist = mm.coMat[cokey]; !exist {
		mat = newCarryObjMat(cokey)
		mm.coMat[cokey] = mat
	}
	var geo *geometry.Geometry
	if geo, exist = mm.coGeo[cokey]; !exist {
		geo = newCarryObjGeo(cokey)
		mm.coGeo[cokey] = geo
	}
	return graphic.NewMesh(geo, mat)
}

func (mm *MeshMaker) GetCarryObj(cokey COKey, x, y int) *graphic.Mesh {
	mm.coInUse[cokey]++
	var mesh *graphic.Mesh
	freeSize := len(mm.coMeshFreeList[cokey])
	if freeSize > 0 {
		mesh = mm.coMeshFreeList[cokey][freeSize-1]
		mm.coMeshFreeList[cokey] = mm.coMeshFreeList[cokey][:freeSize-1]
	} else {
		mesh = mm.newCarryObj(cokey)
	}
	mesh.SetPositionX(float32(x))
	mesh.SetPositionY(float32(y))
	mesh.SetPositionZ(0.5)
	mesh.SetUserData(cokey)
	return mesh
}

func (mm *MeshMaker) PutCarryObj(mesh *graphic.Mesh) {
	cokey := mesh.UserData().(COKey)
	mm.coInUse[cokey]--
	mm.coMeshFreeList[cokey] = append(mm.coMeshFreeList[cokey], mesh)
}
