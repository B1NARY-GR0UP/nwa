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
	"bytes"
	"fmt"
	"github.com/B1NARY-GR0UP/nwa/pkg"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
	"os"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "",
	Long:    ``,
	GroupID: pkg.Common,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// validate skip pattern
		for _, s := range SkipF {
			if !doublestar.ValidatePattern(s) {
				cobra.CheckErr(fmt.Errorf("-skip pattern %v is not valid", s))
			}
		}
		if TmplF == "" {
			tmpl, err := pkg.MatchTmpl(LicenseF)
			if err != nil {
				cobra.CheckErr(err)
			}
			tmplData := &pkg.TmplData{
				Holder: HolderF,
				Year:   YearF,
			}
			renderedTmpl, err := tmplData.RenderTmpl(tmpl)
			if err != nil {
				cobra.CheckErr(err)
			}
			// determine files need to be added
			pkg.PrepareTasks(args, renderedTmpl, pkg.Remove, SkipF, MuteF, TmplF)
		} else {
			content, err := os.ReadFile(TmplF)
			if err != nil {
				cobra.CheckErr(err)
			}
			// TODO: optimize, remove bytes.Buffer
			buf := bytes.NewBuffer(content)
			// add blank line at the end
			_, _ = fmt.Fprintln(buf)
			pkg.PrepareTasks(args, buf.Bytes(), pkg.Remove, SkipF, MuteF, TmplF)
		}
		pkg.ExecuteTasks()
	},
}

func init() {
	setupCommonCmd(removeCmd)
}
