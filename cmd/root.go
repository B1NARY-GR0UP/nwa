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
	"os"
	"time"

	"github.com/B1NARY-GR0UP/nwa/util"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

const (
	name    = "nwa"
	version = "v0.3.1"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   name,
	Short: "A More Powerful License Header Management Tool",
	Long: `
███╗   ██╗██╗    ██╗ █████╗ 
████╗  ██║██║    ██║██╔══██╗
██╔██╗ ██║██║ █╗ ██║███████║
██║╚██╗██║██║███╗██║██╔══██║
██║ ╚████║╚███╔███╔╝██║  ██║
╚═╝  ╚═══╝ ╚══╝╚══╝ ╚═╝  ╚═╝
`,
	Version: version,
}

// Execute executes the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{ .Version }}")
	rootCmd.AddGroup(&cobra.Group{
		ID:    util.Common,
		Title: "Common Mode Commands:",
	}, &cobra.Group{
		ID:    util.Config,
		Title: "Config Mode Commands:",
	})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

type CommonFlags struct {
	Mute    bool
	Holder  string
	Year    string
	License string
	Tmpl    string
	RawTmpl string
	Skip    []string
	SPDXIDs string
}

var defaultCommonFlags = CommonFlags{
	Mute:    false,
	Holder:  "<COPYRIGHT HOLDER>",
	Year:    fmt.Sprint(time.Now().Year()),
	License: "apache",
	Tmpl:    "",
	RawTmpl: "",
	Skip:    []string{},
	SPDXIDs: "",
}

func setupCommonCmd(common *cobra.Command) {
	rootCmd.AddCommand(common)

	common.Flags().BoolVarP(&defaultCommonFlags.Mute, "mute", "m", defaultCommonFlags.Mute, "mute mode")
	common.Flags().StringVarP(&defaultCommonFlags.Holder, "copyright", "c", defaultCommonFlags.Holder, "copyright holder")
	common.Flags().StringVarP(&defaultCommonFlags.Year, "year", "y", defaultCommonFlags.Year, "copyright year")
	common.Flags().StringVarP(&defaultCommonFlags.License, "license", "l", defaultCommonFlags.License, "license type")
	common.Flags().StringVarP(&defaultCommonFlags.Tmpl, "tmpl", "t", defaultCommonFlags.Tmpl, "template file path")
	common.Flags().StringVarP(&defaultCommonFlags.RawTmpl, "rawtmpl", "r", defaultCommonFlags.RawTmpl, "template file path (enable raw template)")
	common.Flags().StringSliceVarP(&defaultCommonFlags.Skip, "skip", "s", defaultCommonFlags.Skip, "skip file path")
	common.Flags().StringVarP(&defaultCommonFlags.SPDXIDs, "spdxids", "i", defaultCommonFlags.SPDXIDs, "spdx ids")

	common.MarkFlagsMutuallyExclusive("tmpl", "rawtmpl")
	// tmpl
	common.MarkFlagsMutuallyExclusive("copyright", "tmpl")
	common.MarkFlagsMutuallyExclusive("year", "tmpl")
	common.MarkFlagsMutuallyExclusive("license", "tmpl")
	// rawtmpl
	common.MarkFlagsMutuallyExclusive("copyright", "rawtmpl")
	common.MarkFlagsMutuallyExclusive("year", "rawtmpl")
	common.MarkFlagsMutuallyExclusive("license", "rawtmpl")

	// SPDX IDs
	common.MarkFlagsMutuallyExclusive("license", "spdxids")
	common.MarkFlagsMutuallyExclusive("tmpl", "spdxids")
	common.MarkFlagsMutuallyExclusive("rawtmpl", "spdxids")
}

func setupConfigCmd(config *cobra.Command) {
	rootCmd.AddCommand(config)
}

func executeCommonCmd(_ *cobra.Command, args []string, flags CommonFlags, operation util.Operation) {
	// validate skip pattern
	for _, s := range flags.Skip {
		if !doublestar.ValidatePattern(s) {
			cobra.CheckErr(fmt.Errorf("--skip (-s) pattern %v is not valid", s))
		}
	}
	// check if enable rawtmpl
	var rawTmpl bool
	if flags.RawTmpl != "" && flags.Tmpl == "" {
		flags.Tmpl = flags.RawTmpl
		rawTmpl = true
	}
	if flags.Tmpl == "" {
		tmpl, err := util.MatchTmpl(flags.License, flags.SPDXIDs != "")
		if err != nil {
			cobra.CheckErr(err)
		}
		tmplData := &util.TmplData{
			Holder:  flags.Holder,
			Year:    flags.Year,
			SPDXIDs: flags.SPDXIDs,
		}
		renderedTmpl, err := tmplData.RenderTmpl(tmpl)
		if err != nil {
			cobra.CheckErr(err)
		}
		util.PrepareTasks(args, renderedTmpl, operation, flags.Skip, flags.Mute, rawTmpl)
	} else {
		content, err := os.ReadFile(flags.Tmpl)
		if err != nil {
			cobra.CheckErr(err)
		}
		buf := bytes.NewBuffer(content)
		if rawTmpl {
			_, _ = fmt.Fprintln(buf)
		}
		util.PrepareTasks(args, buf.Bytes(), operation, flags.Skip, flags.Mute, rawTmpl)
	}
	util.ExecuteTasks(operation, flags.Mute)
}
