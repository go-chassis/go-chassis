# Tooling

## Introduction
go chassis abstract common tool, such as codec, cipher. 
With those plugins you can easily switch implementation.

## Usage
the development of plugin has the same pattern, use codec for example. 

1.Implement and install a new function
```go
// StdJson implement standard json codec
type StdJson struct {
}

func newDefault(opts Options) (codec.Codec, error) {
return &StdJson{}, nil
}
func (s *StdJson) Encode(v any) ([]byte, error) {
return json.Marshal(v)
}

func (s *StdJson) Decode(data []byte, v any) error {
return json.Unmarshal(data, v)
}
...
```
```go
codec.Install("encoding/json", newDefault)
```

2.Configure it in chassis.yaml
```yaml
servicecomb:
  codec:
    plugin: encoding/json
```

3. just call API before you create a resource
```go
data, err := codec.Encode(Person{Name: "a"})
```