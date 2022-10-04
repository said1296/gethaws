# Geth Client For AWS Managed Blockchain

# Installing

```go get github.com/said1296/gethaws```

# Limitations

Only works for HTTP client, will add support for other client types if people request it.

# Usage

Clients are configured using the AWS Config struct found in github.com/aws/aws-sdk-go-v2/config. Further configuration instructions available in:

https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/

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
    err =os.Setenv("AWS_SECRET_ACCESS_KEY", "my_secret_access_key")
    if err != nil {
        panic(err)
    }
    
    // The most common use cases will only use the first returned client, the rpc client is for low level calls not 
    // not implemented by geth.
    // The client can also be created with a manually created AWS Config by calling gethaws.CreateClientsFromConfig
    client, rpcClient, err := gethaws.CreateClients("https://infura.io/api_key", awsConfig)
    if err != nil {
        panic(err)
    }
}

```