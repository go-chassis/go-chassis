# 契约管理
## 概述

go-chassis读取服务契约并将其内容上传至注册中心。

## 配置

契约文件必须为yaml格式文件，契约文件应放置于go-chassis的schema目录。

schema目录位于：

1，conf/{serviceName}/schema，其中conf表示go-chassis的conf文件夹

2，${SCHEMA\_ROOT}

2的优先级高于1。

## API

包路径

```go
import "github.com/ServiceComb/go-chassis/core/config/schema"
```

契约字典，key值为契约文件名，value为契约文件内容

```go
var DefaultSchemaIDsMap map[string]string
```

## 示例

    conf
    `-- myservice
        `-- schema
            |-- myschema1.yaml
            `-- myschema2.yaml



