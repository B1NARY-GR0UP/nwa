![NWA](images/NWA.png)

A Simple Yet Powerful Tool for License Header Management: Effortlessly Add, Check, Update, and Remove License Headers

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/nwa)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/nwa)

## Install

```shell
go install github.com/B1NARY-GR0UP/nwa@latest
```

Do not have a Go environment? Check the [Docker](#docker---run-nwa-through-docker-for-those-do-not-have-a-go-environment) section.

Or 

- [Use NWA in CI](#used-in-ci).
- [Use NWA in pre-commit](#used-in-pre-commit).

## Usage

- **[Flags](#flags)**: Use flags to customize the behavior of NWA
- **[DoubleStar(**) Patterns](#doublestar-patterns)**: Use patterns supported by [doublestar](https://github.com/bmatcuk/doublestar#patterns)
- **[Add](#add---add-license-headers-to-files)**: Add license headers to files
- **[Check](#check---check-license-headers-of-files)**: Check license headers of files
- **[Remove](#remove---remove-licenses-headers-of-files)**: Remove licenses headers of files
- **[Update](#update---update-license-headers-of-files)**: Update license headers of files
- **[Config](#config-mode)**: Edit files according to the configuration file
- **[Supported Licence Templates](#supported-licence-templates)**: Use built-in license templates or use custom templates 
- **[Docker](#docker---run-nwa-through-docker-for-those-do-not-have-a-go-environment)**: Run NWA through docker, for those do not have a Go environment
- **[Used in CI](#used-in-ci)**: Use NWA in CI
- **[Used in pre-commit](#used-in-pre-commit)**: Use NWA in pre-commit

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

| Short | Long        | Default                            | Description                                                                                                    |
|-------|-------------|------------------------------------|----------------------------------------------------------------------------------------------------------------|
| -c    | --copyright | `<COPYRIGHT HOLDER>`               | copyright holder                                                                                               |
| -y    | --year      | `time.Now().Year()` (Current Year) | copyright year                                                                                                 |
| -l    | --license   | `apache`                           | license type                                                                                                   |
| -s    | --skip      | `[]`                               | skip file paths, can use any pattern [supported by doublestar](https://github.com/bmatcuk/doublestar#patterns) |
| -V    | --verbose   | `false` (unspecified)              | verbose mode (Allow log output below the **WARN** level)                                                       |
| -m    | --mute      | `false` (unspecified)              | mute mode (Disable all log output)                                                                             |
| -i    | --spdxids   | `""`                               | SPDX IDs                                                                                                       |
| -t    | --tmpl      | `""`                               | template file path                                                                                             |
| -r    | --rawtmpl   | `""`                               | template file path (enable raw template)                                                                       |
| -f    | --fuzzy     | `false` (unspecified)              | `nwa check` will ignore differences in the **year** within the license header                                  |
| -h    | --help      | null                               | help for command                                                                                               |

### DoubleStar(**) Patterns

Both **`--skip` (`-s`)** and the **working path** support patterns recognized by [doublestar](https://github.com/bmatcuk/doublestar#patterns). 

However, some shells may interpret these patterns (e.g. `**`), which could cause NWA to behave unexpectedly. 

The best way to resolve this issue is to **wrap your paths in double quotes (`""`)**.

**NOTE: Since NWA always operates from the current directory (`.`), parent directories are not visible to NWA. Make sure not to use paths like `../path/to/file`.**

### Add - Add license headers to files

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

- Specify a custom template file using the `--tmpl` (`-t`) or `--rawtmpl` (`-r`) flag. Refer to [Supported Licence Templates](#supported-licence-templates)
- Submit an issue or PR to NWA

NWA will also output logs to inform you if any files already have a license header or if any files are not allowed to be edited (such as code files generated by tools).

### Check - Check license headers of files

- **Usage**

```shell
nwa check [flags] path...
```

- **Example**

```shell
nwa check --tmpl tmpl.txt "**/*.py"
```

The command in the example above **checks** whether the license header of all Python files matches the content specified in the `tmpl.txt` template file.

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

### Remove - Remove licenses headers of files

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

### Update - Update license headers of files

- **Usage**

```shell
nwa update [flags] path...
```

**NOTE: Update identifies the content before the first blank line as a license header; If your file does not meet the requirements, please use `remove` + `add` command.**


- **Example**

```shell
nwa update -l apache -c "BINARY Members" -s "dirA/**" -s "dirB/**/*.py" "**"
```

The command in the example above **updates** the license header of all files in the current directory, **except** for all files under `dirA` and Python files under `dirB`, to:

- License type: `Apache 2.0`
- Copyright holder: `BINARY Members`
- Copyright year: current year

regardless of the previous license header.

### Config Mode

- **Flags**

| Short | Long        | Default                            | Description        |
|-------|-------------|------------------------------------|--------------------|
| -c    | --command   | add                                | command to execute |
| -h    | --help      | null                               | help for config    |

**NOTE: If some configuration are not configured, the default configuration will be used.**

- **Usage**

```shell
nwa config [flags] [path]
```

**NOTE: Path is the configuration file path. If not specified, `.nwa-config.yaml` will be used as the default configuration file path.**

- **Example**

```shell
nwa config -c check
```

The command in the example above reads the `.nwa-config.yaml` configuration file in the current directory 
and checks whether the license headers of the specified files meet the requirements based on its content. 

The structure of the configuration file will be provided below.

- **Sample Configuration file**

**NOTE: If you set the `tmpl` or `rawtmpl` field, the `holder`, `year`, `license` and `spdxids` fields will be ignored.**

```yaml
nwa:
  cmd: "add"                        # Default: "add" Optional: "add", "check", "remove", "update", can be overwritten by --command (-c) flag
  holder: "RHINE LAB.LLC."          # Default: "<COPYRIGHT HOLDER>"
  year: "2077"                      # Default: Current Year
  license: "apache"                 # Default: "apache"
  spdxids: ""                       # Default: ""
  mute: false                       # Default: false (unspecified)
  verbose: false                    # Default: false (unspecified)
  fuzzy: false                      # Default: false (unspecified), only used for "check" command
  path: ["server/**", "example/**"] # Default: []
  skip: ["**.py"]                   # Default: []
  tmpl: "nwa.txt"                   # Default: ""                                                       
  rawtmpl: ""                       # Default: ""
```

### Supported Licence Templates

| License           | Option (ignore case)                                          |
|-------------------|---------------------------------------------------------------|
| Apache-2.0        | `apache`, `apache-2.0`, `apache-2`, `apache 2.0`, `apache2.0` |
| MIT               | `mit`                                                         |
| GPL-3.0 or Later  | `gpl-3.0-or-later`                                            |
| GPL-3.0 Only      | `gpl-3.0-only`                                                |
| AGPL-3.0 or Later | `agpl-3.0-or-later`                                           |
| AGPL-3.0 Only     | `agpl-3.0-only`                                               |

If the license template you need is not available, you could use the `--tmpl` (`-t`) flag, `--rawtmpl` (`-r`) flag or submit an Issue/PR.

**NOTE: The `--tmpl` (`-t`) and `--rawtmpl` (`-r`) flags cannot be used simultaneously.**

#### --tmpl (-t)

- **Example**

```shell
nwa add -t mytmpl.txt "**"
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

#### --rawtmpl (-r)

If you want to use the content of the template file directly to the header of your file without modification, please use the -rawtmpl (-r) flag. 
When the raw template is enabled, NWA will not generate different types of licence headers based on your file type.

- **Example**

```shell
nwa add --rawtmpl myrawtmpl.txt "**"
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

### Used in CI

- **GitHub Action Example**

When there are **mismatched** or **failed** entries in the result of `nwa check`, `nwa` will return a non-zero exit code, causing the CI to fail.

You may refer to the other commands or flags introduced in the [Usage](#usage) section for optimization. 

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
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
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
      rev: v0.5.0
      hooks:
        -   id: nwa-check
            args: ["-c", "BINARY Members", "-f", "-l", "apache"]
```

- **License Header Add Example**

**NOTE: `pre-commit` considers a hook to have failed if it modifies any files. Therefore, after using `add`, `update`, or `remove`, you may need to run `git add` again and commit the updated files.**

```yaml
repos:
  -   repo: https://github.com/B1NARY-GR0UP/nwa
      rev: v0.5.0
      hooks:
        -   id: nwa
            args: ["add", "-c", "BINARY Members", "-l", "apache"]
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