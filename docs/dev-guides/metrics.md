# Metrics Report Plugin

## Introduction
go chassis allows you to install report plugin to receive runtime metrics 

## Usage

1.Implement your reporter

```go
type Reporter func(metrics.Registry) error
```

2.Install reporter to go chassis
```go
func InstallReporter(name string, reporter Reporter) error
```

after above step your plugin is able to receive go chassis runtime metrics  