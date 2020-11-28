# Backend toolkit sets

## Introduction
go chassis allows you to extend backend plugin, like quota management system.

in fact all backend plugin is design to be a common module, more like a toolkit,
which can be used out of go chassis framework. 

what go chassis does is just making module ready to be called by configuration file.

this guide only shows you how to use with go chassis.

## Usage
the development of plugin has the same pattern, use quota for example. 

1.Implement and install a new function
```go
type inMemory struct {
}
func (im *inMemory) GetQuotas(service, domain string) ([]*quota.Quota, error) {
	return []*quota.Quota{
		{ResourceName: "cpu", Used: 10, Limit: 20}, {ResourceName: "mem", Used: 10, Limit: 256},
	}, nil
}
...
```
```go
quota.Install("mock", func() (quota.Manager, error) {
			return &inMemory{}, nil
		})
```

2.Configure it in chassis.yaml
```yaml
servicecomb:
  quota:
    plugin: mock
```

3. just call API before you create a resource
```go
quota.PreCreate("some cloud service", "some user", "cpu", 2)
```