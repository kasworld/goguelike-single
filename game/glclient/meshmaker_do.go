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
	"github.com/kasworld/goguelike-single/enum/dangertype"
)

// manage danger object

var doAttrib = [dangertype.DangerType_Count]struct {
	Co string
}{
	dangertype.None:             {"black"},
	dangertype.BasicAttack:      {"red"},
	dangertype.WideAttack:       {"crimson"},
	dangertype.LongAttack:       {"firebrick"},
	dangertype.RotateLineAttack: {"deeppink"},
	dangertype.MineExplode:      {"orange"},
}

func newDangerObjMat(dt dangertype.DangerType) *material.Standard {
	return material.NewStandard(math32.NewColor(doAttrib[dt].Co))
}

func newDangerObjGeo(dt dangertype.DangerType) *geometry.Geometry {
	return geometry.NewTorus(0.5, 0.1, 16, 8, math.Pi*2)
}

func (mm *MeshMaker) initDangerObj(dataFolder string) {
	// do nothing
}

func (mm *MeshMaker) newDangerObj(dt dangertype.DangerType) *graphic.Mesh {
	var mat *material.Standard
	if mat = mm.doMat[dt]; mat == nil {
		mat = newDangerObjMat(dt)
		mm.doMat[dt] = mat
	}
	var geo *geometry.Geometry
	if geo = mm.doGeo[dt]; geo == nil {
		geo = newDangerObjGeo(dt)
		mm.doGeo[dt] = geo
	}
	return graphic.NewMesh(geo, mat)
}

func (mm *MeshMaker) GetDangerObj(dt dangertype.DangerType, x, y int) *graphic.Mesh {
	mm.doInUse.Inc(dt)
	var mesh *graphic.Mesh
	freeSize := len(mm.doMeshFreeList[dt])
	if freeSize > 0 {
		mesh = mm.doMeshFreeList[dt][freeSize-1]
		mm.doMeshFreeList[dt] = mm.doMeshFreeList[dt][:freeSize-1]
	} else {
		mesh = mm.newDangerObj(dt)
	}
	mesh.SetPositionX(float32(x))
	mesh.SetPositionY(float32(y))
	mesh.SetPositionZ(0.5)
	mesh.SetUserData(dt)
	return mesh
}

func (mm *MeshMaker) PutDangerObj(mesh *graphic.Mesh) {
	dt := mesh.UserData().(dangertype.DangerType)
	mm.doInUse.Dec(dt)
	mm.doMeshFreeList[dt] = append(mm.doMeshFreeList[dt], mesh)
}
