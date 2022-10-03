# Geth Client For AWS Managed Blockchain

# Installing

```go get github.com/said1296/gethaws```

# Limitations

Only works for HTTP client, will add support for other client types if people request it.

# Usage

The client creation functions receive an AWS Config object as defined in github.com/aws/aws-sdk-go-v2/config. Further configuration instructions available in:

https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/


```
import (
	"github.com/aws/aws-sdk-go-v2/config"
)

os.SetEnv("AWS_REGION", "us-east-2")
os.SetEnv("AWS_ACCESS_KEY_ID", "my_access_key_id")
os.SetEnv("AWS_SECRET_ACCESS_KEY", "my_secret_access_key")

awsConfig, err := config.LoadDefaultConfig(ctx)
if err != nil {
    panic(err)
}
awsConfig.HTTPClient = new(http.Client)

// Can also call gethaws.CreateClient or gethaws.CreateRpcClient to just get one type of client
client, rpcClient, err = gethaws.CreateClients(c.EvmProvider, awsConfig)
if err != nil {
    panic(err)
}
```