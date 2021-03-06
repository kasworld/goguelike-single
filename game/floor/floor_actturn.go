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

package floor

import (
	"github.com/kasworld/goguelike-single/config/contagionarea"
	"github.com/kasworld/goguelike-single/config/gameconst"
	"github.com/kasworld/goguelike-single/config/viewportdata"
	"github.com/kasworld/goguelike-single/enum/achievetype"
	"github.com/kasworld/goguelike-single/enum/aotype"
	"github.com/kasworld/goguelike-single/enum/condition"
	"github.com/kasworld/goguelike-single/enum/dangertype"
	"github.com/kasworld/goguelike-single/enum/equipslottype"
	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
	"github.com/kasworld/goguelike-single/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike-single/enum/returncode"
	"github.com/kasworld/goguelike-single/enum/scrolltype"
	"github.com/kasworld/goguelike-single/enum/turnaction"
	"github.com/kasworld/goguelike-single/enum/turnresulttype"
	"github.com/kasworld/goguelike-single/enum/way9type"
	"github.com/kasworld/goguelike-single/game/activeobject/turnresult"
	"github.com/kasworld/goguelike-single/game/aoactreqrsp"
	"github.com/kasworld/goguelike-single/game/cmd2tower"
	"github.com/kasworld/goguelike-single/game/csprotocol"
	"github.com/kasworld/goguelike-single/game/dangerobject"
	"github.com/kasworld/goguelike-single/game/fieldobject"
	"github.com/kasworld/goguelike-single/game/gamei"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/lib/uuidposmani"
)

func (f *Floor) processTurn(TurnCount int) error {

	// wait ai run last turn
	f.aiWG.Wait()

	// prepare to process ao
	ao2ActReqRsp := make(map[gamei.ActiveObjectI]*aoactreqrsp.ActReqRsp, f.aoPosMan.Count())
	aoMapSkipActThisTurn := make(map[string]bool)  // skip this turn
	aoMapLeaveFloorInTurn := make(map[string]bool) // ao to leave floor
	aoListToProcessInTurn := make([]gamei.ActiveObjectI, 0, f.aoPosMan.Count())
	aoAliveInFloorAtStart := make([]gamei.ActiveObjectI, 0, f.aoPosMan.Count())
	for _, v := range f.aoPosMan.GetAllList() {
		ao := v.(gamei.ActiveObjectI)
		aoListToProcessInTurn = append(aoListToProcessInTurn, ao)
		if ao.IsAlive() {
			aoAliveInFloorAtStart = append(aoAliveInFloorAtStart, ao)
			if ao.GetAP() > 0 {
				req := ao.GetClearReq2Handle()
				if req != nil {
					ao2ActReqRsp[ao] = &aoactreqrsp.ActReqRsp{
						Req: *req,
					}
				}
			}
		}
	}
	f.cmdActStat.Add(len(ao2ActReqRsp))

	g2log.Monitor("%v ActiveObj:%v Alive:%v Acted:%v",
		f,
		len(aoListToProcessInTurn),
		len(aoAliveInFloorAtStart),
		len(ao2ActReqRsp),
	)

	for _, ao := range aoListToProcessInTurn {
		ao.PrepareNewTurn(TurnCount)
	}

	// process auto portal and trapteleport steped ao
	for _, ao := range aoListToProcessInTurn {
		if !ao.IsAlive() {
			continue
		}
		if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			g2log.Error("ao not in currentfloor %v %v", f, ao)
			continue
		}

		p, ok := f.foPosMan.Get1stObjAt(aox, aoy).(*fieldobject.FieldObject)
		if !ok {
			continue
		}
		triggered := false
		if p.ActType.AutoTrigger() && p.ActType.TriggerRate() > f.rnd.Float64() {
			triggered = true
		}

		// add found hidden fo to client foman
		f4c := ao.GetFloor4Client(f.GetName())
		f4c.FOPosMan.AddOrUpdateToXY(p.ToPacket_FieldObjClient(aox, aoy), aox, aoy)

		if ao.GetActiveObjType() == aotype.User {
			if p.DisplayType == fieldobjdisplaytype.None {
				f.tower.SendNoti(
					&csprotocol.NotiFoundFieldObj{
						FloorName: f.GetName(),
						FieldObj:  p.ToPacket_FieldObjClient(aox, aoy),
					},
				)
			}
			if p.ActType.TrapNoti() {
				f.tower.SendNoti(
					&csprotocol.NotiActivateTrap{
						FieldObjAct: p.ActType,
						Triggered:   triggered,
					},
				)
			}
		}
		if !triggered {
			continue
		}

		switch p.ActType {
		default:
			fob := fieldobjacttype.GetBuffByFieldObjActType(p.ActType)
			if fob != nil {
				replaced := ao.GetBuffManager().Add(p.ActType.String(), true, true, fob)
				if replaced {
					// need noti?
				}
			}
		case fieldobjacttype.PortalAutoIn:
			p1, p2, err := f.FindUsablePortalPairAt(aox, aoy)
			if err != nil {
				g2log.Error("fail to use portal %v %v %v %v", f, p, ao, err)
				continue
			}
			ao.GetAchieveStat().Inc(achievetype.EnterPortal)
			ao.GetFieldObjActStat().Inc(p2.ActType)
			f.tower.GetCmdCh() <- &cmd2tower.ActiveObjUsePortal{
				SrcFloor:  f,
				ActiveObj: ao,
				P1:        p1,
				P2:        p2,
			}
			aoMapLeaveFloorInTurn[ao.GetUUID()] = true

		case fieldobjacttype.Teleport:
			f.tower.GetCmdCh() <- &cmd2tower.ActiveObjTrapTeleport{
				SrcFloor:     f,
				ActiveObj:    ao,
				DstFloorName: p.DstFloorName,
			}
			aoMapLeaveFloorInTurn[ao.GetUUID()] = true

		case fieldobjacttype.Mine:
			// start explode
			p.CurrentRadius = 0
		}
		if p.ActType.SkipThisTurnAct() {
			aoMapSkipActThisTurn[ao.GetUUID()] = true
			if arr, exist := ao2ActReqRsp[ao]; exist {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Meditate},
					returncode.ActionCanceled)
			}
		}
		if p.ActType.NeedTANoti() {
			ao.SetNeedTANoti()
		}
		ao.GetFieldObjActStat().Inc(p.ActType)
		ao.GetAchieveStat().Inc(achievetype.UseFieldObj)
	}

	// handle sleep condition
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted {
			continue
		}
		if arr.Req.Act.SleepBlockAct() &&
			ao.GetTurnData().Condition.TestByCondition(condition.Sleep) {
			aoMapSkipActThisTurn[ao.GetUUID()] = true
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Meditate},
				returncode.ActionCanceled)
		}
	}
	// handle remain turn2act ao
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted {
			continue
		}
		if ao.GetAP() < 0 {
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Meditate},
				returncode.ActionCanceled)
		}
	}

	// clear dangerobj no remainturn
	if err := f.doPosMan.DelByFilter(func(o uuidposmani.UUIDPosI, x, y int) bool {
		do := o.(*dangerobject.DangerObject)
		return !do.Live1Turn() // del if no remainturn
	}); err != nil {
		g2log.Fatal("fail to delete dangerobject %v", err)
	}

	// add areaattack fieldobj dangerobj
	f.foPosMan.IterAll(func(o uuidposmani.UUIDPosI, foX, foY int) bool {
		fo := o.(*fieldobject.FieldObject)
		switch fo.ActType {
		case fieldobjacttype.RotateLineAttack:
			for wing := 0; wing < fo.WingCount; wing++ {
				for _, v := range fo.GetWingByNum(wing) {
					v.DO.RemainTurn = dangertype.RotateLineAttack.Turn2Live()
					f.doPosMan.AddToXY(
						v.DO,
						foX+v.X, foY+v.Y,
					)
				}
			}
			fo.Degree += fo.DegreePerTurn
		case fieldobjacttype.Mine:
			if fo.CurrentRadius >= gameconst.ViewPortW { // end explode
				fo.CurrentRadius = -1
			}
			if fo.CurrentRadius >= 0 { //  active
				// add do
				for _, v := range fo.GetMineDO() {
					v.DO.RemainTurn = dangertype.MineExplode.Turn2Live()
					f.doPosMan.AddToXY(
						v.DO,
						foX+v.X, foY+v.Y,
					)
				}
				// inc next
				fo.CurrentRadius++
			}
		}
		return false
	})

	// handle attack
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted || !ao.IsAlive() {
			continue
		}
		switch arr.Req.Act {
		case turnaction.Attack:
			f.addBasicAttack(ao, arr)
		case turnaction.AttackWide:
			f.addAttackWide(ao, arr)
		case turnaction.AttackLong:
			f.addAttackLong(ao, arr)
		}
	}

	// handle battle on danger obj
	for _, dstAO := range aoListToProcessInTurn {
		if !dstAO.IsAlive() {
			continue
		}
		dstX, dstY, exist := f.aoPosMan.GetXYByUUID(dstAO.GetUUID())
		if !exist {
			continue
		}
		for _, o := range f.doPosMan.GetObjListAt(dstX, dstY) {
			do := o.(*dangerobject.DangerObject)
			switch do.DangerType {
			default:
				g2log.Fatal("not supported type %v", do.DangerType)
			case dangertype.BasicAttack, dangertype.LongAttack, dangertype.WideAttack:
				owner := do.Owner.(gamei.ActiveObjectI)
				srcTile := f.terrain.GetTiles()[do.OwnerX][do.OwnerY]
				dstTile := f.terrain.GetTiles()[dstX][dstY]
				f.aoAttackActiveObj(owner, dstAO, srcTile, dstTile)
			case dangertype.RotateLineAttack:
				f.foRotateLineAttack(do, dstAO, dstX, dstY)
			case dangertype.MineExplode:
				f.foMineExplodeAttack(do, dstAO, dstX, dstY)
			}
		}
	}

	for _, ao := range aoListToProcessInTurn {
		if ao.ApplyDamageFromDangerObj() { // just killed
			// do nothing here
		}
	}

	// handle ao action except attack
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted || !ao.IsAlive() {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			g2log.Error("ao not in currentfloor %v %v", f, ao)
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Meditate},
				returncode.ActionProhibited)
			continue
		}

		switch arr.Req.Act {
		default:
			g2log.Fatal("unknown aoact %v %v", f, arr)

		case turnaction.Attack, turnaction.AttackWide, turnaction.AttackLong:
			// must be acted
			g2log.Fatal("already acted %v %v", f, arr)

		case turnaction.Meditate:
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Meditate},
				returncode.Success)

		case turnaction.Move:
			mvdir, ec := f.aoAct_Move(ao, arr.Req.Dir, aox, aoy)
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Move, Dir: mvdir},
				ec)

		case turnaction.Pickup:
			if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Pickup, UUID: arr.Req.UUID},
					returncode.ActionProhibited)
				continue
			}
			obj, err := f.poPosMan.GetByXYAndUUID(arr.Req.UUID, aox, aoy)
			if err != nil {
				g2log.Debug("Pickup obj not found %v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Pickup, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			po, ok := obj.(gamei.CarryingObjectI)
			if !ok {
				g2log.Fatal("obj not carryingobject %v", po)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Pickup, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			if err := f.poPosMan.Del(po); err != nil {
				g2log.Fatal("remove po fail %v %v %v", f, po, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Pickup, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			if err := ao.DoPickup(po); err != nil {
				g2log.Error("%v %v %v", f, po, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Pickup, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Pickup, UUID: arr.Req.UUID},
				returncode.Success)

		case turnaction.Drop:
			po := ao.GetInven().GetByUUID(arr.Req.UUID)
			if err := f.aoDropCarryObj(ao, aox, aoy, po); err != nil {
				g2log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Drop, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Drop, UUID: arr.Req.UUID},
				returncode.Success)

		case turnaction.Equip:
			if err := ao.DoEquip(arr.Req.UUID); err != nil {
				g2log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Equip, UUID: arr.Req.UUID},
					returncode.ActionProhibited)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Equip, UUID: arr.Req.UUID},
				returncode.Success)

		case turnaction.UnEquip:
			if err := ao.DoUnEquip(arr.Req.UUID); err != nil {
				g2log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.UnEquip, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.UnEquip, UUID: arr.Req.UUID},
				returncode.Success)

		case turnaction.DrinkPotion:
			po := ao.GetInven().GetByUUID(arr.Req.UUID)
			if po == nil {
				g2log.Error("po not in inventory %v %v", ao, arr.Req.UUID)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.DrinkPotion, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			if err := ao.DoUseCarryObj(arr.Req.UUID); err != nil {
				g2log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.DrinkPotion, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			ao.SetNeedTANoti()
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.DrinkPotion, UUID: arr.Req.UUID},
				returncode.Success)

		case turnaction.ReadScroll:
			po := ao.GetInven().GetByUUID(arr.Req.UUID)
			if po == nil {
				g2log.Error("po not in inventory %v %v", ao, arr.Req.UUID)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.ReadScroll, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			if so, ok := po.(gamei.ScrollI); ok && so.GetScrollType() == scrolltype.Teleport {
				err := f.aoTeleportInFloorRandom(ao)
				if err != nil {
					arr.SetDone(
						aoactreqrsp.Act{Act: turnaction.ReadScroll, UUID: arr.Req.UUID},
						returncode.ActionCanceled)
					g2log.Fatal("fail to teleport %v %v %v", f, ao, err)
					continue
				}
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.ReadScroll, UUID: arr.Req.UUID},
					returncode.Success)
				ao.GetInven().RemoveByUUID(arr.Req.UUID)
				ao.GetAchieveStat().Inc(achievetype.UseCarryObj)
				ao.GetScrollStat().Inc(scrolltype.Teleport)
			} else {
				if err := ao.DoUseCarryObj(arr.Req.UUID); err != nil {
					g2log.Error("%v %v %v", f, ao, err)
					arr.SetDone(
						aoactreqrsp.Act{Act: turnaction.ReadScroll, UUID: arr.Req.UUID},
						returncode.ObjectNotFound)
					continue
				}
				ao.SetNeedTANoti()
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.ReadScroll, UUID: arr.Req.UUID},
					returncode.Success)
			}

		case turnaction.Recycle:
			if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Recycle, UUID: arr.Req.UUID},
					returncode.ActionProhibited)
				continue
			}
			_, ok := f.foPosMan.Get1stObjAt(aox, aoy).(*fieldobject.FieldObject)
			if !ok {
				g2log.Error("not at Recycler FieldObj %v %v", f, ao)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Recycle, UUID: arr.Req.UUID},
					returncode.ActionProhibited)
				continue
			}
			if err := ao.DoRecycleCarryObj(arr.Req.UUID); err != nil {
				g2log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.Recycle, UUID: arr.Req.UUID},
					returncode.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Recycle, UUID: arr.Req.UUID},
				returncode.Success)

		case turnaction.EnterPortal:
			if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.EnterPortal},
					returncode.ActionProhibited)
				continue
			}
			p1, p2, err := f.FindUsablePortalPairAt(aox, aoy)
			if err != nil {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.EnterPortal},
					returncode.ActionProhibited)
				continue
			}
			ao.GetAchieveStat().Inc(achievetype.EnterPortal)
			ao.GetFieldObjActStat().Inc(p1.ActType)
			ao.GetFieldObjActStat().Inc(p2.ActType)
			ao.SetNeedTANoti()
			f.tower.GetCmdCh() <- &cmd2tower.ActiveObjUsePortal{
				SrcFloor:  f,
				ActiveObj: ao,
				P1:        p1,
				P2:        p2,
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.EnterPortal},
				returncode.Success)
			aoMapSkipActThisTurn[ao.GetUUID()] = true
			aoMapLeaveFloorInTurn[ao.GetUUID()] = true
			g2log.Debug("manual in portal %v %v", f, ao)

		case turnaction.ActTeleport:
			if !ao.GetFloor4Client(f.GetName()).Visit.IsComplete() {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.ActTeleport},
					returncode.ActionProhibited)
				continue
			}
			err := f.aoTeleportInFloorRandom(ao)
			if err != nil {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.ActTeleport},
					returncode.ActionCanceled)
				g2log.Fatal("fail to teleport %v %v %v", f, ao, err)
			} else {
				arr.SetDone(
					aoactreqrsp.Act{Act: turnaction.ActTeleport},
					returncode.Success)
			}
		case turnaction.KillSelf:
			ao.ReduceHP(ao.GetHP())
			arr.SetDone(
				aoactreqrsp.Act{Act: turnaction.Meditate},
				returncode.Success)
			aoMapSkipActThisTurn[ao.GetUUID()] = true
		}
	}

	// handle condition greasy Contagion
	for _, ao := range aoListToProcessInTurn {
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			continue
		}

		// contagion can infected from dead ao
		bufname := fieldobjacttype.Contagion.String()
		if ao.GetBuffManager().Exist(bufname) {
			// infact other near
			aoList := f.aoPosMan.GetVPIXYObjByXYLenList(contagionarea.ContagionArea, aox, aoy)
			for _, v := range aoList {
				dstAo := v.O.(gamei.ActiveObjectI)
				if !dstAo.IsAlive() { // skip dead dst
					continue
				}
				if ao.GetUUID() == dstAo.GetUUID() { // skip self
					continue
				}
				if dstAo.GetBuffManager().Exist(bufname) { // skip infected dst
					continue
				}
				if fieldobjacttype.Contagion.TriggerRate() > f.rnd.Float64() {
					fob := fieldobjacttype.GetBuffByFieldObjActType(fieldobjacttype.Contagion)
					dstAo.GetBuffManager().Add(bufname, true, true, fob)

					ao.AppendTurnResult(turnresult.New(turnresulttype.ContagionTo, dstAo, 0))
					dstAo.AppendTurnResult(turnresult.New(turnresulttype.ContagionFrom, ao, 0))

					// fmt.Printf("%v %v to %v\n", bufname, ao, dstAo)
				} else {
					ao.AppendTurnResult(turnresult.New(turnresulttype.ContagionToFail, dstAo, 0))
					dstAo.AppendTurnResult(turnresult.New(turnresulttype.ContagionFromFail, ao, 0))
				}
			}
		}

		if !ao.IsAlive() {
			continue
		}
		if ao.GetTurnData().Condition.TestByCondition(condition.Greasy) && condition.Greasy.Probability() < f.rnd.Float64() {
			eqi := f.rnd.Intn(equipslottype.EquipSlotType_Count)
			co2drop := ao.GetInven().GetEquipSlot()[eqi]
			if co2drop == nil {
				continue
			}
			if err := f.aoDropCarryObj(ao, aox, aoy, co2drop); err != nil {
				g2log.Error("%v %v %v", f, ao, err)
			}
			ao.AppendTurnResult(turnresult.New(turnresulttype.DropCarryObj, co2drop, 0))
		}
	}

	// set ao act result
	for ao, arr := range ao2ActReqRsp {
		ao.SetTurnActReqRsp(arr)
	}

	// apply terrain damage to ao by ao action
	for _, ao := range aoListToProcessInTurn {
		if !ao.IsAlive() {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			g2log.Error("ao not in currentfloor %v %v", f, ao)
			continue
		}
		act := turnaction.Meditate
		dir := way9type.Center
		if arr, exist := ao2ActReqRsp[ao]; exist {
			act = arr.Done.Act
			dir = arr.Done.Dir
		}
		hp, sp := f.terrain.GetTiles()[aox][aoy].ActHPSPCalced(act, dir)
		if hp == 0 && sp == 0 {
			continue
		}
		ao.ApplyHPSPDecByActOnTile(hp, sp)
	}

	f.processCarryObj2floor()

	// apply act result
	for _, ao := range aoListToProcessInTurn {
		ao.ApplyTurnAct() // can die in fn
	}

	// handle ao died in Turn
	for _, ao := range aoAliveInFloorAtStart {
		if ao.IsAlive() {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			g2log.Fatal("ao not in currentfloor %v %v", f, ao)
		}
		if err := f.ActiveObjDropCarryObjByDie(ao, aox, aoy); err != nil {
			g2log.Error("%v %v %v", f, ao, err)
		}
		ao.Death(f) // set rebirth count
		if ao.GetActiveObjType() == aotype.User {
			f.tower.SendNoti(
				&csprotocol.NotiDeath{},
			)
		}
	}

	// for next turn
	// request next turn act for user
	f.sendViewportNoti(TurnCount, aoListToProcessInTurn, aoMapLeaveFloorInTurn)

	// requext next turn act for ai
	for _, ao := range aoListToProcessInTurn {
		if _, exist := aoMapLeaveFloorInTurn[ao.GetUUID()]; exist {
			// skip leaved ao
			continue
		}
		f.aiWG.Add(1)
		go func(ao gamei.ActiveObjectI) {
			ao.RunAI(TurnCount)
			f.aiWG.Done()
		}(ao)
	}

	return nil
}

// ensure po count placed on floor to map script
func (f *Floor) processCarryObj2floor() {
	for _, v := range f.poPosMan.GetAllList() {
		po, ok := v.(gamei.CarryingObjectI)
		if !ok {
			g2log.Fatal("invalid po in poPosMan %v", v)
		}
		if po.DecRemainTurnInFloor() == 0 {
			if err := f.poPosMan.Del(po); err != nil {
				g2log.Warn("remove po fail %v %v %v, maybe already removed",
					f, po, err)
			}
		}
	}

	poNeed := f.terrain.GetCarryObjCount() - f.poPosMan.Count()
	poFailCount := 0
	for i := 0; i < poNeed; i++ {
		if err := f.addNewRandCarryObj2Floor(); err != nil {
			poFailCount++
		}
	}
	if poFailCount > 0 {
		g2log.Monitor("addNewRandCarryObj2Floor fail %v %v/%v", f, poFailCount, poNeed)
	}
}

// send VPTiles, VPObjList noti when need
// request next turn act
func (f *Floor) sendViewportNoti(
	TurnCount int,
	aoListToProcessInTurn []gamei.ActiveObjectI,
	aoMapLeaveFloorInTurn map[string]bool) {

	for _, ao := range aoListToProcessInTurn {
		if _, exist := aoMapLeaveFloorInTurn[ao.GetUUID()]; exist {
			// skip leaved ao
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			g2log.Warn("ao not in currentfloor %v %v, skip tile, obj noti", f, ao)
			continue
		}

		aox, aoy = f.terrain.WrapXY(aox, aoy)
		sightMat := f.terrain.GetViewportCache().GetByCache(aox, aoy)
		sight := float32(ao.GetTurnData().Sight)

		// update ai floor4client info
		vpixyolistsFO := f.foPosMan.GetVPIXYObjByXYLenList(
			viewportdata.ViewportXYLenList, aox, aoy)
		fOs := f.makeViewportFieldObjs2(vpixyolistsFO, sightMat, sight)

		f4c := ao.GetFloor4Client(f.GetName())
		f4c.UpdateObjLists(fOs)

		if ao.GetAndClearNeedTANoti() {
			ao.UpdateVisitAreaBySightMat2(f, aox, aoy, sightMat, sight)
			if ao.GetActiveObjType() == aotype.User {
				f.sendTANoti2Player(ao)
			}
		}
		if ao.GetActiveObjType() == aotype.User {
			f.sendVPObj2Player(ao, TurnCount)
		}
	}
}

// send viewport tiles at ao
// called from processcmd after interfloor move, processturn
func (f *Floor) sendTANoti2Player(ao gamei.ActiveObjectI) {
	aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
	if !exist {
		g2log.Warn("ao not in currentfloor %v %v, skip tile, obj noti", f, ao)
		return
	}
	aox, aoy = f.terrain.WrapXY(aox, aoy)

	sightMat := f.terrain.GetViewportCache().GetByCache(aox, aoy)
	sight := float32(ao.GetTurnData().Sight)
	// make and send NotiTileArea
	f.tower.SendNoti(
		&csprotocol.NotiVPTiles{
			FloorName: f.GetName(),
			VPX:       aox,
			VPY:       aoy,
			VPTiles:   f.makeViewportTiles2(aox, aoy, sightMat, sight),
		},
	)

}

// send viewport object list at ao
// called from processcmd after interfloor move, processturn
func (f *Floor) sendVPObj2Player(ao gamei.ActiveObjectI, TurnCount int) {
	aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
	if !exist {
		g2log.Warn("ao not in currentfloor %v %v, skip tile, obj noti", f, ao)
		return
	}
	aox, aoy = f.terrain.WrapXY(aox, aoy)

	sightMat := f.terrain.GetViewportCache().GetByCache(aox, aoy)
	sight := float32(ao.GetTurnData().Sight)

	vpixyolistsFO := f.foPosMan.GetVPIXYObjByXYLenList(
		viewportdata.ViewportXYLenList, aox, aoy)
	fOs := f.makeViewportFieldObjs2(vpixyolistsFO, sightMat, sight)

	vpixyolistsAO := f.aoPosMan.GetVPIXYObjByXYLenList(
		viewportdata.ViewportXYLenList, aox, aoy)
	aOs := f.makeViewportActiveObjs2(vpixyolistsAO, sightMat, sight)

	vpixyolistsPO := f.poPosMan.GetVPIXYObjByXYLenList(
		viewportdata.ViewportXYLenList, aox, aoy)
	pOs := f.makeViewportCarryObjs2(vpixyolistsPO, sightMat, sight)

	vpixyolistsDO := f.doPosMan.GetVPIXYObjByXYLenList(
		viewportdata.ViewportXYLenList, aox, aoy)
	dOs := f.makeViewportDangerObjs2(vpixyolistsDO, sightMat, sight)

	// make and send NotiVPObjList
	aoContidion := ao.GetTurnData().Condition
	if aoContidion.TestByCondition(condition.Blind) ||
		aoContidion.TestByCondition(condition.Invisible) {
		// if blind, invisible add self
		aOs = append(aOs, ao.ToPacket_ActiveObjClient(aox, aoy))
	}
	f.tower.SendNoti(
		&csprotocol.NotiVPObjList{
			TurnCount:     TurnCount,
			ActiveObj:     ao.ToPacket_PlayerActiveObjInfo(),
			FloorName:     f.GetName(),
			ActiveObjList: aOs,
			CarryObjList:  pOs,
			FieldObjList:  fOs,
			DangerObjList: dOs,
		},
	)
}
