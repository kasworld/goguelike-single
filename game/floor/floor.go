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
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/enum/aotype"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/game/terrain"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/uuidposman_slice"
	"github.com/kasworld/goguelike-single/lib/uuidposmani"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/intervalduration"
)

var _ gamei.FloorI = &Floor{}

func (f *Floor) String() string {
	return fmt.Sprintf("Floor[%v]", f.GetName())
}

type Floor struct {
	rnd *g2rand.G2Rand `prettystring:"hide"`

	tower       gamei.TowerI
	w           int
	h           int
	bias        bias.Bias        `prettystring:"simple"`
	terrain     *terrain.Terrain `prettystring:"simple"`
	initialized bool
	seed        int64

	// async ageing one at a time
	inAgeing int32

	aoPosMan uuidposmani.UUIDPosManI `prettystring:"simple"`
	poPosMan uuidposmani.UUIDPosManI `prettystring:"simple"`
	foPosMan uuidposmani.UUIDPosManI `prettystring:"simple"`
	doPosMan uuidposmani.UUIDPosManI `prettystring:"simple"`

	interDur   *intervalduration.IntervalDuration `prettystring:"simple"`
	cmdActStat *actpersec.ActPerSec               `prettystring:"simple"`

	// for actturn data
	cmdCh chan interface{}

	aiWG sync.WaitGroup // for ai run
}

func New(seed int64, ts []string, tw gamei.TowerI) *Floor {
	f := &Floor{
		tower:      tw,
		seed:       seed,
		rnd:        g2rand.NewWithSeed(seed),
		interDur:   intervalduration.New(""),
		cmdActStat: actpersec.New(),
	}
	f.terrain = terrain.New(f.rnd.Int63(), ts, f.tower.Config().ServerDataFolder)
	return f
}

func (f *Floor) Cleanup() {
	f.aoPosMan.Cleanup()
	f.poPosMan.Cleanup()
	f.terrain.Cleanup()
	f.doPosMan.Cleanup()
}

// Init bi need for randomness
func (f *Floor) Init() error {
	if err := f.terrain.Init(); err != nil {
		return fmt.Errorf("fail to make terrain %v", err)
	}
	if f.terrain.GetName() == "" {
		return nil // skip no name floor terrain
	}
	f.w, f.h = f.GetTerrain().GetXYLen()
	f.aoPosMan = uuidposman_slice.New(f.w, f.h)
	f.poPosMan = uuidposman_slice.New(f.w, f.h)
	f.foPosMan = f.terrain.GetFieldObjPosMan()
	f.doPosMan = uuidposman_slice.New(f.w, f.h)
	f.bias = bias.Bias{
		f.rnd.Float64() - 0.5,
		f.rnd.Float64() - 0.5,
		f.rnd.Float64() - 0.5,
	}.MakeAbsSumTo(gameconst.FloorBaseBiasLen)

	f.initialized = true
	return nil
}

func (f *Floor) Run(ctx context.Context, queuesize int) {
	f.cmdCh = make(chan interface{}, queuesize)
	if f.cmdCh == nil {
		g2log.Fatal("%v fail to make cmdCh %v", f, queuesize)
		return
	}

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerInfoTk.C:
			f.cmdActStat.UpdateLap()
			if len(f.cmdCh) > cap(f.cmdCh)/2 {
				g2log.Fatal("Floor %v cmdch overloaded %v/%v",
					f.GetName(), len(f.cmdCh), cap(f.cmdCh))
			}
			if len(f.cmdCh) >= cap(f.cmdCh) {
				g2log.Fatal("Floor %v cmdch overloaded %v/%v",
					f.GetName(), len(f.cmdCh), cap(f.cmdCh))
				break loop
			}
		}
	}
}

func (f *Floor) Turn(now time.Time) {
	act := f.interDur.BeginAct()
	defer func() {
		act.End()
	}()

	for len(f.cmdCh) > 0 {
		f.processCmd(<-f.cmdCh)
	}
	f.processTurn(now)
	turnPerAge := f.terrain.GetMSPerAgeing() / 1000
	if turnPerAge > 0 && f.interDur.GetCount()%int(turnPerAge) == 0 {
		f.processAgeing()
	}
}

func (f *Floor) processAgeing() {
	if atomic.CompareAndSwapInt32(&f.inAgeing, 0, 1) {
		defer atomic.AddInt32(&f.inAgeing, -1)

		var err error
		err = f.terrain.Ageing()
		if err != nil {
			g2log.Fatal("%v %v", f, err)
			return
		}
		NotiAgeing := f.ToPacket_NotiAgeing()
		for _, v := range f.aoPosMan.GetAllList() {
			ao := v.(gamei.ActiveObjectI)
			ao.SetNeedTANoti()
			// send ageing noti
			if ao.GetActiveObjType() == aotype.User {
				f.tower.SendNoti(
					c2t_idnoti.Ageing,
					NotiAgeing,
				)
			}
		}
	} else {
		g2log.Fatal("processAgeing skipped %v", f)
	}
}
