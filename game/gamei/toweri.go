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

package gamei

import (
	"time"

	"github.com/kasworld/goguelike-single/config/goguelikeconfig"
	"github.com/kasworld/goguelike-single/game/bias"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idnoti"
)

type TowerI interface {
	String() string

	GetCmdCh() chan<- interface{}
	GetTurnCh() chan<- time.Time
	SendNoti(cmd c2t_idnoti.NotiID, body interface{})

	GetBias() bias.Bias

	GetFloorManager() FloorManagerI
	GetExpRanking() []ActiveObjectI

	Config() *goguelikeconfig.GoguelikeConfig
}
