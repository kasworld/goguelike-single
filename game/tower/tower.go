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
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike-single/config/authdata"
	"github.com/kasworld/goguelike-single/config/dataversion"
	"github.com/kasworld/goguelike-single/config/gamedata"
	"github.com/kasworld/goguelike-single/config/towerconfig"
	"github.com/kasworld/goguelike-single/enum/towerachieve_vector"
	"github.com/kasworld/goguelike-single/game/activeobject"
	"github.com/kasworld/goguelike-single/game/aoexpsort"
	"github.com/kasworld/goguelike-single/game/aoid2activeobject"
	"github.com/kasworld/goguelike-single/game/aoid2floor"
	"github.com/kasworld/goguelike-single/game/floormanager"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/game/towerscript"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/loadlines"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_serveconnbyte"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_statapierror"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_statnoti"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_statserveapi"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_version"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/recordduration"
	"github.com/kasworld/uuidstr"
	"github.com/kasworld/version"
	"github.com/kasworld/weblib/retrylistenandserve"
)

var _ gamei.TowerI = &Tower{}

func (tw *Tower) String() string {
	return fmt.Sprintf("Tower[%v Seed:%v %v]",
		tw.sconfig.TowerName,
		tw.seed,
		tw.uuid,
	)
}

type Tower struct {
	doClose func()         `prettystring:"hide"`
	rnd     *g2rand.G2Rand `prettystring:"hide"`
	log     *g2log.LogBase `prettystring:"hide"`

	cmdCh chan interface{}

	sconfig    *towerconfig.TowerConfig
	seed       int64
	uuid       string
	biasFactor [3]int64  `prettystring:"simple"`
	startTime  time.Time `prettystring:"simple"`

	floorMan     *floormanager.FloorManager                  `prettystring:"simple"`
	ao2Floor     *aoid2floor.ActiveObjID2Floor               `prettystring:"simple"`
	id2ao        *aoid2activeobject.ActiveObjID2ActiveObject `prettystring:"simple"`
	aoExpRanking aoexpsort.ByExp                             `prettystring:"simple"`

	serviceInfo *c2t_obj.ServiceInfo
	towerInfo   *c2t_obj.TowerInfo

	// single player
	playerConnection *c2t_serveconnbyte.ServeConnByte
	playerAO         *activeobject.ActiveObject

	towerAchieveStat       *towerachieve_vector.TowerAchieveVector `prettystring:"simple"`
	sendStat               *actpersec.ActPerSec                    `prettystring:"simple"`
	recvStat               *actpersec.ActPerSec                    `prettystring:"simple"`
	protocolStat           *c2t_statserveapi.StatServeAPI          `prettystring:"simple"`
	notiStat               *c2t_statnoti.StatNotification          `prettystring:"simple"`
	errorStat              *c2t_statapierror.StatAPIError          `prettystring:"simple"`
	listenClientPaused     bool
	demuxReq2BytesAPIFnMap [c2t_idcmd.CommandID_Count]func(
		me interface{}, hd c2t_packet.Header, rbody []byte) (
		c2t_packet.Header, interface{}, error) `prettystring:"hide"`

	// tower cmd stats
	cmdActStat *actpersec.ActPerSec `prettystring:"simple"`

	adminWeb  *http.Server `prettystring:"simple"`
	clientWeb *http.Server `prettystring:"simple"`
}

func New(config *towerconfig.TowerConfig) *Tower {
	fmt.Printf("%v\n", config.StringForm())

	if config.BaseLogDir != "" {
		log, err := g2log.NewWithDstDir(
			config.TowerName,
			config.MakeLogDir(),
			logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
			config.LogLevel,
			config.SplitLogLevel,
		)
		if err == nil {
			g2log.GlobalLogger = log
		} else {
			fmt.Printf("%v\n", err)
			g2log.GlobalLogger.SetFlags(
				g2log.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
			g2log.GlobalLogger.SetLevel(
				config.LogLevel)
		}
	} else {
		g2log.GlobalLogger.SetFlags(
			g2log.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
		g2log.GlobalLogger.SetLevel(
			config.LogLevel)
	}

	tw := &Tower{
		uuid:  uuidstr.New(),
		id2ao: aoid2activeobject.New("ActiveObject working"),

		sconfig: config,
		log:     g2log.GlobalLogger,

		sendStat:         actpersec.New(),
		recvStat:         actpersec.New(),
		protocolStat:     c2t_statserveapi.New(),
		notiStat:         c2t_statnoti.New(),
		errorStat:        c2t_statapierror.New(),
		cmdActStat:       actpersec.New(),
		towerAchieveStat: new(towerachieve_vector.TowerAchieveVector),
	}

	tw.seed = int64(config.Seed)
	if tw.seed <= 0 {
		tw.seed = time.Now().UnixNano()
	}
	tw.rnd = g2rand.NewWithSeed(int64(tw.seed))

	tw.doClose = func() {
		tw.log.Fatal("Too early doClose call %v", tw)
	}
	tw.serviceInfo = &c2t_obj.ServiceInfo{
		Version:         version.GetVersion(),
		ProtocolVersion: c2t_version.ProtocolVersion,
		DataVersion:     dataversion.DataVersion,
	}
	authdata.AddAdminKey(config.AdminAuthKey)
	return tw
}

// return implement signalhandle.LoggerI
func (tw *Tower) GetLogger() interface{} {
	return tw.log
}

func (tw *Tower) GetServiceLockFilename() string {
	return tw.sconfig.MakePIDFileFullpath()
}

func (tw *Tower) ServiceInit() error {
	rd := recordduration.New(tw.String())

	tw.log.TraceService("Start ServiceInit %v %v", tw, rd)
	defer func() {
		tw.log.TraceService("End ServiceInit %v %v", tw, rd)
		fmt.Println(rd)
	}()

	tw.log.TraceService("%v", tw.serviceInfo.StringForm())
	tw.log.TraceService("%v", tw.sconfig.StringForm())

	var err error

	gamedata.ActiveObjNameList, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().DataFolder, "ainames.txt"),
	)
	if err != nil {
		tw.log.Fatal("load ainame fail %v", err)
		return err
	}

	gamedata.ChatData, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().DataFolder, "chatdata.txt"),
	)
	if err != nil {
		tw.log.Fatal("load chatdata fail %v", err)
		return err
	}

	tScript, err := towerscript.LoadJSON(
		tw.sconfig.MakeTowerFileFullpath(),
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
	tw.towerInfo = &c2t_obj.TowerInfo{
		StartTime:     tw.startTime,
		UUID:          tw.uuid,
		Name:          tw.sconfig.TowerName,
		Factor:        tw.biasFactor,
		TotalFloorNum: tw.floorMan.GetFloorCount(),
	}

	tw.log.TraceService("%v", tw.towerInfo.StringForm())
	fmt.Printf("%v\n", tw.towerInfo.StringForm())
	fmt.Printf("WebAdmin  : %v:%v id:%v pass:%v\n",
		"http://localhost", tw.sconfig.AdminPort, tw.sconfig.WebAdminID, tw.sconfig.WebAdminPass)
	fmt.Printf("WebClient : %v:%v/\n",
		"http://localhost", tw.sconfig.ServicePort)
	fmt.Printf("WebClient with authkey : %v:%v/?authkey=%v\n",
		"http://localhost", tw.sconfig.ServicePort, tw.sconfig.AdminAuthKey)

	return nil
}

func (tw *Tower) ServiceCleanup() {
	tw.log.TraceService("Start ServiceCleanup %v", tw)
	defer func() { tw.log.TraceService("End ServiceCleanup %v", tw) }()

	tw.id2ao.Cleanup()
	tw.ao2Floor.Cleanup()
	for _, f := range tw.floorMan.GetFloorList() {
		f.Cleanup()
	}
	tw.floorMan.Cleanup()
}

func (tw *Tower) ServiceMain(mainctx context.Context) {
	tw.log.TraceService("Start ServiceMain %v", tw)
	defer func() { tw.log.TraceService("End ServiceMain %v", tw) }()
	ctx, closeCtx := context.WithCancel(mainctx)
	tw.doClose = closeCtx

	defer closeCtx()

	totalaocount := 0
	for _, f := range tw.floorMan.GetFloorList() {
		totalaocount += f.GetTerrain().GetActiveObjCount()
	}
	tw.log.Debug("Total system ActiveObj in tower %v", totalaocount)

	queuesize := totalaocount * 100
	tw.cmdCh = make(chan interface{}, queuesize)
	if tw.cmdCh == nil {
		tw.log.Fatal("fail to make cmdCh %v", queuesize)
		return
	}

	go tw.runTower(ctx)
	for _, f := range tw.floorMan.GetFloorList() {
		go func(f gamei.FloorI) {
			f.Run(ctx, queuesize)
			closeCtx()
		}(f)
	}

	for _, f := range tw.floorMan.GetFloorList() {
		for i := 0; i < f.GetTerrain().GetActiveObjCount(); i++ {
			ao := activeobject.NewSystemActiveObj(tw.rnd.Int63(), f, tw.log, tw.towerAchieveStat)
			if err := tw.ao2Floor.ActiveObjEnterTower(f, ao); err != nil {
				tw.log.Error("%v", err)
				continue
			}
			if err := tw.id2ao.Add(ao); err != nil {
				tw.log.Error("%v", err)
			}
		}
	}

	tw.initAdminWeb()
	tw.initServiceWeb(ctx)

	go retrylistenandserve.RetryListenAndServe(tw.adminWeb, tw.log, "serveAdminWeb")
	go retrylistenandserve.RetryListenAndServe(tw.clientWeb, tw.log, "serveServiceWeb")

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-timerInfoTk.C:
			tw.cmdActStat.UpdateLap()
			tw.sendStat.UpdateLap()
			tw.recvStat.UpdateLap()
			if len(tw.cmdCh) > cap(tw.cmdCh)/2 {
				tw.log.Fatal("Tower cmdch overloaded %v/%v", len(tw.cmdCh), cap(tw.cmdCh))
			}
			if len(tw.cmdCh) >= cap(tw.cmdCh) {
				tw.log.Fatal("Tower cmdch overloaded %v/%v", len(tw.cmdCh), cap(tw.cmdCh))
				break loop
			}
		}
	}
	tw.doClose()
}

func (tw *Tower) runTower(ctx context.Context) {
	tw.log.TraceService("Start Run %v", tw)
	defer func() { tw.log.TraceService("End Run %v", tw) }()

	rankMakeTk := time.NewTicker(1 * time.Second)
	defer rankMakeTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case data := <-tw.cmdCh:
			tw.processCmd(data)

		case <-rankMakeTk.C:
			go tw.makeActiveObjExpRank()

		}
	}
	tw.doClose()
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
