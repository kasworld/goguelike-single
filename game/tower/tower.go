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
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike-single/config/dataversion"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/gamedata"
	"github.com/kasworld/goguelike-single/config/goguelikeconfig"
	"github.com/kasworld/goguelike-single/enum/towerachieve_vector"
	"github.com/kasworld/goguelike-single/game/activeobject"
	"github.com/kasworld/goguelike-single/game/aoexpsort"
	"github.com/kasworld/goguelike-single/game/aoid2activeobject"
	"github.com/kasworld/goguelike-single/game/aoid2floor"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/floormanager"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/game/glclient"
	"github.com/kasworld/goguelike-single/game/towerscript"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/loadlines"
	"github.com/kasworld/intervalduration"
	"github.com/kasworld/version"
	"github.com/kasworld/weblib/retrylistenandserve"
)

var _ gamei.TowerI = &Tower{}

func (tw *Tower) String() string {
	return fmt.Sprintf("Tower[Seed:%v]",
		tw.seed,
	)
}

type Tower struct {
	doClose func() `prettystring:"hide"`

	// cmd to tower channel
	cmdCh chan interface{} `prettystring:"hide"`
	// tower cmd stats
	cmdActStat *actpersec.ActPerSec `prettystring:"simple"`
	// trigger tower turn channel
	turnCh   chan time.Time                     `prettystring:"hide"`
	turnStat *actpersec.ActPerSec               `prettystring:"simple"`
	interDur *intervalduration.IntervalDuration `prettystring:"simple"`

	// init data
	config     *goguelikeconfig.GoguelikeConfig
	seed       int64
	rnd        *g2rand.G2Rand `prettystring:"hide"`
	biasFactor [3]int64       `prettystring:"simple"`
	startTime  time.Time      `prettystring:"simple"`

	// mamagers
	floorMan         *floormanager.FloorManager                  `prettystring:"simple"`
	ao2Floor         *aoid2floor.ActiveObjID2Floor               `prettystring:"simple"`
	id2ao            *aoid2activeobject.ActiveObjID2ActiveObject `prettystring:"simple"`
	aoExpRanking     aoexpsort.ByExp                             `prettystring:"simple"`
	towerAchieveStat *towerachieve_vector.TowerAchieveVector     `prettystring:"simple"`

	gameInfo *csprotocol.GameInfo

	// single player
	playerAO *activeobject.ActiveObject
	// client to tower packet channel
	c2tCh chan *csprotocol.Packet `prettystring:"hide"`
	// tower to client packet channel
	t2cCh chan *csprotocol.Packet `prettystring:"hide"`

	adminWeb *http.Server `prettystring:"simple"`
}

func New(config *goguelikeconfig.GoguelikeConfig) *Tower {
	fmt.Printf("%v\n", config.StringForm())

	tw := &Tower{
		id2ao:            aoid2activeobject.New("ActiveObject working"),
		config:           config,
		turnStat:         actpersec.New(),
		interDur:         intervalduration.New(""),
		cmdActStat:       actpersec.New(),
		towerAchieveStat: new(towerachieve_vector.TowerAchieveVector),
		c2tCh:            make(chan *csprotocol.Packet, gameconst.SendBufferSize),
		t2cCh:            make(chan *csprotocol.Packet, gameconst.SendBufferSize),
		startTime:        time.Now(),
	}

	tw.seed = int64(config.Seed)
	if tw.seed <= 0 {
		tw.seed = time.Now().UnixNano()
	}
	tw.rnd = g2rand.NewWithSeed(int64(tw.seed))

	tw.doClose = func() {
		g2log.Fatal("Too early doClose call %v", tw)
	}

	var err error

	gamedata.ActiveObjNameList, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().ServerDataFolder, "ainames.txt"),
	)
	if err != nil {
		g2log.Fatal("load ainame fail %v", err)
		return nil
	}

	gamedata.ChatData, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().ServerDataFolder, "chatdata.txt"),
	)
	if err != nil {
		g2log.Fatal("load chatdata fail %v", err)
		return nil
	}

	tScript, err := towerscript.LoadJSON(
		tw.config.MakeTowerFileFullpath(),
	)
	if err != nil {
		return nil
	}

	tw.ao2Floor = aoid2floor.New(tw)
	tw.biasFactor = tw.NewRandFactor()

	tw.floorMan = floormanager.New(tScript, tw)
	if err := tw.floorMan.Init(tw.rnd); err != nil {
		g2log.Fatal("floorman init fail %v", err)
		return nil
	}

	tw.gameInfo = &csprotocol.GameInfo{
		Version:       version.GetVersion(),
		DataVersion:   dataversion.DataVersion,
		StartTime:     tw.startTime,
		TowerSeed:     tw.seed,
		TowerName:     tw.config.ScriptFilename,
		Factor:        tw.biasFactor,
		TotalFloorNum: tw.floorMan.GetFloorCount(),
		NickName:      tw.Config().NickName,
	}

	fmt.Printf("%v\n", tw.gameInfo.StringForm())
	fmt.Printf("WebAdmin  : %v:%v id:%v pass:%v\n",
		"http://localhost", tw.config.AdminPort, tw.config.WebAdminID, tw.config.WebAdminPass)

	return tw
}

func (tw *Tower) Run() {
	defer tw.Cleanup()

	ctx, closeCtx := context.WithCancel(context.Background())
	tw.doClose = closeCtx

	defer closeCtx()

	totalaocount := 0
	for _, f := range tw.floorMan.GetFloorList() {
		totalaocount += f.GetTerrain().GetActiveObjCount()
	}
	g2log.Debug("Total system ActiveObj in tower %v", totalaocount)

	queuesize := totalaocount * 2
	if queuesize <= 0 {
		queuesize = 100
	}
	tw.cmdCh = make(chan interface{}, queuesize)
	if tw.cmdCh == nil {
		g2log.Fatal("fail to make cmdCh %v", queuesize)
		return
	}
	tw.turnCh = make(chan time.Time, 100)
	if tw.turnCh == nil {
		g2log.Fatal("fail to make turnCh %v", queuesize)
		return
	}

	// start turn tower
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case now := <-tw.turnCh:
				for len(tw.turnCh) > cap(tw.turnCh)/2 {
					g2log.Warn("remove dup turn req %v/%v", len(tw.turnCh), cap(tw.turnCh))
					now = <-tw.turnCh
				}
				tw.Turn(now)
			}
		}
	}()

	// start floor
	for _, f := range tw.floorMan.GetFloorList() {
		go func(f gamei.FloorI) {
			f.Run(ctx, queuesize)
			closeCtx()
		}(f)
	}

	// add ao to tower/floor
	for _, f := range tw.floorMan.GetFloorList() {
		for i := 0; i < f.GetTerrain().GetActiveObjCount(); i++ {
			ao := activeobject.NewSystemActiveObj(tw.rnd.Int63(), f, tw.towerAchieveStat)
			if err := tw.ao2Floor.ActiveObjEnterTower(f, ao); err != nil {
				g2log.Error("%v", err)
				continue
			}
			if err := tw.id2ao.Add(ao); err != nil {
				g2log.Error("%v", err)
			}
		}
	}

	tw.initAdminWeb()
	go retrylistenandserve.RetryListenAndServe(tw.adminWeb, g2log.GlobalLogger, "serveAdminWeb")

	go tw.handle_c2tch()
	tw.initPlayer()

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()
	rankMakeTk := time.NewTicker(1 * time.Second)
	defer rankMakeTk.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-timerInfoTk.C:
			tw.cmdActStat.UpdateLap()
			tw.turnStat.UpdateLap()
			if len(tw.cmdCh) >= cap(tw.cmdCh) {
				g2log.Fatal("cmdCh full %v/%v", len(tw.cmdCh), cap(tw.cmdCh))
				break loop
			}
			if len(tw.c2tCh) >= cap(tw.c2tCh) {
				g2log.Fatal("c2tCh full %v/%v", len(tw.c2tCh), cap(tw.c2tCh))
				break loop
			}
			if len(tw.t2cCh) >= cap(tw.t2cCh) {
				g2log.Fatal("t2cCh full %v/%v", len(tw.t2cCh), cap(tw.t2cCh))
				break loop
			}
		case <-rankMakeTk.C:
			go tw.makeActiveObjExpRank()
		}
	}
	tw.doClose()
}

func (tw *Tower) Cleanup() {
	tw.id2ao.Cleanup()
	tw.ao2Floor.Cleanup()
	for _, f := range tw.floorMan.GetFloorList() {
		f.Cleanup()
	}
	tw.floorMan.Cleanup()
}

func (tw *Tower) initPlayer() {
	// prepare player ao enter tower
	// new ao
	newAO := activeobject.NewUserActiveObj(
		tw.rnd.Int63(),
		tw.GetFloorManager().GetStartFloor(),
		tw.Config().NickName,
		tw.towerAchieveStat,
	)
	tw.playerAO = newAO
	if err := tw.ao2Floor.ActiveObjEnterTower(tw.playerAO.GetHomeFloor(), tw.playerAO); err != nil {
		g2log.Error("%v", err)
		return
	}
	if err := tw.id2ao.Add(tw.playerAO); err != nil {
		g2log.Fatal("%v", err)
	}
	tw.gameInfo.ActiveObjUUID = tw.playerAO.GetUUID()

	tw.turnCh <- time.Now() // start init turn

	//run client
	go func() {
		cl := glclient.New(tw.config, tw.gameInfo, tw.c2tCh, tw.t2cCh)
		if err := cl.Run(); err != nil {
			g2log.Error("%v", err)
		}
		tw.doClose()
	}()
}

func (tw *Tower) ProcessAllCmds() {
	for len(tw.cmdCh) > 0 {
		tw.processCmd(<-tw.cmdCh)
	}
}
func (tw *Tower) Turn(now time.Time) {
	act := tw.interDur.BeginAct()
	defer func() {
		act.End()
	}()
	tw.turnStat.Inc()

	var ws sync.WaitGroup

	// process turn
	for _, f := range tw.floorMan.GetFloorList() {
		ws.Add(1)
		go func(f gamei.FloorI) {
			f.Turn(now)
			ws.Done()
		}(f)
	}
	ws.Wait()

	// process all cmds
	tw.ProcessAllCmds()
	for _, f := range tw.floorMan.GetFloorList() {
		ws.Add(1)
		go func(f gamei.FloorI) {
			f.ProcessAllCmds()
			ws.Done()
		}(f)
	}
	ws.Wait()

}

func (tw *Tower) makeActiveObjExpRank() {
	rtn := tw.id2ao.GetAllList()
	for _, v := range rtn {
		v.UpdateExpCopy()
	}
	aoexpsort.ByExp(rtn).Sort()
	tw.aoExpRanking = rtn
}

func (tw *Tower) NewRandFactor() [3]int64 {
	// st := 0
	lenPrimes := len(gamedata.Primes)
	rtn := [3]int64{}
loop:
	for i := 0; i < 3; {
		rtn[i] = gamedata.Primes[lenPrimes/2+tw.rnd.Intn(lenPrimes/2)]
		for j := 0; j < i; j++ {
			if rtn[j] == rtn[i] {
				continue loop
			}
		}
		i++
	}
	for i, _ := range rtn {
		if tw.rnd.Intn(2) == 0 {
			rtn[i] = -rtn[i]
		}
	}
	return rtn
}
