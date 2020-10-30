package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// VaultSecret - Encapsulates a single generic secret
type VaultSecret struct {
	vaultObject
	// VaultData specific APIs
	apiRetrieveSecret string
	apiMoveSecret     string
	apiGetChallenge   string

	SecretName              string          `json:"SecretName,omitempty" schema:"secret_name,omitempty"` // User Name
	SecretText              string          `json:"SecretText,omitempty" schema:"secret_text,omitempty"`
	Type                    string          `json:"Type,omitempty" schema:"type,omitempty"`
	FolderID                string          `json:"FolderId,omitempty" schema:"folder_id,omitempty"`
	ParentPath              string          `json:"ParentPath,omitempty" schema:"parent_path,omitempty"`
	DataVaultDefaultProfile string          `json:"DataVaultDefaultProfile" schema:"default_profile_id"` // Default Secret Challenge Profile (used if no conditions matched)
	ChallengeRules          *ChallengeRules `json:"DataVaultRules,omitempty" schema:"challenge_rule,omitempty"`
	Sets                    []string        `json:"Sets,omitempty" schema:"sets,omitempty"`
}

// NewVaultSecret is a VaultSecret constructor
func NewVaultSecret(c *restapi.RestClient) *VaultSecret {
	s := VaultSecret{}
	s.client = c
	s.apiRead = "/ServerManage/GetSecret"
	s.apiCreate = "/ServerManage/AddSecret"
	s.apiDelete = "/ServerManage/DeleteSecret"
	s.apiUpdate = "/ServerManage/UpdateSecret"
	s.apiRetrieveSecret = "/ServerManage/RetrieveSecretContents"
	s.apiMoveSecret = "/ServerManage/MoveSecret"
	s.apiPermissions = "/ServerManage/SetSecretPermissions"
	s.apiGetChallenge = "/ServerManage/GetSecretRightsAndChallenges"

	return &s
}

// Read function fetches a VaultSecret from source, including attribute values. Returns error if any
func (o *VaultSecret) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	LogD.Printf("Response for VaultSecret from tenant: %v", resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New(resp.Message)
	}

	fillWithMap(o, resp.Result)

	// Get challenge profile information
	resp, err = o.client.CallGenericMapAPI(o.apiGetChallenge, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	if v, ok := resp.Result["DataVaultDefaultProfile"]; ok {
		o.DataVaultDefaultProfile = v.(string)
	}

	// Fill challenge rules
	if v, ok := resp.Result["Challenges"]; ok {
		challenges := v.(map[string]interface{})
		if challenges["DataVaultDefaultProfile"] != nil {
			o.DataVaultDefaultProfile = challenges["DataVaultDefaultProfile"].(string)
		}
		if r, ok := challenges["DataVaultRules"]; ok {
			challengerules := &ChallengeRules{}
			fillWithMap(challengerules, r.(map[string]interface{}))
			o.ChallengeRules = challengerules
		}
	}

	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Create function creates a new VaultSecret and returns a map that contains creation result
func (o *VaultSecret) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	queryArg["updateChallenges"] = false

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

// Delete function deletes a VaultSecret and returns a map that contains deletion result
func (o *VaultSecret) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing VaultSecret and returns a map that contains update result
func (o *VaultSecret) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	queryArg["updateChallenges"] = true

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

// MoveSecret function moves an existing VaultSecret to another folder
func (o *VaultSecret) MoveSecret() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["targetFolderId"] = o.FolderID
	//queryArg["updateChallenges"] = true

	LogD.Printf("Generated Map for MoveFolder(): %+v", queryArg)

	reply, err := o.client.CallBoolAPI(o.apiMoveSecret, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Query function returns a single VaultSecret object in map format
func (o *VaultSecret) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM DataVault WHERE 1=1"
	if o.SecretName != "" {
		query += " AND SecretName='" + o.SecretName + "'"
	}
	if o.ParentPath != "" {
		query += " AND ParentPath='" + o.ParentPath + "'"
	}

	return queryVaultObject(o.client, query)
}

// CheckoutSecret checks out account secret from vault
func (o *VaultSecret) CheckoutSecret() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Description"] = "Checkout by Terraform provider"

	reply, err := o.client.CallGenericMapAPI(o.apiRetrieveSecret, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

/*
	API to manage vault secret

	Read Secret
	https://developer.centrify.com/reference#post_servermanage-getsecret

		Request body format
		{
			"ID": "c79f49e7-77c2-4ea1-b84e-fcd49a01d464",
			"RRFormat": true,
			"Args": {
				"PageNumber": 1,
				"Limit": 1,
				"PageSize": 1,
				"Caching": -1
			}
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_encryptkeyid": "ABC0751",
				"_TableName": "datavault",
				"_Timestamp": "/Date(1584413116338)/",
				"WhenContentsReplaced": "/Date(1584413116309)/",
				"ACL": "true",
				"_PartitionKey": "ABC0751",
				"WhenCreated": "/Date(1582558666855)/",
				"_entitycontext": "W/\"datetime'2020-03-17T02%3A45%3A16.3380444Z'\"",
				"_RowKey": "c79f49e7-77c2-4ea1-b84e-fcd49a01d464",
				"WhenUpdated": "/Date(1584413116309)/",
				"ParentPath": "Folder 1\\Folder level 2",
        		"FolderId": "7bded2b8-b481-4302-b2bd-f0a93375953c",
				"Description": "mspadmin@centrify.com.207",
				"SecretName": "Centrify PAS Admin Credential",
				"Type": "Text",
				"_metadata": {
					"Version": 1,
					"IndexingVersion": 1
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

	Add Secret
	https://developer.centrify.com/reference#post_servermanage-addsecret

		Request body format
		{
			"SecretName": "Access key",
			"Description": "AWS access key",
			"SecretText": "987489jkhjkahfdksa980242",
			"Type": "Text",
			"SetID": "4b6caf44-71af-4939-af2e-e9d176e062f4",
			"updateChallenges": false
		}
		or
		{
			"FolderId": "77cc181e-dbab-4662-9be4-67c49d3becf5",
			"SecretName": "Another secret",
			"Description": "Another secret",
			"SecretText": "asfdafsd fass",
			"Type": "Text",
			"updateChallenges": false
		}
		or
		{
			"FolderId": "7bded2b8-b481-4302-b2bd-f0a93375953c",
			"SecretName": "File1",
			"SecretFilePath": "8f051e7f-40c7-41a9-ae9e-e0d4b240211d",
			"SecretFileSize": "38.003 KB",
			"SecretFilePassword": "abc",
			"Type": "File",
			"Description": "",
			"updateChallenges": false
		}

		Respond result
		{
			"success": true,
			"Result": "85cd59ae-0024-456d-97c3-3236e26feb0c",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update Secret
	https://developer.centrify.com/reference#post_servermanage-updatesecret

		Request body format
		{
			"SecretName": "Access key",
			"Description": "AWS access key",
			"SecretText": "987489jkhjkahfdksa980242",
			"Type": "Text",
			"SetID": "4b6caf44-71af-4939-af2e-e9d176e062f4",
			"updateChallenges": true,
			"ID": "85cd59ae-0024-456d-97c3-3236e26feb0c",
			"DataVaultDefaultProfile": "30804754-3b87-4862-a39e-0f042825a3a0"
		}
		or
		{
			"SecretText": "jklkajsldf09890",
			"SecretName": "Random secret",
			"Type": "Text",
			"ID": "361da762-d7da-4d30-9e16-b1c2f40366be"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"ID": "85cd59ae-0024-456d-97c3-3236e26feb0c"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete Secret
	https://developer.centrify.com/reference#post_servermanage-deletesecret

		Request body format
		{
			"ID": "85cd59ae-0024-456d-97c3-3236e26feb0c"
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

	Retrieve Secret content
	https://developer.centrify.com/reference#post_servermanage-retrievesecretcontents

		Request body format
		{
			"ID": "7ea14b7e-f049-469a-bd3f-cebd9e96c77b"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_encryptkeyid": "ABC3434",
				"_TableName": "datavault",
				"_Timestamp": "/Date(1592380339832)/",
				"ACL": "true",
				"_PartitionKey": "ABC3434",
				"WhenCreated": "/Date(1592380339057)/",
				"_entitycontext": "W/\"datetime'2020-06-17T07%3A52%3A19.8321511Z'\"",
				"_RowKey": "7ea14b7e-f049-469a-bd3f-cebd9e96c77b",
				"ParentPath": "",
				"Description": "A random secret",
				"SecretName": "Randon secret",
				"Type": "Text",
				"SecretText": "xxxxxxxxxxx",
				"_metadata": {
					"Version": 1,
					"IndexingVersion": 1
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

	Move Secret to another folder
	https://developer.centrify.com/reference#post_servermanage-movesecret

	Request body format
	{
		"ID": "361da762-d7da-4d30-9e16-b1c2f40366be",
		"targetFolderId": "77cc181e-dbab-4662-9be4-67c49d3becf5"
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
