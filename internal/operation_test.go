// Copyright 2025 BINARY Members
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
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSkipDirNotWalked(t *testing.T) {
	temp := t.TempDir()

	dataDir := filepath.Join(temp, "data")
	if err := os.MkdirAll(filepath.Join(dataDir, "postgres", "18", "docker"), 0755); err != nil {
		t.Fatalf("failed to create test dirs: %v", err)
	}

	// Simulate an unreadable volume directory (permission bits are a no-op on Windows)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(dataDir, 0000); err != nil {
			t.Fatalf("failed to chmod data dir: %v", err)
		}
		t.Cleanup(func() {
			_ = os.Chmod(dataDir, 0755)
		})
	}

	t.Chdir(temp)

	PrepareTasks(&TaskParams{
		Paths: []string{"**/*.py"},
		Skips: []string{"data/**"},
		Op:    OpAdd,
	})

	if counter.failed != 0 {
		t.Errorf("expected failed=0, got %d", counter.failed)
	}
	if counter.skipped != 1 {
		t.Errorf("expected skipped=1 (the data dir), got %d", counter.skipped)
	}
}

func TestIsSkip(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{"empty patterns", "main.py", nil, false},
		{"exact file match", "main.py", []string{"main.py"}, true},
		{"glob file match", "src/main.py", []string{"src/*.py"}, true},
		{"directory pattern matches directory itself", "data", []string{"data/**"}, true},
		{"directory pattern matches nested path", "data/postgres/18/docker", []string{"data/**"}, true},
		{"bare directory name matches directory", "vendor", []string{"vendor"}, true},
		{"bare directory name does not match child", "vendor/lib.go", []string{"vendor"}, false},
		{"second pattern matches", "src/main.py", []string{"docs/**", "src/**"}, true},
		{"no match", "src/main.go", []string{"vendor/**", "docs/**"}, false},
		{"wildcard suffix", "README.md", []string{"*.md"}, true},
		{"nested glob", "internal/operation.go", []string{"internal/*.go"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSkip(tt.path, tt.patterns); got != tt.want {
				t.Errorf("isSkip(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}

func TestPrepareTasksSkipsFile(t *testing.T) {
	temp := t.TempDir()

	files := map[string]string{
		"keep.py": "print('keep')\n",
		"skip.py": "print('skip')\n",
		"note.md": "# note\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(temp, name), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
	}

	t.Chdir(temp)

	PrepareTasks(&TaskParams{
		Paths: []string{"**/*.py"},
		Skips: []string{"skip.py", "**/*.md"},
		Op:    OpAdd,
	})

	if counter.failed != 0 {
		t.Errorf("expected failed=0, got %d", counter.failed)
	}
	// skip.py matches both the main pattern and a skip pattern.
	// note.md does not match the main pattern, so it is never evaluated for skipping.
	if counter.skipped != 1 {
		t.Errorf("expected skipped=1, got %d", counter.skipped)
	}
}

func TestIsRecursiveSkip(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{"recursive pattern matches directory", "data", []string{"data/**"}, true},
		{"recursive pattern matches nested path", "data/postgres/18/docker", []string{"data/**"}, true},
		{"non-recursive wildcard does not prune", "data/sub", []string{"data/*"}, false},
		{"bare directory name does not prune", "vendor", []string{"vendor"}, false},
		{"file pattern does not prune", "foo.py", []string{"*.py"}, false},
		{"nested recursive pattern", "node_modules", []string{"**/node_modules/**"}, true},
		{"trailing file segment is not recursive", "data", []string{"data/**/*"}, false},
		{"non-matching recursive pattern", "src", []string{"data/**"}, false},
		{"second pattern is recursive", "vendor", []string{"vendor/*", "vendor/**"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRecursiveSkip(tt.path, tt.patterns); got != tt.want {
				t.Errorf("isRecursiveSkip(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}

func TestSkipNonRecursiveWalkErrorNotFailed(t *testing.T) {
	temp := t.TempDir()

	subDir := filepath.Join(temp, "data", "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create test dirs: %v", err)
	}

	// Simulate an unreadable directory matched by a non-recursive skip pattern
	// (permission bits are a no-op on Windows)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(subDir, 0000); err != nil {
			t.Fatalf("failed to chmod sub dir: %v", err)
		}
		t.Cleanup(func() { _ = os.Chmod(subDir, 0755) })
	}

	t.Chdir(temp)

	PrepareTasks(&TaskParams{
		Paths: []string{"**/*.py"},
		Skips: []string{"data/*"},
		Op:    OpAdd,
	})

	if counter.failed != 0 {
		t.Errorf("expected failed=0 for a skipped path, got %d", counter.failed)
	}
	if counter.skipped != 0 {
		t.Errorf("expected skipped=0 (walk error is ignored, not counted), got %d", counter.skipped)
	}
}

func TestWalkErrorNotSkippedFails(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are a no-op on Windows")
	}

	temp := t.TempDir()

	subDir := filepath.Join(temp, "data", "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create test dirs: %v", err)
	}
	if err := os.Chmod(subDir, 0000); err != nil {
		t.Fatalf("failed to chmod sub dir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(subDir, 0755) })

	t.Chdir(temp)

	PrepareTasks(&TaskParams{
		Paths: []string{"**/*.py"},
		Skips: []string{"vendor/**"},
		Op:    OpAdd,
	})

	if counter.failed != 1 {
		t.Errorf("expected failed=1 for an unskipped walk error, got %d", counter.failed)
	}
	if counter.skipped != 0 {
		t.Errorf("expected skipped=0, got %d", counter.skipped)
	}
}

func TestRemoveYear(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Remove year from copyright notice",
			input:    []byte("Copyright 2023 Company Name"),
			expected: []byte("Copyright  Company Name"),
		},
		{
			name:     "Remove year from the middle of text",
			input:    []byte("Some text from 2023 with year"),
			expected: []byte("Some text from  with year"),
		},
		{
			name:     "Multiple years in text",
			input:    []byte("From 2022 to 2023 copyright"),
			expected: []byte("From  to 2023 copyright"),
		},
		{
			name:     "No year in text",
			input:    []byte("Text with no year"),
			expected: []byte("Text with no year"),
		},
		{
			name:     "Only year in text",
			input:    []byte("2023"),
			expected: []byte(""),
		},
		{
			name:     "Three-digit number should not be removed",
			input:    []byte("Year 123 not removed"),
			expected: []byte("Year 123 not removed"),
		},
		{
			name:     "Five-digit number should not be removed",
			input:    []byte("Number 12345 not removed"),
			expected: []byte("Number 12345 not removed"),
		},
		{
			name:     "Empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			name:     "Year at beginning of text",
			input:    []byte("2023 is the current year"),
			expected: []byte(" is the current year"),
		},
		{
			name:     "Year at end of text",
			input:    []byte("Copyright until 2023"),
			expected: []byte("Copyright until "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeYear(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("actual: %s, expected: %s", result, tt.expected)
			}
		})
	}
}

func TestDetectLineEnding(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Pure LF",
			input:    []byte("line1\nline2\nline3\n"),
			expected: []byte("\n"),
		},
		{
			name:     "Pure CRLF",
			input:    []byte("line1\r\nline2\r\nline3\r\n"),
			expected: []byte("\r\n"),
		},
		{
			name:     "Empty file",
			input:    []byte(""),
			expected: []byte("\n"),
		},
		{
			name:     "No line break",
			input:    []byte("single line"),
			expected: []byte("\n"),
		},
		{
			name:     "Mixed CRLF majority",
			input:    []byte("line1\r\nline2\r\nline3\n"),
			expected: []byte("\r\n"),
		},
		{
			name:     "Mixed LF majority",
			input:    []byte("line1\nline2\nline3\r\n"),
			expected: []byte("\n"),
		},
		{
			name:     "Single CRLF",
			input:    []byte("line1\r\n"),
			expected: []byte("\r\n"),
		},
		{
			name:     "Single LF",
			input:    []byte("line1\n"),
			expected: []byte("\n"),
		},
		{
			name:     "Equal CRLF and LF count defaults to CRLF",
			input:    []byte("line1\r\nline2\n"),
			expected: []byte("\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectLineEnding(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("actual: %q, expected: %q", result, tt.expected)
			}
		})
	}
}

func TestConvertLineEnding(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		ending   []byte
		expected []byte
	}{
		{
			name:     "LF to CRLF",
			input:    []byte("line1\nline2\n"),
			ending:   []byte("\r\n"),
			expected: []byte("line1\r\nline2\r\n"),
		},
		{
			name:     "CRLF to LF",
			input:    []byte("line1\r\nline2\r\n"),
			ending:   []byte("\n"),
			expected: []byte("line1\nline2\n"),
		},
		{
			name:     "LF to LF no-op",
			input:    []byte("line1\nline2\n"),
			ending:   []byte("\n"),
			expected: []byte("line1\nline2\n"),
		},
		{
			name:     "Mixed with CR to CRLF",
			input:    []byte("line1\r\nline2\rline3\n"),
			ending:   []byte("\r\n"),
			expected: []byte("line1\r\nline2\r\nline3\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertLineEnding(tt.input, tt.ending)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("actual: %q, expected: %q", result, tt.expected)
			}
		})
	}
}
