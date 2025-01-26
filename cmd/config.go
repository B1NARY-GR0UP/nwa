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

package cmd

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/B1NARY-GR0UP/nwa/internal"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "edit the files according to the configuration file",
	Long: `Config Command | Edit files according to the configuration file

EXAMPLE: nwa config -c check config.yaml

NOTE: This command only supports the --command (-c) flag;
You can only specify the path of the configuration file, and everything depends on the configuration file.
If some configuration are not configured, the default configuration will be used.

The command can also be set on the command line.
command priority:
1. --command (-c) flag
2. value configured for cmd in the configuration file
3. default value of cmd in the configuration file (add)

SAMPLE CONFIGURATION FILE (YAML):
nwa:
  holder: "RHINE LAB.LLC."
  year: "2077"
  license: "apache"
  path: ["server/**/*.go", "client/**/*.go", "pkg/**"]
  skip: ["**/*.py"]
`,
	GroupID: _config,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := defaultConfig.readInConfig(args[0]); err != nil {
			cobra.CheckErr(err)
		}

		// command priority:
		// 1. --command (-c) flag
		// 2. value configured for `cmd` in the configuration file
		// 3. default value of `cmd` in the configuration file (add)
		operation := internal.Operation(defaultConfig.Nwa.Cmd)
		if defaultConfigFlags.Command != "" {
			operation = internal.Operation(defaultConfigFlags.Command)
		}

		slog.SetLogLoggerLevel(slog.LevelWarn)
		if defaultConfig.Nwa.Verbose {
			slog.SetLogLoggerLevel(slog.LevelInfo)
		}
		// mute has higher priority
		if defaultConfig.Nwa.Mute {
			slog.SetLogLoggerLevel(_levelMute)
		}

		// validate skip pattern
		for _, s := range defaultConfig.Nwa.Skip {
			if !doublestar.ValidatePattern(s) {
				cobra.CheckErr(fmt.Errorf("-skip pattern %v is not valid", s))
			}
		}
		// validate path pattern
		for _, path := range defaultConfig.Nwa.Path {
			if !doublestar.ValidatePattern(path) {
				cobra.CheckErr(fmt.Errorf("path pattern %v is not valid", path))
			}
		}

		if defaultConfig.Nwa.Tmpl != "" && defaultConfig.Nwa.RawTmpl != "" {
			cobra.CheckErr("tmpl flag should not be used with rawtmpl flag")
		}
		// check if enable rawtmpl
		var rawTmpl bool
		if defaultConfig.Nwa.RawTmpl != "" {
			defaultConfig.Nwa.Tmpl = defaultConfig.Nwa.RawTmpl
			rawTmpl = true
		}
		if defaultConfig.Nwa.Tmpl == "" {
			tmpl, err := internal.MatchTmpl(defaultConfig.Nwa.License, defaultConfig.Nwa.SPDXIDs != "")
			if err != nil {
				cobra.CheckErr(err)
			}
			tmplData := &internal.TmplData{
				Holder:  defaultConfig.Nwa.Holder,
				Year:    defaultConfig.Nwa.Year,
				SPDXIDs: defaultConfig.Nwa.SPDXIDs,
			}
			renderedTmpl, err := tmplData.RenderTmpl(tmpl)
			if err != nil {
				cobra.CheckErr(err)
			}
			// determine files need to be added
			internal.PrepareTasks(defaultConfig.Nwa.Path, renderedTmpl, operation, defaultConfig.Nwa.Skip, rawTmpl, defaultConfig.Nwa.Fuzzy)
		} else {
			content, err := os.ReadFile(defaultConfig.Nwa.Tmpl)
			if err != nil {
				cobra.CheckErr(err)
			}
			buf := bytes.NewBuffer(content)
			if rawTmpl {
				_, _ = fmt.Fprintln(buf)
			}
			internal.PrepareTasks(defaultConfig.Nwa.Path, buf.Bytes(), operation, defaultConfig.Nwa.Skip, rawTmpl, defaultConfig.Nwa.Fuzzy)
		}
		internal.ExecuteTasks(operation, defaultConfig.Nwa.Mute)
	},
}

func init() {
	setupConfigCmd(configCmd)
}

type ConfigFlags struct {
	Command string
}

var defaultConfigFlags = ConfigFlags{
	Command: "", // empty if user not specified
}

type Config struct {
	Nwa NwaConfig `yaml:"nwa"`
}

type NwaConfig struct {
	Cmd     string   `yaml:"cmd"`
	Holder  string   `yaml:"holder"`
	Year    string   `yaml:"year"`
	License string   `yaml:"license"`
	Mute    bool     `yaml:"mute"`
	Verbose bool     `yaml:"verbose"`
	Fuzzy   bool     `yaml:"fuzzy"`
	Path    []string `yaml:"path"`
	Skip    []string `yaml:"skip"`
	SPDXIDs string   `yaml:"spdxids"`
	Tmpl    string   `yaml:"tmpl"`
	RawTmpl string   `yaml:"rawtmpl"`
}

var defaultConfig = &Config{Nwa: NwaConfig{
	Cmd:     "add",
	Holder:  "<COPYRIGHT HOLDER>",
	Year:    fmt.Sprint(time.Now().Year()),
	License: "apache",
	Mute:    false,
	Verbose: false,
	Fuzzy:   false,
	Path:    []string{},
	Skip:    []string{},
	SPDXIDs: "",
	Tmpl:    "",
	RawTmpl: "",
}}

func (cfg *Config) readInConfig(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	// will overwrite default config if some fields is declared
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	return nil
}
