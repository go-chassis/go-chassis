# Cipher
## 概述

---

Go chassis以插件的形式提供加解密组件功能，用户可以自己定制，Chassis默认提供一种华为PAAS成熟的加解密组件，该功能目前以二进制的so文件形式提供。

## 配置

---

#### AES Cipher 配置

1、把so文件拷贝到项目${CHASSIS\_HOME}/lib目录下，或者直接放到系统/usr/lib目录下面。优先读取${CHASSIS\_HOME}/lib，然后再在/usr/lib下面查找。

2、通过环境变量PAAS\_CRYPTO\_PATH指定物料路径（root.key, common\_shared.key）

3、引入aes包，使用加解密方法

#### 自定义Cipher

可通过实现Cipher接口，自定义Cipher

```go
type Cipher interface {
    Encrypt(src string) (string, error)
    Decrypt(src string) (string, error)
}
```

## API

---

#### 加密

```
Encrypt(src string) (string, error)
```

#### 解密

```
Decrypt(src string) (string, error)
```

## 例子

---

#### 使用AES Cipher示例

```go
import (
    _ "github.com/servicecomb/security/plugins/aes"
    "testing"
    "github.com/servicecomb/security"
    "github.com/stretchr/testify/assert"
    "log"
)

func TestAESCipher_Decrypt(t *testing.T) {
    aesFunc := security.CipherPlugins["aes"]
    cipher := aesFunc()
    s, err := cipher.Encrypt("tian")
    assert.NoError(t, err)
    log.Println(s)
    a, _ := cipher.Decrypt(s)
    assert.Equal(t, "tian", a)
}
```

#### 自定义Cipher 示例

```go
package plain

import "github.com/servicecomb/security"

type DefaultCipher struct {
}

// register self-defined plugin in the cipherPlugin map
func init() {
    security.CipherPlugins["default"] = new
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



