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

package c2t_idcmd

func (cmd CommandID) SleepCancel() bool {
	return attrib[cmd].sleepCancel
}

func (cmd CommandID) NeedTurn() float64 {
	return attrib[cmd].needTurn
}

func (cmd CommandID) TriggerTurn() bool {
	return attrib[cmd].triggerTurn
}

var attrib = [CommandID_Count]struct {
	sleepCancel bool
	triggerTurn bool
	needTurn    float64
}{
	Invalid:           {false, false, 0},
	Login:             {false, false, 0},
	Heartbeat:         {false, false, 0},
	Chat:              {false, false, 0},
	AchieveInfo:       {false, false, 0},
	Rebirth:           {false, true, 0},
	MoveFloor:         {false, true, 1}, // need check need turn
	AIPlay:            {false, true, 0},
	VisitFloorList:    {false, false, 0},
	Meditate:          {false, true, 1},
	KillSelf:          {false, true, 1},
	Move:              {true, true, 1},
	Attack:            {true, true, 1.5},
	AttackWide:        {true, true, 3},
	AttackLong:        {true, true, 3},
	Pickup:            {true, true, 1},
	Drop:              {true, true, 1},
	Equip:             {true, true, 1},
	UnEquip:           {true, true, 1},
	DrinkPotion:       {true, true, 1},
	ReadScroll:        {true, true, 1},
	Recycle:           {true, true, 1},
	EnterPortal:       {true, true, 1},
	ActTeleport:       {false, true, 1},
	AdminTowerCmd:     {false, true, 0},
	AdminFloorCmd:     {false, true, 0},
	AdminActiveObjCmd: {false, true, 0},
	AdminFloorMove:    {false, true, 0},
	AdminTeleport:     {false, true, 0},
	AdminAddExp:       {false, true, 0},
	AdminPotionEffect: {false, true, 0},
	AdminScrollEffect: {false, true, 0},
	AdminCondition:    {false, true, 0},
	AdminAddPotion:    {false, true, 0},
	AdminAddScroll:    {false, true, 0},
	AdminAddMoney:     {false, true, 0},
	AdminAddEquip:     {false, true, 0},
	AdminForgetFloor:  {false, true, 0},
	AdminFloorMap:     {false, false, 0},
}
