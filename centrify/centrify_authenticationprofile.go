package centrify

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// AuthenticationProfile - Encapsulates a single Authentication Profile
type AuthenticationProfile struct {
	vaultObject
	UUID              string          `json:"Uuid,omitempty" schema:"uuid,omitempty"`
	DurationInMinutes int             `json:"DurationInMinutes" schema:"pass_through_duration"` // Challenge Pass-Through Duration. Can't omitempty because 0 mean no pass-through
	Challenges        []string        `json:"Challenges,omitempty" schema:"challenges,omitempty"`
	AdditionalData    *AdditionalData `json:"AdditionalData,omitempty" schema:"additional_data,omitempty"`
}

// AdditionalData for AuthenticationProfile
type AdditionalData struct {
	NumberOfQuestions int `json:"NumberOfQuestions,omitempty" schema:"number_of_questions,omitempty"` // Number of questions user must answer
}

// NewAuthenticationProfile is a AuthenticationProfile constructor
func NewAuthenticationProfile(c *restapi.RestClient) *AuthenticationProfile {
	s := AuthenticationProfile{}
	s.client = c
	s.apiRead = "/AuthProfile/GetProfile"
	s.apiCreate = "/AuthProfile/SaveProfile"
	s.apiDelete = "/AuthProfile/DeleteProfile"
	s.apiUpdate = "/AuthProfile/SaveProfile"

	return &s
}

type sliceAPIResponse struct {
	Success bool `json:"success"`
	Result  []interface{}
	Message string
}

// Read function fetches an authentication profile from source, including attribute values. Returns error if any
func (o *AuthenticationProfile) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["Uuid"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	//LogD.Printf("Response for authentication profile from tenant: %v", resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}

	fillWithMap(o, resp.Result)
	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Delete function deletes an authentication profile and returns a map that contains deletion result
func (o *AuthenticationProfile) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("uuid")
}

// Create function creates an authentication profile and returns a map that contains update result
func (o *AuthenticationProfile) Create() (*restapi.GenericMapResponse, error) {
	var queryArg = make(map[string]interface{})

	settings, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	queryArg["settings"] = settings

	LogD.Printf("Generated Map for Create(): %+v", queryArg)

	reply, err := o.client.CallGenericMapAPI(o.apiCreate, queryArg)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Update function updates an existing authentication profile and returns a map that contains update result
func (o *AuthenticationProfile) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	settings, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	queryArg["settings"] = settings

	LogD.Printf("Generated Map for Update(): %+v", queryArg)

	reply, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Query function returns a single authentication profile object
func (o *AuthenticationProfile) Query() (map[string]interface{}, error) {
	var queryArg = make(map[string]interface{})
	args := make(map[string]interface{})
	args["Caching"] = -1
	queryArg["Args"] = args

	// Attempt to read from an upstream API
	resp, err := o.client.CallRawAPI("/AuthProfile/GetProfileList", queryArg)
	if err != nil {
		return nil, err
	}

	reply := &sliceAPIResponse{}
	err = json.Unmarshal(resp, &reply)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal sliceAPIResponse from HTTP response: %v", err)
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	//LogD.Printf("Total returned authentication profile: %d", len(reply.Result))
	// This is the matched list of password profile. There should be only one really
	var autheProfs []keyValue
	//LogD.Printf("Looking for authentication profile: %s", o.Name)
	for _, v := range reply.Result {
		item := v.(map[string]interface{})
		//LogD.Printf("Checking name: %s", item["Name"])
		if item["Name"] == o.Name {
			//LogD.Printf("Found an item: %+v", item)
			autheProfs = append(autheProfs, item)
		}
	}
	if len(autheProfs) == 0 {
		return nil, errors.New("Query returns 0 object")
	}
	if len(autheProfs) > 1 {
		return nil, fmt.Errorf("Query returns too many objects (found %d, expected 1)", len(autheProfs))
	}

	return autheProfs[0], nil
}

/*
	API to manage password profile

	Get authentication profile

		Request body format
		{
			"uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"Uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Name": "LAB Step-up Authentication Profile",
				"DurationInMinutes": 0,
				"Challenges": [
					"OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
				],
				"AdditionalData": {
					"NumberOfQuestions": 1
				}
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}


	Create authentication profile

		Request body format
		{
			"settings": {
				"Name": "test",
				"Challenges": [
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ",
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
				],
				"DurationInMinutes": 30,
				"AdditionalData": {
					"NumberOfQuestions": 1
				}
			}
		}

		Respond result
		{
			"success": true,
			"Result": {
				"Uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Name": "test",
				"DurationInMinutes": 30,
				"Challenges": [
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ",
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
				],
				"AdditionalData": {
					"NumberOfQuestions": 1
				}
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update authentication profile

		Request body format


		Respond result

	Delete authentication profile

		Request body format
		{
			"uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": null,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

*/
