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
	"testing"
)

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
