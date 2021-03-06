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

package activeobject

import (
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/statusoptype"
	"github.com/kasworld/goguelike-single/game/activeobject/aoturndata"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

// apply status effect one time
func (ao *ActiveObject) applyOpArg(
	oldData *aoturndata.ActiveObjTurnData,
	newData *aoturndata.ActiveObjTurnData,
	oparg statusoptype.OpArg) {
	switch oparg.Op {
	default:
		g2log.Error("unknown statusop %v, %v", oparg.Op, oparg.Arg)

	case statusoptype.None:
		// ignore None effect

	case statusoptype.AddHP:
		arg, ok := oparg.Arg.(float64)
		if !ok {
			g2log.Fatal("invalid type arg %v %v %T",
				oparg.Op, oparg.Arg, oparg.Arg)
			return
		}
		ao.hp += arg

	case statusoptype.AddSP:
		arg, ok := oparg.Arg.(float64)
		if !ok {
			g2log.Fatal("invalid type arg %v %v %T",
				oparg.Op, oparg.Arg, oparg.Arg)
			return
		}
		ao.sp += arg

	case statusoptype.AddHPRate:
		arg, ok := oparg.Arg.(float64)
		if !ok {
			g2log.Fatal("invalid type arg %v %v %T",
				oparg.Op, oparg.Arg, oparg.Arg)
			return
		}
		ao.hp += oldData.HPMax * arg

	case statusoptype.AddSPRate:
		arg, ok := oparg.Arg.(float64)
		if !ok {
			g2log.Fatal("invalid type arg %v %v %T",
				oparg.Op, oparg.Arg, oparg.Arg)
			return
		}
		ao.sp += oldData.SPMax * arg

	case statusoptype.RndFaction:
		ao.currentBias = bias.Bias{
			ao.rnd.Float64() - 0.5,
			ao.rnd.Float64() - 0.5,
			ao.rnd.Float64() - 0.5,
		}.MakeAbsSumTo(gameconst.ActiveObjBaseBiasLen)

	case statusoptype.IncFaction:
		arg, ok := oparg.Arg.(int)
		if !ok {
			g2log.Fatal("invalid type arg %v %v %T",
				oparg.Op, oparg.Arg, oparg.Arg)
			return
		}
		ft := ao.currentBias.NearFaction()
		ft = factiontype.FactionType(factiontype.Wraper.WrapSafe(int(ft) + arg))
		ao.currentBias = bias.Bias(
			ft.FactorBase(),
		).MakeAbsSumTo(gameconst.ActiveObjBaseBiasLen)

	case statusoptype.SetFaction:
		switch arg := oparg.Arg.(type) {
		default:
			g2log.Fatal("invalid type arg %v %v %T",
				oparg.Op, oparg.Arg, oparg.Arg)
			return
		case factiontype.FactionType:
			ft := arg
			ao.currentBias = bias.Bias(
				ft.FactorBase(),
			).MakeAbsSumTo(gameconst.ActiveObjBaseBiasLen)
		case int:
			ft := factiontype.FactionType(factiontype.Wraper.WrapSafe(arg))
			ao.currentBias = bias.Bias(
				ft.FactorBase(),
			).MakeAbsSumTo(gameconst.ActiveObjBaseBiasLen)
		}

	case statusoptype.ResetFaction:
		ao.currentBias = bias.Bias(
			ao.bornFaction.FactorBase(),
		).MakeAbsSumTo(gameconst.ActiveObjBaseBiasLen)

	case statusoptype.NegBias:
		ao.currentBias = ao.currentBias.Neg()

	case statusoptype.RotateBiasRight:
		ao.currentBias = ao.currentBias.RotateRight()

	case statusoptype.RotateBiasLeft:
		ao.currentBias = ao.currentBias.RotateLeft()

	case statusoptype.ForgetFloor:
		if err := ao.ForgetFloorByName(ao.currentFloor.GetName()); err != nil {
			g2log.Fatal("%v", err)
		}

	case statusoptype.ForgetOneFloor:
		if err := ao.forgetAnyFloor(); err != nil {
			g2log.Fatal("%v", err)
		}

	case statusoptype.ModSight:
		arg, ok := oparg.Arg.(float64)
		if !ok {
			g2log.Fatal("invalid type arg %v %v %T",
				oparg.Op, oparg.Arg, oparg.Arg)
			return
		}
		newData.Sight += arg

	case statusoptype.SetCondition:
		cnd := oparg.Arg.(condition.Condition)
		newData.Condition.SetByCondition(cnd)
		ao.conditionStat.Inc(cnd)
	}
}
