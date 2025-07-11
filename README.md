![NWA](images/NWA.png)

A Simple Yet Powerful Tool for License Header Management: Effortlessly Add, Check, Update, and Remove License Headers

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/nwa)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/nwa)

## Install

### Homebrew

```shell
brew tap B1NARY-GR0UP/nwa
brew install nwa
```

### Go

```shell
go install github.com/B1NARY-GR0UP/nwa@latest
```

### GitHub Actions

Check the [Used in GitHub Actions](#used-in-github-actions) section.

### pre-commit

Check the [Used in pre-commit](#used-in-pre-commit) section.

### Docker

Check the [Docker](#docker---run-nwa-through-docker-for-those-do-not-have-a-go-environment) section.

## Usage

- **[Flags](#flags)**: Use flags to customize the behavior of NWA
- **[DoubleStar(**) Patterns](#doublestar-patterns)**: Use patterns supported by [doublestar](https://github.com/bmatcuk/doublestar#patterns)
- **[Add](#add---add-license-headers-to-files)**: Add license headers to files
- **[Check](#check---check-license-headers-of-files)**: Check license headers of files
- **[Remove](#remove---remove-licenses-headers-of-files)**: Remove licenses headers of files
- **[Update](#update---update-license-headers-of-files)**: Update license headers of files
- **[Config](#config-mode)**: Edit files according to the configuration file
- **[Built-in License Header Templates and Custom Templates](#built-in-license-header-templates-and-custom-templates)**: Use custom templates 
- **[Used in GitHub Actions](#used-in-github-actions)**: Use NWA in GitHub Actions
- **[Used in pre-commit](#used-in-pre-commit)**: Use NWA in pre-commit
- **[Docker](#docker---run-nwa-through-docker-for-those-do-not-have-a-go-environment)**: Run NWA through docker, for those who do not have a Go environment

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

### Flags

- **Basic**

| Short | Long        | Default              | Description                                                                                                    |
|-------|-------------|----------------------|----------------------------------------------------------------------------------------------------------------|
| -h    | --help      | null                 | help for command                                                                                               |
| -c    | --copyright | `<COPYRIGHT HOLDER>` | copyright holder                                                                                               |
| -y    | --year      | Current Year         | copyright year                                                                                                 |
| -l    | --license   | `apache`             | license type                                                                                                   |
| -i    | --spdxids   | `""`                 | SPDX IDs                                                                                                       |
| -s    | --skip      | `[]`                 | skip file paths, can use any pattern [supported by doublestar](https://github.com/bmatcuk/doublestar#patterns) |

- **Advanced**

| Short | Long       | Default               | Description                                                                                                                                            |
|-------|------------|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| -V    | --verbose  | `false` (unspecified) | verbose mode (Allow log output below the **WARN** level)                                                                                               |
| -m    | --mute     | `false` (unspecified) | mute mode (Disable all log output)                                                                                                                     |
| -t    | --tmpl     | `""`                  | template file path                                                                                                                                     |
| -T    | --tmpltype | `""`                  | template type (`live`, `static`, `raw`)                                                                                                                |
| -f    | --fuzzy    | `false` (unspecified) | commands `check` and `remove` will ignore differences in the **year** within the license header                                                        |
| -k    | --keyword  | `[]`                  | keyword used to confirm the existence of license headers (only used in commands `add` and `update`)                                                    |
| -S    | --style    | `[]`                  | customize the comment style (`line`, `block`, `hash`, `doc`, `starred-block`) for different extensions in the format `extension:style`, e.g.`go:block` |

### DoubleStar(**) Patterns

Both **`--skip` (`-s`)** and the **working path** support patterns recognized by [doublestar](https://github.com/bmatcuk/doublestar#patterns). 

However, some shells may interpret these patterns (e.g. `**`), which could cause NWA to behave unexpectedly. 

The best way to resolve this issue is to **wrap your paths in double quotes (`""`)**.

Additionally, since Windows uses a different path separator (`\`) compared to Linux and macOS, NWA internally converts `\` to `/` when traversing files to match doublestar patterns. Therefore, make sure to **use `/` as the path separator** in both `--skip` (`-s`) and the working path.

> **NOTE: Since NWA always operates from the current directory (`.`), parent directories are not visible to NWA. Make sure not to use paths like `../path/to/file`.**

<details>
<summary><h3>Add - Add license headers to files</h3></summary>

- **Usage**

```shell
nwa add [flags] path...
```
- **Example**

```shell
nwa add -l apache -c "RHINE LAB.LLC." -y 2077 "src/**/*.go" "pkg/pool/**"
```

The command in the example above **adds** a license header to all Go files under the folder with the relative path `src` and to all files under the folder with the relative path `pkg/pool`:

- License type: `Apache 2.0`
- Copyright holder: `RHINE LAB.LLC.`
- Copyright year: `2077`

NWA will generate a corresponding license header based on the file type in the specified paths. For example, `.py` files will use `#` for comments, and `.go` files will use `//`.

If your file type is not supported by NWA, you can:

- Using custom templates. Refer to [Built-in License Header Templates and Custom Templates](#built-in-license-header-templates-and-custom-templates)
- Submit an issue or PR to NWA

NWA will also output logs to inform you if any files already have a license header or if any files are not allowed to be edited (such as code files generated by tools).

</details>

<details>
<summary><h3>Check - Check license headers of files</h3></summary>

- **Usage**

```shell
nwa check [flags] path...
```

> **NOTE: The `check` command uses exact matching to verify whether a file's license header meets the requirements. If you have customized license header needs, you can use the custom template feature. See [Built-in License Header Templates and Custom Templates](#built-in-license-header-templates-and-custom-templates) for more details.**

- **Example**

```shell
nwa check --copyright "BINARY Members" --license apache --fuzzy "**/*.py"
```

The command in the example above **checks** whether the license headers of all Python files comply with the following requirements:

- License type: `Apache 2.0`
- Copyright holder: `BINARY Members`

After the check is complete, NWA will output the results as logs. A sample output is as follows:

```txt
2024/11/24 19:24:29 WARN file does not have a matched header path=dirB\dirC\fileC.go
2024/11/24 19:24:29 WARN file does not have a matched header path=main.go
[NWA SUMMARY] scanned=4 matched=2 mismatched=2 skipped=1 failed=0
```

Verbose mode:

```txt
2024/11/24 19:24:35 INFO skip file path=README.md
2024/11/24 19:24:35 INFO file has a matched header path=dirA\fileA.go
2024/11/24 19:24:35 WARN file does not have a matched header path=dirB\dirC\fileC.go
2024/11/24 19:24:35 WARN file does not have a matched header path=main.go
2024/11/24 19:24:35 INFO file has a matched header path=dirB\fileB.go
[NWA SUMMARY] scanned=4 matched=2 mismatched=2 skipped=1 failed=0
```

Mute mode:

```txt
```

</details>

<details>
<summary><h3>Remove - Remove licenses headers of files</h3></summary>

- **Usage**

```shell
nwa remove [flags] path...
```

- **Example**

```shell
nwa remove -l mit -c "RHINE LAB.LLC." -s "src/vender/**" "src/**"
```
The command in the example above skips all files under the `src/vendor` folder and **removes** the license header for files under the folder with the relative path `src`:

- License type: `MIT`
- Copyright holder: `RHINE LAB.LLC.`
- Copyright year: current year

If a file in the specified path does not have a license header or is not allowed to be edited, NWA will inform you through log output.

</details>

<details>
<summary><h3>Update - Update license headers of files</h3></summary>

- **Usage**

```shell
nwa update [flags] path...
```

> **NOTE: Update identifies the content before the first blank line as a license header; If your file does not meet the requirements, please use `remove` + `add` command.**
>
> If your files use the following types of comment styles, the `update` command may not work properly (refer to [#18](https://github.com/B1NARY-GR0UP/nwa/issues/18)):
> 
> - `<!-- -->`: .html, .md, .vue, etc.
> - `{# #}`: .j2, .twig, etc.
> - `(** *)`: .ml, .mli, etc.
>  
> Please use the `remove` + `add` commands instead of the `update` command.

- **Example**

```shell
nwa update -l apache -c "BINARY Members" -s "dirA/**" -s "dirB/**/*.py" "**"
```

The command in the example above **updates** the license header of all files in the current directory, **except** for all files under `dirA` and Python files under `dirB`, to:

- License type: `Apache 2.0`
- Copyright holder: `BINARY Members`
- Copyright year: current year

regardless of the previous license header.

</details>

### Config Mode

- **Flags**

| Short | Long        | Default                            | Description        |
|-------|-------------|------------------------------------|--------------------|
| -c    | --command   | add                                | command to execute |
| -h    | --help      | null                               | help for config    |

> **NOTE: If some configuration is not configured, the default configuration will be used.**

- **Usage**

```shell
nwa config [flags] [path]
```

> **NOTE: Path is the configuration file path. If not specified, `.nwa-config.yaml` will be used as the default configuration file path.**

- **Example**

```shell
nwa config -c check
```

The command in the example above reads the `.nwa-config.yaml` configuration file in the current directory 
and checks whether the license headers of the specified files meet the requirements based on its content. 

The structure of the configuration file will be provided below.

- **Sample Configuration file**

```yaml
nwa:
  # basic
  cmd: "add"                        # Default: "add"; Can be overwritten by --command (-c) flag; Optional: "add", "check", "remove", "update"
  holder: "RHINE LAB.LLC."          # Default: "<COPYRIGHT HOLDER>"
  year: "2077"                      # Default: Current Year
  license: "apache"                 # Default: "apache"
  spdxids: ""                       # Default: ""
  skip:                             # Default: []
    - "**.py"
  path:                             # Default: []
    - "src/**/*.go"
    - "example/**/*.go"
  # advanced
  mute: false                       # Default: false (unspecified)
  verbose: false                    # Default: false (unspecified)
  fuzzy: false                      # Default: false (unspecified); Used for "check" and "remove" commands
  tmpltype: ""                      # Default: ""; Optional: "live", "static", "raw"
  tmpl: ""                          # Default: ""                                                       
  keyword: []                       # Default: []; Used for "add" and "update" commands
  style: []                         # Default: []
```

### Built-in License Header Templates and Custom Templates

- **Built-in License Header Templates**

| License           | Option (ignore case)                                          |
|-------------------|---------------------------------------------------------------|
| Apache-2.0        | `apache`, `apache-2.0`, `apache-2`, `apache 2.0`, `apache2.0` |
| MIT               | `mit`                                                         |
| GPL-3.0 or Later  | `gpl-3.0-or-later`                                            |
| GPL-3.0 Only      | `gpl-3.0-only`                                                |
| AGPL-3.0 or Later | `agpl-3.0-or-later`                                           |
| AGPL-3.0 Only     | `agpl-3.0-only`                                               |

- **Custom Templates**

> **NOTE: When using a custom template, you must specify both the template type and the template itself (use a file path in common mode, and template content in config mode).**

<details>
<summary><h4>Live Template</h4></summary>

- **Example**

```shell
nwa config -c check
```

`.nwa-config.yaml` is as follows:

```yaml
nwa:
  holder: BINARY Members
  tmpltype: live
  tmpl: |
    Copyright {{ .Year }} {{ .Holder }}

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
  skip:
    - "testdata/**"
  path:
    - "**/*.go"
```

The above command uses the `holder` and `year` specified in the `.nwa-config.yaml` configuration file (defaulting to the current year if not set) to render the complete template, and runs the `check` command with the `-c` option to verify whether the license headers in the files meet the requirements.

</details>

<details>
<summary><h4>Static Template</h4></summary>

- **Example**

```shell
nwa add -T static -t mytmpl.txt "**/*.py"
```

`mytmpl.txt` is as follows:

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

This command example uses the content in `mytmpl.txt` as the license header, and NWA will generate license headers with different comment types based on the file type.

</details>

<details>
<summary><h4>Raw Template</h4></summary>

- **Example**

```shell
nwa add -T raw -t myrawtmpl.txt "**/*.java"
```

`myrawtmpl.txt` is as follows:

```txt
// Copyright 2077 RHINE LAB.LLC.
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
```

This command example uses the content in `myrawtmpl.txt` as the license header. NWA will **not** generate license headers with different comment types based on the file type **but** will instead add the content of `myrawtmpl.txt` directly to each file (include blank line).

</details>

### Used in GitHub Actions

#### Automatic Configuration

For the configuration file example, please refer to the [Config](#config-mode) section.

```yaml
- name: License Header Check
  uses: B1NARY-GR0UP/nwa@main
  with:
    version: latest        # (optional) version of nwa to use; default: latest
    cmd: check             # (optional) command to execute; options: `check`, `add`, `update`, `remove`; default: check
    path: .nwa-config.yaml # (optional) configuration file path; default: .nwa-config.yaml
```

#### Manual Configuration

When there are **mismatched** or **failed** entries in the result of nwa execution, `nwa` will return a non-zero exit code, causing the CI to fail.

You may refer to the other commands or flags introduced in the [Usage](#usage) section for optimization. 

- **Example**

```yaml
name: License Header Check

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: 1.23

      - name: Install NWA
        run: go install github.com/B1NARY-GR0UP/nwa@latest

      - name: Run License Header Check
        run: nwa check -c "BINARY Members" -f -l apache "**/*.go"
```

### Used in pre-commit

Thanks to the support of the [pre-commit](https://pre-commit.com/) tool, NWA can be used as a pre-commit hook.

Currently, NWA offers two hooks, `nwa-check` and `nwa`, as specified in [.pre-commit-hooks.yaml](./.pre-commit-hooks.yaml). 
You can use them by configuring the [.pre-commit-config.yaml](https://pre-commit.com/#adding-pre-commit-plugins-to-your-project) file in your Git repository.

Below are examples of `.pre-commit-config.yaml` file using the `nwa-check` and `nwa` hooks, which you can modify according to your needs:

**NOTE: A Golang runtime environment is required.**

- **License Header Check Example**

```yaml
repos:
  -   repo: https://github.com/B1NARY-GR0UP/nwa
      rev: [version tag]
      hooks:
        -   id: nwa-check
            args: ["-c", "BINARY Members", "-f", "-l", "apache"]
```

- **License Header Add Example**

**NOTE: `pre-commit` considers a hook to have failed if it modifies any files. Therefore, after using `add`, `update`, or `remove`, you may need to run `git add` again and commit the updated files.**

```yaml
repos:
  -   repo: https://github.com/B1NARY-GR0UP/nwa
      rev: [version tag]
      hooks:
        -   id: nwa
            args: ["add", "-c", "BINARY Members", "-l", "apache"]
```

### Docker - Run NWA through docker, for those do not have a Go environment

- **Install**

Install the nwa docker image directly:

```shell
docker pull ghcr.io/b1nary-gr0up/nwa:main
```

**OR**

Build it from source:

```shell
docker build -t ghcr.io/b1nary-gr0up/nwa:main .
```

- **Verify if it can work correctly**

```shell
docker run -it ghcr.io/b1nary-gr0up/nwa:main --version
```

- **Mount the directory you want NWA to work with to `/src` and use the commands mentions in usage**

```shell
docker run -it -v ${PWD}:/src ghcr.io/b1nary-gr0up/nwa:main add -c "RHINE LAB.LLC." -y 2077 "**"
```

## Blogs

- [Add License Headers to Your Code Files](https://dev.to/justlorain/add-license-headers-to-your-code-files-nwa-a-more-powerful-license-statement-management-tool-86o) 

## Credits

Sincere appreciation to the following repositories that made the development of NWA possible.

- [addlicense](https://github.com/google/addlicense)
- [cobra](https://github.com/spf13/cobra)
- [doublestar](https://github.com/bmatcuk/doublestar)
- [viper](https://github.com/spf13/viper)

## License

NWA is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of NWA are explained [here](./licenses).

## ECOLOGY

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BMS.png" alt="BMS"/>
<br/><br/>
NWA is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">Basic Middleware Service</a>
</p>