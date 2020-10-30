package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// VaultDomain - Encapsulates a single Domain
type VaultDomain struct {
	vaultObject
	apiSetAdminAccount string
	apiCanDelete       string

	VerifyDomain bool
	// Policy menu related settings
	DefaultCheckoutTime int `json:"DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)
	// Advanced menu -> Administrative Account Settings
	AdminAccountID               string `json:"AdminAccountID,omitempty" schema:"administrative_account_id,omitempty"`
	AdminAccountDomain           string `json:"AdminAccountDomain,omitempty" schema:"administrative_account_domain,omitempty"`
	AdminAccountPassword         string `json:"AdminAccountPassword,omitempty" schema:"administrative_account_password,omitempty"`
	AdminAccountName             string `json:"AdminAccountName,omitempty" schema:"administrative_account_name,omitempty"`
	AutoDomainAccountMaintenance bool   `json:"AllowAutomaticAccountMaintenance,omitempty" schema:"auto_domain_account_maintenance,omitempty"`     // Enable Automatic Domain Account Maintenance
	AutoLocalAccountMaintenance  bool   `json:"AllowAutomaticLocalAccountMaintenance,omitempty" schema:"auto_local_account_maintenance,omitempty"` // Enable Automatic Local Account Maintenance
	ManualDomainAccountUnlock    bool   `json:"AllowManualAccountUnlock,omitempty" schema:"manual_domain_account_unlock,omitempty"`                // Enable Manual Domain Account Unlock
	ManualLocalAccountUnlock     bool   `json:"AllowManualLocalAccountUnlock,omitempty" schema:"manual_local_account_unlock,omitempty"`            // Enable Manual Local Account Unlock
	// Advanced -> Security Settings
	AllowMultipleCheckouts            bool   `json:"AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts per AD account added for this domain
	AllowPasswordRotation             bool   `json:"AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration            int    `json:"PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin bool   `json:"AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                int    `json:"MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	PasswordProfileID                 string `json:"PasswordProfileID,omitempty" schema:"password_profile_id,omitempty"`                                    // Password Complexity Profile
	// Advanced -> Maintenance Settings
	AllowPasswordHistoryCleanUp    bool `json:"AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`     // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration int  `json:"PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"` // Password history cleanup (days)
	// Advanced -> Domain/Zone Tasks
	AllowRefreshZoneJoined           bool `json:"AllowRefreshZoneJoined,omitempty" schema:"enable_zone_joined_check,omitempty"`             // Enable periodic domain/zone joined check
	RefreshZoneJoinedIntervalMinutes int  `json:"RefreshZoneJoinedIntervalMinutes,omitempty" schema:"zone_joined_check_interval,omitempty"` // Domain/zone joined check interval (minutes)
	AllowZoneRoleCleanup             bool `json:"AllowZoneRoleCleanup,omitempty" schema:"enable_zone_role_cleanup,omitempty"`               // Enable periodic removal of expired zone role assignments
	ZoneRoleCleanupIntervalHours     int  `json:"ZoneRoleCleanupIntervalHours,omitempty" schema:"zone_role_cleanup_interval,omitempty"`     // Expired zone role assignment removal interval (hours)

	// System -> Connectors menu related settings
	ProxyCollectionList string `json:"ProxyCollectionList,omitempty" schema:"connector_list,omitempty"` // List of Connectors used
	// Sets
	//Sets []string `json:"Sets,omitempty" schema:"sets,omitempty"`
}

// NewVaultDomain is a Domain constructor
func NewVaultDomain(c *restapi.RestClient) *VaultDomain {
	s := VaultDomain{}
	s.client = c
	s.apiRead = "/RedRock/query"
	s.apiCreate = "/ServerManage/AddDomain"
	s.apiDelete = "/ServerManage/DeleteDomain"
	s.apiUpdate = "/ServerManage/UpdateDomain"
	s.apiPermissions = "/ServerManage/SetDomainPermissions"
	s.apiSetAdminAccount = "/ServerManage/SetAdministrativeAccounts"
	s.apiCanDelete = "/ServerManage/CanDeleteDomain"

	return &s
}

// Read function fetches a Domain from source, including attribute values. Returns error if any
func (o *VaultDomain) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Script"] = "SELECT * FROM VaultDomain WHERE ID = '" + o.ID + "'"
	queryArg["Args"] = subArgs

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
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
		return errors.New("Domain does not exist in tenant")
	} else if len(results) > 1 {
		// this should never happen
		return errors.New("There are more than one domains with the same ID in tenant")
	}
	var result = results[0].(map[string]interface{})
	// Populate vaultObject struct with map from response
	var row = result["Row"].(map[string]interface{})
	//LogD.Printf("Input map: %+v", row)
	fillWithMap(o, row)

	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Create function creates a new Domain and returns a map that contains creation result
func (o *VaultDomain) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg["Name"] = o.Name
	queryArg["VerifyDomain"] = o.VerifyDomain

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

// Delete function deletes a Domain and returns a map that contains deletion result
func (o *VaultDomain) Delete() (*restapi.BoolResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	// Check if the domain can be deleted
	resp, err := o.client.CallGenericMapAPI(o.apiCanDelete, queryArg)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}
	if resp.Result["can"].(bool) {
		return o.deleteObjectBoolAPI("")
	}

	LogW.Printf("Domain cannot be deleted: %+v", resp.Result["why"])
	return nil, nil
}

// Update function updates an existing Domain and returns a map that contains update result
func (o *VaultDomain) Update() (*restapi.GenericMapResponse, error) {
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

// Query function returns a single Set object in map format
func (o *VaultDomain) Query() (map[string]interface{}, error) {
	query := "SELECT ID, Name FROM VaultDomain WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

func (o *VaultDomain) setAdminAccount() error {
	if o.Name == "" {
		return errors.New("error: Domain name is empty")
	}
	var queryArg = make(map[string]interface{})
	if o.AdminAccountID != "" {
		queryArg["UserUuid"] = o.AdminAccountID
		if o.AdminAccountPassword != "" {
			queryArg["Password"] = o.AdminAccountPassword
		}
		if o.AdminAccountName != "" {
			queryArg["User"] = o.AdminAccountName
		}
	}

	if o.AdminAccountDomain != "" {
		queryArg["Domains"] = []string{o.AdminAccountDomain}
	} else {
		// Try our best to use domain from the Domain itself
		queryArg["Domains"] = []string{o.Name}
	}

	resp, err := o.client.CallGenericMapAPI(o.apiSetAdminAccount, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	return nil
}

/*
	Fetch Domain

		Request body format


		Respond result
		{
			"success": true,
			"Result": {
				"IsAggregate": false,
				"Count": 1,
				"Columns": [
					...
				],
				"FullCount": 1,
				"Results": [
					{
						"Entities": [
							{
								"Type": "VaultDomain",
								"Key": "5f742093-6c3a-4235-b84f-f1b86627ba58",
								"IsForeignKey": false
							}
						],
						"Row": {
							"ID": "5f742093-6c3a-4235-b84f-f1b86627ba58",
							"AllowMultipleCheckouts": null,
							"DefaultCheckoutTime": null,
							"MinimumPasswordAge": null,
							"AllowManualLocalAccountUnlock": true,
							"ZoneRoleWorkflowRoles": "[{\"Name\":\"cfyl-SSH Management/Unix Zone\",\"Unix\":true,\"ZoneDn\":\"CN=Unix Zone,CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\",\"Description\":\"Manage SSHD daemon and configuration\",\"ZoneCanonicalName\":\"centrifylab.aws/Centrify/Zones/Global/Unix Zone\",\"ParentZoneDn\":\"CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\"},{\"Name\":\"cfyl-Super System Admin/Unix Zone\",\"Unix\":true,\"ZoneDn\":\"CN=Unix Zone,CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\",\"Description\":\"Super system admin with root access\",\"ZoneCanonicalName\":\"centrifylab.aws/Centrify/Zones/Global/Unix Zone\",\"ParentZoneDn\":\"CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\"},{\"Name\":\"cfyw-Windows System Admin/Windows Zone\",\"ZoneDn\":\"CN=Windows Zone,CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\",\"Description\":\"Windows system admin for managing local configurations\",\"ZoneCanonicalName\":\"centrifylab.aws/Centrify/Zones/Global/Windows Zone\",\"Windows\":true,\"ParentZoneDn\":\"CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\"}]",
							"AllowAutomaticAccountUnlock": null,
							"ForestID": "5f742093-6c3a-4235-b84f-f1b86627ba58",
							"PasswordHistoryCleanUpDuration": null,
							"AllowHealthCheck": null,
							"HealthStatusError": "Success",
							"UniqueId": null,
							"Reachable": true,
							"LastState": "OK",
							"AllowZoneRoleCleanup": true,
							"ZoneRoleWorkflowApproversList": "[{\"BackupApprover\":{\"PType\":\"Role\",\"ObjectType\":\"Role\",\"Name\":\"LAB Infrastructure Owners\",\"Guid\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"_ID\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1583273484\",\"Principal\":\"LAB Infrastructure Owners\",\"Description\":\"Approver who can approve access to access lab systems.\",\"RoleType\":\"PrincipalList\",\"ReadOnly\":false,\"Type\":\"Role\",\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\"},\"NoManagerAction\":\"useBackup\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1583273484\",\"OptionsSelector\":true,\"Type\":\"Manager\"}]",
							"LastZoneRoleCleanup": "/Date(1595383548740)/",
							"Administrator": "846966fa-dd37-40e9-9958-cab9f5b27b66",
							"SyncFromConnector": true,
							"PasswordProfileID": null,
							"PasswordRotateDuration": null,
							"Description": "",
							"ProxyCollectionList": "7ee6c65b-f000-4300-8e42-d110eb6da6c8",
							"AdministratorDisplayName": "ad_admin (centrifylab.aws)",
							"AllowPasswordRotationAfterCheckin": null,
							"AllowAutomaticAccountMaintenance": true,
							"Name": "centrifylab.aws",
							"AllowRefreshZoneJoined": true,
							"ZoneRoleWorkflowApprovers": "[{\"BackupApprover\":{\"PType\":\"Role\",\"ObjectType\":\"Role\",\"Name\":\"LAB Infrastructure Owners\",\"Guid\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"_ID\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1583273484\",\"Description\":\"Approver who can approve access to access lab systems.\",\"Principal\":\"LAB Infrastructure Owners\",\"RoleType\":\"PrincipalList\",\"ReadOnly\":false,\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",\"Type\":\"Role\"},\"NoManagerAction\":\"useBackup\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1583273484\",\"OptionsSelector\":true,\"Type\":\"Manager\"}]",
							"AllowAutomaticLocalAccountMaintenance": true,
							"LastHealthCheck": "/Date(1595383490412)/",
							"ZoneRoleWorkflowEnabled": true,
							"ReachableError": "Success",
							"HealthCheckInterval": null,
							"ParentID": null,
							"PasswordRotateInterval": null,
							"ZoneRoleCleanupIntervalHours": 6,
							"AllowPasswordHistoryCleanUp": null,
							"RefreshZoneJoinedIntervalMinutes": 1440,
							"DiscoveredTime": null,
							"_MatchFilter": null,
							"LastRefreshZoneJoined": "/Date(1598764755975)/",
							"HealthStatus": "OK",
							"AllowManualAccountUnlock": true,
							"AllowPasswordRotation": null
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

	Create Domain
	https://developer.centrify.com/reference#post_servermanage-adddomain

		Request body format
		{
			"Name": "example.com",
			"VerifyDomain": false
		}

		Respond result
		{
			"success": true,
			"Result": "5857a828-785d-402e-966e-73ca1a198c98",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update Domain
	https://developer.centrify.com/reference#post_servermanage-updatedomain

		Request body format
		{
			"ID": "5857a828-785d-402e-966e-73ca1a198c98",
			"AllowMultipleCheckouts": true,
			"DefaultCheckoutTime": 30,
			"MinimumPasswordAge": 98,
			"AllowManualLocalAccountUnlock": true,
			"PasswordHistoryCleanUpDuration": 100,
			"HealthStatusError": "UnknownError",
			"Reachable": false,
			"LastState": "Unreachable",
			"AllowZoneRoleCleanup": true,
			"PasswordProfileID": "a5143c7e-1195-4f6e-a755-11b5e0213cda",
			"PasswordRotateDuration": 90,
			"Description": "example domain",
			"AdministratorDisplayName": "demo.lab\\ad_admin",
			"AllowPasswordRotationAfterCheckin": true,
			"AllowAutomaticAccountMaintenance": true,
			"Name": "example.com",
			"AllowRefreshZoneJoined": true,
			"AllowAutomaticLocalAccountMaintenance": true,
			"ReachableError": "UnknownError",
			"ZoneRoleCleanupIntervalHours": 6,
			"AllowPasswordHistoryCleanUp": true,
			"RefreshZoneJoinedIntervalMinutes": 1440,
			"HealthStatus": "Unreachable",
			"AllowManualAccountUnlock": true,
			"AllowPasswordRotation": true,
			"newAdministrator": {
				"User": "demo.lab\\ad_admin",
				"PVID": "6ff45b40-7375-4887-bfba-a84849a2250a",
				"Domains": [
					"example.com"
				]
			},
			"_RowKey": "5857a828-785d-402e-966e-73ca1a198c98"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"ID": "5857a828-785d-402e-966e-73ca1a198c98"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete Domain
	https://developer.centrify.com/reference#post_servermanage-deletedomain

		Request body format
		{
			"ID": "ba08a2bb-0667-4737-a80e-13112f2fdb1e"
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

	Set Administrative Account

	Request body format
	{
		"User": "demo.lab\\ad_admin",
		"PVID": "6ff45b40-7375-4887-bfba-a84849a2250a",
		"Domains": [
			"example.com"
		]
	}
	or
	{
		"User": "Administrator@demo.lab",
		"UserUuid": "b5e89f6e-ac09-4258-be70-64dc91ca9054",
		"Password": "xxxxxx",
		"Domains": [
			"demo.lab"
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

	CanDelete Domain?

	Respond result
	{
		"success": true,
		"Result": {
			"can": true,
			"why": null
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}
*/
