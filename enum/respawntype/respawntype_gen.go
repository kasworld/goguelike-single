// Code generated by "genenum.exe -typename=RespawnType -packagename=respawntype -basedir=enum"

package respawntype

import "fmt"

type RespawnType uint8

const (
	ToHomeFloor    RespawnType = iota // respawn to home floor
	ToCurrentFloor                    // respawn to current floor
	ToRandomFloor                     // respawn to random floor in tower
	//

	RespawnType_Count int = iota
)

var _RespawnType2string = [RespawnType_Count][2]string{
	ToHomeFloor:    {"ToHomeFloor", "respawn to home floor"},
	ToCurrentFloor: {"ToCurrentFloor", "respawn to current floor"},
	ToRandomFloor:  {"ToRandomFloor", "respawn to random floor in tower"},
}

func (e RespawnType) String() string {
	if e >= 0 && e < RespawnType(RespawnType_Count) {
		return _RespawnType2string[e][0]
	}
	return fmt.Sprintf("RespawnType%d", uint8(e))
}

func (e RespawnType) CommentString() string {
	if e >= 0 && e < RespawnType(RespawnType_Count) {
		return _RespawnType2string[e][1]
	}
	return ""
}

var _string2RespawnType = map[string]RespawnType{
	"ToHomeFloor":    ToHomeFloor,
	"ToCurrentFloor": ToCurrentFloor,
	"ToRandomFloor":  ToRandomFloor,
}

func String2RespawnType(s string) (RespawnType, bool) {
	v, b := _string2RespawnType[s]
	return v, b
}
