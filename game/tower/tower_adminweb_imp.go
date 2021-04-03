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
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/game/activeobject"
	"github.com/kasworld/goguelike-single/game/cmd2tower"
	"github.com/kasworld/weblib"
)

func (tw *Tower) web_ProtocolStat(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		tw.log.Error("%v", err)
	}
	tw.protocolStat.ToWeb(w, r)
}

func (tw *Tower) web_NotiStat(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		tw.log.Error("%v", err)
	}
	tw.notiStat.ToWeb(w, r)
}

func (tw *Tower) web_ErrorStat(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		tw.log.Error("%v", err)
	}
	tw.errorStat.ToWeb(w, r)
}

func (tw *Tower) web_KickActiveObj(w http.ResponseWriter, r *http.Request) {
	id := weblib.GetStringByName("aoid", "", w, r)
	if id == "" {
		tw.log.Warn("Invalid id")
		http.Error(w, "Invalid id", 404)
		return
	}

	rspCh := make(chan error, 1)
	tw.GetCmdCh() <- &cmd2tower.ActiveObjLeaveTower{
		ActiveObjUUID: id,
		RspCh:         rspCh,
	}
	<-rspCh
}

func (tw *Tower) web_ActiveObjInfo(w http.ResponseWriter, r *http.Request) {
	aoid := weblib.GetStringByName("aoid", "", w, r)
	if aoid == "" {
		tw.log.Warn("Invalid aoid")
		http.Error(w, "Invalid aoid", 404)
		return
	}
	f := tw.ao2Floor.GetFloorByActiveObjID(aoid)
	if f == nil {
		tw.log.Warn("floor not found %v", aoid)
		http.Error(w, "floor not found", 404)
		return
	}

	ao, ok := f.GetActiveObjPosMan().GetByUUID(aoid).(*activeobject.ActiveObject)
	if !ok {
		tw.log.Warn("Invalid aoid %v", aoid)
		http.Error(w, "Invalid aoid", 404)
		return
	}
	if err := weblib.SetFresh(w, r); err != nil {
		tw.log.Error("%v", err)
	}
	ao.Web_ActiveObjInfo(w, r)
}

func (tw *Tower) web_ActiveObjVisitFloorImage(w http.ResponseWriter, r *http.Request) {
	aoid := weblib.GetStringByName("aoid", "", w, r)
	if aoid == "" {
		tw.log.Warn("Invalid aoid")
		http.Error(w, "Invalid aoid", 404)
		return
	}

	visitfloorid := weblib.GetStringByName("floorname", "", w, r)
	if visitfloorid == "" {
		tw.log.Warn("Invalid visitfloorname")
		http.Error(w, "Invalid floor name", 404)
		return
	}

	f := tw.ao2Floor.GetFloorByActiveObjID(aoid)
	if f == nil {
		tw.log.Warn("floor not found %v", aoid)
		http.Error(w, "floor not found", 404)
		return
	}
	ao, ok := f.GetActiveObjPosMan().GetByUUID(aoid).(*activeobject.ActiveObject)
	if !ok {
		tw.log.Warn("Invalid aoid %v", aoid)
		http.Error(w, "Invalid aoid", 404)
		return
	}
	ao.GetFloor4Client(visitfloorid).Visit.Web_Image(w, r)
}

func (tw *Tower) web_ActiveObjRankingList(w http.ResponseWriter, r *http.Request) {
	allActiveObj := tw.aoExpRanking
	page := weblib.GetPage(w, r)
	listActiveObj := allActiveObj.GetPage(page, 40)
	weblib.WebFormBegin("activeobject list", w, r)

	pList := make([]bool, len(allActiveObj)/40+1)

	tplIndex, err := template.New("index").Parse(`
		{{range $i, $v := .}}
		<a href="/ActiveObjRankingList?page={{$i}}">{{$i}}</a>
		{{end}}
	`)
	if err != nil {
		tw.log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, pList); err != nil {
		tw.log.Error("%v", err)
	}

	listActiveObj.ToWebMid(w, r)
	// aolist.ActiveObjList(listActiveObj).ToWebMid(w, r)
	weblib.WebFormEnd(w, r)
}

func (tw *Tower) web_towerStat(w http.ResponseWriter, r *http.Request) {
	err := tw.GetTowerAchieveStat().ToWeb(w, r)
	if err != nil {
		tw.log.Error("%v", err)
	}
}
