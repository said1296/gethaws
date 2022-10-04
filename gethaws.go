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
)

type CreationError struct {
	err error
}

func (e CreationError) Error() string {
	return fmt.Sprintf("Failed to create client(s): %s", e.err.Error())
}

// CreateClients creates an ethclient.Client and an rpc.Client from from a context and an endpoint, AWS config
// is fetched from config.LoadDefaultConfig of AWS package
func CreateClients(ctx context.Context, endpoint string) (*ethclient.Client, *rpc.Client, error) {
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

// CreateClientsFromConfig creates an ethclient.Client and an rpc.Client from an endpoint string and an aws.Config struct
func CreateClientsFromConfig(endpoint string, config aws.Config) (*ethclient.Client, *rpc.Client, error) {
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
