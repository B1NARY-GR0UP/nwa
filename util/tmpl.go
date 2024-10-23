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
//

package util

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
)

var spdxids bool

var (
	errLicenseNotSupported = errors.New("license not supported, please use custom tmpl with --tmpl or -t flag")
)

type TmplData struct {
	Holder  string
	Year    string
	SPDXIDs string
}

func MatchTmpl(license string, useSPDXIDs bool) (string, error) {
	if useSPDXIDs {
		spdxids = true
		return tmplSPDXIDs, nil
	}
	license = strings.ToLower(license)
	switch license {
	case "apache", "apache-2.0", "apache-2", "apache20", "apache 2.0", "apache2.0":
		return tmplApache, nil
	case "mit":
		return tmplMIT, nil
	case "gpl-3.0-or-later":
		return tmplGPLThreeOrLater, nil
	case "gpl-3.0-only":
		return tmplGPLThreeOnly, nil
	default:
		return "", errLicenseNotSupported
	}
}

func (t *TmplData) RenderTmpl(tmpl string) ([]byte, error) {
	if spdxids {
		// if holder is not declared when using spdxids then just don't generate the copyright line
		if t.Holder == "<COPYRIGHT HOLDER>" {
			t.Holder = ""
		}
	}
	buf := bytes.NewBuffer(nil)
	renderedTmpl := template.Must(template.New("nwa-tmpl").Parse(tmpl))
	err := renderedTmpl.Execute(buf, t)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TODO: support more templates
const tmplApache = `Copyright {{ .Year }} {{ .Holder }}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

const tmplMIT = `Copyright {{ .Year }} {{ .Holder }}

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the
Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.`

const tmplGPLThreeOrLater = `Copyright (C) {{ .Year }} {{ .Holder }}

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.`

const tmplGPLThreeOnly = `Copyright (C) {{ .Year }} {{ .Holder }}

This program is free software: you can redistribute it and/or modify 
it under the terms of the GNU General Public License as published by 
the Free Software Foundation, version 3.

This program is distributed in the hope that it will be useful, 
but WITHOUT ANY WARRANTY; without even the implied warranty of 
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the 
GNU General Public License for more details.

You should have received a copy of the GNU General Public License 
along with this program. If not, see <https://www.gnu.org/licenses/>.`

const tmplAGPLThreeOrLater = `Copyright (C) {{ .Year }} {{ .Holder }}

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`

const tmplAGPLThreeOnly = `Copyright (C) {{ .Year }} {{ .Holder }}

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, version 3.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`

const tmplSPDXIDs = `{{ if .Holder }}Copyright{{ if .Year }} {{ .Year }}{{ end }} {{ .Holder }}
{{ end }}SPDX-License-Identifier: {{ .SPDXIDs }}`
