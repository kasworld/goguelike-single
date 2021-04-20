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

package pid2rspfn

import (
	"fmt"
	"sync"

	"github.com/kasworld/goguelike-single/game/csprotocol"
)

type HandleRspFn func(pk *csprotocol.Packet) error
type PID2RspFn struct {
	mutex     sync.Mutex
	pid2rspfn map[int]HandleRspFn
	pid       int
}

func New() *PID2RspFn {
	rtn := &PID2RspFn{
		pid2rspfn: make(map[int]HandleRspFn),
	}
	return rtn
}
func (p2r *PID2RspFn) NewPID(fn HandleRspFn) int {
	p2r.mutex.Lock()
	defer p2r.mutex.Unlock()
	p2r.pid++
	p2r.pid2rspfn[p2r.pid] = fn
	return p2r.pid
}
func (p2r *PID2RspFn) HandleRsp(pk *csprotocol.Packet) error {
	p2r.mutex.Lock()
	if recvfn, exist := p2r.pid2rspfn[pk.PacketID]; exist {
		delete(p2r.pid2rspfn, pk.PacketID)
		p2r.mutex.Unlock()
		return recvfn(pk)
	}
	p2r.mutex.Unlock()
	return fmt.Errorf("pid not found")
}
