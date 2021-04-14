// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package primenum

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"sync"
)

type PrimeIntList []int

func New() PrimeIntList {
	pn := PrimeIntList{2, 3}
	return pn
}

func NewWithCap(n int) PrimeIntList {
	pn := make(PrimeIntList, 0, n)
	pn = append(pn, PrimeIntList{2, 3}...)
	return pn
}

func (pn PrimeIntList) AppendFindTo(n int) PrimeIntList {
	last := pn[len(pn)-1]
	if last >= n {
		return pn
	}
	for i := last + 2; i <= n; i += 2 {
		if pn.CalcPrime(i) {
			pn = append(pn, i)
		}
	}
	return pn
}

func (pn PrimeIntList) FindPos(n int) (int, bool) {
	i := sort.SearchInts(pn, n)
	if i < len(pn) && pn[i] == n {
		// x is present at pn[i]
		return i, true
	} else {
		// x is not present in pn,
		// but i is the index where it would be inserted.
		return i, false
	}
}

func (pn PrimeIntList) MaxCanCheck() int {
	last := pn[len(pn)-1]
	return last * last
}

func (pn PrimeIntList) CanFindIn(n int) bool {
	return pn.MaxCanCheck() > n
}

func (pn PrimeIntList) CalcPrime(n int) bool {
	to := int(math.Sqrt(float64(n)))
	for _, v := range pn {
		if n%v == 0 {
			return false
		}
		if v > to {
			break
		}
	}
	return true
}

func (pn PrimeIntList) MultiAppendFindTo(n int) PrimeIntList {
	lastIndex := len(pn) - 1
	last := pn[lastIndex]
	if last >= n {
		return pn
	}

	if n >= last*last {
		pn = pn.MultiAppendFindTo(n / 2)
		lastIndex = len(pn) - 1
		last = pn[lastIndex]
	}

	bufl := runtime.NumCPU() * 1

	var wgWorker sync.WaitGroup
	var wgAppend sync.WaitGroup

	// recv result
	appendCh := make(chan int, bufl*2)
	wgAppend.Add(1)
	go func() {
		for n := range appendCh {
			pn = append(pn, n)
		}
		wgAppend.Done()
	}()

	// prepare need check data
	argCh := make(chan int, bufl*1000)
	go func() {
		for i := last + 2; i <= n; i += 2 {
			argCh <- i
		}
		close(argCh)
	}()

	// run worker
	for i := 0; i < bufl; i++ {
		wgWorker.Add(1)
		go func() {
			for n := range argCh {
				if pn.CalcPrime(n) {
					appendCh <- n
				}
			}
			wgWorker.Done()
		}()
	}
	wgWorker.Wait()
	close(appendCh)
	wgAppend.Wait()

	sort.Ints(pn[lastIndex+1:])
	return pn
}

func (pn PrimeIntList) MultiAppendFindTo2(n int) PrimeIntList {
	lastIndex := len(pn) - 1
	last := pn[lastIndex]
	if last >= n {
		return pn
	}

	if n >= last*last {
		pn = pn.MultiAppendFindTo2(n / 2)
		lastIndex = len(pn) - 1
		last = pn[lastIndex]
	}

	workerCount := runtime.NumCPU() * 1

	var wgWorker sync.WaitGroup
	var wgAppend sync.WaitGroup

	// recv result
	appendCh := make(chan int, workerCount*2)
	wgAppend.Add(1)
	go func() {
		for n := range appendCh {
			pn = append(pn, n)
		}
		wgAppend.Done()
	}()

	// run worker
	for workerid := 0; workerid < workerCount; workerid++ {
		wgWorker.Add(1)
		go func(workerid int) {
			for i := last + 2 + workerid*2; i <= n; i += workerCount * 2 {
				if pn.CalcPrime(i) {
					appendCh <- i
				}
			}
			wgWorker.Done()
		}(workerid)
	}
	wgWorker.Wait()
	close(appendCh)
	wgAppend.Wait()

	sort.Ints(pn[lastIndex+1:])
	return pn
}

func (pn PrimeIntList) MultiAppendFindTo3(n int) PrimeIntList {
	lastIndex := len(pn) - 1
	last := pn[lastIndex]
	if last >= n {
		return pn
	}

	if n >= last*last {
		pn = pn.MultiAppendFindTo3(n / 2)
		lastIndex = len(pn) - 1
		last = pn[lastIndex]
	}

	workerCount := runtime.NumCPU() * 1

	var wgWorker sync.WaitGroup

	workResult := make([]PrimeIntList, workerCount)
	// run worker
	for workerID := 0; workerID < workerCount; workerID++ {
		wgWorker.Add(1)
		go func(workerID int) {
			var rtn PrimeIntList
			for i := last + 2 + workerID*2; i <= n; i += workerCount * 2 {
				if pn.CalcPrime(i) {
					rtn = append(rtn, i)
				}
			}
			workResult[workerID] = rtn
			wgWorker.Done()
		}(workerID)
	}
	wgWorker.Wait()

	// merge sort
	workerPos := make([]int, workerCount)
	for true {
		minFound := 0
		minWorkerID := 0
		for workerID := 0; workerID < workerCount; workerID++ {
			pos := workerPos[workerID]
			if pos >= len(workResult[workerID]) {
				continue
			}
			if minFound == 0 || workResult[workerID][pos] < minFound {
				minFound = workResult[workerID][pos]
				minWorkerID = workerID
			}
		}
		if minFound != 0 {
			pn = append(pn, minFound)
			workerPos[minWorkerID]++
		} else {
			break
		}
	}
	return pn
}

func (pn PrimeIntList) MultiAppendFindTo4(n int) PrimeIntList {
	lastIndex := len(pn) - 1
	last := pn[lastIndex]
	if last >= n {
		return pn
	}

	if n >= last*last {
		pn = pn.MultiAppendFindTo4(n / 2)
		lastIndex = len(pn) - 1
		last = pn[lastIndex]
	}

	workerCount := runtime.NumCPU() * 1

	var wgWorker sync.WaitGroup

	workResult := make([]PrimeIntList, workerCount)
	workerBufferLen := (n - last) / workerCount / 16
	// run worker
	for workerID := 0; workerID < workerCount; workerID++ {
		wgWorker.Add(1)
		go func(workerID int) {
			rtn := make(PrimeIntList, 0, workerBufferLen)
			for i := last + 2 + workerID*2; i <= n; i += workerCount * 2 {
				if pn.CalcPrime(i) {
					rtn = append(rtn, i)
				}
			}
			workResult[workerID] = rtn
			wgWorker.Done()
		}(workerID)
	}
	wgWorker.Wait()

	// for workerID := 0; workerID < workerCount; workerID++ {
	// 	fmt.Printf("worker %v buflen %v datalen %v\n", workerID, workerBufferLen, len(workResult[workerID]))
	// }

	// merge sort
	workerPos := make([]int, workerCount)
	for true {
		minFound := 0
		minWorkerID := 0
		for workerID := 0; workerID < workerCount; workerID++ {
			pos := workerPos[workerID]
			if pos >= len(workResult[workerID]) {
				continue
			}
			if minFound == 0 || workResult[workerID][pos] < minFound {
				minFound = workResult[workerID][pos]
				minWorkerID = workerID
			}
		}
		if minFound != 0 {
			pn = append(pn, minFound)
			workerPos[minWorkerID]++
		} else {
			break
		}
	}
	return pn
}

func (pn PrimeIntList) Save(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("err in create %v", err)
	}
	defer fd.Close()
	enc := gob.NewEncoder(fd)
	err = enc.Encode(pn)
	if err != nil {
		return err
	}
	return nil
}

func LoadPrimeIntList(filename string) (PrimeIntList, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Fail to %v", err)
	}
	defer fd.Close()
	var rtn PrimeIntList
	dec := gob.NewDecoder(fd)
	err = dec.Decode(&rtn)
	if err != nil {
		return nil, err
	}
	return rtn, nil
}
