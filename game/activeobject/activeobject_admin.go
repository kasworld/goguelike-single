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

package activeobject

import (
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_error"
)

func (ao *ActiveObject) DoAdminCmd(cmd, args string) c2t_error.ErrorCode {
	g2log.AdminAudit("$c %v %v", ao, cmd, args)

	ec := c2t_error.None

	switch cmd {
	default:
		g2log.Error("unknown admin ao cmd %v %v", cmd, args)
		ec = c2t_error.ActionProhibited

	}
	return ec
}
