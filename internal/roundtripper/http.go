package roundtripper

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"io"
	"net/http"
	"time"
)

// httpRoundtripper holds the AWS config and satisfies the http.RoundTripper interface so that it can be used with an
// http.Client
type httpRoundtripper struct {
	config aws.Config
}

// NewHttpRoundTripper creates an http.Client using an http.RoundTripper that handles AWS Managed Blockchain requests
func NewHttpRoundTripper(cfg aws.Config) http.RoundTripper {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = new(http.Client)
	}

	return httpRoundtripper{
		config: cfg,
	}
}

// RoundTrip gets AWS credentials from the aws.Config, adds the necessary AWS heahders to the request, performs the request,
// and finally parses the gzip response from Managed Blockchain
func (h httpRoundtripper) RoundTrip(request *http.Request) (*http.Response, error) {
	credentials, err := h.config.Credentials.Retrieve(request.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	internalRequest := request.Clone(request.Context())

	hash, err := requestDataHash(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request for v4 signature: %w", err)
	}

	signer := v4.NewSigner()
	err = signer.SignHTTP(
		context.Background(),
		credentials,
		internalRequest,
		hash,
		"managedblockchain",
		h.config.Region,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed v4 sign request: %w", err)
	}

	response, err := h.config.HTTPClient.Do(internalRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to perform Managed Blockchain request: %w", err)
	}

	if response.Header.Get("Content-Type") == "gzip" {
		gzipReader, err := gzip.NewReader(base64.NewDecoder(base64.StdEncoding, response.Body))
		if err != nil {
			return nil, fmt.Errorf("failed decode gzip content-type response: %w", err)
		}

		request.Header.Set("Content-Type", "application/json")

		response.Body = gzipReader
	}

	return response, nil
}

// hashRequest gets the sha256 hash of the request body so that it can be signed for AWS V4 signature authentication
func requestDataHash(req *http.Request) (string, error) {
	var requestData []byte
	if req.Body != nil {
		requestBody, err := req.GetBody()
		if err != nil {
			return "", err
		}
		defer requestBody.Close()

		requestData, err = io.ReadAll(io.LimitReader(requestBody, 1<<20))
		if err != nil {
			return "", err
		}
	}

	return getSha256(requestData)
}

// getSha256 returns the sha256 hash of a byte array as a string
func getSha256(input []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(input); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
