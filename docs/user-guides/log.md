# Log
## 概述

用户可配置微服务的运行日志的相关属性，比如输出方式，日志级别，文件路径以及日志转储相关属性。

## 配置

日志配置文件为lager.yaml，配置模板如下：

- logLevel: 由低到高分别为 DEBUG, INFO, WARN, ERROR, FATAL 共5个级别，这里设置的级别是日志输出的最低级别，只有不低于该级别的日志才会输出。
- logWriters: 表示日志的输出方式，默认为文件和标准输出。
- logFile: 日志路径
- logFormatText: 默认为false，即设定日志的输出格式为 json。若为true则输出格式为plaintext，类似log4j。建议使用json格式输出的日志。
- logRotateDisable: 是否开启日志绕接.
- logRotateCompress: 是否压缩旧的日志
- logRotateAge: 日志rotate时间配置，单位"day"，范围为(0, 10)。
- logRotateSize: 日志rotate文件大小配置，单位"MB",范围为(0,50)。
- logBackupCount: 日志最大存储数量，单位“个”,范围为[0,100)。

```yaml
---
logWriters: file,stdout
# LoggerLevel: |DEBUG|INFO|WARN|ERROR|FATAL
logLevel: DEBUG
logFile: log/chassis.log
logFormatText: false
logRotateDisable: false
#log rotate and backup settings
logRotateAge: 1 # after n days
logRotateSize: 10 # megabytes
logBackupCount: 7
```
