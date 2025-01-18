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
	"bufio"
	"bytes"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bmatcuk/doublestar/v4"
)

// TODO: refactor

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
var counter = struct {
	scanned    int // files have been read
	matched    int // files license headers matched as required
	mismatched int // files license headers do not match as required
	modified   int // files have been modified (e.g. add, update, remove license header)
	skipped    int // file paths match the skip pattern
	failed     int // unexpected error occurred
}{}

func walkDir(pattern string, tmpl []byte, operation Operation, skips []string, raw, fuzzy bool) {
	// TODO: add doc notes for walkDir started from root (.)
	if err := filepath.WalkDir(_root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			counter.failed++
			slog.Error("walk dir error", slog.String("path", path), slog.String("err", err.Error()))
			return nil
		}

		// match glob pattern
		match, err := doublestar.Match(pattern, path)
		if err != nil {
			counter.failed++
			slog.Error("match doublestar pattern error", slog.String("path", path), slog.String("err", err.Error()))
			return nil
		}
		if !match {
			return nil
		}

		// ignore dir
		if d.IsDir() {
			return nil
		}

		// determine if this file needs to be skipped
		if isSkip(path, skips) {
			counter.skipped++
			slog.Info("skip file", slog.String("path", path))
			return nil
		}

		// generate header or use tmpl
		header := tmpl
		if !raw {
			// generate header according to the file type
			// NOTE: The file has not been modified yet
			header, err = generateHeader(path, tmpl)
			if err != nil {
				counter.failed++
				slog.Warn(err.Error(), slog.String("path", path))
				return nil
			}
		}

		// submit task
		switch operation {
		case Add:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareAdd(path, d, header)
			}()
		case Update:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareUpdate(path, d, header)
			}()
		case Remove:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareRemove(path, d, header)
			}()
		case Check:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareCheck(path, header, fuzzy)
			}()
		default:
			slog.Warn("not a valid operation")
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func prepareCheck(path string, header []byte, fuzzy bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		if isGenerated(content) {
			slog.Warn("file is generated, won't check", slog.String("path", path))
			return
		}

		// fuzzy matching
		if fuzzy {
			header = removeYear(header)
			content = removeYear(content)
		}

		// get the first index of the header in the file
		idx := bytes.Index(content, header)
		// matched
		if idx != -1 {
			counter.matched++
			slog.Info("file has a matched header", slog.String("path", path))
			return
		}
		// mismatched
		counter.mismatched++
		slog.Warn("file does not have a matched header", slog.String("path", path))
	}
}

func prepareUpdate(path string, d fs.DirEntry, header []byte) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		if !hasHeader(content) || isGenerated(content) {
			counter.failed++
			slog.Warn("file does not have a header or is generated", slog.String("path", path))
			return
		}

		// get the first line of the special file
		line := matchShebang(content)
		file, err := os.Open(path)
		if err != nil {
			counter.failed++
			slog.Error("open file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			l := scanner.Bytes()
			if len(l) == 0 {
				break
			}
		}
		afterBlankLine := make([]byte, 0)
		// NOTE: scanner will not scan from the beginning
		for scanner.Scan() {
			afterBlankLine = append(afterBlankLine, scanner.Bytes()...)
			afterBlankLine = append(afterBlankLine, '\n')
		}
		err = file.Close()
		if err != nil {
			slog.Error("file close error")
		}
		// assemble license header and modify the file
		b := assemble(line, header, afterBlankLine, true)
		err = os.WriteFile(path, b, d.Type())
		if err != nil {
			counter.failed++
			slog.Error("write file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		counter.modified++
		slog.Info("file has been modified", slog.String("path", path))
	}
}

func prepareRemove(path string, d fs.DirEntry, header []byte) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		if isGenerated(content) {
			counter.failed++
			slog.Warn("file is generated", slog.String("path", path))
			return
		}

		// get the first index of the header in the file
		idx := bytes.Index(content, header)
		if idx == -1 {
			counter.failed++
			slog.Warn("file does not have a matched header", slog.String("path", path))
			return
		}

		// remove the header of the file
		content = append(content[:idx], content[idx+len(header):]...)
		// modify the file
		err = os.WriteFile(path, content, d.Type())
		if err != nil {
			counter.failed++
			slog.Error("write file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		counter.modified++
		slog.Info("file has been modified", slog.String("path", path))
	}
}

func prepareAdd(path string, d fs.DirEntry, header []byte) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		// TODO: split hasHeader and isGenerated
		// TODO: do not count hasHeader and isGenerated as failed
		if hasHeader(content) || isGenerated(content) {
			counter.failed++
			slog.Warn("file already has a header or is generated", slog.String("path", path))
			return
		}

		// get the first line of the special file
		line := matchShebang(content)
		// assemble license header and modify the file
		b := assemble(line, header, content, false)
		err = os.WriteFile(path, b, d.Type())
		if err != nil {
			counter.failed++
			slog.Error("write file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		counter.modified++
		slog.Info("file has been modified", slog.String("path", path))
	}
}

func isSkip(path string, pattern []string) bool {
	for _, p := range pattern {
		if match, err := doublestar.Match(p, path); match {
			if err != nil {
				slog.Error("skip pattern match error", slog.String("path", path), slog.String("pattern", p))
				return false
			}
			return true
		}
	}
	return false
}

func removeYear(b []byte) []byte {
	return regexp.MustCompile(`\b\d{4}\b`).ReplaceAll(b, []byte{})
}
