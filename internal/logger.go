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

	"github.com/fatih/color"
)

// Output format:
//
//	[LABEL] TAG msg: path        (msg + path)
//	[LABEL] TAG path             (path only)
//	[LABEL] TAG msg              (msg only)
//	[NWA SUMMARY] scanned=...    (summary)
//
// Labels: [ERROR] [WARN] [INFO] [DRY-RUN]
// Tags:   ADD UPDATE REMOVE CHECK SKIP
//
// Dry-run output goes to stdout, log output goes to stderr.

var (
	ErrorLabelColor  = color.New(color.FgRed, color.Bold)    // [ERROR]
	WarnLabelColor   = color.New(color.FgYellow, color.Bold) // [WARN]
	InfoLabelColor   = color.New(color.FgCyan, color.Bold)   // [INFO]
	DryRunLabelColor = color.New(color.FgCyan, color.Bold)   // [DRY-RUN]
)

var (
	AddTagColor    = color.New(color.FgGreen, color.Bold)   // ADD
	UpdateTagColor = color.New(color.FgMagenta, color.Bold) // UPDATE
	RemoveTagColor = color.New(color.FgRed, color.Bold)     // REMOVE
	CheckTagColor  = color.New(color.FgCyan, color.Bold)    // CHECK
	SkipTagColor   = color.New(color.FgYellow, color.Bold)  // SKIP
)

var (
	PathColor    = color.New(color.FgHiWhite, color.Bold) // file path
	SummaryColor = color.New(color.FgBlue, color.Bold)    // [NWA SUMMARY] header
)

const (
	TagAdd    = "ADD  "
	TagUpdate = "UPDATE"
	TagRemove = "REMOVE"
	TagCheck  = "CHECK"
	TagSkip   = "SKIP "
)

type Level int

const (
	LvlError  Level = iota // [ERROR]  → stderr
	LvlWarn                // [WARN]   → stderr
	LvlInfo                // [INFO]   → stderr
	LvlDryRun              // [DRY-RUN]→ stdout
	LvlMute   Level = -1   // silence stderr
)

var logLevel = LvlWarn // default: Error + Warn

// SetLevel sets the minimum log level for stderr output.
func SetLevel(level Level) { logLevel = level }

// IsMuted reports whether stderr output is completely silenced.
func IsMuted() bool { return logLevel == LvlMute }

// opColor returns the tag color for an operation.
func opColor(op Operation) *color.Color {
	switch op {
	case OpAdd:
		return AddTagColor
	case OpUpdate:
		return UpdateTagColor
	case OpRemove:
		return RemoveTagColor
	case OpCheck:
		return CheckTagColor
	default:
		panic("unknown operation")
	}
	return nil
}

func labelColor(level Level) *color.Color {
	switch level {
	case LvlError:
		return ErrorLabelColor
	case LvlWarn:
		return WarnLabelColor
	case LvlInfo:
		return InfoLabelColor
	case LvlDryRun:
		return DryRunLabelColor
	default:
		panic("unknown log level")
	}
}

func labelString(level Level) string {
	switch level {
	case LvlError:
		return "ERROR"
	case LvlWarn:
		return "WARN"
	case LvlInfo:
		return "INFO"
	case LvlDryRun:
		return "DRY-RUN"
	default:
		panic("unknown log level")
	}
}

// Logf writes a colored log line to stdout (LvlDryRun) or stderr (all others).
func Logf(level Level, tag string, tagColor *color.Color, path, format string, args ...any) {
	if level == LvlDryRun {
		printLine(os.Stdout, level, tag, tagColor, path, fmt.Sprintf(format, args...))
		return
	}
	if level > logLevel {
		return
	}
	printLine(os.Stderr, level, tag, tagColor, path, fmt.Sprintf(format, args...))
}

// Log writes a colored log line with a pre-formatted message.
func Log(level Level, tag string, tagColor *color.Color, path, msg string) {
	if level == LvlDryRun {
		printLine(os.Stdout, level, tag, tagColor, path, msg)
		return
	}
	if level > logLevel {
		return
	}
	printLine(os.Stderr, level, tag, tagColor, path, msg)
}

func printLine(w *os.File, level Level, tag string, tagColor *color.Color, path, msg string) {
	lc := labelColor(level)
	if lc == nil {
		lc = color.New()
	}
	ls := labelString(level)

	var tagPart string
	if tagColor != nil && tag != "" {
		tagPart = tagColor.Sprint(tag)
	} else if tag != "" {
		tagPart = tag
	}

	var body string
	switch {
	case path != "" && msg != "":
		body = msg + ": " + PathColor.Sprint(path)
	case path != "":
		body = PathColor.Sprint(path)
	case msg != "":
		body = msg
	}

	if tagPart != "" {
		_, _ = fmt.Fprintf(w, "%s %s %s\n", lc.Sprintf("[%s]", ls), tagPart, body)
	} else {
		_, _ = fmt.Fprintf(w, "%s %s\n", lc.Sprintf("[%s]", ls), body)
	}
}
