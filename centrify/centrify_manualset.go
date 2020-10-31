package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// ManualSet - Encapsulates a single Generic ManualSet
type ManualSet struct {
	vaultObject
	apiUpdateMembers     string
	apiMemberPermissions string

	ObjectType        string `json:"ObjectType,omitempty" schema:"type,omitempty"`
	SubObjectType     string `json:"SubObjectType,omitempty" schema:"subtype,omitempty"`
	CollectionType    string `json:"CollectionType,omitempty" schema:"collection_type,omitempty"`
	MemberPermissions []Permission
}

// setMember is a simply struct used for performing Set related API call
type setMember struct {
	MemberType string
	Table      string
	Key        string
}

// NewManualSet is a ManualSet constructor
func NewManualSet(c *restapi.RestClient) *ManualSet {
	s := ManualSet{}
	s.CollectionType = "ManualBucket"
	s.client = c
	s.apiRead = "/Collection/GetCollection"
	s.apiCreate = "/Collection/CreateManualCollection"
	s.apiDelete = "/Collection/DeleteCollection"
	s.apiUpdate = "/Collection/UpdateCollection"
	s.apiUpdateMembers = "/Collection/UpdateMembersCollection"
	s.apiPermissions = "/Collection/SetCollectionPermissions"

	return &s
}

// Read function fetches a ManualSet from source, including attribute values. Returns error if any
func (o *ManualSet) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)

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

// Create function creates a new Manual Set and returns a map that contains creation result
func (o *ManualSet) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg["Name"] = o.Name
	queryArg["ObjectType"] = o.ObjectType
	queryArg["SubObjectType"] = o.SubObjectType
	queryArg["CollectionType"] = o.CollectionType
	if o.Description != "" {
		queryArg["Description"] = o.Description
	}

	reply, err := o.client.CallStringAPI(o.apiCreate, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Delete function deletes a Manual Set and returns a map that contains deletion result
func (o *ManualSet) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing Manual Set and returns a map that contains update result
func (o *ManualSet) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	if o.Name != "" {
		queryArg["Name"] = o.Name
	}
	if o.Description != "" {
		queryArg["Description"] = o.Description
	}

	reply, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Query function returns a single Set object in map format
func (o *ManualSet) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Sets WHERE 1=1"
	if o.ObjectType != "" {
		query += " AND ObjectType='" + o.ObjectType + "'"
	}
	if o.CollectionType != "" {
		query += " AND CollectionType='" + o.CollectionType + "'"
	}
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// UpdateSetMembers adds or removes members from the ManualSet
func (o *ManualSet) UpdateSetMembers(ids []string, action string) (*restapi.StringResponse, error) {
	var reply *restapi.StringResponse
	var err error
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	if action != "add" && action != "remove" {
		return nil, errors.New("error: Update Set action must be either 'add' or 'remove'")
	}
	if o.ObjectType == "" {
		return nil, errors.New("error: Set object type is empty")
	}

	for _, v := range ids {
		var members []setMember
		var queryArg = make(map[string]interface{})
		queryArg["id"] = o.ID // This is the ID of Set
		setData := setMember{
			Key:        v,
			MemberType: "Row",
			Table:      o.ObjectType,
		}
		members = append(members, setData)
		queryArg[action] = members
		reply, err = o.client.CallStringAPI(o.apiUpdateMembers, queryArg)
		if err != nil {
			return nil, err
		}
		if !reply.Success {
			return nil, errors.New(reply.Message)
		}
	}

	return reply, nil
}

// SetMemberPermissions sets member permissions. isRemove indicates whether to remove all permissions instead of setting permissions
func (o *ManualSet) SetMemberPermissions(isRemove bool) (*restapi.GenericMapResponse, error) {
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
		queryArg["RowKey"] = o.ID
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

func (o *ManualSet) getMemberPerms() map[string]string {
	memberPerms := setPermissions
	switch o.ObjectType {
	case "Server":
		memberPerms = systemPermissions
		o.apiMemberPermissions = "/ServerManage/SetResourceCollectionPermissions"
	case "VaultDatabase":
		memberPerms = databasePermissions
		o.apiMemberPermissions = "/ServerManage/SetDatabaseCollectionPermissions"
	case "VaultDomain":
		memberPerms = domainPermissions
		o.apiMemberPermissions = "/ServerManage/SetDomainCollectionPermissions"
	case "VaultAccount":
		memberPerms = accountPermissions
		o.apiMemberPermissions = "/ServerManage/SetAccountCollectionPermissions"
	case "DataVault":
		memberPerms = secretPermissions
		o.apiMemberPermissions = "/ServerManage/SetSecretCollectionPermissions"
	case "SshKeys":
		memberPerms = sshkeyPermissions
		o.apiMemberPermissions = "/ServerManage/SetSSHKeyCollectionPermissions"
	case "Subscriptions":
		memberPerms = servicePermissions
		o.apiMemberPermissions = "/SetSubscriptionCollectionPermissions"
	case "Application":
		memberPerms = appPermissions
		o.apiMemberPermissions = "/SaasManage/SetApplicationCollectionPermissions"
	case "ResourceProfiles":
		memberPerms = setPermissions
		o.apiMemberPermissions = "/ResourceProfile/SetResourceProfileCollectionPermissions"
	}

	return memberPerms
}

/*
	API to manage set

	Get Set
	https://developer.centrify.com/reference#post_collection-getcollection

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_encryptkeyid": "XXXXX",
				"_TableName": "collections",
				"_Timestamp": "/Date(1569301295080)/",
				"ObjectType": "Server",
				"ACL": "true",
				"Name": "LAB Systems",
				"Filters": ...
				"_PartitionKey": "XXXXX",
				"WhenCreated": "/Date(1569301294975)/",
				"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"_entitycontext": "W/\"datetime'2019-09-24T05%3A01%3A35.0809483Z'\"",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"ParentPath": "",
				"CollectionType": "ManualBucket",
				"MembersFile": "/sys/buckets/exxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx.json",
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

	Create Set
	https://developer.centrify.com/reference#post_collection-createmanualcollection

		Request body format
		{
			"ObjectType": "Server",
			"Name": "Test Set 1",
			"Description": "Test Set 1",
			"CollectionType": "ManualBucket",
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

	Update Set
	https://developer.centrify.com/reference#post_collection-updatecollection

		Request body format
		{
			"Rights": "View, Edit, Delete, Grant",
			"_TableName": "collections",
			"_encryptkeyid": "XXXXX",
			"ObjectType": "Server",
			"ACL": "true",
			"Name": "Test Set 1",
			"_PartitionKey": "XXXXX",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"_entitycontext": "W/\"datetime'2020-08-01T09%3A49%3A55.7338109Z'\"",
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"CollectionType": "ManualBucket",
			"Description": "Test Set 2",
			"MembersFile": "/sys/buckets/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx.json",
			"_metadata": {
				"Version": 1,
				"IndexingVersion": 1
			},
			"sql": "undefined"
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

	Delete Set
	https://developer.centrify.com/reference#post_collection-deletecollection

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Updates set members
	https://developer.centrify.com/reference#post_collection-updatememberscollection
	Request body format
	{
		"id": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"add": [
			{
				"MemberType": "Row",
				"Table": "Server",
				"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			}
		]
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

	Set Collection Permissions
	https://developer.centrify.com/reference#post_collection-setcollectionpermissions

		Request body format
		{
			"Grants": [
				{
					"Principal": "LAB Infrastructure Admins",
					"PType": "Role",
					"Rights": "Grant,View,Edit,Delete",
					"PrincipalId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
				}
			],
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"PVID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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
