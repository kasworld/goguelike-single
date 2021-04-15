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
	"github.com/kasworld/goguelike-single/game/floormanager"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/game/glclient"
	"github.com/kasworld/goguelike-single/game/towerscript"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/loadlines"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_version"
	"github.com/kasworld/recordduration"
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
	doClose func()         `prettystring:"hide"`
	rnd     *g2rand.G2Rand `prettystring:"hide"`

	cmdCh chan interface{}

	config     *goguelikeconfig.GoguelikeConfig
	seed       int64
	biasFactor [3]int64  `prettystring:"simple"`
	startTime  time.Time `prettystring:"simple"`

	floorMan     *floormanager.FloorManager                  `prettystring:"simple"`
	ao2Floor     *aoid2floor.ActiveObjID2Floor               `prettystring:"simple"`
	id2ao        *aoid2activeobject.ActiveObjID2ActiveObject `prettystring:"simple"`
	aoExpRanking aoexpsort.ByExp                             `prettystring:"simple"`

	gameInfo *c2t_obj.GameInfo

	// single player
	playerAO *activeobject.ActiveObject

	turnStat         *actpersec.ActPerSec                    `prettystring:"simple"`
	towerAchieveStat *towerachieve_vector.TowerAchieveVector `prettystring:"simple"`

	// tower cmd stats
	cmdActStat *actpersec.ActPerSec `prettystring:"simple"`

	adminWeb *http.Server `prettystring:"simple"`

	// client to tower packet channel
	c2tCh chan *c2t_packet.Packet
	// tower to client packet channel
	t2cCh chan *c2t_packet.Packet
}

func New(config *goguelikeconfig.GoguelikeConfig) *Tower {
	fmt.Printf("%v\n", config.StringForm())

	tw := &Tower{
		id2ao:            aoid2activeobject.New("ActiveObject working"),
		config:           config,
		turnStat:         actpersec.New(),
		cmdActStat:       actpersec.New(),
		towerAchieveStat: new(towerachieve_vector.TowerAchieveVector),
	}

	tw.seed = int64(config.Seed)
	if tw.seed <= 0 {
		tw.seed = time.Now().UnixNano()
	}
	tw.rnd = g2rand.NewWithSeed(int64(tw.seed))

	tw.doClose = func() {
		g2log.Fatal("Too early doClose call %v", tw)
	}
	return tw
}

func (tw *Tower) ServiceInit() error {
	rd := recordduration.New(tw.String())

	g2log.TraceService("Start ServiceInit %v %v", tw, rd)
	defer func() {
		g2log.TraceService("End ServiceInit %v %v", tw, rd)
		fmt.Println(rd)
	}()

	g2log.TraceService("%v", tw.config.StringForm())

	var err error

	gamedata.ActiveObjNameList, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().ServerDataFolder, "ainames.txt"),
	)
	if err != nil {
		g2log.Fatal("load ainame fail %v", err)
		return err
	}

	gamedata.ChatData, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().ServerDataFolder, "chatdata.txt"),
	)
	if err != nil {
		g2log.Fatal("load chatdata fail %v", err)
		return err
	}

	tScript, err := towerscript.LoadJSON(
		tw.config.MakeTowerFileFullpath(),
	)
	if err != nil {
		return err
	}

	tw.ao2Floor = aoid2floor.New(tw)
	tw.biasFactor = tw.NewRandFactor()

	tw.floorMan = floormanager.New(tScript, tw)
	if err := tw.floorMan.Init(tw.rnd); err != nil {
		return err
	}
	tw.startTime = time.Now()

	tw.gameInfo = &c2t_obj.GameInfo{
		Version:         version.GetVersion(),
		ProtocolVersion: c2t_version.ProtocolVersion,
		DataVersion:     dataversion.DataVersion,

		StartTime:     tw.startTime,
		TowerSeed:     tw.seed,
		TowerName:     tw.config.ScriptFilename,
		Factor:        tw.biasFactor,
		TotalFloorNum: tw.floorMan.GetFloorCount(),

		NickName: tw.Config().NickName,
	}

	g2log.TraceService("%v", tw.gameInfo.StringForm())
	fmt.Printf("%v\n", tw.gameInfo.StringForm())
	fmt.Printf("WebAdmin  : %v:%v id:%v pass:%v\n",
		"http://localhost", tw.config.AdminPort, tw.config.WebAdminID, tw.config.WebAdminPass)

	return nil
}

func (tw *Tower) ServiceCleanup() {
	g2log.TraceService("Start ServiceCleanup %v", tw)
	defer func() { g2log.TraceService("End ServiceCleanup %v", tw) }()

	tw.id2ao.Cleanup()
	tw.ao2Floor.Cleanup()
	for _, f := range tw.floorMan.GetFloorList() {
		f.Cleanup()
	}
	tw.floorMan.Cleanup()
}

func (tw *Tower) ServiceMain(mainctx context.Context) {
	g2log.TraceService("Start ServiceMain %v", tw)
	defer func() { g2log.TraceService("End ServiceMain %v", tw) }()
	ctx, closeCtx := context.WithCancel(mainctx)
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

	tw.c2tCh = make(chan *c2t_packet.Packet, gameconst.SendBufferSize)
	tw.t2cCh = make(chan *c2t_packet.Packet, gameconst.SendBufferSize)

	// start tower
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case data := <-tw.cmdCh:
				tw.processCmd(data)
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

	go tw.handle_c2tch()

	//run client
	go func() {
		time.Sleep(time.Second)
		cl := glclient.New(tw.config, tw.gameInfo, tw.c2tCh, tw.t2cCh)
		if err := cl.Run(); err != nil {
			g2log.Error("%v", err)
		}
		tw.doClose()
	}()

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
			if len(tw.cmdCh) > cap(tw.cmdCh)/2 {
				g2log.Fatal("Tower cmdch overloaded %v/%v", len(tw.cmdCh), cap(tw.cmdCh))
			}
			if len(tw.cmdCh) >= cap(tw.cmdCh) {
				g2log.Fatal("Tower cmdch overloaded %v/%v", len(tw.cmdCh), cap(tw.cmdCh))
				break loop
			}
		case <-rankMakeTk.C:
			go tw.makeActiveObjExpRank()
		}
	}
	tw.doClose()
}

func (tw *Tower) Turn(now time.Time) {
	tw.turnStat.Inc()

	var ws sync.WaitGroup
	for _, f := range tw.floorMan.GetFloorList() {
		ws.Add(1)
		go func(f gamei.FloorI) {
			f.Turn(now)
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
