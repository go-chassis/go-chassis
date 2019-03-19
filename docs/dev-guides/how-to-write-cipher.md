# Cipher
## 概述

---

Go chassis以插件的形式提供加解密组件功能，用户可以自己定制

## Configuration

you can use cipher in SK and passphase decryption

```yaml
cse:
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

## API

---
可通过实现Cipher接口，自定义Cipher

```go
type Cipher interface {
    Encrypt(src string) (string, error)
    Decrypt(src string) (string, error)
}
```
#### 加密

```
Encrypt(src string) (string, error)
```

#### 解密

```
Decrypt(src string) (string, error)
```

## Example

---

#### Protect your SK
in real environment, you may not want anyone to touch your real SK.
if anyone hack into your micro service system, 
then he can leverage AK SK to do anything to your cloud resource.

in other hand, by distributing encrypted SK to developer of micro services
can control the security risk. developer is only able to develop services, not able 
to touch any cloud resources. 


1.Implement and install cipher to go chassis
```go
package plain

import sec "github.com/go-chassis/go-chassis/security"
import "github.com/go-chassis/foundation/security"

type DefaultCipher struct {
}

// register plugin 
func init() {
    sec.InstallCipherPlugin("custom" ,new)
}

// define a method of newing a plugin object, and register this method
func new() security.Cipher {
    return &DefaultCipher{
    }
}

// implement the Encrypt(string) method, to encrypt the clear text
func (c *DefaultCipher)Encrypt(src string) (string, error) {
    return src, nil
}

// implement the Decrypt(string) method, to decrypt the cipher text
func (c *DefaultCipher)Decrypt(src string) (string, error) {
    return src, nil
}
```
2.encrypt your SK with this cipher implementation
(for example you can develop a simple command tool for encryption)

3.Set cipher name in auth.yaml to decrypt SK
```yaml
cse:
  credentials:
    accessKey: xxx
    secretKey: xxx #ecrypted
    akskCustomCipher: custom #used to decrypt sk if it is encrypted
```