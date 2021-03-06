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
	"html/template"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/goguelike-single/enum/tile_flag"
	"github.com/kasworld/goguelike-single/enum/towerachieve_vector_float64"
	"github.com/kasworld/goguelike-single/game/activeobject"
	"github.com/kasworld/goguelike-single/game/aoid2activeobject"
	"github.com/kasworld/goguelike-single/game/carryingobject"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/dangerobject"
	"github.com/kasworld/goguelike-single/game/fieldobject"
	"github.com/kasworld/goguelike-single/game/terrain/room"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/version"
	"github.com/kasworld/weblib"
	"github.com/kasworld/weblib/webprofile"
	"github.com/kasworld/wrapper"
)

func (tw *Tower) web_FaviconIco(w http.ResponseWriter, r *http.Request) {
}

func (tw *Tower) runAdminWeb() {
	authData := weblib.NewAuthData("tower")
	authData.ReLoadUserData([][2]string{
		{tw.config.WebAdminID, tw.config.WebAdminPass},
	})
	webMux := weblib.NewAuthMux(authData, g2log.GlobalLogger)

	if !version.IsRelease() {
		webprofile.AddWebProfile(webMux)
	}

	webMux.HandleFunc("/favicon.ico", tw.web_FaviconIco)

	webMux.HandleFunc("/addturn", tw.web_AddTurn)

	webMux.HandleFuncAuth("/", tw.web_TowerInfo)
	webMux.HandleFuncAuth("/floor", tw.web_FloorInfo)
	webMux.HandleFuncAuth("/floorimagezoom", tw.web_FloorImageZoom)
	webMux.HandleFuncAuth("/floorimageautozoom", tw.web_FloorImageAutoZoom)
	webMux.HandleFuncAuth("/floortile", tw.web_FloorTile)

	webMux.HandleFuncAuth("/ActiveObj", tw.web_ActiveObjInfo)
	webMux.HandleFuncAuth("/ActiveObjVisitImgae", tw.web_ActiveObjVisitFloorImage)
	webMux.HandleFuncAuth("/ActiveObjRankingList", tw.web_ActiveObjRankingList)
	webMux.HandleFuncAuth("/towerStat", tw.web_towerStat)

	webMux.HandleFuncAuth("/terrain", tw.web_TerrainInfo)
	webMux.HandleFuncAuth("/terrainimagezoom", tw.web_TerrainImageZoom)
	webMux.HandleFuncAuth("/terrainimageautozoom", tw.web_TerrainImageAutoZoom)
	webMux.HandleFuncAuth("/terraintile", tw.web_TerrainTile)

	authData.AddAllActionName(tw.config.WebAdminID)

	tw.adminWeb = &http.Server{
		Handler: webMux,
		Addr:    fmt.Sprintf(":%v", tw.config.AdminPort),
	}
	for {
		err := tw.adminWeb.ListenAndServe()
		if err != nil {
			operr, ok := err.(*net.OpError)
			if ok {
				syscallerr, ok := operr.Err.(*os.SyscallError)
				if ok && syscallerr.Err == syscall.EADDRINUSE {
					g2log.Warn("Retry adminWeb %v", err) // bind: address already in use
					time.Sleep(time.Second)
					continue
				}
			}
			break
		}
	}
}

func (tw *Tower) BuildDate() time.Time {
	return version.GetBuildDate()
}

func (tw *Tower) NumGoroutine() int {
	return runtime.NumGoroutine()
}

func (tw *Tower) TileCacheCount() int {
	return tile_flag.TileCacheCount()
}

func (tw *Tower) WrapInfo() string {
	return wrapper.G_WrapperInfo()
}

func (tw *Tower) GetTowerAchieveStat() *towerachieve_vector_float64.TowerAchieveVector_float64 {
	return tw.towerAchieveStat
}

func (tw *Tower) GetStartTime() time.Time {
	return tw.startTime
}

func (tw *Tower) GetID2ActiveObj() *aoid2activeobject.ActiveObjID2ActiveObject {
	return tw.id2ao
}

func (tw *Tower) GetTurnStat() *actpersec.ActPerSec {
	return tw.turnStat
}
func (tw *Tower) GetGameInfo() *csprotocol.GameInfo {
	return tw.gameInfo
}

func (tw *Tower) GetTowerCmdActStat() *actpersec.ActPerSec {
	return tw.cmdActStat
}

func (tw *Tower) SysAOID() string {
	return activeobject.SysAOIDMaker.String()
}
func (tw *Tower) EquipID() string {
	return carryingobject.EquipIDMaker.String()
}
func (tw *Tower) MoneyID() string {
	return carryingobject.MoneyIDMaker.String()
}
func (tw *Tower) PotionID() string {
	return carryingobject.PotionIDMaker.String()
}
func (tw *Tower) ScrollID() string {
	return carryingobject.ScrollIDMaker.String()
}
func (tw *Tower) DOID() string {
	return dangerobject.DOIDMaker.String()
}
func (tw *Tower) FOID() string {
	return fieldobject.FOIDMaker.String()
}
func (tw *Tower) RoomID() string {
	return room.RoomIDMaker.String()
}

func (tw *Tower) GetPlayerAO() *activeobject.ActiveObject {
	return tw.playerAO
}

func (tw *Tower) CmdChState() string {
	return fmt.Sprintf("cmdCh %v/%v", len(tw.cmdCh), cap(tw.cmdCh))
}
func (tw *Tower) TurnChInfo() string {
	return fmt.Sprintf("turnCh %v/%v", len(tw.turnCh), cap(tw.turnCh))
}

func (tw *Tower) C2tChInfo() string {
	return fmt.Sprintf("c2tCh %v/%v", len(tw.c2tCh), cap(tw.c2tCh))
}
func (tw *Tower) T2cChInfo() string {
	return fmt.Sprintf("t2cCh %v/%v", len(tw.t2cCh), cap(tw.t2cCh))
}

func (tw *Tower) web_TowerInfo(w http.ResponseWriter, r *http.Request) {
	tplIndex, err := template.New("index").Parse(`
	<html> <head>
	<title>Tower {{.GetGameInfo.TowerName}} admin</title>
	</head>
	<body>

	BuildDate : {{.BuildDate.Format "2006-01-02T15:04:05Z07:00"}}
	<br/>
	Version: {{.GetGameInfo.Version}}
	<br/>
	DataVersion : {{.GetGameInfo.DataVersion}}
	<hr/>
	{{.}}
	<br/>
	Factor : {{.GetGameInfo.Factor}}
	<br/>
	Current Bias : {{.GetBias}}
	<br/>
	TotalFloor : {{.GetFloorManager.GetFloorCount}}
	<br/>
	Tile2Discover : {{.GetFloorManager.CalcTiles2Discover}}
	<br/>
	Max Exp From Discover : {{.GetFloorManager.CalcFullDiscoverExp}} 
	<br/>
	Max Level From Discover :  {{.GetFloorManager.CalcFullDiscoverLevel}} 
	<br/>
	SysAOID : {{.SysAOID}}
	<br/>
	EquipID : {{.EquipID}}
	<br/>
	MoneyID : {{.MoneyID}}
	<br/>
	PotionID : {{.PotionID}}
	<br/>
	ScrollID : {{.ScrollID}}
	<br/>
	DOID : {{.DOID}}
	<br/>
	FOID : {{.FOID}}
	<br/>
	RoomID : {{.RoomID}}
	<br/>
	goroutine : {{.NumGoroutine}}	
	<br/>
	global wrapper : {{.WrapInfo}}	
	<br/>
	TileCache : {{.TileCacheCount}}
	<br/>
	TurnStat : {{.GetTurnStat}} {{.TurnChInfo}} <a href= "/addturn" target="_blank">Add</a> 
	<br/>
	{{.GetInterDur.GetCount}} {{.GetInterDur.GetInterval}} {{.GetInterDur.GetDuration}}
	<br/>
	<a href= "/towerStat" target="_blank">Tower Achieve</a>
    <br/>
	TowerCmd act : {{.CmdChState}} {{.GetTowerCmdActStat}}
    <br/>
	{{.C2tChInfo}}
    <br/>
	{{.T2cChInfo}}
    <br/>

	<hr/>
	<a href="/ActiveObjRankingList?page=0" target="_blank">{{.GetID2ActiveObj}}</a>
    <br/>
	<a href="/ActiveObj?aoid={{.GetPlayerAO.GetUUID}}" target="_blank">{{.GetPlayerAO}}</a>
    <br/>
    <br/>

	<table border=1 style="border-collapse:collapse;">
	` + floor_HTML_header + `
	{{range $i, $v := .GetFloorManager.GetFloorList}}
		{{if $v}}
	` + floor_HTML_row + `
		{{end}}
	{{end}}
	` + floor_HTML_header + `
	</table>
	<br/>
	<pre>{{.GetGameInfo.StringForm}}</pre>
	<pre>{{.Config.StringForm}}</pre>
	<br/>
	</body> </html> 
	`)
	if err != nil {
		g2log.Error("%v", err)
	}
	if err := weblib.SetFresh(w, r); err != nil {
		g2log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, tw); err != nil {
		g2log.Error("%v", err)
	}
}

const (
	floor_HTML_header = `<tr>
	<td>Floor</td>
	<td>FloorCmd act</td>
	<td>Terrain Ageing</td>
	<td>Faction</td>
	<td>W/H</td>
	<td>ActiveObj/CarryObj</td>
	<td>ActTurnJitter</td>
	<td>Viewport Cache</td>
	</tr>`
	floor_HTML_row = `<tr>
	<td>
	<a href= "/floor?floorname={{$v.GetName}}" target="_blank">
		{{$v.GetName}}
	</a>
	</td>
	<td>{{$v.CmdChState}} {{$v.GetCmdFloorActStat}}
	</td>
	<td>
	<a href= "/terrain?floorname={{$v.GetName}}" target="_blank">
		{{$v.GetTerrain.AgeingCount}}/{{$v.GetTerrain.GetResetAfterNAgeing}}
	</a>
	</td>
	<td>{{$v.GetBias.NearFaction}}</td>
	<td>{{$v.GetWidth}}/{{$v.GetHeight}}</td>
	<td>{{$v.TotalActiveObjCount}} / {{$v.TotalCarryObjCount}}</td>
	<td>{{$v.GetInterDur.GetCount}} {{$v.GetInterDur.GetInterval}} {{$v.GetInterDur.GetDuration}}</td>
	<td>{{$v.GetTerrain.GetViewportCache.HitRate}}</td>
	</tr>`
)
