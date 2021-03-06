// Code generated by "genenum.exe -typename=TurnResultType -packagename=turnresulttype -basedir=enum"

package turnresulttype

import "fmt"

type TurnResultType uint8

const (
	None                       TurnResultType = iota // empty make error
	AttackTo                                         //
	AttackedFrom                                     //
	Kill                                             //
	KilledBy                                         //
	DamagedByTile                                    //
	DeadByTile                                       //
	HPDamageFromTrap                                 //
	SPDamageFromTrap                                 //
	DropCarryObj                                     //
	DropMoney                                        //
	DropMoneyInsteadOfCarryObj                       //
	ContagionTo                                      // success
	ContagionFrom                                    // success
	ContagionToFail                                  // fail
	ContagionFromFail                                // fail

	TurnResultType_Count int = iota
)

var _TurnResultType2string = [TurnResultType_Count][2]string{
	None:                       {"None", "empty make error"},
	AttackTo:                   {"AttackTo", ""},
	AttackedFrom:               {"AttackedFrom", ""},
	Kill:                       {"Kill", ""},
	KilledBy:                   {"KilledBy", ""},
	DamagedByTile:              {"DamagedByTile", ""},
	DeadByTile:                 {"DeadByTile", ""},
	HPDamageFromTrap:           {"HPDamageFromTrap", ""},
	SPDamageFromTrap:           {"SPDamageFromTrap", ""},
	DropCarryObj:               {"DropCarryObj", ""},
	DropMoney:                  {"DropMoney", ""},
	DropMoneyInsteadOfCarryObj: {"DropMoneyInsteadOfCarryObj", ""},
	ContagionTo:                {"ContagionTo", "success"},
	ContagionFrom:              {"ContagionFrom", "success"},
	ContagionToFail:            {"ContagionToFail", "fail"},
	ContagionFromFail:          {"ContagionFromFail", "fail"},
}

func (e TurnResultType) String() string {
	if e >= 0 && e < TurnResultType(TurnResultType_Count) {
		return _TurnResultType2string[e][0]
	}
	return fmt.Sprintf("TurnResultType%d", uint8(e))
}

func (e TurnResultType) CommentString() string {
	if e >= 0 && e < TurnResultType(TurnResultType_Count) {
		return _TurnResultType2string[e][1]
	}
	return ""
}

var _string2TurnResultType = map[string]TurnResultType{
	"None":                       None,
	"AttackTo":                   AttackTo,
	"AttackedFrom":               AttackedFrom,
	"Kill":                       Kill,
	"KilledBy":                   KilledBy,
	"DamagedByTile":              DamagedByTile,
	"DeadByTile":                 DeadByTile,
	"HPDamageFromTrap":           HPDamageFromTrap,
	"SPDamageFromTrap":           SPDamageFromTrap,
	"DropCarryObj":               DropCarryObj,
	"DropMoney":                  DropMoney,
	"DropMoneyInsteadOfCarryObj": DropMoneyInsteadOfCarryObj,
	"ContagionTo":                ContagionTo,
	"ContagionFrom":              ContagionFrom,
	"ContagionToFail":            ContagionToFail,
	"ContagionFromFail":          ContagionFromFail,
}

func String2TurnResultType(s string) (TurnResultType, bool) {
	v, b := _string2TurnResultType[s]
	return v, b
}
