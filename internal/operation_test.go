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
