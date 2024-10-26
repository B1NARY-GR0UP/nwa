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
//

package util

import (
	"bufio"
	"bytes"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

func walkDir(start string, tmpl []byte, operation Operation, skipF []string, muteF bool, rawTmpl bool) {
	_ = filepath.WalkDir(start, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("walk dir error", slog.String("path", path), slog.String("err", err.Error()))
			return nil
		}
		if d.IsDir() {
			return nil
		}
		// determine if this file needs to be skipped
		if isSkip(path, skipF) {
			if !muteF {
				slog.Info("skip file", slog.String("path", path))
			}
			return nil
		}
		header := tmpl
		if !rawTmpl {
			// generate header according to the file type
			// NOTE: The file has not been modified yet
			header, err = generateHeader(path, tmpl)
			if err != nil {
				slog.Warn(err.Error(), slog.String("path", path))
				return nil
			}
		}
		switch operation {
		case Add:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareAdd(path, d, header, muteF)
			}()
		case Update:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareUpdate(path, d, header, muteF)
			}()
		case Remove:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareRemove(path, d, header, muteF)
			}()
		case Check:
			wg.Add(1)
			go func() {
				defer wg.Done()
				taskC <- prepareCheck(path, header, muteF)
			}()
		default:
			slog.Warn("not a valid operation")
		}
		return nil
	})
}

func prepareCheck(path string, header []byte, muteF bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		if isGenerated(content) {
			slog.Warn("file is generated", slog.String("path", path))
			return
		}
		// get the first index of the header in the file
		idx := bytes.Index(content, header)
		// not matched
		if idx != -1 && !muteF {
			slog.Info("file has a matched header", slog.String("path", path))
			return
		}
		if !muteF {
			slog.Info("file does not have a matched header", slog.String("path", path))
		}
	}
}

func prepareUpdate(path string, d fs.DirEntry, header []byte, muteF bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		if !hasHeader(content) || isGenerated(content) {
			slog.Warn("file does not have a header or is generated", slog.String("path", path))
			return
		}
		// get the first line of the special file
		line := matchShebang(content)
		file, err := os.Open(path)
		if err != nil {
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
			slog.Error("write file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		if !muteF {
			slog.Info("file has been modified", slog.String("path", path))
		}
	}
}

func prepareRemove(path string, d fs.DirEntry, header []byte, muteF bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		if isGenerated(content) {
			slog.Warn("file is generated", slog.String("path", path))
			return
		}
		// get the first index of the header in the file
		idx := bytes.Index(content, header)
		if idx == -1 {
			slog.Warn("file does not have a matched header", slog.String("path", path))
			return
		}
		// remove the header of the file
		content = append(content[:idx], content[idx+len(header):]...)
		// modify the file
		err = os.WriteFile(path, content, d.Type())
		if err != nil {
			slog.Error("write file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		if !muteF {
			slog.Info("file has been modified", slog.String("path", path))
		}
	}
}

func prepareAdd(path string, d fs.DirEntry, header []byte, muteF bool) func() {
	return func() {
		content, err := os.ReadFile(path)
		if err != nil {
			slog.Error("read file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		if hasHeader(content) || isGenerated(content) {
			slog.Warn("file already has a header or is generated", slog.String("path", path))
			return
		}
		// get the first line of the special file
		line := matchShebang(content)
		// assemble license header and modify the file
		b := assemble(line, header, content, false)
		err = os.WriteFile(path, b, d.Type())
		if err != nil {
			slog.Error("write file error", slog.String("path", path), slog.String("err", err.Error()))
			return
		}
		if !muteF {
			slog.Info("file has been modified", slog.String("path", path))
		}
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
