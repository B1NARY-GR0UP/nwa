// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"fmt"
	"os"
	"sync"
)

const _size = 1000

var (
	taskC  chan func()
	taskWG sync.WaitGroup
)

// lock-free because of serial
//
// Add, Update, Remove:
// - scanned
// - modified
// - skipped
// - failed
//
// Check:
// - scanned
// - matched
// - mismatched
// - skipped
// - failed
var counter struct {
	scanned    int // files have been read
	matched    int // files license headers matched as required
	mismatched int // files license headers do not match as required
	modified   int // files have been modified (e.g. add, update, remove license header)
	skipped    int // file paths match the skip pattern
	failed     int // unexpected error occurred
}

// PrepareTasks walk through the dir and add tasks into task chan
// TODO: optimize params
func PrepareTasks(paths []string, tmpl []byte, operation Operation, skips, keywords, styles []string, raw, fuzzy bool) {
	counter = struct {
		scanned    int
		matched    int
		mismatched int
		modified   int
		skipped    int
		failed     int
	}{}
	taskC = make(chan func(), _size)

	for _, path := range paths {
		walkDir(path, tmpl, operation, skips, keywords, styles, raw, fuzzy)
	}

	go func() {
		taskWG.Wait()
		close(taskC)
	}()
}

func ExecuteTasks(operation Operation, muteF bool) {
	for task := range taskC {
		task()
	}

	switch operation {
	case Add, Update, Remove:
		if !muteF {
			fmt.Printf("[NWA SUMMARY] scanned=%d modified=%d skipped=%d failed=%d\n", counter.scanned, counter.modified, counter.skipped, counter.failed)
		}
		if counter.failed > 0 {
			os.Exit(1)
		}
	case Check:
		if !muteF {
			fmt.Printf("[NWA SUMMARY] scanned=%d matched=%d mismatched=%d skipped=%d failed=%d\n", counter.scanned, counter.matched, counter.mismatched, counter.skipped, counter.failed)
		}
		// exit 1 to fail ci check
		if counter.mismatched > 0 || counter.failed > 0 {
			os.Exit(1)
		}
	default:
		panic("not a valid operation")
	}
}
