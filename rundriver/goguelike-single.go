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

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kasworld/argdefault"
	profile "github.com/kasworld/go-profile"
	"github.com/kasworld/goguelike-single/config/dataversion"
	"github.com/kasworld/goguelike-single/config/towerconfig"
	"github.com/kasworld/goguelike-single/game/tower"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_version"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func printVersion() {
	fmt.Println("Goguelike-single")
	fmt.Println("Build     ", version.GetVersion())
	fmt.Println("Data      ", dataversion.DataVersion)
	fmt.Println("Protocol  ", c2t_version.ProtocolVersion)
	fmt.Println()
}

func main() {
	printVersion()

	profile.AddArgs()

	ads := argdefault.New(&towerconfig.TowerConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &towerconfig.TowerConfig{}
	ads.SetDefaultToNonZeroField(config)
	ads.ApplyFlagTo(config)

	if profile.IsCpu() {
		fn := profile.StartCPUProfile()
		defer fn()
	}

	g2log.GlobalLogger.SetFlags(g2log.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
	g2log.GlobalLogger.SetLevel(config.LogLevel)
	if config.BaseLogDir != "" {
		log, err := g2log.NewWithDstDir(
			config.ScriptFilename,
			config.MakeLogDir(),
			logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
			config.LogLevel,
			config.SplitLogLevel,
		)
		if err == nil {
			g2log.GlobalLogger = log
		} else {
			fmt.Printf("%v\n", err)
		}
	}

	tw := tower.New(config)
	if err := tw.ServiceInit(); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		tw.ServiceMain(context.Background())
		tw.ServiceCleanup()
	}

	if profile.IsMem() {
		profile.WriteHeapProfile()
	}
}
