# Log
## 概述

用户可配置微服务的运行日志的相关属性，比如输出方式，日志级别，文件路径以及日志转储相关属性。

## 配置

日志配置文件为lager.yaml，配置模板如下：

- logger_level表示日志级别，由低到高分别为 DEBUG, INFO, WARN, ERROR, FATAL 共5个级别，这里设置的级别是日志输出的最低级别，只有不低于该级别的日志才会输出。
- writers表示日志的输出方式，默认为文件和标准输出。
- logger_file表示日志输出文件。
- log_format_text: 默认为false，即设定日志的输出格式为 json。若为true则输出格式为plaintext，类似log4j。建议使用json格式输出的日志。
- rollingPolicy: 默认为size，即根据大小进行日志rotate操作；若配置为daily则基于事件做日志rotate。
- log_rotate_date: 日志rotate时间配置，单位"day"，范围为(0, 10)。
- log_rotate_size: 日志rotate文件大小配置，单位"MB",范围为(0,50)。
- log_backup_count: 日志最大存储数量，单位“个”,范围为[0,100)。

```yaml
---
writers: file,stdout
# LoggerLevel: |DEBUG|INFO|WARN|ERROR|FATAL
logger_level: DEBUG
logger_file: log/chassis.log
log_format_text: false

#rollingPolicy daily/size
rollingPolicy: size
#log rotate and backup settings
log_rotate_date: 1
log_rotate_size: 10
log_backup_count: 7
```

## API

通过配置lager.yaml，go-chassis会自动为服务加载日志模块，并默认输出日志到文件或标准输出。

##### Debug和Info级别的日志

```go
lager.Logger.Info(action string, data ...Data)
lager.Logger.Debug(action string, data ...Data)

lager.Logger.Debugf(format string, args ...interface{})
lager.Logger.Infof(format string, args ...interface{})
```

##### Warn Error Fatal级别的日志

```go
lager.Logger.Warn(action string, err error, data ...Data)
lager.Logger.Error(action string, err error, data ...Data)
lager.Logger.Fatal(action string, err error, data ...Data)

lager.Logger.Warnf(err error, format string, args ...interface{})
lager.Logger.Errorf(err error, format string, args ...interface{})
lager.Logger.Fatalf(err error, format string, args ...interface{})
```

