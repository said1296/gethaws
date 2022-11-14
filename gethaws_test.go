package gethaws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEvmProviderType(t *testing.T) {
	for _, tt := range []struct {
		rpcUrl   string
		expected Provider
	}{
		{
			rpcUrl:   "https://mainnet.alchemy.io/api/aksdnosknadoas",
			expected: ProviderRegular,
		},
		{
			rpcUrl:   "https://ujinaidnfjdsa.ethereum.managedblockchain.us-east-1.amazonaws.com",
			expected: ProviderAws,
		},
	} {
		pt := GetEvmProviderType(tt.rpcUrl)
		assert.Equal(t, tt.expected, pt)
	}
}

func TestGenerateProviderClients(t *testing.T) {
	ctx := context.Background()

	for _, tt := range []struct {
		rpcUrl   string
		config   *aws.Config
		expected Provider
	}{
		{
			rpcUrl:   "https://mainnet.alchemy.io/api/aksdnosknadoas",
			config:   nil,
			expected: ProviderRegular,
		},
		{
			rpcUrl:   "https://ujinaidnfjdsa.ethereum.managedblockchain.us-east-1.amazonaws.com",
			config:   nil,
			expected: ProviderAws,
		},
		{
			rpcUrl:   "https://ujinaidnfjdsa.ethereum.managedblockchain.us-east-1.amazonaws.com",
			config:   aws.NewConfig(),
			expected: ProviderAws,
		},
	} {
		_, _, err := CreateClients(ctx, tt.rpcUrl, tt.config)
		assert.NoError(t, err)
	}
}
