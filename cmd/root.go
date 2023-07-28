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

package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	Name    = "NWA"
	version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   Name,
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
	rootCmd.SetVersionTemplate("{{.Version}}")
	rootCmd.AddGroup(&cobra.Group{
		ID:    common,
		Title: "Common Mode Commands:",
	}, &cobra.Group{
		ID:    config,
		Title: "Config Mode Commands:",
	})
}

const (
	common = "common"
	config = "config"
)

var (
	mute    bool
	holder  string
	year    int
	license string
	skip    string
	tmpl    string
)

func setupCommonCmd(common *cobra.Command) {
	rootCmd.AddCommand(common)

	common.Flags().BoolVarP(&mute, "mute", "m", false, "mute mode")
	common.Flags().StringVarP(&holder, "copyright", "c", "[copyright holder]", "copyright holder")
	common.Flags().IntVarP(&year, "year", "y", time.Now().Year(), "copyright year")
	common.Flags().StringVarP(&license, "license", "l", "apache", "license type")
	// Note: use spaces to separate sections and use double quotation to enclose all the sections
	common.Flags().StringVarP(&skip, "skip", "s", "", "skip file")
	common.Flags().StringVarP(&tmpl, "tmpl", "t", "", "template file path")

	common.MarkFlagsMutuallyExclusive("copyright", "tmpl")
	common.MarkFlagsMutuallyExclusive("year", "tmpl")
	common.MarkFlagsMutuallyExclusive("license", "tmpl")
	common.MarkFlagsMutuallyExclusive("ignore", "tmpl")
}

func setupConfigCmd(config *cobra.Command) {
	rootCmd.AddCommand(config)
}
