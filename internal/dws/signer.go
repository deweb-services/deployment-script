package dws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Signer struct {
	v4 *v4.Signer
}

type SignerOpt func(signer *Signer)

type CredentialsOpt func() *credentials.Credentials

func NewSigner(credentialsOpt CredentialsOpt, opts ...SignerOpt) *Signer {
	cred := credentialsOpt()
	signer := &Signer{v4: v4.NewSigner(cred)}

	for _, opt := range opts {
		opt(signer)
	}

	return signer
}

func (s *Signer) SignRequest(req *http.Request, body io.ReadSeeker) error {
	// TODO: normal fix
	originalURL := req.URL
	req.URL, _ = url.Parse(strings.Replace(req.URL.String(), "/api", "", -1))
	if err := s.signRequest(req, body); err != nil {
		return err
	}
	req.URL = originalURL

	return nil
}

func (s *Signer) signRequest(req *http.Request, body io.ReadSeeker) error {
	_, err := s.v4.Sign(req, body, "terraform", "global", time.Now())
	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}
	return nil
}
