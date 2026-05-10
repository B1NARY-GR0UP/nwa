package internal

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// --- Label colors (level prefix [ERROR] [WARN] [INFO] [DRY-RUN]) ---
var (
	ErrorLabelColor  = color.New(color.FgRed, color.Bold)  // [ERROR]
	WarnLabelColor   = color.New(color.FgYellow)           // [WARN]
	InfoLabelColor   = color.New(color.FgCyan, color.Bold) // [INFO]
	DryRunLabelColor = color.New(color.FgCyan, color.Bold) // [DRY-RUN]
)

// --- Tag colors (operation/action: ADD UPDATE REMOVE CHECK SKIP) ---
var (
	AddTagColor    = color.New(color.FgGreen, color.Bold)   // ADD
	UpdateTagColor = color.New(color.FgMagenta, color.Bold) // UPDATE
	RemoveTagColor = color.New(color.FgRed, color.Bold)     // REMOVE
	CheckTagColor  = color.New(color.FgCyan, color.Bold)    // CHECK
	SkipTagColor   = color.New(color.FgYellow)              // SKIP
)

// --- Other ---
var (
	PathColor    = color.New(color.FgWhite)            // file path
	SummaryColor = color.New(color.FgBlue, color.Bold) // [NWA SUMMARY] header
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
		body = PathColor.Sprint(path) + ": " + msg
	case path != "":
		body = PathColor.Sprint(path)
	case msg != "":
		body = msg
	}

	if tagPart != "" {
		fmt.Fprintf(w, "%s %s %s\n", lc.Sprintf("[%s]", ls), tagPart, body)
	} else {
		fmt.Fprintf(w, "%s %s\n", lc.Sprintf("[%s]", ls), body)
	}
}
