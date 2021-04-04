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

	"github.com/kasworld/argdefault"
	"github.com/kasworld/configutil"
	"github.com/kasworld/goguelike-single/config/textclientconfig"
	"github.com/kasworld/goguelike-single/game/clientai"
	"github.com/kasworld/goguelike-single/lib/g2log"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func main() {
	configurl := flag.String("i", "", "config file or url")

	ads := argdefault.New(&textclientconfig.TextClientConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &textclientconfig.TextClientConfig{
		LogLevel: g2log.LL_All,
	}
	ads.SetDefaultToNonZeroField(config)
	if *configurl != "" {
		if err := configutil.LoadIni(*configurl, config); err != nil {
			g2log.Error("%v", err)
		}
	}
	ads.ApplyFlagTo(config)

	app := clientai.New(config)
	app.Run(context.Background())
	g2log.Error("%v", app.GetRunResult())
}
