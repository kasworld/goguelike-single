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

package tower

import (
	"fmt"
	"time"

	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/goguelikeconfig"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/gamei"
)

// managers

func (tw *Tower) GetFloorManager() gamei.FloorManagerI {
	return tw.floorMan
}

func (tw *Tower) GetExpRanking() []gamei.ActiveObjectI {
	return tw.aoExpRanking
}

// attribute get/set

func (tw *Tower) GetCmdCh() chan<- interface{} {
	return tw.cmdCh
}

func (tw *Tower) CmdChState() string {
	return fmt.Sprintf("%v/%v", len(tw.cmdCh), cap(tw.cmdCh))
}

func (tw *Tower) GetUUID() string {
	return tw.uuid
}

func (tw *Tower) GetRunDur() time.Duration {
	return time.Now().Sub(tw.startTime)
}

func (tw *Tower) GetBias() bias.Bias {
	rtn := bias.MakeBiasByProgress(tw.biasFactor, tw.GetRunDur().Seconds(), gameconst.TowerBaseBiasLen)
	return rtn
}

func (tw *Tower) Config() *goguelikeconfig.GoguelikeConfig {
	return tw.sconfig
}

func (tw *Tower) PauseListenClient() {
	tw.listenClientPaused = true
}

func (tw *Tower) ResumeListenClient() {
	tw.listenClientPaused = false
}

func (tw *Tower) IsListenClientPaused() bool {
	return tw.listenClientPaused
}
