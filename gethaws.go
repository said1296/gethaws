package gethaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
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

// CreateClient creates an ethclient.Client from an endpoint string and an aws.Config struct
func CreateClient(endpoint string, config aws.Config) (*ethclient.Client, error) {
	c, _, err := CreateClients(endpoint, config)
	if err != nil {
		return nil, CreationError{err}
	}

	return c, nil
}

// CreateRpcClient creates an rpc.Client from an endpoint string and an aws.Config struct
func CreateRpcClient(endpoint string, config aws.Config) (*rpc.Client, error) {
	hc := new(http.Client)
	hc.Transport = roundtripper.NewHttpRoundTripper(config)
	return rpc.DialHTTPWithClient(endpoint, hc)
}

// CreateClients creates an ethclient.Client and an rpc.Client from an endpoint string and an aws.Config struct
func CreateClients(endpoint string, config aws.Config) (*ethclient.Client, *rpc.Client, error) {
	r, err := CreateRpcClient(endpoint, config)
	if err != nil {
		return nil, nil, CreationError{err}
	}

	return ethclient.NewClient(r), r, nil
}
