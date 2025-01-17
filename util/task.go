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

package util

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

const _size = 1000

var taskC = make(chan func(), _size)

var wg sync.WaitGroup

// PrepareTasks walk through the dir and add tasks into task chan
func PrepareTasks(paths []string, tmpl []byte, operation Operation, skips []string, raw, fuzzy bool) {
	for _, path := range paths {
		walkDir(path, tmpl, operation, skips, raw, fuzzy)
	}
	go func() {
		wg.Wait()
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
			fmt.Printf("[NWA SUMMARY] modified=%d skipped=%d failed=%d\n", counter.modified, counter.skipped, counter.failed)
		}
		// exit 1 to fail pre-commit
		if counter.failed > 0 {
			os.Exit(1)
		}
	case Check:
		if !muteF {
			fmt.Printf("[NWA SUMMARY] matched=%d mismatched=%d skipped=%d failed=%d\n", counter.matched, counter.mismatched, counter.skipped, counter.failed)
		}
		// exit 1 to fail ci check
		if counter.mismatched > 0 || counter.failed > 0 {
			os.Exit(1)
		}
	default:
		slog.Warn("not a valid operation")
	}
}
