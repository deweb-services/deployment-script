package dws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/deweb-services/deployment-script/internal/types"
	"go.uber.org/zap"
)

type DWSClient struct {
	Config *types.DWSProviderConfiguration
	client *http.Client
	signer *Signer
	url    string
	log    *zap.SugaredLogger
}

type ClientOpt func(c *DWSClient)

func NewClient(ctx context.Context, configuration *types.DWSProviderConfiguration, log *zap.SugaredLogger, opts ...ClientOpt) *DWSClient {
	signerOpts := make([]CredentialsOpt, 0)

	fmt.Println(configuration)
	if configuration.AccessKey != "" && configuration.SecretAccessKey != "" {
		signerOpts = append(signerOpts, WithStaticCredentials(configuration.AccessKey, configuration.SecretAccessKey))
	}

	if len(signerOpts) == 0 {
		signerOpts = append(signerOpts, WithAnonymousCredentials())
	}

	signer := NewSigner(signerOpts[len(signerOpts)-1], WithDebugLogger(&types.DebugLogger{Context: ctx, SugaredLogger: log}))

	c := &DWSClient{
		Config: configuration,
		client: &http.Client{},
		signer: signer,
		url:    types.APIURL,
		log:    log,
	}

	if configuration.Timeout == 0 {
		c.client.Timeout = types.DefaultTimeoutSeconds * time.Second
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithStaticCredentials(accessKey, secretKey string) CredentialsOpt {
	return func() *credentials.Credentials {
		return credentials.NewStaticCredentials(accessKey, secretKey, "")
	}
}

func WithAnonymousCredentials() CredentialsOpt {
	return func() *credentials.Credentials {
		return credentials.AnonymousCredentials
	}
}

func WithDebugLogger(logger aws.Logger) SignerOpt {
	return func(signer *Signer) {
		signer.v4.Debug = aws.LogDebugWithSigning
		signer.v4.Logger = logger
	}
}

func ClientOptWithURL(url string) ClientOpt {
	return func(c *DWSClient) {
		c.url = url
	}
}

func (c *DWSClient) CreateGPU(ctx context.Context, gpu *types.GPUCreateConfig) (*types.GPUCreateConfigResponse, error) {
	errPrefix := "failed to create gpu: %w"
	b, err := json.Marshal(gpu)
	if err != nil {
		return nil, fmt.Errorf(errPrefix, err)
	}

	responseBody, err := c.DoSignedRequest(ctx, http.MethodPost, c.url+types.GPUEndpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf(errPrefix, err)
	}

	resp := types.GPUCreateConfigResponse{}
	if err = json.Unmarshal(responseBody, &resp); err != nil {
		return nil, fmt.Errorf(errPrefix, err)
	}

	return &resp, nil
}

func (c *DWSClient) DeleteGPU(ctx context.Context, id string) error {
	c.log.Debugf("Delete GPU by id: %s", id)

	_, err := c.DoSignedRequest(ctx, http.MethodDelete, fmt.Sprintf(c.url+types.GPUEndpoint+"/%s", id), nil)
	if err != nil {
		return fmt.Errorf("failed to delete GPU: %w", err)
	}

	return nil
}

func (c *DWSClient) GetGPU(ctx context.Context, id string) (*types.RentedGpuInfoResponse, error) {
	c.log.Debugf("Get GPU by id: %s", id)

	responseBody, err := c.DoSignedRequest(ctx, http.MethodGet, fmt.Sprintf(c.url+types.GPUEndpoint+"/%s", id), nil)
	c.log.Debugf("Get GPU responseBody: %s", string(responseBody))
	if err != nil {
		return nil, fmt.Errorf("failed to get GPU: %w", err)
	}

	gpu := types.RentedGpuInfoResponse{}
	err = json.Unmarshal(responseBody, &gpu)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal get deployment response body: %w", err)
	}

	return &gpu, nil
}

func (c *DWSClient) DoSignedRequest(ctx context.Context, method string, endpoint string, body io.ReadSeeker) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c.log.Debug("signing request")
	if err = c.signer.SignRequest(req, body); err != nil {
		return nil, err
	}

	return c.DoRequest(req)
}

func (c *DWSClient) DoRequest(req *http.Request) ([]byte, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	c.log.Debugf("ERROR: %v, SIGN: %s", err, string(b))
	err = checkResponse(res)
	if err != nil {
		return nil, fmt.Errorf("external API returned an error code: %w", err)
	}

	return b, nil
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 400 && res.StatusCode <= 599 {
		return fmt.Errorf("request failed, status code: %d", res.StatusCode)
	}

	return nil
}
