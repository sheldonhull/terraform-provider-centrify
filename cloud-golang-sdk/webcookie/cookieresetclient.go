package webcookie

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/centrify/terraform-provider/cloud-golang-sdk/util"
)

// GetClient creates REST client
func (c *WebCookie) GetClient() (*restapi.RestClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// Munge on the service a little bit, force it to have no trailing / and always start with https://
	url, err := url.Parse(c.Service)
	if err != nil {
		return nil, err
	}
	url.Scheme = "https"
	url.Path = ""
	c.Service = url.String()

	if c.SkipCertVerify {
		// Ignore certificate error for on-prem deployment
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.Client = &http.Client{Transport: tr}
	} else {
		c.Client = &http.Client{}
	}
	c.Client.Jar = jar

	util.LogD.Printf("Start authentication...\n")
	authResp, err := c.startAuthentication()
	if err != nil {
		return nil, err
	}

	if !authResp.Success {
		return nil, fmt.Errorf("Failed to start authentication")
	}

	authMechs, err := c.advanceAuthentication(authResp)
	if err != nil {
		return nil, err
	}

	token, err := c.doAuthentication(authResp, authMechs)
	if err != nil {
		return nil, err
	}
	util.LogD.Printf("Token: %v\n", token)

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

	restClient.Headers["Authorization"] = "Bearer " + token
	return restClient, nil
}
