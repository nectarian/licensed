# USAGE

## Prerequisites

* go 1.13+
* make

## clone

```
git clone https://github.com/nectarian/licensed.git
```

## build

```
make
```

## generate key pair

```
./dist/keymaker/keymaker -b 4096 -o output
private key generated : output/rsa_private.pem
 public key generated : output/rsa_public.pem

# check
ls -alh ./output/
```

## Collect the unique identifier of the client

> This operation is performed on the customer client machine (except `text` type)

```
# network MAC Address
./dist/fingerprint/fingerprint -o ./output/mac.pem -t mac -n etho
client ID generated : ./output/mac.pem

# machine-id
./dist/fingerprint/fingerprint -o ./output/machine-id.pem -t machine-id
client ID generated : ./output/machine-id.pem

# baseboard , we should Use sudo to get permissions
sudo ./dist/fingerprint/fingerprint -o ./output/baseboard.pem -t baseboard
client ID generated : ./output/baseboard.pem

ls -alh ./output/mac.pem ./output/machine-id.pem ./output/baseboard.pem
```

## signature

```
./dist/postmark/postmark -e 2021-12-31 -k ./output/rsa_private.pem -o ./output/license.key -i ./output/machine-id.pem -m '{"corporation":"corporation-name","machine":"192.168.0.x"}'   
license generated : ./output/license.key
```

## client control

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

> `license.key` and `rsa_public.pem` should be delivered with the program.
