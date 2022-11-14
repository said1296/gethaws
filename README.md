# Geth Client For AWS Managed Blockchain

Generates RPC clients for regular JSON-RPC clients such as Infura or Alchemy, and for Managed Blockchain which uses a 
custom authentication mechanism.

# Installing

```go get github.com/said1296/gethaws```

# Limitations

Only works for HTTP client, will add support for other client types if people request it.

# Usage

In this example the clients are configured loading the configuration from env variables. Further config instructions found at:

https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/

You can also pass an aws.Config struct instead of nil to CreateClients in order to use a custom aws.Config.

```
package main

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/config"
    "net/http"
    "os"
)

func main() {
    err := os.Setenv("AWS_REGION", "us-east-2")
    if err != nil {
        panic(err)
    }
    err = os.Setenv("AWS_ACCESS_KEY_ID", "my_access_key_id")
    if err != nil {
        panic(err)
    }
    err = os.Setenv("AWS_SECRET_ACCESS_KEY", "my_secret_access_key")
    if err != nil {
        panic(err)
    }
    
    // The most common use cases will only use the first returned client, the rpc client is for low level calls not 
    // implemented by geth.
    // The clients can also be created with a manually created aws.Config by passing it instead of nil.
    client, rpcClient, err := gethaws.CreateClients(ctx.Background(), "https://ethereum.managedblockchain/1jsj1i23213nk32mo1", nil)
    if err != nil {
        panic(err)
    }
}

```