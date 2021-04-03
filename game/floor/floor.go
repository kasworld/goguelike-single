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
	return fmt.Sprintf("Floor[%v Seed:%v]",
		f.GetName(),
		f.seed,
	)
}

type Floor struct {
	rnd *g2rand.G2Rand `prettystring:"hide"`
	log *g2log.LogBase `prettystring:"hide"`

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

	interDur          *intervalduration.IntervalDuration `prettystring:"simple"`
	statPacketObjOver *actpersec.ActPerSec               `prettystring:"simple"`
	cmdActStat        *actpersec.ActPerSec               `prettystring:"simple"`

	// for actturn data
	recvRequestCh chan interface{}

	aiWG sync.WaitGroup // for ai run
}

func New(seed int64, ts []string, tw gamei.TowerI) *Floor {
	f := &Floor{
		log:               tw.Log(),
		tower:             tw,
		seed:              seed,
		rnd:               g2rand.NewWithSeed(seed),
		interDur:          intervalduration.New(""),
		statPacketObjOver: actpersec.New(),
		cmdActStat:        actpersec.New(),
	}
	f.terrain = terrain.New(f.rnd.Int63(), ts, f.tower.Config().DataFolder, f.log)
	return f
}

func (f *Floor) Cleanup() {
	f.log.TraceService("Start Cleanup Floor %v", f.GetName())
	defer func() { f.log.TraceService("End Cleanup Floor %v", f.GetName()) }()

	f.aoPosMan.Cleanup()
	f.poPosMan.Cleanup()
	f.terrain.Cleanup()
	f.doPosMan.Cleanup()
}

// Init bi need for randomness
func (f *Floor) Init() error {
	f.log.TraceService("Start Init %v", f)
	defer func() { f.log.TraceService("End Init %v", f) }()

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
	f.log.TraceService("start Run %v", f)
	defer func() { f.log.TraceService("End Run %v", f) }()

	f.recvRequestCh = make(chan interface{}, queuesize)

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerInfoTk.C:
			f.statPacketObjOver.UpdateLap()
			f.cmdActStat.UpdateLap()
			if len(f.recvRequestCh) > cap(f.recvRequestCh)/2 {
				f.log.Fatal("Floor %v %v reqch overloaded %v/%v",
					f.terrain.Name, f.GetName(),
					len(f.recvRequestCh), cap(f.recvRequestCh))
			}
			if len(f.recvRequestCh) >= cap(f.recvRequestCh) {
				break loop
			}

		case data := <-f.recvRequestCh:
			f.processCmd(data)

		}
	}
}

func (f *Floor) processAgeing() {
	if atomic.CompareAndSwapInt32(&f.inAgeing, 0, 1) {
		defer atomic.AddInt32(&f.inAgeing, -1)

		var err error
		err = f.terrain.Ageing()
		if err != nil {
			f.log.Fatal("%v %v", f, err)
			return
		}
		NotiAgeing := f.ToPacket_NotiAgeing()
		for _, v := range f.aoPosMan.GetAllList() {
			ao := v.(gamei.ActiveObjectI)
			ao.SetNeedTANoti()
			// send ageing noti
			if aoconn := ao.GetClientConn(); aoconn != nil {
				if err := aoconn.SendNotiPacket(c2t_idnoti.Ageing,
					NotiAgeing,
				); err != nil {
					f.log.Error("%v %v %v", f, ao, err)
				}
			}
		}
	} else {
		f.log.Fatal("processAgeing skipped %v", f)
	}
}
