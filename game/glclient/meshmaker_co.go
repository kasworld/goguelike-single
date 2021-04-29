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
	"github.com/kasworld/goguelike-single/config/moneycolor"
	"github.com/kasworld/goguelike-single/enum/carryingobjecttype"
	"github.com/kasworld/goguelike-single/enum/equipslottype"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/htmlcolors"
)

// manage carry object

type ShiftInfo struct {
	X float32
	Y float32
	Z float32
}

// equipped shift, around ao
var aoEqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	// center
	equipslottype.Helmet:   {0.5, 0.00, 0.5},
	equipslottype.Amulet:   {0.5, 0.33, 0.5},
	equipslottype.Armor:    {0.5, 0.66, 0.5},
	equipslottype.Footwear: {0.5, 1.00, 0.5},

	// right
	equipslottype.Weapon:   {1.00, 0.33, 0.5},
	equipslottype.Gauntlet: {1.00, 0.66, 0.5},

	// left
	equipslottype.Shield: {0.00, 0.33, 0.5},
	equipslottype.Ring:   {0.00, 0.66, 0.5},
}

func cokey2ShiftInfo(cokey COKey) ShiftInfo {
	switch cokey.CT {
	default:
		return otherCarryObjShift[cokey.CT]
	case carryingobjecttype.Equip:
		return eqPosShift[cokey.ET]
	}
}

// on floor in tile
var eqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {0.0, 0.0, 0.0},
	equipslottype.Amulet: {0.75, 0.0, 0.0},

	equipslottype.Weapon: {0.0, 0.25, 0.0},
	equipslottype.Shield: {0.75, 0.25, 0.0},

	equipslottype.Ring:     {0.0, 0.50, 0.0},
	equipslottype.Gauntlet: {0.75, 0.50, 0.0},

	equipslottype.Armor:    {0.0, 0.75, 0.0},
	equipslottype.Footwear: {0.75, 0.75, 0.0},
}

var otherCarryObjShift = [carryingobjecttype.CarryingObjectType_Count]ShiftInfo{
	carryingobjecttype.Money:  {0.33, 0.0, 0.0},
	carryingobjecttype.Potion: {0.33, 0.33, 0.0},
	carryingobjecttype.Scroll: {0.33, 0.66, 0.0},
}

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

func cokey2Color(cokey COKey) string {
	var v htmlcolors.Color24
	switch cokey.CT {
	case carryingobjecttype.Equip:
		v = cokey.FT.Color24()
	case carryingobjecttype.Money:
		if cokey.Val >= len(moneycolor.Attrib) {
			v = moneycolor.Attrib[len(moneycolor.Attrib)-1].Color
		} else {
			v = moneycolor.Attrib[cokey.Val].Color
		}
	case carryingobjecttype.Potion:
		v = cokey.PT.Color24()
	case carryingobjecttype.Scroll:
		v = cokey.ST.Color24()
	}
	return v.NearNamedColor24().String()
}

// TODO update
func newCarryObjMat(cokey COKey) *material.Standard {
	return material.NewStandard(math32.NewColor(cokey2Color(cokey)))
}

// TODO update
func newCarryObjGeo(cokey COKey) *geometry.Geometry {
	return geometry.NewBox(0.1, 0.1, 0.1)
}

func (mm *MeshMaker) initCarryObj(dataFolder string) {
	// do nothing
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
	si := cokey2ShiftInfo(cokey)
	mesh.SetPositionX(float32(x) + si.X - 0.5)
	mesh.SetPositionY(float32(y) + si.Y - 0.5)
	mesh.SetPositionZ(0.5 + si.Z)
	mesh.SetUserData(cokey)
	return mesh
}

func (mm *MeshMaker) PutCarryObj(mesh *graphic.Mesh) {
	cokey := mesh.UserData().(COKey)
	mm.coInUse[cokey]--
	mm.coMeshFreeList[cokey] = append(mm.coMeshFreeList[cokey], mesh)
}
