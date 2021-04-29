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

// Package activeobject manage actor object
package activeobject

import (
	"sync/atomic"

	"github.com/kasworld/goguelike-single/enum/aiplan"
	"github.com/kasworld/goguelike-single/lib/g2log"
)

func (ao *ActiveObject) IsAIUse() bool {
	return ao.isAIInUse
}
func (ao *ActiveObject) SetUseAI(b bool) {
	ao.isAIInUse = b
	if aio := ao.ai; aio != nil {
		go ao.ResetPlan(aio)
	}
}

func (ao *ActiveObject) RunAI(TurnCount int) {
	if aio := ao.ai; aio != nil && ao.IsAIUse() {
		ao.ActTurn(aio, TurnCount)
	}
}

func (ao *ActiveObject) ActTurn(sai *ServerAIState, TurnCount int) {
	if atomic.CompareAndSwapInt32(&sai.isAIRunning, 0, 1) {
		defer atomic.AddInt32(&sai.isAIRunning, -1)
		ao.actTurn(sai, TurnCount)
	} else {
		g2log.Warn("skip ai Turn %v %v", ao, sai.isAIRunning)
	}
}

func (ao *ActiveObject) actTurn(sai *ServerAIState, TurnCount int) {
	if !ao.IsAlive() {
		return
	}
	if ao.NeedChangePlan(sai, ao.GetTurnActReqRsp()) {
		sai.planRemainCount = 0
	}

	if ao.GetAP() < 0 { // skip
		return
	}

	act := sai.InterDur.BeginAct()
	defer func() {
		act.End()
	}()

	sai.TurnCount = TurnCount

	aox, aoy, exist := ao.currentFloor.GetActiveObjPosMan().GetXYByUUID(ao.GetUUID())
	if !exist {
		// g2log.Error("ao not in currentfloor %v %v", ao, ao.currentFloor)
		return
	}
	sai.aox, sai.aoy = aox, aoy

	// attacked?
	if sai.RunningPlanList.GetCurrentPlan() != aiplan.Attack &&
		sai.RunningPlanList.GetCurrentPlan() != aiplan.Revenge &&
		ao.aoAttackLast(sai) != nil {

		sai.RunningPlanList.Move2Front(aiplan.Revenge)
		ao.selectPlan(sai)
	} else {
		// need select new plan?
		if sai.planRemainCount <= 0 {
			ao.selectPlan(sai)
		}
	}
	if sai.planRemainCount > 0 {
		continuePlan := allPlanList[sai.RunningPlanList.GetCurrentPlan()].ActFn(ao, sai)
		if continuePlan {
			sai.planRemainCount--
		} else {
			sai.planRemainCount = 0
		}
	}
}
