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
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func generateHeader(path string, tmpl *bytes.Buffer) []byte {
	var header []byte
	base := strings.ToLower(filepath.Base(path))
	ext := base
	// use file name if there is no extension
	if e := filepath.Ext(base); e != "" {
		ext = e
	}
	switch ext {
	case ".c", ".h", ".gv", ".java", ".scala", ".kt", ".kts":
		header = doGenerate(tmpl, "/*", " * ", " */")
	case ".js", ".mjs", ".cjs", ".jsx", ".tsx", ".css", ".scss", ".sass", ".ts":
		header = doGenerate(tmpl, "/**", " * ", " */")
	case ".cc", ".cpp", ".cs", ".go", ".hcl", ".hh", ".hpp", ".m", ".mm", ".proto", ".rs", ".swift", ".dart", ".groovy", ".v", ".sv":
		header = doGenerate(tmpl, "", "// ", "//")
	case ".py", ".sh", ".yaml", ".yml", ".dockerfile", "dockerfile", ".rb", "gemfile", ".tcl", ".tf", ".bzl", ".pl", ".pp", "build", ".build", ".toml":
		header = doGenerate(tmpl, "", "# ", "")
	case ".el", ".lisp":
		header = doGenerate(tmpl, "", ";; ", "")
	case ".erl":
		header = doGenerate(tmpl, "", "% ", "")
	case ".hs", ".sql", ".sdl":
		header = doGenerate(tmpl, "", "-- ", "")
	case ".html", ".xml", ".vue", ".wxi", ".wxl", ".wxs":
		header = doGenerate(tmpl, "<!--", " ", "-->")
	case ".php":
		header = doGenerate(tmpl, "", "// ", "")
	case ".ml", ".mli", ".mll", ".mly":
		header = doGenerate(tmpl, "(**", "   ", "*)")
	default:
		if base == "cmakelists.txt" || strings.HasSuffix(base, ".cmake.in") || strings.HasSuffix(base, ".cmake") {
			header = doGenerate(tmpl, "", "# ", "")
		}
	}
	return header
}

// TODO: test needed
func doGenerate(tmpl *bytes.Buffer, top, mid, bot string) []byte {
	buf := bytes.NewBuffer(nil)
	if top != "" {
		_, _ = fmt.Fprintln(buf, top)
	}
	s := bufio.NewScanner(tmpl)
	for s.Scan() {
		_, _ = fmt.Fprintln(buf, strings.TrimRightFunc(mid+s.Text(), unicode.IsSpace))
	}
	if bot != "" {
		_, _ = fmt.Fprintln(buf, bot)
	}
	_, _ = fmt.Fprintln(buf)
	return buf.Bytes()
}

func isGenerated(b []byte) bool {
	// go generate: ^// Code generated .* DO NOT EDIT\.$
	goGenerated := regexp.MustCompile(`(?m)^.{1,2} Code generated .* DO NOT EDIT\.$`)
	// cargo raze: ^DO NOT EDIT! Replaced on runs of cargo-raze$
	cargoRazeGenerated := regexp.MustCompile(`(?m)^DO NOT EDIT! Replaced on runs of cargo-raze$`)
	return goGenerated.Match(b) || cargoRazeGenerated.Match(b)
}

func hasHeader(b []byte) bool {
	n := 1000
	if len(b) < n {
		n = len(b)
	}
	return bytes.Contains(bytes.ToLower(b[:n]), []byte("copyright")) ||
		bytes.Contains(bytes.ToLower(b[:n]), []byte("mozilla public")) ||
		bytes.Contains(bytes.ToLower(b[:n]), []byte("spdx-license-identifier"))
}

func matchFirstLine(b []byte) []byte {
	var line []byte
	fls := []string{
		"#!",                       // shell script
		"<?xml",                    // XML declaratioon
		"<!doctype",                // HTML doctype
		"# encoding:",              // Ruby encoding
		"# frozen_string_literal:", // Ruby interpreter instruction
		"<?php",                    // PHP opening tag
		"# escape",                 // Dockerfile directive https://docs.docker.com/engine/reference/builder/#parser-directives
		"# syntax",                 // Dockerfile directive https://docs.docker.com/engine/reference/builder/#parser-directives
	}
	// get headline of the file
	for _, c := range b {
		line = append(line, c)
		if c == '\n' {
			break
		}
	}
	first := strings.ToLower(string(line))
	// match first line
	for _, fl := range fls {
		if strings.HasPrefix(first, fl) {
			return line
		}
	}
	return nil
}

func assemble(line, header, content []byte) []byte {
	if line != nil {
		// get content exclude the first line
		content = content[len(line):]
		// add \n if the first line do not end with \n
		if line[len(line)-1] != '\n' {
			line = append(line, '\n')
		}
		header = append(line, header...)
	}
	header = append(header, content...)
	return header
}
