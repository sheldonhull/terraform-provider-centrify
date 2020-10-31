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
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
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
				"_encryptkeyid": "XXXXXX",
				"_TableName": "datavault",
				"_Timestamp": "/Date(1584413116338)/",
				"WhenContentsReplaced": "/Date(1584413116309)/",
				"ACL": "true",
				"_PartitionKey": "XXXXXX",
				"WhenCreated": "/Date(1582558666855)/",
				"_entitycontext": "W/\"datetime'2020-03-17T02%3A45%3A16.3380444Z'\"",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"WhenUpdated": "/Date(1584413116309)/",
				"ParentPath": "Folder 1\\Folder level 2",
        		"FolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Description": "admin@example.com",
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
			"SecretText": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Type": "Text",
			"SetID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": false
		}
		or
		{
			"FolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"SecretName": "Another secret",
			"Description": "Another secret",
			"SecretText": "xxxxxxxxxxxxx",
			"Type": "Text",
			"updateChallenges": false
		}
		or
		{
			"FolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxc",
			"SecretName": "File1",
			"SecretFilePath": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"SecretFileSize": "38.003 KB",
			"SecretFilePassword": "xxxxxxx",
			"Type": "File",
			"Description": "",
			"updateChallenges": false
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

	Update Secret
	https://developer.centrify.com/reference#post_servermanage-updatesecret

		Request body format
		{
			"SecretName": "Access key",
			"Description": "AWS access key",
			"SecretText": "xxxxxxxxxxxxx",
			"Type": "Text",
			"SetID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": true,
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"DataVaultDefaultProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}
		or
		{
			"SecretText": "xxxxxxxxxxxxx",
			"SecretName": "Random secret",
			"Type": "Text",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Delete Secret
	https://developer.centrify.com/reference#post_servermanage-deletesecret

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

	Retrieve Secret content
	https://developer.centrify.com/reference#post_servermanage-retrievesecretcontents

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_encryptkeyid": "XXXXXX",
				"_TableName": "datavault",
				"_Timestamp": "/Date(1592380339832)/",
				"ACL": "true",
				"_PartitionKey": "XXXXXX",
				"WhenCreated": "/Date(1592380339057)/",
				"_entitycontext": "W/\"datetime'2020-06-17T07%3A52%3A19.8321511Z'\"",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
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
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"targetFolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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
