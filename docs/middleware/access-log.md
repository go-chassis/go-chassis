# how to use this model

1. add this in provider chain, and as the first handler
```
cse:
  service:
    registry:
      disabled: true
      registry: manual
  protocols: # what kind of server you want to launch
    rest: #launch a http server
      listenAddress: 127.0.0.1:5001
  handler:
    chain:
      Provider:
        default: access-log
```

2. add a config in lager.yaml
```
# can be a file path or stdout
# a file path: record access log in this file, recommend access file path' dir is same as log file'dir
# stdout: access log will record in console stdout
access_log_file: xxx
```

3. import access log package
```
// should import after import go-chassis
	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/middleware/accesslog"
```
