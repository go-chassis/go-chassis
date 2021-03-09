# Cipher
## 概述
Go chassis以插件的形式提供加解密组件功能，用户可以自己定制
in real environment, you may not want anyone to touch your real SK.
if anyone hack into your micro service system, 
then he can leverage AK SK to do anything to your cloud resource.

in other hand, by distributing encrypted SK to developer of micro services
can control the security risk. developer is only able to develop services, not able 
to touch any cloud resources. 
## Configuration

you can use cipher in SK and passphase decryption, or specify cipher plugin for generic usage 

```yaml
servicecomb:
  credentials:
    accessKey: xxx
    secretKey: xxx #ecrypted
    akskCustomCipher: default #used to decrypt sk if it is encrypted
```
```yaml
ssl:
  rest.Consumer.cipherPlugin: default
  rest.Consumer.certPwdFile: /path/to/passphase
  ...
```
```yaml
servicecomb:
  cipher:
    plugin: default
  ...
```
## Example

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
cipher.InstallCipherPlugin("default", new)
```

#### 加密

```
d, _ := cipher.Encrypt("ok")
```

#### 解密

```
```go
d, _ := cipher.Decrypt("ok")
```
```
