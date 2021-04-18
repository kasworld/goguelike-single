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

// Package aoscore score data for ground server
package aoscore

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/kasworld/goguelike-single/config/leveldata"
	"github.com/kasworld/goguelike-single/enum/achievetype_vector"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/turnaction_vector"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/lib/idu64str"
)

func init() {
	gob.Register(&ActiveObjScore{})
}

var ScoreIDMaker = idu64str.New("ScoreID")

type ActiveObjScore struct {
	TowerUUID  string
	TowerName  string
	RecordTime time.Time `prettystring:"simple"`

	UUID        string
	NickName    string
	AchieveStat achievetype_vector.AchieveTypeVector `prettystring:"simple"`
	ActionStat  turnaction_vector.TurnActionVector   `prettystring:"simple"`
	Exp         float64
	Wealth      float64
	BornFaction factiontype.FactionType
	CurrentBias bias.Bias `prettystring:"simple"`
}

// NewActiveObjScoreByLevel make dummy data for ground server
func NewActiveObjScoreByLevel(lv int, bornFaction factiontype.FactionType, CurrentBias bias.Bias) *ActiveObjScore {
	exp := leveldata.BaseExp(lv)
	aos := &ActiveObjScore{
		TowerUUID:  "NoUUID",
		TowerName:  "Noname",
		RecordTime: time.Now(),

		UUID:        ScoreIDMaker.New(),
		NickName:    fmt.Sprintf("Unnamed%v", lv),
		Exp:         exp,
		BornFaction: bornFaction,
		CurrentBias: CurrentBias,
	}
	return aos
}
