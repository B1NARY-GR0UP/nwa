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
	"os"
	"time"

	"github.com/B1NARY-GR0UP/nwa/internal"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

const (
	Name    = "nwa"
	Version = "v0.7.8"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: "A Simple Yet Powerful Tool for License Header Management",
	Long: `
++.     :*@@@@@@.#@@@@@@..@@@@@@*.      =@*.           
+@%:      =@@@@@.#@@@@*.  .*@@@@*.    :*@@@%:          
+@@@=      :%@@@.#@@%-      -%@@*.   =%@@@@@@*.        
+@@@@%:      =@@.#@*:        :*@*.  *@@@@@@@@@%:       
+@@@@@@=      :@.*=            =+.:@@@@@@@@@@@@@*     
`,
	Version: Version,
}

// Execute executes the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const (
	_useAdd    = "add"
	_useCheck  = "check"
	_useUpdate = "update"
	_useRemove = "remove"
	_useConfig = "config"
)

const (
	_modeCommon = "common"
	_modeConfig = "config"
)

const (
	_tmplLive   = "live"
	_tmplStatic = "static"
	_tmplRaw    = "raw"
)

const _levelMute = 12

func init() {
	rootCmd.SetVersionTemplate("{{ .Version }}")
	rootCmd.AddGroup(&cobra.Group{
		ID:    _modeCommon,
		Title: "Common Mode Commands:",
	}, &cobra.Group{
		ID:    _modeConfig,
		Title: "Config Mode Commands:",
	})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

type CommonFlags struct {
	// Basic Flags
	Holder  string
	Year    string
	License string
	SPDXIDs string
	Skip    []string

	// Advanced Flags
	Mute     bool
	DryRun   bool
	Verbose  bool
	Fuzzy    bool
	TmplType string
	Tmpl     string // template file path
	Keyword  []string
	Style    []string
}

var defaultCommonFlags = CommonFlags{
	Holder:   "<COPYRIGHT HOLDER>",
	Year:     fmt.Sprint(time.Now().Year()),
	License:  "apache",
	SPDXIDs:  "",
	Skip:     []string{},
	Mute:     false,
	DryRun:   false,
	Verbose:  false,
	Fuzzy:    false,
	TmplType: "",
	Tmpl:     "",
	Keyword:  []string{},
	Style:    []string{},
}

// ResetCommonFlags resets the common flags to their default values.
// This is needed for testing because cobra does not reset flags between
// consecutive calls to Execute.
func ResetCommonFlags() {
	defaultCommonFlags = CommonFlags{
		Holder:   "<COPYRIGHT HOLDER>",
		Year:     fmt.Sprint(time.Now().Year()),
		License:  "apache",
		SPDXIDs:  "",
		Skip:     []string{},
		Mute:     false,
		DryRun:   false,
		Verbose:  false,
		Fuzzy:    false,
		TmplType: "",
		Tmpl:     "",
		Keyword:  []string{},
		Style:    []string{},
	}
}

func setupCommonCmd(common *cobra.Command) {
	rootCmd.AddCommand(common)

	// basic
	common.Flags().StringVarP(&defaultCommonFlags.Holder, "copyright", "c", defaultCommonFlags.Holder, "copyright holder")
	common.Flags().StringVarP(&defaultCommonFlags.Year, "year", "y", defaultCommonFlags.Year, "copyright year")
	common.Flags().StringVarP(&defaultCommonFlags.License, "license", "l", defaultCommonFlags.License, "license type")
	common.Flags().StringVarP(&defaultCommonFlags.SPDXIDs, "spdxids", "i", defaultCommonFlags.SPDXIDs, "spdx ids")
	common.Flags().StringSliceVarP(&defaultCommonFlags.Skip, "skip", "s", defaultCommonFlags.Skip, "skip file path")

	// advanced
	common.Flags().BoolVarP(&defaultCommonFlags.Mute, "mute", "m", defaultCommonFlags.Mute, "mute mode")
	common.Flags().BoolVarP(&defaultCommonFlags.Verbose, "verbose", "V", defaultCommonFlags.Verbose, "verbose mode")
	common.Flags().BoolVarP(&defaultCommonFlags.Fuzzy, "fuzzy", "f", defaultCommonFlags.Fuzzy, "fuzzy matching")
	common.Flags().StringVarP(&defaultCommonFlags.TmplType, "tmpltype", "T", defaultCommonFlags.TmplType, "template type (live, static, raw)")
	common.Flags().StringVarP(&defaultCommonFlags.Tmpl, "tmpl", "t", defaultCommonFlags.Tmpl, "template file path")
	common.Flags().StringSliceVarP(&defaultCommonFlags.Keyword, "keyword", "k", defaultCommonFlags.Keyword, "keyword used to confirm the existence of license headers")
	common.Flags().StringSliceVarP(&defaultCommonFlags.Style, "style", "S", defaultCommonFlags.Style, "comment style `extension:style`, e.g. go:block")

	// flag rules
	common.MarkFlagsMutuallyExclusive("mute", "verbose")
	common.MarkFlagsRequiredTogether("tmpl", "tmpltype")
	common.MarkFlagsMutuallyExclusive("license", "tmpl")
	common.MarkFlagsMutuallyExclusive("license", "spdxids")
	common.MarkFlagsMutuallyExclusive("style", "tmpl")

	// for dry-run mode
	if common.Use != _useCheck {
		common.Flags().BoolVarP(&defaultCommonFlags.DryRun, "dry-run", "D", defaultCommonFlags.DryRun, "dry-run mode: print operations without modifying files")
		common.MarkFlagsMutuallyExclusive("dry-run", "mute")
		common.MarkFlagsMutuallyExclusive("dry-run", "verbose")
	}
}

func setupConfigCmd(config *cobra.Command) {
	rootCmd.AddCommand(config)

	config.Flags().StringVarP(&defaultConfigFlags.Command, "command", "c", defaultConfigFlags.Command, "command to execute")
	config.Flags().BoolVarP(&defaultConfigFlags.DryRun, "dry-run", "D", defaultConfigFlags.DryRun, "dry-run mode: print operations without modifying files")
}

func executeCommonCmd(_ *cobra.Command, args []string, flags CommonFlags, operation internal.Operation) {
	slog.SetLogLoggerLevel(slog.LevelWarn)
	if flags.Verbose {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
	// dry-run mode uses stdout to print infos, mute slog to avoid duplicate outputs
	if flags.Mute || flags.DryRun {
		slog.SetLogLoggerLevel(_levelMute)
	}

	// validate skip pattern
	for _, s := range flags.Skip {
		if !doublestar.ValidatePattern(s) {
			cobra.CheckErr(fmt.Errorf("--skip (-s) pattern %v is not valid", s))
		}
	}
	// validate path pattern
	for _, arg := range args {
		if !doublestar.ValidatePattern(arg) {
			cobra.CheckErr(fmt.Errorf("path pattern %v is not valid", arg))
		}
	}

	if flags.Tmpl == "" {
		tmpl, err := internal.MatchTmpl(flags.License, flags.SPDXIDs != "")
		if err != nil {
			cobra.CheckErr(err)
		}

		tmplData := &internal.TmplData{
			Holder:  flags.Holder,
			Year:    flags.Year,
			SPDXIDs: flags.SPDXIDs,
		}

		renderedTmpl, err := tmplData.RenderTmpl(tmpl)
		if err != nil {
			cobra.CheckErr(err)
		}

		internal.PrepareTasks(&internal.TaskParams{
			Paths:    args,
			Skips:    flags.Skip,
			Keywords: flags.Keyword,
			Styles:   flags.Style,
			Raw:      false,
			Fuzzy:    flags.Fuzzy,
			Tmpl:     renderedTmpl,
			Op:       operation,
			DryRun:   flags.DryRun,
		})
	} else {
		// use customize template
		content, err := os.ReadFile(flags.Tmpl)
		if err != nil {
			cobra.CheckErr(err)
		}

		params := &internal.TaskParams{
			Paths:    args,
			Skips:    flags.Skip,
			Keywords: flags.Keyword,
			Styles:   flags.Style,
			Fuzzy:    flags.Fuzzy,
			Op:       operation,
			DryRun:   flags.DryRun,
		}

		switch flags.TmplType {
		case _tmplLive:
			tmplData := &internal.TmplData{
				Holder:  flags.Holder,
				Year:    flags.Year,
				SPDXIDs: flags.SPDXIDs,
			}

			renderedTmpl, err := tmplData.RenderTmpl(string(content))
			if err != nil {
				cobra.CheckErr(err)
			}

			params.Tmpl = renderedTmpl
			params.Raw = false

			internal.PrepareTasks(params)
		case _tmplStatic:
			params.Tmpl = content
			params.Raw = false

			internal.PrepareTasks(params)
		case _tmplRaw:
			params.Tmpl = content
			params.Raw = true

			internal.PrepareTasks(params)
		default:
			cobra.CheckErr(fmt.Errorf("invalid template type: %v", flags.TmplType))
		}
	}

	// mute flag and dry-run flag used to handle SUMMARY output
	internal.ExecuteTasks(operation, flags.Mute, flags.DryRun)
}
