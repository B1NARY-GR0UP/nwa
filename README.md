# NWA

> A More Powerful License Header Management Tool

## Install

```shell
go install github.com/B1NARY-GR0UP/nwa@latest
```

## Usage

```shell
nwa add
nwa update
nwa remove
nwa config
```

nwa cmd [flags] path 

--help -h => 帮助信息
--copyright -c => 指定 copyright holder 默认 {Copyright Holder}
--year -y => 指定年份 默认 {Current Year}
--license -l => 指定开源许可证类型 默认 {Apache-2.0}
--mute -m => 静默模式(即，不输出修改的文件信息等) 默认 {false}
--tmpl -t => 指定模板文件路径 默认 {./nwa.tmpl}
--path -p => 指定文件路径 默认 {./} TODO: 待定
--style -s => 指定注释风格 默认 {line} 可选 {line, block} TODO: 如果要支持不同语言的话可能会有更多的选项

config 模式下，所有配置以配置文件为准，即配置文件的优先级高于命令行参数

一个完整的配置文件如下（待定）：

```yaml
nwa:
    cmd: "add"
    holder: "B1NARY-GR0UP"
    year: "2023"
    license: "Apache-2.0"
    mute: false
    tmpl: "./nwa.tmpl"
    path: "./"
```

## Roadmap

- [ ] 支持多种类型的文件