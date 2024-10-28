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
	"github.com/B1NARY-GR0UP/nwa/util"
	"github.com/spf13/cobra"
)

// TODO: support fuzzy matching (e.g. ignore year, only check copyright and license type)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check license headers of files",
	Long: `Common Command | Check license headers of files
EXAMPLE: nwa check -t tmpl.txt .
NOTE: Do not use --mute (-m) flag with the command`,
	GroupID: util.Common,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeCommonCmd(cmd, args, defaultCommonFlags, util.Check)
	},
}

func init() {
	setupCommonCmd(checkCmd)
}
