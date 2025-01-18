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
	"github.com/B1NARY-GR0UP/nwa/internal"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update license headers of files",
	Long: `Common Command | Update license headers of files
EXAMPLE: nwa update -l mit -c Anmory "**/*.py"
NOTE: Update identifies the content before the first blank line as a license header;
If your file does not meet the requirements, please use remove + add`,
	GroupID: _common,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeCommonCmd(cmd, args, defaultCommonFlags, internal.Update)
	},
}

func init() {
	setupCommonCmd(updateCmd)
}
