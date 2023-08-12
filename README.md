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

Local Flags (add, update, remove):

--mute -m => 静默模式(即，不输出修改的文件信息等) 默认 {false}
--skip -s => 指定忽略的文件 默认 {}
--copyright -c => 指定 copyright holder 默认 {Copyright Holder}
--year -y => 指定年份 默认 {Current Year}
--license -l => 指定开源许可证类型 默认 {Apache-2.0}
--tmpl -t => 指定模板文件路径 默认 {} **和 -c -y -l 互斥**

**config 模式下，所有配置以配置文件为准，即配置文件的优先级高于命令行参数**

一个完整的配置文件如下（待定）：

```yaml
nwa:
  cmd: "add"
  holder: "RHINE LAB.LLC."
  year: "2023"
  license: "apache"
  mute: false
  path: ["server/*.go"]
  skip: ["example/**", "test/**"]
  tmpl: "./nwa.tmpl"
```

## Related Projects

- [VIOLIN](https://github.com/B1NARY-GR0UP/violin) | VIOLIN worker/connection pool | `go` `connection-pool` `worker-pool`

## License

NWA is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of NWA are explained [here](./licenses).

## ECOLOGY

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BINARY-WEB-ECO.png" alt="BINARY-WEB-ECO"/>
<br/><br/>
NWA is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">BINARY WEB ECOLOGY</a>
</p>