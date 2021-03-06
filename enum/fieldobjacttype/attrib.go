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

package fieldobjacttype

import (
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/statusoptype"
	"github.com/kasworld/htmlcolors"
)

func (v FieldObjActType) Color24() htmlcolors.Color24 {
	return attrib[v].color24
}

func (v FieldObjActType) Rune() string {
	return attrib[v].runeStr
}

func (v FieldObjActType) TrapNoti() bool {
	return attrib[v].trapNoti
}

func (v FieldObjActType) AutoTrigger() bool {
	return attrib[v].autoTrigger
}

func (v FieldObjActType) TriggerRate() float64 {
	return attrib[v].triggerRate
}

func (v FieldObjActType) SkipThisTurnAct() bool {
	return attrib[v].skipThisTurnAct
}
func (v FieldObjActType) NeedTANoti() bool {
	return attrib[v].needTANoti
}

func (v FieldObjActType) MustCharPlaceable() bool {
	return attrib[v].mustCharPlaceable
}

var attrib = [FieldObjActType_Count]struct {
	runeStr           string
	trapNoti          bool // send noti on step
	autoTrigger       bool
	triggerRate       float64 // if AutoTrigger true
	skipThisTurnAct   bool
	needTANoti        bool // if pos changed
	mustCharPlaceable bool // fatal if placed at noCharPlaceable
	color24           htmlcolors.Color24
}{
	None: {"?", false, false, 1.0, false, false, false, htmlcolors.Black},

	PortalInOut:     {"?", false, false, 0.0, false, false, true, htmlcolors.MediumVioletRed},
	PortalIn:        {"?", false, false, 0.0, false, false, true, htmlcolors.MediumVioletRed},
	PortalOut:       {"?", false, false, 0.0, false, false, false, htmlcolors.MediumVioletRed},
	PortalAutoIn:    {"?", false, true, 1.0, true, true, true, htmlcolors.MediumVioletRed},
	RecycleCarryObj: {"?", false, false, 0.0, false, false, false, htmlcolors.Green},
	Teleport:        {"?", true, true, 0.1, true, true, false, htmlcolors.Red},

	ForgetFloor:    {"?", true, true, 0.2, false, true, false, htmlcolors.OrangeRed},
	ForgetOneFloor: {"?", true, true, 0.3, false, true, false, htmlcolors.OrangeRed},
	AlterFaction:   {"?", true, true, 0.5, false, false, false, htmlcolors.Red},
	AllFaction:     {"?", true, true, 0.5, false, false, false, htmlcolors.Red},
	Bleeding:       {"?", true, true, 0.2, false, false, false, htmlcolors.Crimson},
	Chilly:         {"?", true, true, 0.2, false, false, false, htmlcolors.DarkTurquoise},

	Blind:     {"?", true, true, 0.2, false, false, false, condition.Blind.Color()},
	Invisible: {"?", true, true, 0.5, false, false, false, condition.Invisible.Color()},
	Burden:    {"?", true, true, 0.2, false, false, false, condition.Burden.Color()},
	Float:     {"?", true, true, 0.3, false, false, false, condition.Float.Color()},
	Greasy:    {"?", true, true, 0.5, false, false, false, condition.Greasy.Color()},
	Drunken:   {"?", true, true, 0.5, false, false, false, condition.Drunken.Color()},
	Sleepy:    {"?", true, true, 0.1, false, false, false, condition.Sleep.Color()},
	Contagion: {"?", true, true, 0.1, false, false, false, condition.Contagion.Color()},
	Slow:      {"?", true, true, 0.1, false, false, false, condition.Slow.Color()},
	Haste:     {"?", true, true, 0.1, false, false, false, condition.Haste.Color()},

	RotateLineAttack: {"?", false, false, 0.0, false, false, false, htmlcolors.Lavender},
	Mine:             {"?", true, true, 1.0, false, false, false, htmlcolors.Orange},
}

// try act on fieldobj
var ClientData = [FieldObjActType_Count]struct {
	ActOn bool
	Text  string
}{
	None:             {true, ""},
	PortalInOut:      {true, "portal in/out"},
	PortalIn:         {true, "portal oneway"},
	PortalOut:        {true, "portal out only"},
	PortalAutoIn:     {false, "portal auto in oneway"},
	RecycleCarryObj:  {true, "recycle carryobj to money"},
	Teleport:         {false, "teleport somewhere"},
	ForgetFloor:      {false, "forget current floor"},
	ForgetOneFloor:   {false, "forget some floor you visited"},
	AlterFaction:     {false, "change faction randomly"},
	AllFaction:       {false, "rotate all faction"},
	Bleeding:         {false, "hp damage"},
	Chilly:           {false, "sp damage"},
	Blind:            {false, "sight 0"},
	Invisible:        {false, "other cannot see you"},
	Burden:           {false, "overload limit reduced"},
	Float:            {false, "float in air"},
	Greasy:           {false, "greasy body"},
	Drunken:          {false, "random direction"},
	Sleepy:           {false, "cannot act"},
	Contagion:        {false, "make contagion other, die or heal randomly"},
	RotateLineAttack: {false, "rotate line of dangerobj"},
	Mine:             {false, "explode on step"},
}

func GetBuffByFieldObjActType(at FieldObjActType) []statusoptype.OpArg {
	return foAct2BuffList[at]
}

var foAct2BuffList = [FieldObjActType_Count][]statusoptype.OpArg{
	// immediate effect
	AlterFaction: {
		statusoptype.OpArg{statusoptype.RndFaction, nil},
	},
	AllFaction: statusoptype.RepeatShift(260, 10,
		statusoptype.OpArg{statusoptype.IncFaction, 1},
	),

	ForgetFloor: {
		statusoptype.OpArg{statusoptype.ForgetFloor, nil},
	},
	ForgetOneFloor: {
		statusoptype.OpArg{statusoptype.ForgetOneFloor, nil},
	},

	// statusop debuff
	Bleeding: statusoptype.RepeatShift(200, 10,
		statusoptype.OpArg{statusoptype.AddHPRate, -0.05},
	),
	Chilly: statusoptype.RepeatShift(200, 10,
		statusoptype.OpArg{statusoptype.AddSPRate, -0.05},
	),

	// condition debuff
	Blind: statusoptype.RepeatShift(200, 2,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Blind},
	),
	Invisible: statusoptype.RepeatShift(200, 2,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Invisible},
	),
	Burden: statusoptype.RepeatShift(100, 1,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Burden},
	),
	Float: statusoptype.RepeatShift(200, 1,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Float},
	),
	Greasy: statusoptype.RepeatShift(400, 1,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Greasy},
	),
	Drunken: statusoptype.RepeatShift(200, 2,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Drunken},
	),
	Sleepy: statusoptype.RepeatShift(200, 4,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Sleep},
	),
	Contagion: statusoptype.RepeatShift(400, 4,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Contagion},
	),
	Slow: statusoptype.RepeatShift(200, 1,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Slow},
	),
	Haste: statusoptype.RepeatShift(200, 1,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Haste},
	),
}
