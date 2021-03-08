# Cipher

## Introduction
go chassis allows you to extend cipher plugin.

## Usage

1.Implement and install a new cipher
```go
//DefaultCipher is a struct
type DefaultCipher struct {
}

func new() security.Cipher {
	return &DefaultCipher{}
}

//Encrypt is method used for encryption
func (c *DefaultCipher) Encrypt(src string) (string, error) {
	return src, nil
}

//Decrypt is method used for decryption
func (c *DefaultCipher) Decrypt(src string) (string, error) {
	return  src, nil
}
```
```go
cipher.InstallCipherPlugin("noop", new)
```

2.Configure it in chassis.yaml
```yaml
servicecomb:
  cipher:
    plugin: noop
```

3. call API
```go
d, _ := cipher.Decrypt("ok")
```