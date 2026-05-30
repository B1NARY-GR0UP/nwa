// Copyright 2026 BINARY Members
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
	"io"
	"slices"
	"strings"

	"github.com/fatih/color"
)

var (
	diffDelColor  = color.New(color.FgRed)
	diffAddColor  = color.New(color.FgGreen)
	diffHunkColor = color.New(color.FgCyan)
	diffHintColor = color.New(color.FgHiWhite, color.Bold)
)

// extractHeader returns the license header block from the beginning of content.
// It finds the first blank line (\n\n) and returns everything before it,
// consistent with doUpdate's header boundary detection.
// If no blank line is found, content is capped at limit lines.
func extractHeader(content []byte, limit int) []byte {
	if idx := bytes.Index(content, []byte("\n\n")); idx != -1 {
		return content[:idx]
	}
	lines := 0
	for i, c := range content {
		if c == '\n' {
			lines++
			if lines >= limit {
				return content[:i]
			}
		}
	}
	return content
}

type diffEdit struct {
	typ  byte // ' ', '-', '+'
	line string
}

// DiffReport computes and writes a colored unified diff to w.
func DiffReport(w io.Writer, expected, actual []byte) {
	a := splitLines(expected)
	b := splitLines(actual)
	dp := lcs(a, b)
	edits := backtrack(dp, a, b)

	// Compute visible width for adaptive separator
	maxWidth := 0
	for _, e := range edits {
		lw := 1 + len(e.line) // prefix + content
		if lw > maxWidth {
			maxWidth = lw
		}
	}
	if maxWidth < 12 {
		maxWidth = 12
	}
	sep := strings.Repeat("─", maxWidth) // ─ (horizontal box drawing)

	// Top separator
	_, _ = fmt.Fprintln(w, sep)

	// Hint
	_, _ = fmt.Fprintf(w, "%s\n", diffHintColor.Sprint("--- expected"))
	_, _ = fmt.Fprintf(w, "%s\n", diffHintColor.Sprint("+++ actual"))

	// Hunk header — both headers are compared from their first line,
	// so the start offset is always 1.
	//
	//   @@ -1,13 +1,13 @@
	//       │  │  │  │
	//       │  │  │  └─ actual line count (newLen)
	//       │  │  └───── actual start line
	//       │  └─────────── expected line count (oldLen)
	//       └─────────────── expected start line
	//
	// oldLen: count of edits where typ != '+'
	//   = common lines ( ) + removed lines (-)
	//   — lines present in the expected header.
	// newLen: count of edits where typ != '-'
	//   = common lines ( ) + added lines (+)
	//   — lines present in the actual header.
	oldLen, newLen := 0, 0
	for _, e := range edits {
		if e.typ != '+' {
			oldLen++
		}
		if e.typ != '-' {
			newLen++
		}
	}
	_, _ = fmt.Fprintf(w, "%s\n", diffHunkColor.Sprintf("@@ -1,%d +1,%d @@", oldLen, newLen))

	// Edit lines
	for _, e := range edits {
		switch e.typ {
		case '-':
			_, _ = fmt.Fprintf(w, "%s\n", diffDelColor.Sprint("-"+e.line))
		case '+':
			_, _ = fmt.Fprintf(w, "%s\n", diffAddColor.Sprint("+"+e.line))
		default:
			_, _ = fmt.Fprintf(w, " %s\n", e.line)
		}
	}

	// Bottom separator
	_, _ = fmt.Fprintln(w, sep)
}

func splitLines(b []byte) []string {
	if len(b) == 0 {
		return nil
	}
	s := string(b)
	s, _ = strings.CutSuffix(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func lcs(a, b []string) [][]int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp
}

func backtrack(dp [][]int, a, b []string) []diffEdit {
	var edits []diffEdit
	i, j := len(a), len(b)
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			edits = append(edits, diffEdit{' ', a[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			edits = append(edits, diffEdit{'+', b[j-1]})
			j--
		} else {
			edits = append(edits, diffEdit{'-', a[i-1]})
			i--
		}
	}
	slices.Reverse(edits)
	return edits
}
