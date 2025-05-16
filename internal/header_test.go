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
	"testing"
)

func TestHasHeader(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "empty content",
			content:  []byte(""),
			expected: false,
		},
		{
			name:     "copyright word present",
			content:  []byte("This file has Copyright 2023 notice"),
			expected: true,
		},
		{
			name:     "COPYRIGHT uppercase present",
			content:  []byte("THIS FILE HAS COPYRIGHT 2023 NOTICE"),
			expected: true,
		},
		{
			name:     "copyright symbol present",
			content:  []byte("This file has © 2023 notice"),
			expected: true,
		},
		{
			name:     "copr. abbreviation present",
			content:  []byte("This file has Copr. 2023 notice"),
			expected: true,
		},
		{
			name:     "mozilla public license present",
			content:  []byte("Mozilla Public License Version 2.0"),
			expected: true,
		},
		{
			name:     "spdx identifier present",
			content:  []byte("SPDX-License-Identifier: Apache-2.0"),
			expected: true,
		},
		{
			name:     "copyright beyond first 1000 bytes",
			content:  []byte(string(make([]byte, 1100)) + "Copyright 2023"),
			expected: false,
		},
		{
			name:     "mixed case with copyright",
			content:  []byte("this File has cOpYriGht 2023 notice"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := hasHeader(tt.content, []string{})
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestHasHeaderWithKeywords(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		keywords []string
		expected bool
	}{
		{
			name:     "empty content with keywords",
			content:  []byte(""),
			keywords: []string{"license", "copyright"},
			expected: false,
		},
		{
			name:     "content with matching keyword",
			content:  []byte("This file has a license notice"),
			keywords: []string{"license", "copyright"},
			expected: true,
		},
		{
			name:     "content with matching uppercase keyword",
			content:  []byte("This file has a LICENSE notice"),
			keywords: []string{"license"},
			expected: true,
		},
		{
			name:     "content with matching mixed case keyword",
			content:  []byte("This file has a LiCeNsE notice"),
			keywords: []string{"license"},
			expected: true,
		},
		{
			name:     "content without matching keywords",
			content:  []byte("This file has no matching words"),
			keywords: []string{"license", "copyright"},
			expected: false,
		},
		{
			name:     "keyword beyond first 1000 bytes",
			content:  []byte(string(make([]byte, 1100)) + "license"),
			keywords: []string{"license"},
			expected: false,
		},
		{
			name:     "multiple keywords with one match",
			content:  []byte("This file has an apache notice"),
			keywords: []string{"mit", "apache", "gpl"},
			expected: true,
		},
		{
			name:     "empty keywords list falls back to default checks",
			content:  []byte("This file has Copyright 2023 notice"),
			keywords: []string{},
			expected: true,
		},
		{
			name:     "empty keywords list with no default matches",
			content:  []byte("This file has no info"),
			keywords: []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := hasHeader(tt.content, tt.keywords)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
