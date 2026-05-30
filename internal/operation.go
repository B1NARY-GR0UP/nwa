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
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bmatcuk/doublestar/v4"
)

type Operation string

const (
	OpAdd    Operation = "ADD"
	OpUpdate Operation = "UPDATE"
	OpRemove Operation = "REMOVE"
	OpCheck  Operation = "CHECK"
)

const _root = "."

func walkDir(pattern string, tmpl []byte, operation Operation, skips, keywords, styles []string, raw, fuzzy, dryRun, diff bool) {
	tag := string(operation)
	tagColor := opColor(operation)

	if err := filepath.WalkDir(_root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			counter.failed++
			Logf(LvlError, tag, tagColor, path, "walk dir error: %v", err)
			return nil
		}

		// convert Windows path separators
		path = filepath.ToSlash(path)

		// match glob pattern
		match, err := doublestar.Match(pattern, path)
		if err != nil {
			counter.failed++
			Logf(LvlError, tag, tagColor, path, "match doublestar pattern error: %v", err)
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
			Log(LvlInfo, tag, tagColor, path, "skip file")
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
				Log(LvlWarn, tag, tagColor, path, err.Error())
				return nil
			}
		}

		// submit task
		switch operation {
		case OpAdd:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doAdd(path, d, header, keywords, dryRun)
			}()
		case OpUpdate:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doUpdate(path, d, header, keywords, dryRun)
			}()
		case OpRemove:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doRemove(path, d, header, fuzzy, dryRun)
			}()
		case OpCheck:
			taskWG.Add(1)
			go func() {
				defer taskWG.Done()
				taskC <- doCheck(path, header, fuzzy, diff)
			}()
		default:
			Logf(LvlWarn, tag, tagColor, "", "not a valid operation")
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func doCheck(path string, header []byte, fuzzy, diff bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			Logf(LvlError, TagCheck, CheckTagColor, path, "read file error: %v", err)
			return
		}

		counter.scanned++

		if isGenerated(content) {
			Log(LvlWarn, TagCheck, CheckTagColor, path, "generated file, won't be checked")
			return
		}

		// standardize line separator
		header = standardizeLineSeparator(header)
		content = standardizeLineSeparator(content)

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
			Log(LvlInfo, TagCheck, CheckTagColor, path, "file has a matched header")
			return
		}
		// mismatched
		counter.mismatched++

		Log(LvlWarn, TagCheck, CheckTagColor, path, "file does not have a matched header")
		if diff {
			// +1 converts newline count to line count, +5 provides context margin
			currentHeader := extractHeader(content, bytes.Count(header, []byte("\n"))+1+5)
			DiffReport(os.Stderr, header, currentHeader)
			_, _ = fmt.Fprintln(os.Stderr)
		}
	}
}

func doUpdate(path string, d fs.DirEntry, header []byte, keywords []string, dryRun bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			Logf(LvlError, TagUpdate, UpdateTagColor, path, "read file error: %v", err)
			return
		}

		counter.scanned++

		// check generated first
		if isGenerated(content) {
			if dryRun {
				Log(LvlDryRun, TagSkip, SkipTagColor, path, "generated file")
			} else {
				Log(LvlWarn, TagUpdate, UpdateTagColor, path, "generated file, won't be modified")
			}
			return
		}
		if !hasHeader(content, keywords) {
			if dryRun {
				Log(LvlDryRun, TagSkip, SkipTagColor, path, "no header found")
			} else {
				Log(LvlWarn, TagUpdate, UpdateTagColor, path, "file does not have a header")
			}
			return
		}

		// detect line ending style before standardizing
		ending := detectLineEnding(content)

		// standardize line separator for internal processing
		stdContent := standardizeLineSeparator(content)

		// get the shebang of the special file
		shebang := matchShebang(stdContent)
		// check if file has a utf8BOM
		hasBOM := matchBOM(stdContent)

		// find the first blank line using bytes.Index on standardized content
		blankLineIdx := bytes.Index(stdContent, []byte("\n\n"))
		var afterBlankLine []byte
		if blankLineIdx == -1 {
			// no blank line found, entire content is header
			afterBlankLine = nil
		} else {
			afterBlankLine = stdContent[blankLineIdx+2:]
		}

		// add a blank line at the end of the header
		header = append(header, '\n')

		// assemble license header and modify the file
		b := assemble(shebang, header, afterBlankLine, hasBOM, true)

		// restore original line ending style
		b = convertLineEnding(b, ending)

		if dryRun {
			counter.wouldModify++
			Log(LvlDryRun, TagUpdate, UpdateTagColor, path, "")
			return
		}
		err = os.WriteFile(path, b, d.Type())
		if err != nil {
			counter.failed++
			Logf(LvlError, TagUpdate, UpdateTagColor, path, "write file error: %v", err)
			return
		}
		counter.modified++
		Log(LvlInfo, TagUpdate, UpdateTagColor, path, "file has been modified")
	}
}

func doRemove(path string, d fs.DirEntry, header []byte, fuzzy, dryRun bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			Logf(LvlError, TagRemove, RemoveTagColor, path, "read file error: %v", err)
			return
		}

		counter.scanned++

		if isGenerated(content) {
			if dryRun {
				Log(LvlDryRun, TagSkip, SkipTagColor, path, "generated file")
			} else {
				Log(LvlWarn, TagRemove, RemoveTagColor, path, "generated file, won't be modified")
			}
			return
		}

		// detect line ending style before standardizing
		ending := detectLineEnding(content)

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
			if dryRun {
				Log(LvlDryRun, TagSkip, SkipTagColor, path, "no matched header")
			} else {
				Log(LvlWarn, TagRemove, RemoveTagColor, path, "file does not have a matched header")
			}
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

		// restore original line ending style
		content = convertLineEnding(content, ending)

		if dryRun {
			counter.wouldModify++
			Log(LvlDryRun, TagRemove, RemoveTagColor, path, "")
			return
		}
		// modify the file
		err = os.WriteFile(path, content, d.Type())
		if err != nil {
			counter.failed++
			Logf(LvlError, TagRemove, RemoveTagColor, path, "write file error: %v", err)
			return
		}
		counter.modified++
		Log(LvlInfo, TagRemove, RemoveTagColor, path, "file has been modified")
	}
}

func doAdd(path string, d fs.DirEntry, header []byte, keywords []string, dryRun bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			counter.failed++
			Logf(LvlError, TagAdd, AddTagColor, path, "read file error: %v", err)
			return
		}

		counter.scanned++

		// check generated first
		if isGenerated(content) {
			if dryRun {
				Log(LvlDryRun, TagSkip, SkipTagColor, path, "generated file")
			} else {
				Log(LvlWarn, TagAdd, AddTagColor, path, "generated file, won't be modified")
			}
			return
		}
		if hasHeader(content, keywords) {
			if dryRun {
				Log(LvlDryRun, TagSkip, SkipTagColor, path, "already has header")
			} else {
				Log(LvlWarn, TagAdd, AddTagColor, path, "file already has a header")
			}
			return
		}

		// detect line ending style before standardizing
		ending := detectLineEnding(content)

		// standardize content for shebang matching
		stdContent := standardizeLineSeparator(content)

		// get the shebang of the special file
		shebang := matchShebang(stdContent)
		// check if file has a utf8BOM
		hasBOM := matchBOM(stdContent)

		// add a blank line at the end of the header
		header = append(header, '\n')

		// assemble license header and modify the file
		b := assemble(shebang, header, stdContent, hasBOM, false)

		// restore original line ending style
		b = convertLineEnding(b, ending)

		if dryRun {
			counter.wouldModify++
			Log(LvlDryRun, TagAdd, AddTagColor, path, "")
			return
		}
		err = os.WriteFile(path, b, d.Type())
		if err != nil {
			counter.failed++
			Logf(LvlError, TagAdd, AddTagColor, path, "write file error: %v", err)
			return
		}
		counter.modified++
		Log(LvlInfo, TagAdd, AddTagColor, path, "file has been modified")
	}
}

func isSkip(path string, pattern []string) bool {
	for _, p := range pattern {
		if match, err := doublestar.Match(p, path); match {
			if err != nil {
				Logf(LvlError, "", nil, path, "skip pattern match error: %s", p)
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

// detectLineEnding detects the dominant line ending style in the content.
// It returns "\r\n" for CRLF or "\n" for LF.
// Empty files or files with no line breaks default to "\n".
func detectLineEnding(b []byte) []byte {
	crlfCount := bytes.Count(b, []byte("\r\n"))
	lfOnly := bytes.Count(bytes.ReplaceAll(b, []byte("\r\n"), []byte("")), []byte("\n"))
	if crlfCount == 0 {
		return []byte("\n")
	}
	if crlfCount >= lfOnly {
		return []byte("\r\n")
	}
	return []byte("\n")
}

// convertLineEnding converts all line endings in the content to the target style.
// It first standardizes to LF, then replaces LF with the target ending.
func convertLineEnding(b []byte, ending []byte) []byte {
	b = standardizeLineSeparator(b)
	if string(ending) == "\n" {
		return b
	}
	return bytes.ReplaceAll(b, []byte("\n"), ending)
}

var yearRE = regexp.MustCompile(`\b\d{4}\b`)

func removeYear(b []byte) []byte {
	loc := yearRE.FindIndex(b)
	if loc == nil {
		return b
	}
	return append(b[:loc[0]], b[loc[1]:]...)
}
