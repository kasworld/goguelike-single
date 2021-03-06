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

package floor

import (
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/slippperydata"
	"github.com/kasworld/goguelike-single/enum/achievetype"
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/dangertype"
	"github.com/kasworld/goguelike-single/enum/returncode"
	"github.com/kasworld/goguelike-single/enum/tile_flag"
	"github.com/kasworld/goguelike-single/enum/turnaction"
	"github.com/kasworld/goguelike-single/enum/turnresulttype"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/activeobject/turnresult"
	"github.com/kasworld/goguelike-single/game/aoactreqrsp"
	"github.com/kasworld/goguelike-single/game/dangerobject"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

func (f *Floor) checkAttackSrc(ao gamei.ActiveObjectI, arr *aoactreqrsp.ActReqRsp) (int, int, way9type.Way9Type) {
	atkdir := arr.Req.Dir
	aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
	if !exist {
		g2log.Error("ao not in currentfloor %v %v", f, ao)
		arr.SetDone(aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
			returncode.ActionProhibited)
		return aox, aoy, atkdir
	}
	if ao.GetTurnData().Condition.TestByCondition(condition.Drunken) {
		turnmod := slippperydata.Drunken[f.rnd.Intn(len(slippperydata.Drunken))]
		atkdir = atkdir.TurnDir(turnmod)
	}
	// add dopoaman near attack

	// check valid attack
	if !atkdir.IsValid() || atkdir == way9type.Center {
		arr.SetDone(aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
			returncode.InvalidDirection)
		return aox, aoy, atkdir
	}
	if f.terrain.GetTileWrapped(aox, aoy).NoBattle() {
		arr.SetDone(aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
			returncode.ActionProhibited)
		return aox, aoy, atkdir
	}
	return aox, aoy, atkdir
}

func (f *Floor) addAttackWide(ao gamei.ActiveObjectI, arr *aoactreqrsp.ActReqRsp) {
	aox, aoy, atkdir := f.checkAttackSrc(ao, arr)
	if arr.Acted {
		return
	}

	for _, dir := range []way9type.Way9Type{atkdir.TurnDir(-1), atkdir, atkdir.TurnDir(1)} {
		dstX, dstY := f.terrain.WrapXY(aox+dir.Dx(), aoy+dir.Dy())
		if f.terrain.GetTiles()[dstX][dstY].NoBattle() {
			continue
		}
		if err := f.doPosMan.AddToXY(
			dangerobject.NewAOAttact(ao, dangertype.WideAttack, aox, aoy),
			dstX, dstY); err != nil {
			g2log.Fatal("fail to AddToXY %v", err)
			continue
		}
	}
	arr.SetDone(
		aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
		returncode.Success)
}

func (f *Floor) addAttackLong(ao gamei.ActiveObjectI, arr *aoactreqrsp.ActReqRsp) {
	aox, aoy, atkdir := f.checkAttackSrc(ao, arr)
	if arr.Acted {
		return
	}

	for i := 1; i < gameconst.AttackLongLen; i++ {
		dstX, dstY := f.terrain.WrapXY(aox+atkdir.Dx()*i, aoy+atkdir.Dy()*i)
		if f.terrain.GetTiles()[dstX][dstY].NoBattle() {
			continue
		}
		if err := f.doPosMan.AddToXY(
			dangerobject.NewAOAttact(ao, dangertype.LongAttack, aox, aoy),
			dstX, dstY); err != nil {
			g2log.Fatal("fail to AddToXY %v", err)
			continue
		}
	}
	arr.SetDone(
		aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
		returncode.Success)
}

func (f *Floor) addBasicAttack(ao gamei.ActiveObjectI, arr *aoactreqrsp.ActReqRsp) {
	aox, aoy, atkdir := f.checkAttackSrc(ao, arr)
	if arr.Acted {
		return
	}
	dstX, dstY := f.terrain.WrapXY(aox+atkdir.Dx(), aoy+atkdir.Dy())
	if f.terrain.GetTiles()[dstX][dstY].NoBattle() {
		arr.SetDone(aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
			returncode.ActionProhibited)
		return
	}
	if err := f.doPosMan.AddToXY(
		dangerobject.NewAOAttact(ao, dangertype.BasicAttack, aox, aoy),
		dstX, dstY); err != nil {
		g2log.Fatal("fail to AddToXY %v", err)
		arr.SetDone(aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
			returncode.ActionCanceled)
		return
	}
	arr.SetDone(
		aoactreqrsp.Act{Act: turnaction.Attack, Dir: atkdir},
		returncode.Success)
}

func (f *Floor) aoAttackActiveObj(src, dst gamei.ActiveObjectI, srcTile, dstTile tile_flag.TileFlag) {

	// attack to invisible ao miss 50%
	if dst.GetTurnData().Condition.TestByCondition(condition.Invisible) && f.rnd.Intn(2) == 0 {
		src.GetAchieveStat().Inc(achievetype.AttackMiss)
		return
	}

	// blind ao attack miss 50%
	if src.GetTurnData().Condition.TestByCondition(condition.Blind) && f.rnd.Intn(2) == 0 {
		src.GetAchieveStat().Inc(achievetype.AttackMiss)
		return
	}

	envbias := f.GetEnvBias()
	srcbias := src.GetTurnData().AttackBias.Add(envbias)
	dstbias := dst.GetTurnData().DefenceBias.Add(envbias)

	atkMod := srcTile.AtkMod()
	defMod := dstTile.DefMod()
	atkValue := srcbias.SelectSkill(f.rnd.Intn(3))
	defValue := dstbias.SelectSkill(f.rnd.Intn(3))
	diffValue := atkValue*atkMod - defValue*defMod +
		src.GetTurnData().Level - dst.GetTurnData().Level +
		f.rnd.NormFloat64Range(gameconst.ActiveObjBaseBiasLen, 0)
	atkSuccess := diffValue > 0
	atkCritical := false
	rndValue := f.rnd.Intn(20)
	if rndValue == 0 {
		atkSuccess = false
	} else if rndValue == 19 {
		atkSuccess = true
	}
	if atkSuccess && f.rnd.Intn(20) == 19 {
		atkCritical = true
	}

	src.GetAchieveStat().Inc(achievetype.Attack)
	dst.GetAchieveStat().Inc(achievetype.Attacked)

	if !atkSuccess {
		src.GetAchieveStat().Inc(achievetype.AttackMiss)
		return
	}

	src.GetAchieveStat().Inc(achievetype.AttackHit)
	if diffValue < 0 {
		diffValue = -diffValue
	}

	damage := diffValue

	if atkCritical {
		damage *= 2
		src.GetAchieveStat().Inc(achievetype.AttackCritical)
	}

	src.AppendTurnResult(turnresult.New(turnresulttype.AttackTo, dst, damage))
	dst.AppendTurnResult(turnresult.New(turnresulttype.AttackedFrom, src, damage))

	src.GetAchieveStat().Add(achievetype.DamageTotalGive, damage)
	src.GetAchieveStat().SetIfGt(achievetype.DamageMaxGive, damage)
	dst.GetAchieveStat().Add(achievetype.DamageTotalRecv, damage)
	dst.GetAchieveStat().SetIfGt(achievetype.DamageMaxRecv, damage)

	src.AddBattleExp(damage * gameconst.ActiveObjExp_Damage)
}

func (f *Floor) foRotateLineAttack(do *dangerobject.DangerObject, dstao gamei.ActiveObjectI, dstx, dsty int) {
	hpdamage := do.AffectRate * dstao.GetTurnData().HPMax
	dstao.AppendTurnResult(turnresult.New(turnresulttype.AttackedFrom, do.Owner, hpdamage))
}

func (f *Floor) foMineExplodeAttack(do *dangerobject.DangerObject, dstao gamei.ActiveObjectI, dstx, dsty int) {
	hpdamage := do.AffectRate * dstao.GetTurnData().HPMax
	dstao.AppendTurnResult(turnresult.New(turnresulttype.AttackedFrom, do.Owner, hpdamage))
}
