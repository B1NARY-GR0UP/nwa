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
	"fmt"
	"log/slog"
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

NOTE: You can specify the path of the configuration file. 
If not specified, .nwa-config.yaml will be used as the default configuration file path. 
The behavior of NWA depends entirely on the configuration file. 
If some configurations are not set, the default configurations will be used.  

The cmd field value can also be set by --command (-c) flag.

Priority:
1. --command (-c) flag
2. value configured for cmd in the configuration file
3. default value of cmd in the configuration file (add)`,
	Example: `nwa config -c check sample-config.yaml

SAMPLE CONFIGURATION FILE (YAML):
nwa:
  holder: BINARY Members
  license: apache-2.0
  verbose: true
  fuzzy: true
  skip:
    - "testdata/**"
  path:
    - "**/*.go"`,
	GroupID: _config,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// read config
		var path string
		if len(args) != 0 {
			path = args[0]
		}
		if err := defaultConfig.readInConfig(path); err != nil {
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
		// mute has a higher priority
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

		if (defaultConfig.Nwa.TmplType == "") != (defaultConfig.Nwa.Tmpl == "") {
			cobra.CheckErr("tmpltype and tmpl must be set together")
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

			internal.PrepareTasks(&internal.TaskParams{
				Paths:    defaultConfig.Nwa.Path,
				Skips:    defaultConfig.Nwa.Skip,
				Keywords: defaultConfig.Nwa.Keyword,
				Styles:   defaultConfig.Nwa.Style,
				Raw:      false,
				Fuzzy:    defaultConfig.Nwa.Fuzzy,
				Tmpl:     renderedTmpl,
				Op:       operation,
			})
		} else {
			// use customize template

			params := &internal.TaskParams{
				Paths:    defaultConfig.Nwa.Path,
				Skips:    defaultConfig.Nwa.Skip,
				Keywords: defaultConfig.Nwa.Keyword,
				Styles:   defaultConfig.Nwa.Style,
				Fuzzy:    defaultConfig.Nwa.Fuzzy,
				Op:       operation,
			}

			switch defaultConfig.Nwa.TmplType {
			case _live:
				tmplData := &internal.TmplData{
					Holder:  defaultConfig.Nwa.Holder,
					Year:    defaultConfig.Nwa.Year,
					SPDXIDs: defaultConfig.Nwa.SPDXIDs,
				}

				renderedTmpl, err := tmplData.RenderTmpl(defaultConfig.Nwa.Tmpl)
				if err != nil {
					cobra.CheckErr(err)
				}

				params.Tmpl = renderedTmpl
				params.Raw = false

				internal.PrepareTasks(params)
			case _static:
				params.Tmpl = []byte(defaultConfig.Nwa.Tmpl)
				params.Raw = false

				internal.PrepareTasks(params)
			case _raw:
				params.Tmpl = []byte(defaultConfig.Nwa.Tmpl)
				params.Raw = true

				internal.PrepareTasks(params)
			default:
				cobra.CheckErr(fmt.Errorf("invalid template type: %v", defaultConfig.Nwa.TmplType))
			}
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
	// basic
	Cmd     string   `yaml:"cmd"`
	Holder  string   `yaml:"holder"`
	Year    string   `yaml:"year"`
	License string   `yaml:"license"`
	SPDXIDs string   `yaml:"spdxids"`
	Skip    []string `yaml:"skip"`
	Path    []string `yaml:"path"`

	// advanced
	Mute     bool     `yaml:"mute"`
	Verbose  bool     `yaml:"verbose"`
	Fuzzy    bool     `yaml:"fuzzy"`
	TmplType string   `yaml:"tmpltype"`
	Tmpl     string   `yaml:"tmpl"`
	Keyword  []string `yaml:"keyword"`
	Style    []string `yaml:"style"`
}

var defaultConfig = &Config{Nwa: NwaConfig{
	Cmd:      "add",
	Holder:   "<COPYRIGHT HOLDER>",
	Year:     fmt.Sprint(time.Now().Year()),
	License:  "apache",
	SPDXIDs:  "",
	Skip:     []string{},
	Path:     []string{},
	Mute:     false,
	Verbose:  false,
	Fuzzy:    false,
	TmplType: "",
	Tmpl:     "",
	Keyword:  []string{},
	Style:    []string{},
}}

func (cfg *Config) readInConfig(path string) error {
	if path == "" {
		// default configuration path: `./.nwa-config.yaml`
		viper.SetConfigName(".nwa-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	// will overwrite the default config if some fields are declared
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	return nil
}
