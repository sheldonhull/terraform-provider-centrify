package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// DesktopApp - Encapsulates a single Generic DesktopApp
type DesktopApp struct {
	vaultObject

	TemplateName             string            `json:"TemplateName,omitempty" schema:"template_name,omitempty"`
	DesktopAppRunHostID      string            `json:"DesktopAppRunHostId,omitempty" schema:"application_host_id,omitempty"`         // Application host
	DesktopAppRunAccountType string            `json:"DesktopAppRunAccountType,omitempty" schema:"login_credential_type,omitempty"`  // Host login credential type: ADCredential, SetByUser, AlternativeAccount, SharedAccount
	DesktopAppRunAccountID   string            `json:"DesktopAppRunAccountUuid,omitempty" schema:"application_account_id,omitempty"` // Host login credential account
	DesktopAppProgramName    string            `json:"DesktopAppProgramName,omitempty" schema:"application_alias,omitempty"`         // Application alias
	DesktopAppCmdline        string            `json:"DesktopAppCmdlineTemplate,omitempty" schema:"command_line,omitempty"`          // Command line
	DesktopAppParams         []DesktopAppParam `json:"DesktopAppParams,omitempty" schema:"command_parameter,omitempty"`
	DefaultAuthProfile       string            `json:"DefaultAuthProfile" schema:"default_profile_id"`
	ChallengeRules           *ChallengeRules   `json:"AuthRules,omitempty" schema:"challenge_rule,omitempty"`
	PolicyScript             string            `json:"PolicyScript,omitempty" schema:"policy_script,omitempty"` // Use script to specify authentication rules (configured rules are ignored)
	WorkflowEnabled          bool              `json:"WorkflowEnabled,omitempty" schema:"workflow_enabled,omitempty"`
}

// DesktopAppParam - desktop app command line parameters
type DesktopAppParam struct {
	ParamName      string `json:"ParamName,omitempty" schema:"name,omitempty"`
	ParamType      string `json:"ParamType,omitempty" schema:"type,omitempty"` // int, date, string, User, Role, Device, Server, VaultAccount, VaultDomain, VaultDatabase, Subscriptions, DataVault, SshKeys, system_profile
	ParamValue     string `json:"ParamValue,omitempty" schema:"value,omitempty"`
	TargetObjectID string `json:"TargetObjectId,omitempty" schema:"target_object_id,omitempty"`
}

// NewDesktopApp is a esktopApp constructor
func NewDesktopApp(c *restapi.RestClient) *DesktopApp {
	s := DesktopApp{}
	s.client = c
	s.apiRead = "/SaasManage/GetApplication"
	s.apiCreate = "/SaasManage/ImportAppFromTemplate"
	s.apiDelete = "/SaasManage/DeleteApplication"
	s.apiUpdate = "/SaasManage/UpdateApplicationDE"
	s.apiPermissions = "/SaasManage/SetApplicationPermissions"

	return &s
}

// Read function fetches a DesktopApp from source, including attribute values. Returns error if any
func (o *DesktopApp) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["_RowKey"] = o.ID

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

// Create function creates a new DesktopApp and returns a map that contains creation result
func (o *DesktopApp) Create() (*restapi.SliceResponse, error) {
	var queryArg = make(map[string]interface{})

	queryArg["ID"] = []string{o.TemplateName}
	LogD.Printf("Generated Map for Create(): %+v", queryArg)

	resp, err := o.client.CallSliceAPI(o.apiCreate, queryArg)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

// Update function updates an existing DesktopApp and returns a map that contains update result
func (o *DesktopApp) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	queryArg["_RowKey"] = o.ID

	LogD.Printf("Generated Map for Update(): %+v", queryArg)

	resp, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

// Delete function deletes a DesktopApp and returns a map that contains deletion result
func (o *DesktopApp) Delete() (*restapi.SliceResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["_RowKey"] = []string{o.ID}

	resp, err := o.client.CallSliceAPI(o.apiDelete, queryArg)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}
	return resp, nil
}

// Query function returns a single DesktopApp object in map format
func (o *DesktopApp) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Application WHERE 1=1 AND AppType='Desktop'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

/*
Fetch desktop app
https://developer.centrify.com/reference#post_saasmanage-getapplication

	Request body format
	{
		"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
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
			"IsTestApp": false,
			"Icon": "/vfs/Application/Icons/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"DesktopAppRunAccountUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"DisplayName": "AirWatch ONE UEM",
			"UseDefaultSigningCert": true,
			"_entitycontext": "W/\"datetime'2020-06-17T09%3A24%3A09.4903113Z'\"",
			"_TableName": "application",
			"Generic": true,
			"LocalizationMappings": [
				...
			],
			"State": "Active",
			"RegistrationLinkMessage": null,
			"DesktopAppCmdlineTemplate": "--ini=ini\\web_airwatch_webdriver.ini --username={login.Description} --password={login.SecretText}",
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"DesktopAppProgramName": "pas_desktopapp",
			"_encryptkeyid": "XXXXXXX",
			"RemoteDesktopHostName": "member2",
			"DesktopAppRunAccountType": "SharedAccount",
			"_PartitionKey": "XXXXXX",
			"RemoteDesktopUser": "shared_account (demo.lab)",
			"CertificateSubjectName": "CN=Centrify Customer Application Signing Certificate",
			"ParentDisplayName": null,
			"_metadata": {
				"Version": 1,
				"IndexingVersion": 1
			},
			"Description": "This template allows you to provide single sign-on to a custom desktop application.",
			"DesktopAppType": "CommandLine",
			"AuthRules": {
				"_UniqueKey": "Condition",
				"_Value": [],
				"Enabled": true,
				"_Type": "RowSet"
			},
			"AppType": "Desktop",
			"Name": "AirWatch ONE UEM",
			"Thumbprint": "XXXXXXXXXXXXXXXXX",
			"TemplateName": "GenericDesktopApplication",
			"Handler": "Server.Cloud;Centrify.Server.DesktopApp.GenericDesktopAppHandler",
			"DefaultAuthProfile": "AlwaysAllowed",
			"AppTypeDisplayName": "Desktop",
			"AuthChallengeDefinitionId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RegistrationMessage": null,
			"DesktopAppRunHostId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"_Timestamp": "/Date(1592385849490)/",
			"ProvCapable": false,
			"AdminTag": "Other",
			"ProvSettingEnabled": false,
			"ProvConfigured": false,
			"ACL": "true",
			"Category": "Other",
			"LocalizationEnabled": false,
			"DesktopAppParams": [
				{
					"_encryptkeyid": "XXXXXXX",
					"_TableName": "applicationparams",
					"TargetObjectId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
					"_Timestamp": "/Date(1592385849244)/",
					"ParamName": "login",
					"ApplicationId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
					"ParamValue": "AirWatch Workspace ONE UEM Login",
					"_PartitionKey": "XXXXXX",
					"_entitycontext": "*",
					"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
					"ParamType": "DataVault",
					"_metadata": {
						"Version": 1,
						"IndexingVersion": 1
					}
				}
			]
		},
		"IsSoftError": false
	}

Create desktop app
https://developer.centrify.com/reference#post_saasmanage-importappfromtemplate

	Request body format
	{
		"ID": [
			"GenericDesktopApplication"      |"Ssms"|"Toad"|"VpxClient"
		]
	}

	Respond result
	{
		"success": true,
		"Result": [
			{
				"success": true,
				"ID": "GenericDesktopApplication",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			}
		],
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Update desktop app
https://developer.centrify.com/reference#post_saasmanage-updateapplicationde

	Request body format
	{
		"LocalizationEnabled": false,
		"LocalizationMappings": [
			...
		],
		"Name": "AirWatch ONE UEM",
		"Description": "This template allows you to provide single sign-on to a custom desktop application.",
		"Icon": "/vfs/Application/Icons/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Handler": "Server.Cloud;Centrify.Server.DesktopApp.GenericDesktopAppHandler",
		"IconUri": "/vfs/Application/Icons/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"DesktopAppRunHostId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"DesktopAppRunAccountUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"DesktopAppAccountContainerId": "",
		"DesktopAppAccountUuid": "",
		"RemoteDesktopHostName": "member2",
		"DesktopAppRunAccountType": "SharedAccount",
		"RemoteDesktopUser": "shared_account (demo.lab)",
		"DesktopAppProgramName": "pas_desktopapp",
		"DesktopAppParams": [
			{
				"_encryptkeyid": "XXXXXXX",
				"_TableName": "applicationparams",
				"TargetObjectId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"_Timestamp": "2020-06-17T09:24:09.244Z",
				"ParamName": "login",
				"ApplicationId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"ParamValue": "AirWatch Workspace ONE UEM Login",
				"_PartitionKey": "XXXXXX",
				"_entitycontext": "*",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"ParamType": "DataVault",
				"_metadata": {
					"Version": 1,
					"IndexingVersion": 1
				}
			}
		],
		"DesktopAppCmdlineTemplate": "--ini=ini\\web_airwatch_webdriver.ini --username={login.Description} --password={login.SecretText}",
		"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": {
			"State": 0
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Delete desktop app
https://developer.centrify.com/reference#post_saasmanage-deleteapplication

	Request body format
	{
		"_RowKey": [
			"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		]
	}

	Respond result
	{
		"success": true,
		"Result": [
			{
				"success": true,
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			}
		],
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}
*/
