// Code generated by "genprotocol.exe -ver=fc37e02b6858cffd9591410bf9ff4f28fcf1782014d44a7d0e102918f2b1f57d -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_obj

import (
	"time"

	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/enum/achievetype_vector"
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/condition_vector"
	"github.com/kasworld/goguelike-single/enum/equipslottype"
	"github.com/kasworld/goguelike-single/enum/factiontype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype_vector"
	"github.com/kasworld/goguelike-single/enum/potiontype"
	"github.com/kasworld/goguelike-single/enum/potiontype_vector"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/enum/scrolltype_vector"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/tilearea"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd_stats"
)

////////////////////////////////////////////////////////////////////
// commnad protocol

// Invalid make empty packet error
type ReqInvalid_data struct {
	Dummy uint8
}
type RspInvalid_data struct {
	Dummy uint8
}

// Login
type ReqLogin_data struct {
	SessionUUID string
	AuthKey     string
}
type RspLogin_data struct {
	ServiceInfo *ServiceInfo
	AccountInfo *AccountInfo
}

// Heartbeat
type ReqHeartbeat_data struct {
	Time time.Time `prettystring:"simple"`
}
type RspHeartbeat_data struct {
	Time time.Time `prettystring:"simple"`
}

// Chat
type ReqChat_data struct {
	Chat string
}
type RspChat_data struct {
	Dummy uint8
}

// AchieveInfo
type ReqAchieveInfo_data struct {
	Dummy uint8
}
type RspAchieveInfo_data struct {
	AchieveStat   achievetype_vector.AchieveTypeVector         `prettystring:"simple"`
	PotionStat    potiontype_vector.PotionTypeVector           `prettystring:"simple"`
	ScrollStat    scrolltype_vector.ScrollTypeVector           `prettystring:"simple"`
	FOActStat     fieldobjacttype_vector.FieldObjActTypeVector `prettystring:"simple"`
	AOActionStat  c2t_idcmd_stats.CommandIDStat                `prettystring:"simple"`
	ConditionStat condition_vector.ConditionVector             `prettystring:"simple"`
}

// Rebirth
type ReqRebirth_data struct {
	Dummy uint8
}
type RspRebirth_data struct {
	Dummy uint8
}

// MoveFloor tower cmd
type ReqMoveFloor_data struct {
	UUID string
}
type RspMoveFloor_data struct {
	Dummy uint8
}

// AIPlay
type ReqAIPlay_data struct {
	On bool
}
type RspAIPlay_data struct {
	Dummy uint8
}

// VisitFloorList floor info of visited
type ReqVisitFloorList_data struct {
	Dummy uint8 // change as you need
}

// VisitFloorList floor info of visited
type RspVisitFloorList_data struct {
	FloorList []*FloorInfo
}

////////////////////////////////////////////////////////////////////////////
// ao act

type ReqMeditate_data struct {
	Dummy uint8
}
type RspMeditate_data struct {
	Dummy uint8
}

type ReqKillSelf_data struct {
	Dummy uint8
}
type RspKillSelf_data struct {
	Dummy uint8
}

type ReqMove_data struct {
	Dir way9type.Way9Type
}
type RspMove_data struct {
	Dir way9type.Way9Type
}

type ReqAttack_data struct {
	Dir way9type.Way9Type
}
type RspAttack_data struct {
	Dummy uint8
}

type ReqAttackWide_data struct {
	Dir way9type.Way9Type
}
type RspAttackWide_data struct {
	Dummy uint8
}

type ReqAttackLong_data struct {
	Dir way9type.Way9Type
}
type RspAttackLong_data struct {
	Dummy uint8
}

type ReqPickup_data struct {
	UUID string
}
type RspPickup_data struct {
	Dummy uint8
}

type ReqDrop_data struct {
	UUID string
}
type RspDrop_data struct {
	Dummy uint8
}

type ReqEquip_data struct {
	UUID string
}
type RspEquip_data struct {
	Dummy uint8
}

type ReqUnEquip_data struct {
	UUID string
}
type RspUnEquip_data struct {
	Dummy uint8
}

type ReqDrinkPotion_data struct {
	UUID string
}
type RspDrinkPotion_data struct {
	Dummy uint8
}

type ReqReadScroll_data struct {
	UUID string
}
type RspReadScroll_data struct {
	Dummy uint8
}

type ReqRecycle_data struct {
	UUID string
}
type RspRecycle_data struct {
	Dummy uint8
}

type ReqEnterPortal_data struct {
	Dummy uint8
}
type RspEnterPortal_data struct {
	Dummy uint8
}

type ReqActTeleport_data struct {
	Dummy uint8
}
type RspActTeleport_data struct {
	Dummy uint8
}

////////////////////////////////////////////////////////////////////
// admin

// AdminTowerCmd generic cmd
type ReqAdminTowerCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminTowerCmd_data struct {
	Dummy uint8
}

// AdminFloorCmd generic cmd
type ReqAdminFloorCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminFloorCmd_data struct {
	Dummy uint8
}

// AdminActiveObjCmd generic cmd
type ReqAdminActiveObjCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminActiveObjCmd_data struct {
	Dummy uint8
}

// AdminFloorMove Next Before floorUUID
type ReqAdminFloorMove_data struct {
	Floor string
}
type RspAdminFloorMove_data struct {
	Dummy uint8
}

// AdminTeleport random pos in floor
type ReqAdminTeleport_data struct {
	X int
	Y int
}
type RspAdminTeleport_data struct {
	Dummy uint8
}

// AdminAddExp  add arg to battle exp
type ReqAdminAddExp_data struct {
	Exp int
}

// AdminAddExp  add arg to battle exp
type RspAdminAddExp_data struct {
	Dummy uint8 // change as you need
}

// AdminPotionEffect buff by arg potion type
type ReqAdminPotionEffect_data struct {
	Potion potiontype.PotionType
}

// AdminPotionEffect buff by arg potion type
type RspAdminPotionEffect_data struct {
	Dummy uint8 // change as you need
}

// AdminScrollEffect buff by arg Scroll type
type ReqAdminScrollEffect_data struct {
	Scroll scrolltype.ScrollType
}

// AdminScrollEffect buff by arg Scroll type
type RspAdminScrollEffect_data struct {
	Dummy uint8 // change as you need
}

// AdminCondition add arg condition for 100 turn
type ReqAdminCondition_data struct {
	Condition condition.Condition
}

// AdminCondition add arg condition for 100 turn
type RspAdminCondition_data struct {
	Dummy uint8 // change as you need
}

// AdminAddPotion add arg potion to inven
type ReqAdminAddPotion_data struct {
	Potion potiontype.PotionType
}

// AdminAddPotion add arg potion to inven
type RspAdminAddPotion_data struct {
	Dummy uint8 // change as you need
}

// AdminAddScroll add arg scroll to inven
type ReqAdminAddScroll_data struct {
	Scroll scrolltype.ScrollType
}

// AdminAddScroll add arg scroll to inven
type RspAdminAddScroll_data struct {
	Dummy uint8 // change as you need
}

// AdminAddMoney add arg money to inven
type ReqAdminAddMoney_data struct {
	Money int
}

// AdminAddMoney add arg money to inven
type RspAdminAddMoney_data struct {
	Dummy uint8 // change as you need
}

// AdminAddEquip add random equip to inven
type ReqAdminAddEquip_data struct {
	Faction factiontype.FactionType
	Equip   equipslottype.EquipSlotType
}

// AdminAddEquip add random equip to inven
type RspAdminAddEquip_data struct {
	Dummy uint8 // change as you need
}

// AdminForgetFloor forget current floor map
type ReqAdminForgetFloor_data struct {
	Dummy uint8 // change as you need
}

// AdminForgetFloor forget current floor map
type RspAdminForgetFloor_data struct {
	Dummy uint8 // change as you need
}

// AdminFloorMap complete current floor map
type ReqAdminFloorMap_data struct {
	Dummy uint8 // change as you need
}

// AdminFloorMap complete current floor map
type RspAdminFloorMap_data struct {
	Dummy uint8 // change as you need
}

/////////////////////////////////////////////////////////////////
// noti

type NotiInvalid_data struct {
	Dummy uint8
}

type NotiEnterTower_data struct {
	TowerInfo *TowerInfo
}
type NotiLeaveTower_data struct {
	TowerInfo *TowerInfo
}

type NotiEnterFloor_data struct {
	FI *FloorInfo
}
type NotiLeaveFloor_data struct {
	FI *FloorInfo
}

type NotiAgeing_data struct {
	FloorName string
}

type NotiDeath_data struct {
	Dummy uint8
}

type NotiReadyToRebirth_data struct {
	Dummy uint8
}
type NotiRebirthed_data struct {
	Dummy uint8
}

type NotiBroadcast_data struct {
	Msg string
}

type NotiVPObjList_data struct {
	Time          time.Time `prettystring:"simple"`
	FloorName     string
	ActiveObj     *PlayerActiveObjInfo
	ActiveObjList []*ActiveObjClient
	CarryObjList  []*CarryObjClientOnFloor
	FieldObjList  []*FieldObjClient
	DangerObjList []*DangerObjClient
}

// NotiVPTiles_data contains tile info center from pos
type NotiVPTiles_data struct {
	FloorName string
	VPX       int // viewport center X
	VPY       int // viewport center Y
	VPTiles   *viewportdata.ViewportTileArea2
}

// NotiFloorTiles_data used for floor map, reconnect client
type NotiFloorTiles_data struct {
	FI    *FloorInfo
	X     int // X start position, not center
	Y     int // Y start position, not center
	Tiles tilearea.TileArea
}

// FieldObjList    // for rebuild known floor
type NotiFieldObjList_data struct {
	FI     *FloorInfo
	FOList []*FieldObjClient
}

type NotiFoundFieldObj_data struct {
	FloorName string
	FieldObj  *FieldObjClient
}

type NotiForgetFloor_data struct {
	FloorName string
}

type NotiActivateTrap_data struct {
	FieldObjAct fieldobjacttype.FieldObjActType
	Triggered   bool
}
