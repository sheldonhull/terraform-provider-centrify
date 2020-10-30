package oauth

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/centrify/terraform-provider/cloud-golang-sdk/util"
)

// GetClient creates REST client
func (c *OauthClient) GetClient() (*restapi.RestClient, error) {
	// Use an oauth client to get our bearer token, currently always via confidential client flow
	token, err := c.GetOauthToken()
	if err != nil {
		return nil, err
	}

	// Then get rest client and set it up to use our token
	restClient, err := c.GetRestClient(token)
	if err != nil {
		return nil, err
	}

	return restClient, nil
}

// GetOauthToken obtains OAuth token string
func (c *OauthClient) GetOauthToken() (*TokenResponse, error) {
	//oclient, err := oauth.GetNewConfidentialClient(c.URL, c.Username, c.Password, nil)
	var clientFactory HttpClientFactory = func() *http.Client {
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
	oclient, err := GetNewConfidentialClient(c.Service, c.ClientID, c.ClientSecret, clientFactory)

	if err != nil {
		return nil, fmt.Errorf("Failed to get confidential client: %v", err)
	}
	oclient.SourceHeader = restapi.SourceHeader
	token, failure, err := oclient.ClientCredentials(c.AppID, c.Scope)

	if err != nil {
		return nil, fmt.Errorf("Failed to get confidential client token: %v", err)
	}

	if failure != nil {
		return nil, fmt.Errorf("Failed to get oauth token, failure: %v", failure)
	}

	util.LogD.Printf("Client token established - type: %s expires in: %d", token.TokenType, token.ExpiresIn)
	return token, nil
}

// GetRestClient returns rest client directly with oauth token
func (c *OauthClient) GetRestClient(token *TokenResponse) (*restapi.RestClient, error) {
	//restClient, err := restapi.GetNewRestClient(c.URL, nil)
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
	restClient, err := restapi.GetNewRestClient(c.Service, clientFactory)
	if err != nil {
		return nil, err
	}

	restClient.Headers["Authorization"] = token.TokenType + " " + token.AccessToken
	return restClient, nil
}
