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
	"reflect"
	"strings"
	"testing"
)

func TestExtractHeader(t *testing.T) {
	tests := []struct {
		name    string
		content string
		limit   int
		want    string
	}{
		{
			name:    "header separated by blank line",
			content: "// line1\n// line2\n\npackage main\n",
			limit:   10,
			want:    "// line1\n// line2",
		},
		{
			name:    "no blank line, within limit",
			content: "// line1\n// line2\npackage main\n",
			limit:   3,
			want:    "// line1\n// line2\npackage main",
		},
		{
			name:    "no blank line, exceeds limit",
			content: "// line1\n// line2\n// line3\npackage main\n",
			limit:   2,
			want:    "// line1\n// line2",
		},
		{
			name:    "no blank line, at exact limit",
			content: "// line1\n// line2\npackage main\n",
			limit:   2,
			want:    "// line1\n// line2",
		},
		{
			name:    "empty content",
			content: "",
			limit:   5,
			want:    "",
		},
		{
			name:    "blank line at start",
			content: "\n\n// header\n\npackage main\n",
			limit:   10,
			want:    "",
		},
		{
			name:    "content shorter than limit, no blank line",
			content: "only one line",
			limit:   10,
			want:    "only one line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractHeader([]byte(tt.content), tt.limit)
			if string(got) != tt.want {
				t.Errorf("extractHeader() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "multi-line with trailing newline",
			input: "a\nb\nc\n",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "multi-line without trailing newline",
			input: "a\nb\nc",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "single line",
			input: "hello",
			want:  []string{"hello"},
		},
		{
			name:  "single line with trailing newline",
			input: "hello\n",
			want:  []string{"hello"},
		},
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name:  "only newline",
			input: "\n",
			want:  nil,
		},
		{
			name:  "blank line between content",
			input: "a\n\nb\n",
			want:  []string{"a", "", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLines([]byte(tt.input))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitLines() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestLCS(t *testing.T) {
	tests := []struct {
		name string
		a, b []string
		want int // expected LCS length (last cell of DP table)
	}{
		{name: "both empty", a: nil, b: nil, want: 0},
		{name: "a empty", a: nil, b: []string{"x", "y"}, want: 0},
		{name: "b empty", a: []string{"x", "y"}, b: nil, want: 0},
		{name: "identical", a: []string{"a", "b", "c"}, b: []string{"a", "b", "c"}, want: 3},
		{name: "completely different", a: []string{"a", "b"}, b: []string{"x", "y"}, want: 0},
		{name: "single change in middle", a: []string{"a", "b", "c"}, b: []string{"a", "x", "c"}, want: 2},
		{name: "one added", a: []string{"a", "c"}, b: []string{"a", "b", "c"}, want: 2},
		{name: "one deleted", a: []string{"a", "b", "c"}, b: []string{"a", "c"}, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := lcs(tt.a, tt.b)
			got := dp[len(tt.a)][len(tt.b)]
			if got != tt.want {
				t.Errorf("lcs() length = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBacktrack(t *testing.T) {
	// Verify backtrack produces a valid edit script: applying it to 'a' yields 'b'
	a := []string{"a", "b", "c"}
	b := []string{"a", "x", "c"}
	dp := lcs(a, b)
	edits := backtrack(dp, a, b)

	var rebuilt []string
	for _, e := range edits {
		if e.typ != '-' {
			rebuilt = append(rebuilt, e.line)
		}
	}
	if !reflect.DeepEqual(rebuilt, b) {
		t.Errorf("applying edits to 'a' should yield 'b': got %#v, want %#v", rebuilt, b)
	}
}

func TestBacktrackIdentical(t *testing.T) {
	a := []string{"line1", "line2"}
	dp := lcs(a, a)
	edits := backtrack(dp, a, a)

	for _, e := range edits {
		if e.typ != ' ' {
			t.Errorf("identical lines should produce only ' ' edits, got %c", e.typ)
		}
	}
}

func TestBacktrackCompletelyDifferent(t *testing.T) {
	a := []string{"old1", "old2"}
	b := []string{"new1", "new2"}
	dp := lcs(a, b)
	edits := backtrack(dp, a, b)

	// All old lines should be deleted, all new lines added
	delCount, addCount := 0, 0
	for _, e := range edits {
		switch e.typ {
		case '-':
			delCount++
		case '+':
			addCount++
		}
	}
	if delCount != len(a) {
		t.Errorf("expected %d deletions, got %d", len(a), delCount)
	}
	if addCount != len(b) {
		t.Errorf("expected %d additions, got %d", len(b), addCount)
	}
}

func TestDiffReport(t *testing.T) {
	expected := []byte("// line1\n// line2\n")
	actual := []byte("// line1\n// lineX\n")

	var buf bytes.Buffer
	DiffReport(&buf, expected, actual)
	output := buf.String()

	// Must contain key elements
	checks := []string{
		"--- expected",
		"+++ actual",
		"@@ ",
		"-// line2",
		"+// lineX",
	}
	for _, c := range checks {
		if !strings.Contains(output, c) {
			t.Errorf("DiffReport output missing %q", c)
		}
	}
}

func TestDiffReportIdentical(t *testing.T) {
	data := []byte("// line1\n// line2\n")
	var buf bytes.Buffer
	DiffReport(&buf, data, data)
	output := buf.String()

	if strings.Contains(output, "-//") || strings.Contains(output, "+//") {
		t.Errorf("identical inputs should have no diffs: %s", output)
	}
}
