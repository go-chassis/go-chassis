# Circuit breaker

## Introduction
go chassis allows you to define how to handle errors, if a remote call fails

## Usage

1.Implement your fallback logic

```go
//Fallback defines how to return response if remote call fails
//a implementation should return a closure to handle the error
//in this closure, if you fallback logic should handle the original error,
//you can return a fallback error to replace the original error
//you can assemble invocation.Response on demand
//in summary the closure defines, if err happens, how to handle it.
type Fallback func(inv *invocation.Invocation, finish chan *invocation.Response) func(error) error
```

2.Register function to go chassis
```go
circuit.RegisterFallback("your_fallback", f)
```

3 operate circuit_breaker.yaml to use custom fallback

```yaml
cse:
  fallbackpolicy:
    Consumer:
      policy: your_fallback
```