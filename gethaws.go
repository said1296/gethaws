package gethaws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/said1296/gethaws/internal/roundtripper"
	"net/http"
	"strings"
)

type (
	Provider      string
	CreationError struct {
		err error
	}
)

const (
	ProviderAws     Provider = "aws"
	ProviderRegular Provider = "regular"
)

func (e CreationError) Error() string {
	return fmt.Sprintf("Failed to create client(s): %s", e.err.Error())
}

// CreateClients determines the Provider type from the rpcUrl passed. If config is different than nil and the provider
// type is ProviderAws, then CreateAwsClientsFromConfig is used to generate the clients, otherwise CreateAwsClients is used.
func CreateClients(ctx context.Context, rpcUrl string, config *aws.Config) (client *ethclient.Client, rpcClient *rpc.Client, err error) {
	switch GetEvmProviderType(rpcUrl) {
	case ProviderAws:
		if config != nil {
			client, rpcClient, err = CreateAwsClientsFromConfig(rpcUrl, *config)
		} else {
			client, rpcClient, err = CreateAwsClients(ctx, rpcUrl)
		}
		if err != nil {
			return nil, nil, err
		}
	case ProviderRegular:
		rpcClient, err = rpc.DialContext(ctx, rpcUrl)
		if err != nil {
			return nil, nil, err
		}

		client = ethclient.NewClient(rpcClient)
	}

	return client, rpcClient, nil
}

// CreateAwsClients creates an ethclient.Client and an rpc.Client from from a context and an endpoint, AWS config
// is fetched from config.LoadDefaultConfig of aws package which uses environment variables to get credentials. Refer
// to: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
func CreateAwsClients(ctx context.Context, endpoint string) (*ethclient.Client, *rpc.Client, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, err
	}

	r, err := createRpcClient(endpoint, awsConfig)
	if err != nil {
		return nil, nil, CreationError{err}
	}

	return ethclient.NewClient(r), r, nil
}

// CreateAwsClientsFromConfig creates an ethclient.Client and an rpc.Client from an endpoint string and an aws.Config struct
func CreateAwsClientsFromConfig(endpoint string, config aws.Config) (*ethclient.Client, *rpc.Client, error) {
	r, err := createRpcClient(endpoint, config)
	if err != nil {
		return nil, nil, CreationError{err}
	}

	return ethclient.NewClient(r), r, nil
}

// createRpcClient creates an rpc.Client from an endpoint string and an aws.Config struct
func createRpcClient(endpoint string, config aws.Config) (*rpc.Client, error) {
	hc := new(http.Client)
	hc.Transport = roundtripper.NewHttpRoundTripper(config)
	return rpc.DialHTTPWithClient(endpoint, hc)
}

// GetEvmProviderType determines if the provider url is for an AWS or regular JSON-RPC
func GetEvmProviderType(rpcUrl string) Provider {
	if strings.Contains(strings.ToLower(rpcUrl), "ethereum.managedblockchain") {
		return ProviderAws
	}
	return ProviderRegular
}
