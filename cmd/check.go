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

package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/B1NARY-GR0UP/nwa/util"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

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
		// validate skip pattern
		for _, s := range SkipF {
			if !doublestar.ValidatePattern(s) {
				cobra.CheckErr(fmt.Errorf("-skip pattern %v is not valid", s))
			}
		}
		// check if enable rawtmpl
		var rawTmpl bool
		if RawTmplF != "" && TmplF == "" {
			TmplF = RawTmplF
			rawTmpl = true
		}
		if TmplF == "" {
			tmpl, err := util.MatchTmpl(LicenseF, SPDXIDsF != "")
			if err != nil {
				cobra.CheckErr(err)
			}
			tmplData := &util.TmplData{
				Holder:  HolderF,
				Year:    YearF,
				SPDXIDs: SPDXIDsF,
			}
			renderedTmpl, err := tmplData.RenderTmpl(tmpl)
			if err != nil {
				cobra.CheckErr(err)
			}
			// determine files need to be added
			util.PrepareTasks(args, renderedTmpl, util.Check, SkipF, MuteF, rawTmpl)
		} else {
			content, err := os.ReadFile(TmplF)
			if err != nil {
				cobra.CheckErr(err)
			}
			buf := bytes.NewBuffer(content)
			if rawTmpl {
				_, _ = fmt.Fprintln(buf)
			}
			util.PrepareTasks(args, buf.Bytes(), util.Check, SkipF, MuteF, rawTmpl)
		}
		util.ExecuteTasks()
	},
}

func init() {
	setupCommonCmd(checkCmd)
}
