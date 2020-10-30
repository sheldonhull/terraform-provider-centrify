package webcookie

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/centrify/terraform-provider/cloud-golang-sdk/util"
	"golang.org/x/crypto/ssh/terminal"
)

// WebCookie represents a stateful web cookie client
type WebCookie struct {
	restapi.RestClient
	ClientID       string
	ClientSecret   string
	SkipCertVerify bool
}

func (c *WebCookie) startAuthentication() (*AuthResponse, error) {
	method := "/Security/StartAuthentication"
	args := make(map[string]interface{})
	args["User"] = c.ClientID
	args["Version"] = "1.0"

	body, err := c.postAndGetBody(method, args)
	if err != nil {
		return nil, err
	}
	//util.LogD.Printf("body: %+v\n", string(body))

	response, err := NewAuthResponse(body)
	util.LogD.Printf("response: %+v\n", response)
	if err != nil {
		return nil, err
	}
	if response.Success == false {
		return nil, fmt.Errorf("Failed to initiate authentication: %+v", response.Message)
	}

	return response, nil

}

func (c *WebCookie) advanceAuthentication(authResp *AuthResponse) ([]AuthMechanism, error) {
	challenges := authResp.Result.Challenges

	var authMechs []AuthMechanism
	for i, challenge := range challenges {
		util.LogD.Printf("Challenge %d -->\n", i)
		mechanisms := challenge.Mechanisms
		var authMech AuthMechanism
		if len(mechanisms) > 1 {
			fmt.Print("\n\n")
			// Display mechanisms
			for j, mechanism := range mechanisms {
				displayNum := j + 1
				fmt.Printf("%d. %s\n", displayNum, mechanism.PromptSelectMech)
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Please choose an authentication mechanism: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n")
			choice, err := strconv.Atoi(input)
			choice = choice - 1
			if choice < 0 || choice > len(mechanisms)-1 || err != nil {
				return nil, fmt.Errorf("Invalid choice")
			}
			authMech = mechanisms[choice]
			util.LogD.Printf("selected mech: %+v\n", authMech)
		} else {
			// Only one mechanism so go ahead to ask for credential
			authMech = mechanisms[0]

		}
		// Enter credential
		switch authMech.Name {
		case "UP":
			fmt.Print("Enter Password: ")
		case "SQ":
			fmt.Printf("%s : ", authMech.Question)
		case "OATH", "SMS", "EMAIL", "PF":
			fmt.Print("Enter Verification Code: ")
		default:
			fmt.Print("Enter Credential: ")
		}
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password := strings.TrimSpace(string(bytePassword))

		authMech.Credential = password
		authMechs = append(authMechs, authMech)
	}

	return authMechs, nil
}

func (c *WebCookie) doAuthentication(authResp *AuthResponse, authMechs []AuthMechanism) (string, error) {
	method := "/Security/AdvanceAuthentication"
	args := make(map[string]interface{})
	args["TenantId"] = authResp.Result.TenantID
	args["SessionId"] = authResp.Result.SessionID

	if len(authMechs) > 1 {
		var ops []map[string]interface{}
		for _, authMech := range authMechs {
			subargs := make(map[string]interface{})
			subargs["MechanismId"] = authMech.MechanismID
			subargs["Action"] = "Answer"
			subargs["Answer"] = authMech.Credential
			ops = append(ops, subargs)
		}
		args["MultipleOperations"] = ops
	} else if len(authMechs) == 1 {
		authMech := authMechs[0]
		args["MechanismId"] = authMech.MechanismID
		args["Action"] = "Answer"
		args["Answer"] = authMech.Credential
	}
	util.LogD.Printf("Auth post args: %+v\n", args)
	httpresp, err := c.postAndGetResp(method, args)
	if err != nil {
		return "", err
	}
	if httpresp.StatusCode != 200 {
		return "", fmt.Errorf("Bad http status code %v", httpresp.StatusCode)
	}

	body, _ := ioutil.ReadAll(httpresp.Body)
	resp, err := NewAdvanceAuthResponse(body)
	if err != nil {
		return "", fmt.Errorf("Error process respond body: %v", err)
	}

	if !resp.Success {
		return "", fmt.Errorf("Authentication failed: %s", resp.Message)
	}
	// Get auth cookie
	cookie := httpresp.Cookies()
	util.LogD.Printf("Cookies: %+v\n", getCookieByName(cookie, ".ASPXAUTH"))
	return getCookieByName(cookie, ".ASPXAUTH"), nil
}

func (c *WebCookie) postAndGetResp(method string, args map[string]interface{}) (*http.Response, error) {
	service := strings.TrimSuffix(c.Service, "/")
	method = strings.TrimPrefix(method, "/")
	postdata, _ := json.Marshal(args)
	util.LogD.Printf("Post json: %+v", bytes.NewBuffer(postdata))
	postreq, err := http.NewRequest("POST", service+"/"+method, bytes.NewBuffer(postdata))

	if err != nil {
		return nil, err
	}

	postreq.Header.Add("Content-Type", "application/json")
	postreq.Header.Add("X-CENTRIFY-NATIVE-CLIENT", "Yes")
	postreq.Header.Add("X-CFY-SRC", restapi.SourceHeader)

	for k, v := range c.Headers {
		postreq.Header.Add(k, v)
	}

	httpresp, err := c.Client.Do(postreq)
	if err != nil {
		c.ResponseHeaders = nil
		return nil, err
	}

	// save response heasder
	c.ResponseHeaders = httpresp.Header

	return httpresp, nil
}

func (c *WebCookie) postAndGetBody(method string, args map[string]interface{}) ([]byte, error) {
	httpresp, err := c.postAndGetResp(method, args)
	if err != nil {
		return nil, err
	}
	defer httpresp.Body.Close()

	if httpresp.StatusCode == 200 {
		return ioutil.ReadAll(httpresp.Body)
	}

	body, _ := ioutil.ReadAll(httpresp.Body)
	return nil, fmt.Errorf("POST to %s failed with code %d, body: %s", method, httpresp.StatusCode, body)
}

func getCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}
