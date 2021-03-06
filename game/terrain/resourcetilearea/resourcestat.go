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

package resourcetilearea

import "github.com/kasworld/goguelike-single/enum/resourcetype_vector_int"

func (rta ResourceTileArea) CalcStat() *resourcetype_vector_int.ResourceTypeVector_int {
	var rtn resourcetype_vector_int.ResourceTypeVector_int
	for _, xv := range rta {
		for _, yv := range xv {
			for i, v := range yv {
				rtn[i] += int(v)
			}
		}
	}
	return &rtn
}
