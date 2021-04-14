// Code generated by "genenum.exe -typename=StatusOpType -packagename=statusoptype -basedir=enum"

package statusoptype

import "fmt"

type StatusOpType uint8

const (
	None      StatusOpType = iota //
	AddHP                         //
	AddSP                         //
	AddHPRate                     //
	AddSPRate                     //
	//
	ModSight // change sight for 1 turn
	//
	RndFaction   // change to random faction
	IncFaction   // change to next faction
	SetFaction   // set faction to arg
	ResetFaction // reset to born faction
	//
	NegBias         // neg bias value
	RotateBiasRight // rotate bias right
	RotateBiasLeft  // rotate bias left
	//
	SetCondition // set condition for 1 turn
	//
	ForgetFloor    // forget this floor tiles
	ForgetOneFloor // forget visited floor tiles

	StatusOpType_Count int = iota
)

var _StatusOpType2string = [StatusOpType_Count][2]string{
	None:            {"None", ""},
	AddHP:           {"AddHP", ""},
	AddSP:           {"AddSP", ""},
	AddHPRate:       {"AddHPRate", ""},
	AddSPRate:       {"AddSPRate", ""},
	ModSight:        {"ModSight", "change sight for 1 turn"},
	RndFaction:      {"RndFaction", "change to random faction"},
	IncFaction:      {"IncFaction", "change to next faction"},
	SetFaction:      {"SetFaction", "set faction to arg"},
	ResetFaction:    {"ResetFaction", "reset to born faction"},
	NegBias:         {"NegBias", "neg bias value"},
	RotateBiasRight: {"RotateBiasRight", "rotate bias right"},
	RotateBiasLeft:  {"RotateBiasLeft", "rotate bias left"},
	SetCondition:    {"SetCondition", "set condition for 1 turn"},
	ForgetFloor:     {"ForgetFloor", "forget this floor tiles"},
	ForgetOneFloor:  {"ForgetOneFloor", "forget visited floor tiles"},
}

func (e StatusOpType) String() string {
	if e >= 0 && e < StatusOpType(StatusOpType_Count) {
		return _StatusOpType2string[e][0]
	}
	return fmt.Sprintf("StatusOpType%d", uint8(e))
}

func (e StatusOpType) CommentString() string {
	if e >= 0 && e < StatusOpType(StatusOpType_Count) {
		return _StatusOpType2string[e][1]
	}
	return ""
}

var _string2StatusOpType = map[string]StatusOpType{
	"None":            None,
	"AddHP":           AddHP,
	"AddSP":           AddSP,
	"AddHPRate":       AddHPRate,
	"AddSPRate":       AddSPRate,
	"ModSight":        ModSight,
	"RndFaction":      RndFaction,
	"IncFaction":      IncFaction,
	"SetFaction":      SetFaction,
	"ResetFaction":    ResetFaction,
	"NegBias":         NegBias,
	"RotateBiasRight": RotateBiasRight,
	"RotateBiasLeft":  RotateBiasLeft,
	"SetCondition":    SetCondition,
	"ForgetFloor":     ForgetFloor,
	"ForgetOneFloor":  ForgetOneFloor,
}

func String2StatusOpType(s string) (StatusOpType, bool) {
	v, b := _string2StatusOpType[s]
	return v, b
}
