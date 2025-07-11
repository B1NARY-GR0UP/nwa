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
	"bufio"
	"bytes"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bmatcuk/doublestar/v4"
)

type Operation string

const (
	Add    Operation = "add"
	Update Operation = "update"
	Remove Operation = "remove"
	Check  Operation = "check"
)

const _root = "."

func walkDir(pattern string, tmpl []byte, operation Operation, skips, keywords, styles []string, raw, fuzzy bool) {
	if err := filepath.WalkDir(_root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			counter.failed++
			slog.Error("walk dir error", slog.String("path", path), slog.String("err", err.Error()))
			return nil
		}

		// convert Windows path separators
		path = filepath.ToSlash(path)

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
			header, err = generateHeader(path, tmpl, styles)
			if err != nil {
				counter.failed++
				slog.Warn(err.Error(), slog.String("path", path))
				return nil
			}
		}

		// submit task
		switch operation {
		case Add:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doAdd(path, d, header, keywords)
			}()
		case Update:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doUpdate(path, d, header, keywords)
			}()
		case Remove:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doRemove(path, d, header, fuzzy)
			}()
		case Check:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doCheck(path, header, fuzzy)
			}()
		default:
			slog.Warn("not a valid operation")
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func doCheck(path string, header []byte, fuzzy bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		if isGenerated(content) {
			slog.Warn("file is generated, won't be checked", slog.String("path", path))
			return
		}

		// standardize line separator
		content = standardizeLineSeparator(content)
		header = standardizeLineSeparator(header)

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

func doUpdate(path string, d fs.DirEntry, header []byte, keywords []string) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		// check generated first
		if isGenerated(content) {
			slog.Warn("file is generated, won't be modified", slog.String("path", path))
			return
		}
		if !hasHeader(content, keywords) {
			slog.Warn("file does not have a header", slog.String("path", path))
			return
		}

		// get the shebang of the special file
		shebang := matchShebang(content)
		// check if file has a utf8BOM
		hasBOM := matchBOM(content)

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

		// add a blank line at the end of the header
		header = append(header, '\n')

		// assemble license header and modify the file
		b := assemble(shebang, header, afterBlankLine, hasBOM, true)

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

func doRemove(path string, d fs.DirEntry, header []byte, fuzzy bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		if isGenerated(content) {
			slog.Warn("file is generated, won't be modified", slog.String("path", path))
			return
		}

		// standardize line separator
		content = standardizeLineSeparator(content)
		header = standardizeLineSeparator(header)

		// fuzzy matching
		if fuzzy {
			header = removeYear(header)
			content = removeYear(content)
		}

		// get the first index of the header in the file
		idx := bytes.Index(content, header)
		if idx == -1 {
			counter.failed++
			slog.Warn("file does not have a matched header", slog.String("path", path))
			return
		}

		// if exist a blank line after the header, remove it
		headerIdx := idx + len(header)
		if headerIdx < len(content) {
			nextNewLineIdx := bytes.IndexByte(content[headerIdx:], '\n')
			if nextNewLineIdx == 0 {
				headerIdx++
			}
		}

		// remove the header of the file
		content = append(content[:idx], content[headerIdx:]...)
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

func doAdd(path string, d fs.DirEntry, header []byte, keywords []string) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}

		counter.scanned++

		// check generated first
		if isGenerated(content) {
			slog.Warn("file is generated, won't be modified", slog.String("path", path))
			return
		}
		if hasHeader(content, keywords) {
			slog.Warn("file already has a header", slog.String("path", path))
			return
		}

		// get the shebang of the special file
		shebang := matchShebang(content)
		// check if file has a utf8BOM
		hasBOM := matchBOM(content)

		// add a blank line at the end of the header
		header = append(header, '\n')

		// assemble license header and modify the file
		b := assemble(shebang, header, content, hasBOM, false)

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

// convert CR and CRLF line separator to LF
func standardizeLineSeparator(b []byte) []byte {
	// CRLF => LF
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	// CR => LF
	b = bytes.ReplaceAll(b, []byte("\r"), []byte("\n"))
	return b
}

var yearRE = regexp.MustCompile(`\b\d{4}\b`)

func removeYear(b []byte) []byte {
	loc := yearRE.FindIndex(b)
	if loc == nil {
		return b
	}
	return append(b[:loc[0]], b[loc[1]:]...)
}
