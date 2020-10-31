package centrify

import (
	"errors"
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// PasswordProfile - Encapsulates a single Password Profile
type PasswordProfile struct {
	vaultObject
	ProfileFeature string `json:"ProfileFeature,omitempty" schema:"profile_feature,omitempty"`
	ProfileType    string `json:"ProfileType,omitempty" schema:"profile_type,omitempty"` // UserDefined, CheckPointGaia
	// password requirements
	MinimumPasswordLength              int    `json:"MinimumPasswordLength" schema:"minimum_password_length"`
	MaximumPasswordLength              int    `json:"MaximumPasswordLength" schema:"maximum_password_length"`
	AtLeastOneLowercase                bool   `json:"AtLeastOneLowercase" schema:"at_least_one_lowercase"`                                    // At least one lower-case alpha character
	AtLeastOneUppercase                bool   `json:"AtLeastOneUppercase" schema:"at_least_one_uppercase"`                                    // At least one upper-case alpha character
	AtLeastOneDigit                    bool   `json:"AtLeastOneDigit" schema:"at_least_one_digit"`                                            // At least one digit
	ConsecutiveCharRepeatAllowed       bool   `json:"ConsecutiveCharRepeatAllowed,omitempty" schema:"no_consecutive_repeated_char,omitempty"` // No consecutive repeated characters
	AtLeastOneSpecial                  bool   `json:"AtLeastOneSpecial" schema:"at_least_one_special_char"`                                   // At least one special character
	MaximumCharOccurrenceCount         int    `json:"MaximumCharOccurrenceCount,omitempty" schema:"maximum_char_occurrence_count,omitempty"`
	SpecialCharSet                     string `json:"SpecialCharSet,omitempty" schema:"special_charset,omitempty"`                                            // Special Characters
	FirstCharacterType                 string `json:"FirstCharacterType,omitempty" schema:"first_character_type,omitempty"`                                   // AlphaOnly or AlphaNumericOnly
	LastCharacterType                  string `json:"LastCharacterType,omitempty" schema:"last_character_type,omitempty"`                                     // AlphaOnly or AlphaNumericOnly
	MinimumAlphabeticCharacterCount    int    `json:"MinimumAlphabeticCharacterCount,omitempty" schema:"minimum_alphabetic_character_count,omitempty"`        // Min number of alpha characters
	MinimumNonAlphabeticCharacterCount int    `json:"MinimumNonAlphabeticCharacterCount,omitempty" schema:"minimum_non_alphabetic_character_count,omitempty"` // Min number of non-alpha characters
}

// NewPasswordProfile is a PasswordProfile constructor
func NewPasswordProfile(c *restapi.RestClient) *PasswordProfile {
	s := PasswordProfile{}
	s.client = c
	s.apiRead = "/ServerManage/GetPasswordProfiles"
	s.apiCreate = "/ServerManage/AddPasswordProfile"
	s.apiDelete = "/ServerManage/DeletePasswordProfile"
	s.apiUpdate = "/ServerManage/UpdatePasswordProfile"

	return &s
}

// Read function fetches an password profile from source, including attribute values. Returns error if any
func (o *PasswordProfile) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	args := make(map[string]interface{})
	args["Caching"] = -1
	queryArg["ProfileTypes"] = "All"
	queryArg["Args"] = args

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	// Loop through respond results and grab the matched record
	var results = resp.Result["Results"].([]interface{})
	//LogD.Printf("Total returned password profile: %d", len(results))
	// This is the matched list of password profile. There should be only one really
	var pwdpfs []keyValue
	//LogD.Printf("Looking for password profile: %s", o.ID)
	for _, v := range results {
		item := v.(map[string]interface{})
		row := item["Row"].(map[string]interface{})
		//LogD.Printf("Checking name: %s, profiletype: %s", row["Name"], row["ProfileType"])
		if row["ID"] == o.ID {
			LogD.Printf("Found an item: %+v", row)
			// If ProfileType is defined, then compare it
			if o.ProfileType == "" {
				pwdpfs = append(pwdpfs, row)
			} else if row["ProfileType"] == o.ProfileType {
				pwdpfs = append(pwdpfs, row)
			}
		}
	}
	if len(pwdpfs) == 0 {
		return errors.New("Query returns 0 object")
	}
	if len(pwdpfs) > 1 {
		return fmt.Errorf("Query returns too many objects (found %d, expected 1)", len(pwdpfs))
	}

	fillWithMap(o, pwdpfs[0])
	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Query function returns a single password profile object
func (o *PasswordProfile) Query() (map[string]interface{}, error) {
	var queryArg = make(map[string]interface{})
	args := make(map[string]interface{})
	args["Caching"] = -1
	queryArg["ProfileTypes"] = "All"
	queryArg["Args"] = args

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}
	// Loop through respond results and grab the matched record
	var results = resp.Result["Results"].([]interface{})

	//LogD.Printf("Total returned password profile: %d", len(results))
	// This is the matched list of password profile. There should be only one really
	var pwdpfs []keyValue
	//LogD.Printf("Looking for password profile: %s", o.Name)
	for _, v := range results {
		item := v.(map[string]interface{})
		row := item["Row"].(map[string]interface{})
		//LogD.Printf("Checking name: %s, profiletype: %s", row["Name"], row["ProfileType"])
		if row["Name"] == o.Name {
			LogD.Printf("Found an item: %+v", row)
			// If ProfileType is defined, then compare it
			if o.ProfileType == "" {
				pwdpfs = append(pwdpfs, row)
			} else if row["ProfileType"] == o.ProfileType {
				pwdpfs = append(pwdpfs, row)
			}
		}
	}
	if len(pwdpfs) == 0 {
		return nil, errors.New("Query returns 0 object")
	}
	if len(pwdpfs) > 1 {
		return nil, fmt.Errorf("Query returns too many objects (found %d, expected 1)", len(pwdpfs))
	}

	return pwdpfs[0], nil
}

// Delete function deletes an password profile and returns a map that contains deletion result
func (o *PasswordProfile) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Create function creates an password profile and returns a map that contains update result
func (o *PasswordProfile) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})

	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}

	LogD.Printf("Generated Map for Create(): %+v", queryArg)

	reply, err := o.client.CallStringAPI(o.apiCreate, queryArg)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Update function updates an existing password profile and returns a map that contains update result
func (o *PasswordProfile) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}

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

/*
	API to manage password profile

	Get password profile

		Request body format
		{
			"RRFormat": true,
			"ProfileTypes": "All",
			"Args": {
				"PageNumber": 1,
				"PageSize": 100000,
				"Limit": 100000,
				"SortBy": "",
				"direction": "False",
				"Caching": -1
			}
		}

		Respond result
		{
			"success": true,
			"Result": {
				"IsAggregate": false,
				"Count": 22,
				"Columns": [
					...}\
				],
				"FullCount": 22,
				"Results": [
					{
						"Entities": [
							{
								"Type": "PasswordProfile",
								"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
								"IsForeignKey": false
							}
						],
						"Row": {
							"LastCharacterType": "AnyChar",
							"_entitycontext": "*",
							"FirstCharacterType": "AlphaOnly",
							"SpecialCharSet": "!@#$%&()+,-./:;<=>?[\\]^_{|}~",
							"Name": "Check Point Gaia Profile",
							"ConsecutiveCharRepeatAllowed": false,
							"_metadata": {
								"Version": 1,
								"IndexingVersion": 1
							},
							"_PartitionKey": "centrify",
							"_Timestamp": "/Date(1596273224728)/",
							"AtLeastOneSpecial": true,
							"_encryptkeyid": "RootIndex:1",
							"AtLeastOneLowercase": true,
							"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"Description": "Default profile for Check Point Gaia systems",
							"MaximumPasswordLength": 32,
							"ProfileFeature": "Infrastructure",
							"AtLeastOneDigit": true,
							"_TableName": "pvpasswordprofile",
							"MinimumPasswordLength": 12,
							"AtLeastOneUppercase": true,
							"ProfileType": "CheckPointGaia"
						}
					},
				],
				"ReturnID": ""
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Create password profile

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Test PW Profile 1",
			"Description": "Test PW Profile 1",
			"MinimumPasswordLength": 12,
			"MaximumPasswordLength": 32,
			"AtLeastOneLowercase": true,
			"AtLeastOneUppercase": true,
			"AtLeastOneDigit": true,
			"ConsecutiveCharRepeatAllowed": true,
			"AtLeastOneSpecial": true,
			"jsutil-checkbox-8128-inputEl": true,
			"MaximumCharOccurrenceCount": 2,
			"SpecialCharSet": "!#$%&()*+,-./:;<=>?@[\\]^_{|}~",
			"jsutil-checkbox-8138-inputEl": true,
			"FirstCharacterType": "AlphaOnly",
			"jsutil-checkbox-8142-inputEl": true,
			"LastCharacterType": "AlphaNumericOnly",
			"MinimumAlphabeticCharacterCount": 1,
			"MinimumNonAlphabeticCharacterCount": 1
		}

		Respond result
		{
			"success": true,
			"Result": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update password profile

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Test PW Profile",
			"Description": "testa afasdf jhjkhjkh",
			"MinimumPasswordLength": 12,
			"MaximumPasswordLength": 32,
			"AtLeastOneLowercase": true,
			"AtLeastOneUppercase": true,
			"AtLeastOneDigit": true,
			"ConsecutiveCharRepeatAllowed": true,
			"AtLeastOneSpecial": true,
			"jsutil-checkbox-4100-inputEl": false,
			"SpecialCharSet": "!#$%&()*+,-./:;<=>?@[\\]^_{|}~",
			"jsutil-checkbox-4110-inputEl": false,
			"jsutil-checkbox-4114-inputEl": false,
			"LastCharacterType": "AnyChar",
			"FirstCharacterType": "AnyChar"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete password profile

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": true,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

*/
