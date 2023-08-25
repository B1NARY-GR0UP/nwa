# NWA

> mending

![NWA](images/NWA.png)

A More Powerful License Header Management Tool

## Install

```shell
go install github.com/B1NARY-GR0UP/nwa@latest
```

## Usage

- **[Add](#add)**: Add license headers to files
- **[Check](#check)**: Check license headers of files
- **[Remove](#remove)**: Remove licenses headers of files
- **[Update](#update)**: Update license headers of files
- **[Config](#config)**: Edit files according to the configuration file

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

### Add

[EXAMPLE]()

```shell
Common Command | Add license headers to files
EXAMPLE: nwa add -l apache -c Lorain -m .

Usage:
  nwa add [flags]

Flags:
  -c, --copyright string   copyright holder (default "<COPYRIGHT HOLDER>")
  -h, --help               help for add
  -l, --license string     license type (default "apache")
  -m, --mute               mute mode
  -s, --skip strings       skip file path
  -t, --tmpl string        template file path
  -y, --year string        copyright year (default "2023")
```

### Check

[EXAMPLE]()

```shell
Common Command | Check license headers of files   
EXAMPLE: nwa check -t tmpl.txt .                  
NOTE: Do not use --mute (-m) flag with the command

Usage:                                                                    
  nwa check [flags]                                                       
                                                                          
Flags:                                                                    
  -c, --copyright string   copyright holder (default "<COPYRIGHT HOLDER>")
  -h, --help               help for check                                 
  -l, --license string     license type (default "apache")                
  -m, --mute               mute mode                                      
  -s, --skip strings       skip file path                                 
  -t, --tmpl string        template file path                             
  -y, --year string        copyright year (default "2023")   
```

### Remove

[EXAMPLE]()

```shell
Common Command | Remove licenses headers of files
EXAMPLE: nwa remove -l mit -c Anmory .           

Usage:                                                                    
  nwa remove [flags]                                                      
                                                                          
Flags:                                                                    
  -c, --copyright string   copyright holder (default "<COPYRIGHT HOLDER>")
  -h, --help               help for remove                                
  -l, --license string     license type (default "apache")                
  -m, --mute               mute mode                                      
  -s, --skip strings       skip file path                                 
  -t, --tmpl string        template file path                             
  -y, --year string        copyright year (default "2023") 
```

### Update

[EXAMPLE]()

```shell
Common Command | Update license headers of files                                    
EXAMPLE: nwa update -l mit -c Anmory .                                              
NOTE: Update identifies the content before the first blank line as a license header;
If your file does not meet the requirements, please use remove + add                

Usage:                                                                    
  nwa update [flags]                                                      
                                                                          
Flags:                                                                    
  -c, --copyright string   copyright holder (default "<COPYRIGHT HOLDER>")
  -h, --help               help for update                                
  -l, --license string     license type (default "apache")                
  -m, --mute               mute mode                                      
  -s, --skip strings       skip file path                                 
  -t, --tmpl string        template file path                             
  -y, --year string        copyright year (default "2023") 
```

### Config

[EXAMPLE]()

```shell
Config Command | Edit files according to the configuration file                                           
EXAMPLE: nwa config config.yaml                                                                           
NOTE: This command does not have any flag;                                                                
You can only specify the path of the configuration file, and everything depends on the configuration file;
If some configuration are not configured, the default configuration will be used                          
SAMPLE CONFIGURATION FILE(YAML):                                                                          
nwa:                                                                                                      
  cmd: "add"                                                                                              
  holder: "RHINE LAB.LLC."                                                                                
  year: "2077"                                                                                            
  license: "apache"                                                                                       
  mute: false                                                                                             
  path: ["server", "client", "pkg"]                                                                       
  skip: ["**.py"]                                                                                         
  tmpl: "nwa.txt"                                                                                         

Usage:                        
  nwa config [flags]          
                              
Flags:                        
  -h, --help   help for config
```

--mute -m => 静默模式(即，不输出修改的文件信息等) 默认 {false} 不会显示 Info 级别以下的信息，但是会显示 Warn 和 Error
--skip -s => 指定忽略的文件 默认 {}
--copyright -c => 指定 copyright holder 默认 {Copyright Holder}
--year -y => 指定年份 默认 {Current Year}
--license -l => 指定开源许可证类型 默认 {Apache-2.0}
--tmpl -t => 指定模板文件路径 默认 {} **和 -c -y -l 互斥** 完全自定义，允许你使用不同风格的注释

一个完整的配置文件如下（待定）：

```yaml
nwa:
  cmd: "add"
  holder: "RHINE LAB.LLC."
  year: "2077"
  license: "apache"
  mute: false
  path: ["server", "client", "pkg"]
  skip: ["**.py"]
  tmpl: "nwa.txt"
```

## License

NWA is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of NWA are explained [here](./licenses).

## Acknowledgement

## ECOLOGY

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BMS.png" alt="BMS"/>
<br/><br/>
NWA is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">Basic Middleware Service</a>
</p>