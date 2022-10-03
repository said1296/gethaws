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

type RoundTripperError struct {
	err error
}

func (e RoundTripperError) Error() string {
	return fmt.Sprintf("Rountripper error: %s", e.err.Error())
}

type httpRoundtripper struct {
	config aws.Config
}

func NewHttpRoundTripper(cfg aws.Config) http.RoundTripper {
	return httpRoundtripper{
		config: cfg,
	}
}

func (h httpRoundtripper) RoundTrip(request *http.Request) (*http.Response, error) {
	credentials, err := h.config.Credentials.Retrieve(request.Context())
	if err != nil {
		return nil, RoundTripperError{err}
	}

	internalRequest := request.Clone(request.Context())

	hash, err := requestDataHash(request)
	if err != nil {
		return nil, RoundTripperError{err}
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
		return nil, RoundTripperError{err}
	}

	response, err := h.config.HTTPClient.Do(internalRequest)
	if err != nil {
		return nil, RoundTripperError{err}
	}

	if response.Header.Get("Content-Type") == "gzip" {
		gzipReader, err := gzip.NewReader(base64.NewDecoder(base64.StdEncoding, response.Body))
		if err != nil {
			return nil, RoundTripperError{err}
		}

		request.Header.Set("Content-Type", "application/json")

		response.Body = gzipReader
	}

	return response, nil
}

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

func getSha256(input []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(input); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
