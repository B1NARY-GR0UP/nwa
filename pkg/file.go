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

package pkg

import (
	"bytes"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
)

// PrepareAddTasks walk through the dir and add tasks into task chan
// TODO: replace *bytes.Buffer with []byte
// TODO: optimize function args
func PrepareAddTasks(paths []string, tmpl *bytes.Buffer, skipF []string, muteF bool, tmplF string) {
	for _, path := range paths {
		walkDir(path, tmpl, skipF, muteF, tmplF)
	}
}

func walkDir(start string, tmpl *bytes.Buffer, skipF []string, muteF bool, tmplF string) {
	_ = filepath.WalkDir(start, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
				"err":  err,
			}).Errorln("walk dir error")
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
		taskC <- func() {
			header := tmpl.Bytes()
			if tmplF == "" {
				// generate header according to the file type
				// NOTE: The file has not been modified yet
				header = generateHeader(path, tmpl)
			}
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
			line := matchFirstLine(content)
			// modify the file
			err = os.WriteFile(path, assemble(line, header, content), d.Type())
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
		return nil
	})
}

func isMatch(path string, pattern []string) bool {
	for _, p := range pattern {
		if match, _ := doublestar.Match(p, path); match {
			return true
		}
	}
	return false
}
