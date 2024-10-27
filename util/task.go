// Copyright 2024 BINARY Members
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
	"log/slog"
	"strconv"
	"sync"
)

const _size = 1000

var taskC = make(chan func(), _size)

var wg sync.WaitGroup

// PrepareTasks walk through the dir and add tasks into task chan
func PrepareTasks(paths []string, tmpl []byte, operation Operation, skipF []string, muteF bool, rawTmpl bool) {
	for _, path := range paths {
		walkDir(path, tmpl, operation, skipF, muteF, rawTmpl)
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
		if muteF {
			return
		}
		slog.Info("Summary",
			slog.String("modified", strconv.Itoa(counter.modified)),
			slog.String("skipped", strconv.Itoa(counter.skipped)),
			slog.String("failed", strconv.Itoa(counter.failed)),
		)
	case Check:
		if muteF {
			return
		}
		slog.Info("Summary",
			slog.String("matched", strconv.Itoa(counter.matched)),
			slog.String("mismatched", strconv.Itoa(counter.mismatched)),
			slog.String("skipped", strconv.Itoa(counter.skipped)),
			slog.String("failed", strconv.Itoa(counter.failed)),
		)
	default:
		slog.Warn("not a valid operation")
	}
}
