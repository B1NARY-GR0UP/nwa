![NWA](images/NWA.png)

一款更强大的许可声明管理工具

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/nwa)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/nwa)

[English](README.md) | 中文

## 安装

### 通过 Go 安装（如果您有 Go 环境）

```shell
go install github.com/B1NARY-GR0UP/nwa@latest
```

执行以下命令以验证是否安装成功：

```shell
nwa --version
```
### 通过 Docker 安装（如果您没有 Go 环境）

1. **安装**

直接获取 NWA 的 docker 镜像：

```shell
docker pull ghcr.io/b1nary-gr0up/nwa:main
```
或者自己从源码构建镜像，示例：

```shell
docker build -t ghcr.io/b1nary-gr0up/nwa:main .
```
2. **验证是否安装成功**

```shell
docker run -it ghcr.io/b1nary-gr0up/nwa:main --version
```

3. **将希望让 NWA 进行操作的目录挂载到 `/src` 目录下，以下是一个示例，您可以通过阅读 [使用](#使用) 章节了解 NWA 的用法后再来体验：**

```shell
docker run -it -v ${PWD}:/src ghcr.io/b1nary-gr0up/nwa:main add -c "RHINE LAB.LLC." -y 2077 .
```
## 使用

为了帮助您更好的上手 NWA 的功能使用，开发者提供了 [**nwa-examples**](https://github.com/rainiring/nwa-examples) 来帮助您以实践的方式学习 NWA 的使用。您可以 clone 这个项目并结合这个教程一起食用或者也可以直接查看本教程，并通过 nwa-example 来获取更多的使用示例。

NWA 是一个基于 [cobra](https://github.com/spf13/cobra) 搭建的命令行工具，以下是 NWA 的命令总览：

```shell
Usage:         
  nwa [command]

Common Mode Commands:
  add         add license headers to files
  check       check license headers of files
  remove      remove license headers of files
  update      update license headers of files

Config Mode Commands:
  config      edit the files according to the configuration file

Additional Commands:
  help        Help about any command

Flags:
  -h, --help      help for nwa
  -v, --version   version for nwa

Use "nwa [command] --help" for more information about a command.
```

**命令列表**

- **[Add](#add---为文件添加许可声明)**: 为文件添加许可声明；
- **[Remove](#remove---移除文件的许可声明)**: 移除文件的许可声明；
- **[Update](#update---更新文件的许可声明)**: 更新文件的许可声明；
- **[Check](#check---检查文件的许可声明)**: 检查文件的许可声明；
- **[Config](#config---以配置文件的方式编辑文件的许可声明)**: 以配置文件的方式编辑文件的许可声明；

### Add - 为文件添加许可声明

#### 使用

```shell
nwa add [flags] path...
```
通过在 `nwa` 后面指定 `add` 命令来执行添加许可声明的操作，`[flags]` 表示可以添加 0 个或多个可选的 flag，`path...` 表示您可以指定一个或多个要添加许可声明的文件的路径。

#### 使用示例

```shell
nwa add -l apache -c "RHINE LAB.LLC." -y 2077 ./server ./utils/bufferpool
```

以上示例的命令表示为相对路径为 `./server` 和 `./utils/bufferpool` 文件夹下的所有文件**添加**

- 许可证类型： `Apache 2.0`
- 版权归属（copyright holder）： `RHINE LAB.LLC.`
- 版权年份： `2077`

的许可声明。

并且 NWA 会根据指定路径下的文件类型生成对应的许可声明，比如 `.py` 文件使用 `#` 注释，`.go` 文件使用 `//` 注释；

如果您的文件不在 NWA 内置的文件类型中，您可以：

- 通过 `-t`（`--tmpl`）flag 指定模板文件；
- 为 NWA 提一个 issue 或者 PR；

NWA 也会以输出日志的形式提示您如果有文件已经具有了许可声明或者是不被允许编辑的文件（例如一些工具生成的代码文件）。

#### Flags

`add` 命令目前支持的 flags 如下表所示：

| 短标志 | 长标志         | 默认值                        | 介绍                |
|-----|-------------|----------------------------|-------------------|
| -c  | --copyright | `<COPYRIGHT HOLDER>`       | 版权归属              |
| -l  | --license   | `apache`                   | 许可证类型             |
| -i  | --spdxids   | `""`                       | SPDX IDs          |
| -m  | --mute      | `false` (不指定)              | 静默模式              |
| -s  | --skip      | `[]`                       | 跳过的文件路径           |
| -t  | --tmpl      | `""`                       | 模板文件路径            |
| -r  | --rawtmpl   | `""`                       | 模板文件路径（启用 raw 模式） |
| -y  | --year      | `time.Now().Year()` (当前年份) | 版权年份              |
| -h  | --help      | 无                          | 帮助                |

大部分的 flag 都很容易理解，值得一提的是 `-s` （`-skip`）和 `-t`（`-tmpl`）;

- `-s` 的默认值是一个空数组，所以您可以使用 0 个或多个 `-s` 来指定要跳过的文件路径，并且这个路径支持 [doublestar](https://github.com/bmatcuk/doublestar) 语法，所以非常灵活，比如可以这样使用：

```shell
nwa add -s **.py -s /example/**/*.txt -c Lorain .
```

- `-t` 允许您指定一个模板文件的路径，NWA 会读取这个模板文件的内容并为要进行添加许可声明的文件添加上模板中所定义的内容；

**注意：如果您指定了模板，则 NWA 不会根据文件类型生成对应的许可声明，而是对所有路径下的文件添加模板中的内容，指定了模板也会让 NWA 忽略通过 `-c`，`-l` 和 `-y` 指定的配置。**

一个模板文件的示例如下：

- `example-tmpl.txt`

```text
Copyright 2077 RHINE LAB.LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

更多的使用示例请您参考 [nwa-examples](https://github.com/rainiring/nwa-examples)。

### Remove - 移除文件的许可声明

#### 使用

```shell
nwa remove [flags] path...
```

通过在 `nwa` 后面指定 `remove` 命令来执行添加许可声明的操作，`[flags]` 表示可以添加 0 个或多个可选的 flag，`path...` 表示您可以指定一个或多个要移除许可声明的文件的路径。

#### 使用示例

```shell
nwa remove -l mit -c "RHINE LAB.LLC." -s **.py pkg
```

以上示例的命令表示为相对路径为 `pkg` 文件夹下的所有 **非 Python** 文件 **移除**

- 许可证类型： `MIT`
- 版权归属（copyright holder）： `RHINE LAB.LLC.`
- 版权年份： `2023`

的许可声明。

如果路径下的文件没有许可声明或者是不可进行编辑的文件，NWA 会以输出日志的形式提示您。

#### Flags

`remove` 命令目前支持的 flags 如下表所示：

| 短标志 | 长标志         | 默认值                        | 介绍                |
|-----|-------------|----------------------------|-------------------|
| -c  | --copyright | `<COPYRIGHT HOLDER>`       | 版权归属              |
| -l  | --license   | `apache`                   | 许可证类型             |
| -i  | --spdxids   | `""`                       | SPDX IDs          |
| -m  | --mute      | `false` (不指定)              | 静默模式              |
| -s  | --skip      | `[]`                       | 跳过的文件路径           |
| -t  | --tmpl      | `""`                       | 模板文件路径            |
| -r  | --rawtmpl   | `""`                       | 模板文件路径（启用 raw 模式） |
| -y  | --year      | `time.Now().Year()` (当前年份) | 版权年份              |
| -h  | --help      | 无                          | 帮助                |

**注意：如果您指定了模板，则 NWA 不会根据文件类型来移除对应的许可声明，而是移除路径下所有符合模板文件内容的许可声明。**

更多的使用示例请您参考 [nwa-examples](https://github.com/rainiring/nwa-examples)。

### Update - 更新文件的许可声明

#### 使用

```shell
nwa update [flags] path...
```

通过在 `nwa` 后面指定 `update` 命令来执行更新许可声明的操作，`[flags]` 表示可以添加 0 个或多个可选的 flag，`path...` 表示您可以指定一个或多个要更新许可声明的文件的路径。

#### 使用示例

```shell
nwa update -l apache -c "BINARY Members" .
```

以上示例的命令表示将当前路径下的所有文件的许可声明 **更新** 为

- 许可证类型： `Apache 2.0`
- 版权归属（copyright holder）： `BINARY Members`
- 版权年份： `2023`

而不管之前的版权声明是什么。

**注意：Update 还是 Remove + Add**

`update` 命令的执行逻辑为将文件中**第一个空行之前**的内容识别为许可声明（如果有 shebang 则会保留），`update` 会根据指定的 flags 生成新的许可声明并进行替换。

如果您的文件第一个空行之前的内容包含了除了许可声明之外的内容，我建议您使用 `remove` + `add` 的形式来更新许可声明，这可以达到相同的效果，并且在 `remove` 中和模板文件 flag 一起使用可以达到更好的移除效果。

#### Flags

`update` 命令目前支持的 flags 如下表所示：

| 短标志 | 长标志         | 默认值                        | 介绍                |
|-----|-------------|----------------------------|-------------------|
| -c  | --copyright | `<COPYRIGHT HOLDER>`       | 版权归属              |
| -l  | --license   | `apache`                   | 许可证类型             |
| -i  | --spdxids   | `""`                       | SPDX IDs          |
| -m  | --mute      | `false` (不指定)              | 静默模式              |
| -s  | --skip      | `[]`                       | 跳过的文件路径           |
| -t  | --tmpl      | `""`                       | 模板文件路径            |
| -r  | --rawtmpl   | `""`                       | 模板文件路径（启用 raw 模式） |
| -y  | --year      | `time.Now().Year()` (当前年份) | 版权年份              |
| -h  | --help      | 无                          | 帮助                |

**注意：如果您指定了模板，则 NWA 不会根据文件类型来更新对应的许可声明，而是更新路径下所有文件第一个空行前的内容为模板文件中指定的内容。**

更多的使用示例请您参考 [nwa-examples](https://github.com/rainiring/nwa-examples)。

### Check - 检查文件的许可声明

#### 使用

```shell
nwa check [flags] path...
```

通过在 `nwa` 后面指定 `check` 命令来执行检查许可声明的操作，`[flags]` 表示可以添加 0 个或多个可选的 flag，`path...` 表示您可以指定一个或多个要检查许可声明的文件的路径。

#### 使用示例

```shell
nwa check --tmpl tmpl.txt ./client
```

以上示例的命令表示 **检查** 相对路径为 `./client` 文件夹下的所有文件许可声明是否与 `tmpl.txt` 模板文件中指定的内容相符。

检查完毕后，NWA 会以输出日志的形式来告知您检查的结果，一个样例输出如下：

```text
time="2023-08-25T23:01:49+08:00" level=info msg="skip file" path=README.md
time="2023-08-25T23:01:49+08:00" level=info msg="file has a matched header" path="dirA\\fileA.go"                
time="2023-08-25T23:01:49+08:00" level=info msg="file does not have a matched header" path="dirB\\dirC\\fileC.go"
time="2023-08-25T23:01:49+08:00" level=info msg="file has a matched header" path="dirB\\fileB.go"                
time="2023-08-25T23:01:49+08:00" level=info msg="file does not have a matched header" path=main.go 
```

#### Flags

`check` 命令目前支持的 flags 如下表所示：

| 短标志 | 长标志         | 默认值                        | 介绍                |
|-----|-------------|----------------------------|-------------------|
| -c  | --copyright | `<COPYRIGHT HOLDER>`       | 版权归属              |
| -l  | --license   | `apache`                   | 许可证类型             |
| -i  | --spdxids   | `""`                       | SPDX IDs          |
| -m  | --mute      | `false` (不指定)              | 静默模式              |
| -s  | --skip      | `[]`                       | 跳过的文件路径           |
| -t  | --tmpl      | `""`                       | 模板文件路径            |
| -r  | --rawtmpl   | `""`                       | 模板文件路径（启用 raw 模式） |
| -y  | --year      | `time.Now().Year()` (当前年份) | 版权年份              |
| -h  | --help      | 无                          | 帮助                |

**注意：在使用 `check` 命令时不要与 `-m` flag 一起使用，因为这会让 NWA 无法输出检查的结果。**

更多的使用示例请您参考 [nwa-examples](https://github.com/rainiring/nwa-examples)。

### Config - 以配置文件的方式编辑文件的许可声明

#### 使用

```shell
nwa config [flags] path
```

通过在 `nwa` 后面指定 `config` 命令来执行以配置文件的方式编辑文件的许可声明的操作，`[flags]` 表示可以添加 0 个或多个可选的 flag，`path` 表示您 **必须指定** 配置文件所在的路径。

#### 使用示例

```shell
nwa config config.yaml
```

以上的示例命令表示读取 `config.yaml` 配置文件并根据其内容来编辑指定文件的许可声明，配置文件的结构将在下面提供。

#### Flags

由于所有的配置都位于配置文件，所以只有一个帮助 flag。

| 短标志 | 长标志    | 默认值 | 介绍 |
|-----|--------|-----|----|
| -h  | --help | 无   | 帮助 |

#### 配置文件

一个完整的配置文件示例如下所示：

```yaml
nwa:
  cmd: "add"                        # 默认值: "add" 可选值: "add", "check", "remove", "update" 
  holder: "RHINE LAB.LLC."          # 默认值: "<COPYRIGHT HOLDER>"
  year: "2077"                      # 默认值: Current Year
  license: "apache"                 # 默认值: "apache"
  spdxids: ""                       # 默认值: ""
  mute: false                       # 默认值: false (unspecified)
  path: ["server", "client", "pkg"] # 默认值: []
  skip: ["**.py"]                   # 默认值: []
  tmpl: "nwa.txt"                   # 默认值: ""                                                       
  rawtmpl: ""                       # 默认值: ""                                                       
```

如果您不指定一些配置文件字段，NWA 将会使用默认值，比如如果您不指定许可证类型，则会默认使用 Apache 2.0 协议。

**注意：如果您设置了 `tmpl` 字段或者 `rawtmpl`，那么 `holder`，`year`，`license` 和 `spdxids` 字段将被 NWA 忽略。**

更多的使用示例请您参考 [nwa-examples](https://github.com/rainiring/nwa-examples)。

## 支持的协议模板

| 协议                | 选项 (忽略大小写)                                                    |
|-------------------|---------------------------------------------------------------|
| Apache-2.0        | `apache`, `apache-2.0`, `apache-2`, `apache 2.0`, `apache2.0` |
| MIT               | `mit`                                                         |
| GPL-3.0 or Later  | `gpl-3.0-or-later`                                            |
| GPL-3.0 Only      | `gpl-3.0-only`                                                |
| AGPL-3.0 or Later | `agpl-3.0-or-later`                                           |
| AGPL-3.0 Only     | `agpl-3.0-only`                                               |

如果没有您需要的模板，您可以使用 `--tmpl` (`-t`) 选项或者提交一个 Issue/PR。

如果你想要原封不动的把模板文件里的内容添加到你的文件头，请使用 `--rawtmpl` (`-r`) 选项，启用 raw template 后 NWA 不会根据你的文件类型生成不同类型（注释）的协议头。

**注意：`--tmpl` (`-t`) 和 `--rawtmpl` (`-r`) 不能同时使用。**

## 相关项目

- [nwa-examples](https://github.com/rainiring/nwa-examples): NWA 的使用示例；

## 博客

- [为您的代码文件添加许可声明：NWA一款更强大的许可声明管理工具](https://juejin.cn/post/7272025820967976996)

## 鸣谢

衷心感谢以下使 NWA 的开发成为可能的项目：

- [addlicense](https://github.com/google/addlicense)
- [cobra](https://github.com/spf13/cobra)
- [doublestar](https://github.com/bmatcuk/doublestar)
- [viper](https://github.com/spf13/viper)

## 许可证

NWA 使用 [Apache License 2.0](./LICENSE) 进行分发。NWA 的第三方依赖项的许可证说明在[此处](./licenses)。

## 生态

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BMS.png" alt="BMS"/>
<br/><br/>
NWA 是 <a href="https://github.com/B1NARY-GR0UP"> 基础中间件服务 </a> 的一个子项目
</p>