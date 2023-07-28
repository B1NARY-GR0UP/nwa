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
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "",
	Long:    ``,
	GroupID: common,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		// 校验路径参数，即 args 代表添加的路径
		// 查看是否设置了 tmpl，如果设置了则忽略 holder year license 参数
		// 查看 holder year license 参数，如果没有设置则使用默认值
		// 加载对应的 license 模板并使用 holder 和 year 参数渲染
		// 查看 skip 参数，和 args 一起决定需要进行修改的文件路径（列表）
		// 查看是否使用 mute 参数，没有启用则日志输出修改文件列表
		// 将文件修改任务添加到 chan 中
		// 使用 worker pool 消费任务
	},
}

func init() {
	setupCommonCmd(addCmd)
}
