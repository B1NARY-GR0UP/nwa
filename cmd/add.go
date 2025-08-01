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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "add license headers to files",
	Long:    `Common Command | Add license headers to files`,
	Example: `nwa add --skip "**/*.py" --license apache --copyright Lorain "**/*.go"`,
	GroupID: _common,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeCommonCmd(cmd, args, defaultCommonFlags, internal.Add)
	},
}

func init() {
	setupCommonCmd(addCmd)
}
