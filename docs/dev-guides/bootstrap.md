# BootstrapPlugin

### Introduction

go-chassis gives your a way to load any custom plugin you write before start go chassis, so that you don't need to add code in go-chassis project

you can use bootstrap plugin to load your custom logic to manipulate go-chassis modules, for example change Registry, router etc.


### Usage
```go
func InstallPlugin(name string, plugin BootstrapPlugin)
```