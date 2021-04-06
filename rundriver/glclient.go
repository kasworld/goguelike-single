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
	"github.com/kasworld/goguelike-single/config/glclientconfig"
	"github.com/kasworld/goguelike-single/game/glclient"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func main() {

	ads := argdefault.New(&glclientconfig.GLClientConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &glclientconfig.GLClientConfig{}
	ads.SetDefaultToNonZeroField(config)
	ads.ApplyFlagTo(config)

	g2log.GlobalLogger.SetFlags(g2log.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
	g2log.GlobalLogger.SetLevel(config.LogLevel)
	if config.BaseLogDir != "" {
		log, err := g2log.NewWithDstDir(
			"glclient",
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
	app := glclient.New(config)
	err := app.Run(context.Background())
	app.Cleanup()
	if err != nil {
		g2log.Error("%v", err)
	}
}
