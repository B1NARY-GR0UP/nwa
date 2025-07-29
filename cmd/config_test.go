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

package cmd

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	err := defaultConfig.readInConfig("../testdata/config/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(defaultConfig)
}

func TestMultilineTmpl(t *testing.T) {
	cfg := &Config{Nwa: NwaConfig{}}

	err := cfg.readInConfig("../testdata/tmpl/.static-tmpl-test.yaml")
	if err != nil {
		t.Fatalf("read config file failed: %v", err)
	}

	expectedTmpl := `Copyright 2025 BINARY Members

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

	if cfg.Nwa.Tmpl != expectedTmpl {
		t.Errorf("template field content doesn't match\nExpected:\n%s\n\nActual:\n%s", expectedTmpl, cfg.Nwa.Tmpl)
	}
}
