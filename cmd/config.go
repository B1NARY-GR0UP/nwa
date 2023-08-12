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
	"github.com/spf13/viper"
	"time"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "",
	Long:    ``,
	GroupID: config,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
		// 校验路径参数，即 args 代表配置文件的路径
		// 读取配置文件
		if err := defaultConfig.readInConfig(""); err != nil {
		}
		// 查看是否配置了 tmpl，如果配置了就忽略配置文件中 holder year license 的值
		// 查看 skip 和 path 一起决定需要进行修改的文件路径（列表）注意：path 为必须参数，如果没有返回错误
		// 查看是否使用 mute 参数，没有启用则日志输出修改文件列表
		// 将文件修改任务添加到 chan 中
		// 使用 worker pool 消费任务
	},
}

func init() {
	setupConfigCmd(configCmd)
}

type Config struct {
	Nwa NwaConfig `yaml:"nwa"`
}

type NwaConfig struct {
	Cmd     string   `yaml:"cmd"`
	Holder  string   `yaml:"holder"`
	Year    string   `yaml:"year"`
	License string   `yaml:"license"`
	Mute    bool     `yaml:"mute"`
	Path    []string `yaml:"path"`
	Skip    []string `yaml:"skip"`
	Tmpl    string   `yaml:"tmpl"`
}

var defaultConfig = &Config{Nwa: NwaConfig{
	Cmd:     "add",
	Holder:  "<COPYRIGHT HOLDER>",
	Year:    fmt.Sprint(time.Now().Year()),
	License: "apache",
	Mute:    false,
	Path:    []string{},
	Skip:    []string{},
	Tmpl:    "",
}}

func (cfg *Config) readInConfig(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	return nil
}
