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

package glclientconfig

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/prettystring"
)

type GLClientConfig struct {
	BaseLogDir        string        `default:"" argname:""`
	LogLevel          g2log.LL_Type `default:"7" argname:""`
	SplitLogLevel     g2log.LL_Type `default:"0" argname:""`
	ConnectToTower    string        `default:"localhost:14101" argname:""`
	DisconnectOnDeath bool          `default:"false" argname:""`
}

func (config *GLClientConfig) StringForm() string {
	return prettystring.PrettyString(config, 4)
}

func (config *GLClientConfig) MakeLogDir() string {
	rstr := filepath.Join(config.BaseLogDir,
		"goguelike_glclient.logfiles",
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}
