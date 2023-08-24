// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/bmatcuk/doublestar/v4"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
)

// PrepareTasks walk through the dir and add tasks into task chan
// TODO: optimize function args
func PrepareTasks(paths []string, tmpl []byte, operation Operation, skipF []string, muteF bool, tmplF string) {
	for _, path := range paths {
		walkDir(path, tmpl, operation, skipF, muteF, tmplF)
	}
}

func walkDir(start string, tmpl []byte, operation Operation, skipF []string, muteF bool, tmplF string) {
	_ = filepath.WalkDir(start, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("walk dir error")
			return nil
		}
		if d.IsDir() {
			return nil
		}
		// determine if this file needs to be skipped
		if isMatch(path, skipF) {
			if !muteF {
				logrus.WithFields(logrus.Fields{
					"path": path,
				}).Infoln("skip file")
			}
			return nil
		}
		header := tmpl
		if tmplF == "" {
			// generate header according to the file type
			// NOTE: The file has not been modified yet
			header = generateHeader(path, tmpl)
		}
		switch operation {
		case Add:
			prepareAdd(path, d, header, muteF)
		case Update:
			prepareUpdate(path, d, header, muteF)
		case Remove:
			prepareRemove(path, d, header, muteF)
		case Check:
			prepareCheck(path, header, muteF)
		default:
			logrus.Errorln("no matched operation")
		}
		return nil
	})
}

func prepareCheck(path string, header []byte, muteF bool) {
	taskC <- func() {
		content, err := os.ReadFile(path)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("read file error")
			return
		}
		if isGenerated(content) {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Warnln("file is generated")
			return
		}
		// get the first index of the header in the file
		idx := bytes.Index(content, header)
		// not matched
		if idx != -1 && !muteF {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Infoln("file does have a matched header")
		}
		if !muteF {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Infoln("file does not have a matched header")
		}
	}
}

func prepareUpdate(path string, d fs.DirEntry, header []byte, muteF bool) {
	taskC <- func() {
		content, err := os.ReadFile(path)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("read file error")
			return
		}
		if !hasHeader(content) || isGenerated(content) {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Infoln("file does not have a header or is generated")
			return
		}
		// get the first line of the special file
		line := matchFirstLine(content)
		file, err := os.Open(path)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("open file error")
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
			logrus.Warnln("file close error")
		}
		// assemble license header and modify the file
		b := assemble(line, header, afterBlankLine, true)
		err = os.WriteFile(path, b, d.Type())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("write file error")
			return
		}
		if !muteF {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Infoln("file has been modified")
		}
	}
}

func prepareRemove(path string, d fs.DirEntry, header []byte, muteF bool) {
	taskC <- func() {
		content, err := os.ReadFile(path)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("read file error")
			return
		}
		if isGenerated(content) {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Warnln("file is generated")
			return
		}
		// get the first index of the header in the file
		idx := bytes.Index(content, header)
		if idx == -1 {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Warnln("file does not have a matched header")
			return
		}
		// remove the header of the file
		content = append(content[:idx], content[idx+len(header):]...)
		// modify the file
		err = os.WriteFile(path, content, d.Type())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("write file error")
			return
		}
		if !muteF {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Infoln("file has been modified")
		}
	}
}

func prepareAdd(path string, d fs.DirEntry, header []byte, muteF bool) {
	taskC <- func() {
		content, err := os.ReadFile(path)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("read file error")
			return
		}
		if hasHeader(content) || isGenerated(content) {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Warnln("file already has a header or is generated")
			return
		}
		// get the first line of the special file
		line := matchFirstLine(content)
		// assemble license header and modify the file
		b := assemble(line, header, content, false)
		err = os.WriteFile(path, b, d.Type())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("write file error")
			return
		}
		if !muteF {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Infoln("file has been modified")
		}
	}
}

func isMatch(path string, pattern []string) bool {
	for _, p := range pattern {
		if match, _ := doublestar.Match(p, path); match {
			return true
		}
	}
	return false
}
