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

package glclient

import (
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/game/bias"
)

func (app *GLClient) GetArg() interface{} {
	return app.config
}

func (app *GLClient) TowerBias() bias.Bias {
	if app.OLNotiData == nil {
		return bias.Bias{}
	}
	ft := app.TowerInfo.Factor
	dur := app.OLNotiData.Time.Sub(app.TowerInfo.StartTime)
	return bias.MakeBiasByProgress(ft, dur.Seconds(), gameconst.TowerBaseBiasLen)
}

func (app *GLClient) GetPlayerXY() (int, int) {
	ao := app.playerActiveObjClient
	if ao != nil {
		return ao.X, ao.Y
	}
	return 0, 0
}
