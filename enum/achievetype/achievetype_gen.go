// Code generated by "genenum.exe -typename=AchieveType -packagename=achievetype -basedir=enum -vectortype=float64"

package achievetype

import "fmt"

type AchieveType uint8

const (
	Admin           AchieveType = iota //
	Kill                               //
	Death                              //
	Turn                               //
	Move                               //
	EnterPortal                        //
	PickupCarryObj                     //
	EquipCarryObj                      //
	UnEquipCarryObj                    //
	UseCarryObj                        //
	DropCarryObj                       //
	UseFieldObj                        //
	Attack                             //
	AttackHit                          //
	AttackCritical                     //
	AttackMiss                         //
	DamageTotalGive                    //
	DamageMaxGive                      //
	Attacked                           //
	DamageTotalRecv                    //
	DamageMaxRecv                      //
	MaxExp                             //
	MoneyGet                           //

	AchieveType_Count int = iota
)

var _AchieveType2string = [AchieveType_Count][2]string{
	Admin:           {"Admin", ""},
	Kill:            {"Kill", ""},
	Death:           {"Death", ""},
	Turn:            {"Turn", ""},
	Move:            {"Move", ""},
	EnterPortal:     {"EnterPortal", ""},
	PickupCarryObj:  {"PickupCarryObj", ""},
	EquipCarryObj:   {"EquipCarryObj", ""},
	UnEquipCarryObj: {"UnEquipCarryObj", ""},
	UseCarryObj:     {"UseCarryObj", ""},
	DropCarryObj:    {"DropCarryObj", ""},
	UseFieldObj:     {"UseFieldObj", ""},
	Attack:          {"Attack", ""},
	AttackHit:       {"AttackHit", ""},
	AttackCritical:  {"AttackCritical", ""},
	AttackMiss:      {"AttackMiss", ""},
	DamageTotalGive: {"DamageTotalGive", ""},
	DamageMaxGive:   {"DamageMaxGive", ""},
	Attacked:        {"Attacked", ""},
	DamageTotalRecv: {"DamageTotalRecv", ""},
	DamageMaxRecv:   {"DamageMaxRecv", ""},
	MaxExp:          {"MaxExp", ""},
	MoneyGet:        {"MoneyGet", ""},
}

func (e AchieveType) String() string {
	if e >= 0 && e < AchieveType(AchieveType_Count) {
		return _AchieveType2string[e][0]
	}
	return fmt.Sprintf("AchieveType%d", uint8(e))
}

func (e AchieveType) CommentString() string {
	if e >= 0 && e < AchieveType(AchieveType_Count) {
		return _AchieveType2string[e][1]
	}
	return ""
}

var _string2AchieveType = map[string]AchieveType{
	"Admin":           Admin,
	"Kill":            Kill,
	"Death":           Death,
	"Turn":            Turn,
	"Move":            Move,
	"EnterPortal":     EnterPortal,
	"PickupCarryObj":  PickupCarryObj,
	"EquipCarryObj":   EquipCarryObj,
	"UnEquipCarryObj": UnEquipCarryObj,
	"UseCarryObj":     UseCarryObj,
	"DropCarryObj":    DropCarryObj,
	"UseFieldObj":     UseFieldObj,
	"Attack":          Attack,
	"AttackHit":       AttackHit,
	"AttackCritical":  AttackCritical,
	"AttackMiss":      AttackMiss,
	"DamageTotalGive": DamageTotalGive,
	"DamageMaxGive":   DamageMaxGive,
	"Attacked":        Attacked,
	"DamageTotalRecv": DamageTotalRecv,
	"DamageMaxRecv":   DamageMaxRecv,
	"MaxExp":          MaxExp,
	"MoneyGet":        MoneyGet,
}

func String2AchieveType(s string) (AchieveType, bool) {
	v, b := _string2AchieveType[s]
	return v, b
}
