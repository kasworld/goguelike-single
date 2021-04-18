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

package aoactreqrsp

import (
	"fmt"

	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/condition_flag"
	"github.com/kasworld/goguelike-single/enum/returncode"
	"github.com/kasworld/goguelike-single/enum/turnaction"
	"github.com/kasworld/goguelike-single/enum/way9type"
)

type Act struct {
	Act  turnaction.TurnAction
	Dir  way9type.Way9Type
	UUID string
}

func (act Act) CalcAPByActAndCondition(cndflag condition_flag.ConditionFlag) float64 {
	turn2need := act.Act.NeedTurn()
	if cndflag.TestByCondition(condition.Slow) {
		turn2need *= 2
	}
	if cndflag.TestByCondition(condition.Haste) {
		turn2need /= 2
	}
	if act.Act == turnaction.Move && act.Dir != way9type.Center {
		turn2need *= act.Dir.Len()
	}
	return turn2need
}

type ActReqRsp struct {
	Req   Act // act requested
	Done  Act // act done
	Acted bool

	Error returncode.ReturnCode
}

func (arr ActReqRsp) IsSuccess() bool {
	return arr.Acted && arr.Error == returncode.Success
}

func (arr *ActReqRsp) SetDone(done Act, Error returncode.ReturnCode) {
	if arr.Acted {
		fmt.Printf("already Acted %v %+v", Error, arr)
	}
	arr.Done = done
	arr.Error = Error
	arr.Acted = true
}
