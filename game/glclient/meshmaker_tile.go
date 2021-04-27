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
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/kasworld/goguelike-single/enum/tile"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

// mange tile

var tileAttrib = [tile.Tile_Count]struct {
	tranparent bool
	zPos       float32
	height     float32
}{
	tile.Swamp:  {false, -0.2, 0.2},
	tile.Soil:   {false, -0.2, 0.2},
	tile.Stone:  {false, -0.2, 0.2},
	tile.Sand:   {false, -0.2, 0.2},
	tile.Sea:    {false, -0.3, 0.2},
	tile.Magma:  {false, -0.3, 0.2},
	tile.Ice:    {false, -0.1, 0.1},
	tile.Grass:  {false, -0.0, 0.2},
	tile.Tree:   {false, -0.0, 0.3},
	tile.Road:   {false, -0.0, 0.1},
	tile.Room:   {false, -0.0, 0.1},
	tile.Wall:   {false, -0.0, 1.0},
	tile.Window: {true, -0.0, 1.0},
	tile.Door:   {true, -0.0, 1.0},
	tile.Fog:    {true, 0.1, 1.0},
	tile.Smoke:  {true, 0.1, 1.0},
}

func loadTileTexture(texFilename string) *texture.Texture2D {
	// Open image file
	file, err := os.Open(texFilename)
	if err != nil {
		g2log.Fatal("Error loading texture: %s", err)
		return nil
	}
	defer file.Close()

	// Decodes image
	img, _, err := image.Decode(file)
	if err != nil {
		g2log.Fatal("Error loading texture: %s", err)
		return nil
	}

	// Converts image to RGBA format
	texSize := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{64, 64},
	}
	rgba := image.NewRGBA(texSize)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		err := fmt.Errorf("unsupported stride")
		g2log.Fatal("Error loading texture: %s", err)
		return nil
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	tex := texture.NewTexture2DFromRGBA(rgba)
	tex.SetFromRGBA(rgba)

	// tex.SetWrapS(gls.REPEAT)
	// tex.SetWrapT(gls.REPEAT)
	// tex.SetRepeat(fw/128, fh/128)
	return tex
}

func (mm *MeshMaker) initTile(dataFolder string, initSize int) {
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
}

func (mm *MeshMaker) newTile(tl tile.Tile) *graphic.Mesh {
	mesh := graphic.NewMesh(mm.tileGeo[tl], mm.tileMat[tl])
	return mesh
}

func (mm *MeshMaker) GetTile(tl tile.Tile, x, y int) *graphic.Mesh {
	mm.tileInUse.Inc(tl)
	var mesh *graphic.Mesh
	freeSize := len(mm.tileMeshFreeList[tl])
	if freeSize > 0 {
		mesh = mm.tileMeshFreeList[tl][freeSize-1]
		mm.tileMeshFreeList[tl] = mm.tileMeshFreeList[tl][:freeSize-1]
	} else {
		mesh = mm.newTile(tl)
	}
	mesh.SetPositionX(float32(x))
	mesh.SetPositionY(float32(y))
	mesh.SetPositionZ(tileAttrib[tl].zPos)
	return mesh
}

func (mm *MeshMaker) PutTile(tl tile.Tile, mesh *graphic.Mesh) {
	mm.tileInUse.Dec(tl)
	mm.tileMeshFreeList[tl] = append(mm.tileMeshFreeList[tl], mesh)
}
