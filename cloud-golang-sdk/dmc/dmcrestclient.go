package dmc

import (
	"crypto/tls"
	"net/http"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// DMC represents a stateful dmc client
type DMC struct {
	restapi.RestClient
	Scope          string // Delegated Machine Credential scope definition
	Token          string // DMC Oauth token. If this is provided, then no need to make LRPC call
	SkipCertVerify bool
}

// GetClient creates REST client
func (c *DMC) GetClient() (*restapi.RestClient, error) {
	var clientFactory restapi.HttpClientFactory = func() *http.Client {
		return &http.Client{}
	}
	if c.SkipCertVerify {
		// Ignore certificate error for on-prem deployment
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		clientFactory = func() *http.Client {
			return &http.Client{Transport: tr}
		}
	}

	if c.Token == "" {
		rpc := NewLRPC2()
		token, err := rpc.GetToken(c.Scope)
		if err != nil {
			return nil, err
		}
		c.Token = token
		//LogD.Printf("dmc token: %v\n", c.Token)
	}

	restClient, err := restapi.GetNewRestClient(c.Service, clientFactory)
	if err != nil {
		return nil, err
	}

	//restClient.SourceHeader = restapi.SourceHeader
	restClient.Headers["Authorization"] = "Bearer " + c.Token
	return restClient, nil
}
