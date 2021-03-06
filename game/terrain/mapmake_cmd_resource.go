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

package terrain

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/goguelike-single/enum/resourcetype"
	"github.com/kasworld/goguelike-single/game/terrain/resourcetile"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/maze2"
	"github.com/kasworld/goguelike-single/lib/scriptparse"
	"github.com/kasworld/walk2d"
)

func cmdResourceAt(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, x, y int
	if err := ca.GetArgs(&rsctl, &amount, &x, &y); err != nil {
		return err
	}
	tr.resourceTileArea[x][y][rsctl] = resourcetile.ResourceValue(amount)
	return nil
}

func cmdResourceHLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, x, w, y int
	if err := ca.GetArgs(&rsctl, &amount, &x, &w, &y); err != nil {
		return err
	}

	fn := func(ax, ay int) bool {
		tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
		return false
	}
	walk2d.HLine(x, x+w, y, fn)
	return nil
}

func cmdResourceVLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, x, y, h int
	if err := ca.GetArgs(&rsctl, &amount, &x, &y, &h); err != nil {
		return err
	}

	fn := func(ax, ay int) bool {
		tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
		return false
	}
	walk2d.VLine(y, y+h, x, fn)
	return nil
}

func cmdResourceLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, x1, y1, x2, y2 int
	if err := ca.GetArgs(&rsctl, &amount, &x1, &y1, &x2, &y2); err != nil {
		return err
	}

	fn := func(ax, ay int) bool {
		tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
		return false
	}
	walk2d.Line(x1, y1, x2, y2, fn)
	return nil
}

func cmdResourceRect(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, x, w, y, h int
	if err := ca.GetArgs(&rsctl, &amount, &x, &w, &y, &h); err != nil {
		return err
	}

	fn := func(ax, ay int) bool {
		tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
		return false
	}
	walk2d.Rect(x, y, x+w, y+h, fn)
	return nil
}

func cmdResourceFillRect(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, x, w, y, h int
	if err := ca.GetArgs(&rsctl, &amount, &x, &w, &y, &h); err != nil {
		return err
	}
	fn := func(ax, ay int) bool {
		tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
		return false
	}
	walk2d.FillHV(x, y, x+w, y+h, fn)
	return nil
}

func cmdResourceFillEllipses(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, x, w, y, h int
	if err := ca.GetArgs(&rsctl, &amount, &x, &w, &y, &h); err != nil {
		return err
	}
	fn := func(ax, ay int) bool {
		tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
		return false
	}
	walk2d.Ellipses(x, y, x+w, y+h, fn)
	return nil
}

func cmdResourceRand(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var mean, stddev, repeat int
	if err := ca.GetArgs(&rsctl, &mean, &stddev, &repeat); err != nil {
		return err
	}
	for i := 0; i < repeat; i++ {
		xpos := tr.rnd.Intn(tr.Xlen)
		ypos := tr.rnd.Intn(tr.Ylen)
		amount := tr.rnd.NormIntRange(mean, stddev)
		tr.resourceTileArea[xpos][ypos][rsctl] = resourcetile.ResourceValue(amount)
	}
	return nil
}

func cmdResourceMazeWall(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, xn, yn int
	var conerFill bool
	var maX, maY, maW, maH int
	if err := ca.GetArgs(&rsctl, &amount, &maX, &maY, &maW, &maH, &xn, &yn, &conerFill); err != nil {
		return err
	}

	m := maze2.New(tr.rnd, xn, yn)
	ma, err := m.ToBoolMatrix(maW, maH, conerFill)
	if err != nil {
		return fmt.Errorf("tr %v %v", tr, err)
	}
	for x, xv := range ma {
		for y, yv := range xv {
			if yv {
				ax, ay := tr.XWrapper.WrapSafe(maX+x), tr.YWrapper.WrapSafe(maY+y)
				tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
			}
		}
	}
	return nil
}

func cmdResourceMazeWalk(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var rsctl resourcetype.ResourceType
	var amount, xn, yn int
	var conerFill bool
	var maX, maY, maW, maH int
	if err := ca.GetArgs(&rsctl, &amount, &maX, &maY, &maW, &maH, &xn, &yn, &conerFill); err != nil {
		return err
	}

	m := maze2.New(tr.rnd, xn, yn)
	ma, err := m.ToBoolMatrix(tr.GetXLen(), tr.GetYLen(), conerFill)
	if err != nil {
		return fmt.Errorf("tr %v %v", tr, err)
	}
	for x, xv := range ma {
		for y, yv := range xv {
			if !yv {
				ax, ay := tr.XWrapper.WrapSafe(maX+x), tr.YWrapper.WrapSafe(maY+y)
				tr.resourceTileArea[ax][ay][rsctl] = resourcetile.ResourceValue(amount)
			}
		}
	}
	return nil
}

func cmdResourceFromPNG(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var name string
	if err := ca.GetArgs(&name); err != nil {
		return err
	}
	if err := tr.resourceTileArea.FromImage(filepath.Join(tr.dataDir, name)); err != nil {
		g2log.Fatal("%v %v", tr, err)
		return err
	}
	return nil
}

func cmdAgeing(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var initrun, msper, resetaftern int
	if err := ca.GetArgs(&initrun, &msper, &resetaftern); err != nil {
		return err
	}
	tr.MSPerAgeing = int64(msper)
	tr.ResetAfterNAgeing = int64(resetaftern)
	tr.resourceTileArea.Ageing(tr.rnd.Intn, initrun)
	return nil
}
