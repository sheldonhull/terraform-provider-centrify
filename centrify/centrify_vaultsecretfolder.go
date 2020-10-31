package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// VaultSecretFolder - Encapsulates a single generic secret folder
type VaultSecretFolder struct {
	vaultObject
	// VaultSecretFolder specific APIs
	apiGetChallenge      string
	apiMoveFolder        string
	apiMemberPermissions string

	Type                            string          `json:"Type,omitempty" schema:"type,omitempty"`        // Can only be Folder
	ParentID                        string          `json:"Parent,omitempty" schema:"parent_id,omitempty"` // ID of parent folder
	ParentPath                      string          `json:"ParentPath,omitempty" schema:"parent_path,omitempty"`
	CollectionMembersDefaultProfile string          `json:"CollectionMembersDefaultProfile" schema:"default_profile_id"` // Default Secret Challenge Profile (used if no conditions matched)
	ChallengeRules                  *ChallengeRules `json:"CollectionMembersRules,omitempty" schema:"challenge_rule,omitempty"`
	MemberPermissions               []Permission
}

// NewVaultSecretFolder is a VaultSecretFolder constructor
func NewVaultSecretFolder(c *restapi.RestClient) *VaultSecretFolder {
	s := VaultSecretFolder{}
	s.client = c
	s.apiRead = "/ServerManage/GetSecretFolder"
	s.apiCreate = "/ServerManage/AddSecretsFolder"
	s.apiDelete = "/ServerManage/DeleteSecretsFolder"
	s.apiUpdate = "/ServerManage/UpdateSecretsFolder"
	s.apiGetChallenge = "/ServerManage/GetSecretsFolderRightsAndChallenges"
	s.apiMoveFolder = "/ServerManage/MoveFolder"
	s.apiPermissions = "/ServerManage/SetSecretsFolderPermissions"
	s.apiMemberPermissions = "/ServerManage/SetSecretCollectionPermissions"
	s.Type = "Folder"

	return &s
}

// Read function fetches a VaultSecretFolder from source, including attribute values. Returns error if any
func (o *VaultSecretFolder) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	LogD.Printf("Response for VaultSecretFolder from tenant: %v", resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New(resp.Message)
	}

	// Loop through respond results and grab the first record
	var results = resp.Result["Results"].([]interface{})
	if len(results) < 1 {
		// Make sure error message contains "not exist"
		LogD.Printf("Returning error: VaultSecretFolder does not exist in tenant")
		return errors.New("VaultSecretFolder does not exist in tenant")
	} else if len(results) > 1 {
		// this should never happen
		return errors.New("There are more than one VaultSecretFolder with the same ID in tenant")
	}
	var result = results[0].(map[string]interface{})
	// Populate vaultObject struct with map from response
	var row = result["Row"].(map[string]interface{})

	fillWithMap(o, row)

	// Get challenge profile
	resp, err = o.client.CallGenericMapAPI(o.apiGetChallenge, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	//LogD.Printf("Challenges result: %+v", resp)
	if v, ok := resp.Result["Challenges"]; ok {
		challenges := v.(map[string]interface{})
		if challenges["CollectionMembersDefaultProfile"] != nil {
			o.CollectionMembersDefaultProfile = challenges["CollectionMembersDefaultProfile"].(string)
		}
		// Fill challenge rules
		if r, ok := challenges["CollectionMembersRules"]; ok {
			challengerules := &ChallengeRules{}
			fillWithMap(challengerules, r.(map[string]interface{}))
			o.ChallengeRules = challengerules
		}
	}

	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Create function creates a new VaultSecretFolder and returns a map that contains creation result
func (o *VaultSecretFolder) Create() (*restapi.StringResponse, error) {
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

// Delete function deletes a VaultSecretFolder and returns a map that contains deletion result
func (o *VaultSecretFolder) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing VaultSecretFolder and returns a map that contains update result
func (o *VaultSecretFolder) Update() (*restapi.GenericMapResponse, error) {
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

// MoveFolder function moves an existing VaultSecretFolder to another folder
func (o *VaultSecretFolder) MoveFolder() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["targetFolderId"] = o.ParentID
	//queryArg["updateChallenges"] = true

	LogD.Printf("Generated Map for MoveFolder(): %+v", queryArg)

	reply, err := o.client.CallBoolAPI(o.apiMoveFolder, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Query function returns a single VaultSecretFolder object in map format
func (o *VaultSecretFolder) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Sets WHERE 1=1 AND ObjectType='DataVault' AND CollectionType='Phantom'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.ParentPath != "" {
		query += " AND ParentPath='" + o.ParentPath + "'"
	}

	return queryVaultObject(o.client, query)
}

// SetMemberPermissions sets member permissions. isRemove indicates whether to remove all permissions instead of setting permissions
func (o *VaultSecretFolder) SetMemberPermissions(isRemove bool) (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}

	var permissions []map[string]interface{}
	for _, v := range o.MemberPermissions {
		var permission = make(map[string]interface{})
		permission, err := generateRequestMap(v)
		if isRemove {
			permission["Rights"] = "None"
		}
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if len(permissions) > 0 {
		var queryArg = make(map[string]interface{})
		queryArg["ID"] = o.ID
		queryArg["PVID"] = o.ID
		queryArg["Grants"] = permissions
		LogD.Printf("Generated Map for SetMemberPermissions(): %+v", queryArg)
		resp, err := o.client.CallGenericMapAPI(o.apiMemberPermissions, queryArg)
		if err != nil {
			return nil, err
		}
		if !resp.Success {
			return nil, errors.New(resp.Message)
		}
		return resp, nil
	}
	return nil, nil
}

/*
	API to manage vault secret

	Read folder
	https://developer.centrify.com/reference#post_servermanage-getsecretfolder

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RRFormat": true
		}

		Respond result
		{
			"success": true,
			"Result": {
				"IsAggregate": false,
				"Count": 1,
				"Columns": [...],
				"FullCount": 1,
				"Results": [
					{
						"Entities": [
							{
								"Type": "DataVaultFolder",
								"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
								"IsForeignKey": false
							}
						],
						"Row": {
							"_entitycontext": "W/\"datetime'2020-08-16T02%3A08%3A49.301252Z'\"",
							"Name": "Test Folder",
							"_metadata": {
								"Version": 1,
								"IndexingVersion": 1
							},
							"ACL": "true",
							"_PartitionKey": "XXXXXX",
							"WhenCreated": "/Date(1597543728366)/",
							"_Timestamp": "/Date(1597543728399)/",
							"_encryptkeyid": "XXXXXX",
							"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"Description": "Test Folder",
							"Type": "Folder",
							"ObjectType": "DataVault",
							"SecretName": "Test Folder",
							"CollectionType": "Phantom",
							"_TableName": "collections",
							"Filters": "FolderId",
							"ParentPath": ""
						}
					}
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

	Add folder
	https://developer.centrify.com/reference#post_servermanage-addsecretsfolder

		Request body format
		{
			"Name": "Test Folder",
			"Description": "Test Folder",
			"Type": "Folder",
			"updateChallenges": false
		}
		or
		{
			"Parent": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Sub folder",
			"Description": "Sub folder",
			"Type": "Folder",
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

	Update folder
	https://developer.centrify.com/reference#post_servermanage-updatesecretsfolder

		Request body format
		{
			"_entitycontext": "W/\"datetime'2020-08-16T02%3A17%3A58.6860435Z'\"",
			"Name": "Sub folder",
			"_metadata": {
				"Version": 1,
				"IndexingVersion": 1
			},
			"ACL": "true",
			"_PartitionKey": "XXXXXX",
			"_encryptkeyid": "XXXXXX",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Description": "Sub folder of something",
			"Type": "Folder",
			"ObjectType": "DataVault",
			"SecretName": "Sub folder",
			"CollectionType": "Phantom",
			"_TableName": "collections",
			"Filters": "FolderId",
			"Parent": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"ParentPath": "Folder 1",
			"Rights": "View, Edit, Delete, Grant, Add",
			"CollectionMembersDefaultProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": true
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

	Delete folder
	https://developer.centrify.com/reference#post_servermanage-deletesecretsfolder

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": false,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}
		or
		{
			"success": false,
			"Result": {
				"ChallengeId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": "You must respond to a challenge to proceed.",
			"MessageID": null,
			"Exception": null,
			"ErrorID": "ChallengeRequired",
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Get Rights and Challenges

	Request body format
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": {
			"Rights": "View, Edit, Delete, Grant, Add",
			"Challenges": {
				"CollectionMembersRules": {
					"_UniqueKey": "Condition",
					"_Value": [],
					"Enabled": true,
					"_Type": "RowSet"
				},
				"CollectionMembersDefaultProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Move folder to another folder
	https://developer.centrify.com/reference#post_servermanage-movefolder

	Request body format
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"targetFolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": false,
		"Result": {
			"ChallengeId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		},
		"Message": "You must respond to a challenge to proceed.",
		"MessageID": null,
		"Exception": null,
		"ErrorID": "ChallengeRequired",
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}
	or
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
