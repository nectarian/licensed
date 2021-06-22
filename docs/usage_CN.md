# 使用

## 前置条件 

* go 1.13+
* make

## 克隆

```
git clone https://github.com/nectarian/licensed.git
```

## 构建

```
make
```

## 生成密钥对

```
./dist/keymaker/keymaker -b 4096 -o output
private key generated : output/rsa_private.pem
 public key generated : output/rsa_public.pem

# 检查
ls -alh ./output/
```

## 收集客户端ID

> 此操作在客户端执行(text类型除外)

```
# network MAC Address
./dist/fingerprint/fingerprint -o ./output/mac.pem -t mac -n etho
client ID generated : ./output/mac.pem

# machine-id
./dist/fingerprint/fingerprint -o ./output/machine-id.pem -t machine-id
client ID generated : ./output/machine-id.pem

# baseboard , we should use sudo to get permissions
sudo ./dist/fingerprint/fingerprint -o ./output/baseboard.pem -t baseboard
client ID generated : ./output/baseboard.pem

ls -alh ./output/mac.pem ./output/machine-id.pem ./output/baseboard.pem
```

## 生成签名

```
./dist/postmark/postmark -e 2021-12-31 -k ./output/rsa_private.pem -o ./output/license.key -i ./output/machine-id.pem -m '{"corporation":"corporation-name","machine":"192.168.0.x"}'
license generated : ./output/license.key
```

## 在交付程序中添加控制代码

```
package main

import (
	"context"

	"github.com/nectarian/licensed/client"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	licenseFilePath := "./license.key"
	publicKeyPath := "./rsa_public.pem"
	go client.WatchLicense(ctx, "machine-id", licenseFilePath, publicKeyPath)

	// Do something
}
```

> 在交付时，将`license.key` 和 `rsa_public.pem` 随交付程序一起交付即可。