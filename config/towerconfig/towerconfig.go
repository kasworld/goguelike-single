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

package towerconfig

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/prettystring"
)

type TowerConfig struct {
	// player
	NickName string `default:"Player" argname:""`

	// common to all tower
	LogLevel         g2log.LL_Type `default:"7" argname:""`
	SplitLogLevel    g2log.LL_Type `default:"0" argname:""`
	BaseLogDir       string        `default:"" argname:""`
	DataFolder       string        `default:"./serverdata" argname:""`
	ClientDataFolder string        `default:"./clientdata" argname:""`
	WebAdminID       string        `default:"root" argname:""`
	WebAdminPass     string        `default:"password" argname:"" prettystring:"hidevalue"`

	// config for each tower
	Seed           int    `default:"0" argname:""` // <=0 time seed
	ServicePort    int    `default:"14101" argname:""`
	AdminPort      int    `default:"14201" argname:""`
	ScriptFilename string `default:"start" argname:""`
}

func (config *TowerConfig) MakeLogDir() string {
	rstr := filepath.Join(config.BaseLogDir,
		fmt.Sprintf("goguelike_tower_%v.logfiles",
			config.ScriptFilename),
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *TowerConfig) MakeTowerFileFullpath() string {
	rstr := filepath.Join(config.DataFolder,
		fmt.Sprintf("%v.tower", config.ScriptFilename),
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *TowerConfig) StringForm() string {
	return prettystring.PrettyString(config, 4)
}

func (config *TowerConfig) ConnectToTower() string {
	return fmt.Sprintf("localhost:%v", config.ServicePort)
}
