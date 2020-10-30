package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// VaultAccount - Encapsulates a single generic VaultAccount
type VaultAccount struct {
	vaultObject
	// VaultAccount specific APIs
	apiUpdatePassword   string
	apiCheckDelete      string
	apiGetChallenge     string
	apiCheckoutPassword string
	apiCheckinPassword  string
	apiSetAdminAccount  string

	// Settings menu
	User           string `json:"User,omitempty" schema:"name,omitempty"` // User Name
	Password       string `json:"Password,omitempty" schema:"password,omitempty"`
	Host           string `json:"Host,omitempty" schema:"host_id,omitempty"`
	SSHKeyID       string `json:"SshKeyId,omitempty" schema:"sshkey_id,omitempty"`
	DomainID       string `json:"DomainID,omitempty" schema:"domain_id,omitempty"`
	DatabaseID     string `json:"DatabaseID,omitempty" schema:"database_id,omitempty"`
	CredentialType string `json:"CredentialType,omitempty" schema:"credential_type,omitempty"`

	// Policy menu
	UseWheel                       bool            `json:"UseWheel,omitempty" schema:"use_proxy_account,omitempty"` // Use proxy account
	IsManaged                      bool            `json:"IsManaged,omitempty" schema:"managed,omitempty"`          // manage this credential
	Description                    string          `json:"Description,omitempty" schema:"description,omitempty"`
	Status                         string          `json:"Status,omitempty" schema:"status,omitempty"`
	DefaultCheckoutTime            int             `json:"DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)
	PasswordCheckoutDefaultProfile string          `json:"PasswordCheckoutDefaultProfile" schema:"default_profile_id"`         // Default Password Checkout Profile (used if no conditions matched)
	ChallengeRules                 *ChallengeRules `json:"PasswordCheckoutRules,omitempty" schema:"challenge_rule,omitempty"`

	IsAdminAccount bool `json:"IsAdminAccount,omitempty" schema:"is_admin_account,omitempty"`
}

// NewVaultAccount is a VaultAccount constructor
func NewVaultAccount(c *restapi.RestClient) *VaultAccount {
	s := VaultAccount{}
	s.client = c
	s.apiRead = "/ServerManage/GetAllAccountInformation"
	s.apiCreate = "/ServerManage/AddAccount"
	s.apiDelete = "/ServerManage/DeleteAccount"
	s.apiUpdate = "/ServerManage/UpdateAccount"
	s.apiUpdatePassword = "/ServerManage/UpdatePassword"
	s.apiCheckDelete = "/ServerManage/CanDeleteAccount"
	s.apiGetChallenge = "/ServerManage/GetAccountChallenges"
	s.apiCheckoutPassword = "/ServerManage/CheckoutPassword"
	s.apiCheckinPassword = "/ServerManage/CheckinPassword"
	s.apiPermissions = "/ServerManage/SetAccountPermissions"
	s.apiSetAdminAccount = "/ServerManage/SetAdministrativeAccounts"

	return &s
}

// Read function fetches a VaultAccount from source, including attribute values. Returns error if any
func (o *VaultAccount) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	LogD.Printf("Response for VaultAccount from tenant: %v", resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New(resp.Message)
	}

	var va = resp.Result["VaultAccount"].(map[string]interface{})
	var row = va["Row"].(map[string]interface{})

	fillWithMap(o, row)

	// Get password checkout profile information
	resp, err = o.client.CallGenericMapAPI(o.apiGetChallenge, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	if v, ok := resp.Result["PasswordCheckoutDefaultProfile"]; ok {
		o.PasswordCheckoutDefaultProfile = v.(string)
	}

	// Fill challenge rules
	if v, ok := resp.Result["PasswordCheckoutRules"]; ok {
		challengerules := &ChallengeRules{}
		fillWithMap(challengerules, v.(map[string]interface{}))
		o.ChallengeRules = challengerules
	}

	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Create function creates a new VaultAccount and returns a map that contains creation result
func (o *VaultAccount) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	// Special handling of password checkout profile
	//if queryArg["PasswordCheckoutDefaultProfile"] != "" {
	queryArg["updateChallenges"] = true
	//}

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

// Delete function deletes a VaultAccount and returns a map that contains deletion result
func (o *VaultAccount) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing VaultAccount and returns a map that contains update result
func (o *VaultAccount) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	// Special handling of password checkout profile
	//if queryArg["PasswordCheckoutDefaultProfile"] != "" {
	queryArg["updateChallenges"] = true
	//}

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

// ChangePassword function updates an existing VaultAccount password and returns a map that contains update result
func (o *VaultAccount) ChangePassword() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	// Mandatory attributes
	queryArg["ID"] = o.ID
	queryArg["Password"] = o.Password

	reply, err := o.client.CallBoolAPI(o.apiUpdatePassword, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// ValidateCredentialType checks credential type matches password or sshkey setting
func (o *VaultAccount) ValidateCredentialType() error {
	if o.CredentialType == "Password" && o.Password == "" {
		return errors.New("Credential type is 'Password' but password isn't set")
	}
	if o.CredentialType == "SshKey" && o.SSHKeyID == "" {
		return errors.New("Credential type is 'SSHKey' but Sshkey_id isn't set")
	}
	return nil
}

// Query function returns a single VaultAccount object in map format
func (o *VaultAccount) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM VaultAccount WHERE 1=1"
	if o.User != "" {
		query += " AND User='" + o.User + "'"
	}
	if o.Host != "" {
		query += " AND Host='" + o.Host + "'"
	}
	if o.DatabaseID != "" {
		query += " AND DatabaseID='" + o.DatabaseID + "'"
	}
	if o.DomainID != "" {
		query += " AND DomainID='" + o.DomainID + "'"
	}

	return queryVaultObject(o.client, query)
}

// CheckoutPassword checks out account password from vault
func (o *VaultAccount) CheckoutPassword() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Description"] = "Checkout by Terraform provider"

	reply, err := o.client.CallGenericMapAPI(o.apiCheckoutPassword, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// CheckinPassword checks in an checked out account password
func (o *VaultAccount) CheckinPassword(coid string) (*restapi.BoolResponse, error) {
	if coid == "" {
		return nil, errors.New("error: COID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = coid

	reply, err := o.client.CallBoolAPI(o.apiCheckinPassword, queryArg)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

func (o *VaultAccount) getPerms() map[string]string {
	perms := accountPermissions
	if o.Host != "" {
		perms = accountPermissions
	} else if o.DomainID != "" {
		perms = domainaccountPermissions
	} else if o.DatabaseID != "" {
		perms = dbaccountPermissions
	}

	return perms
}

func (o *VaultAccount) setAdminAccount(enable bool) error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	if enable {
		queryArg["PVID"] = o.ID
	}
	queryArg["Systems"] = []string{o.Host}
	reply, err := o.client.CallGenericMapAPI(o.apiSetAdminAccount, queryArg)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}

	return nil
}

/*
	API to manage vault account

	Read Account

		Request body format
		{
			"ID": "aff00452-71d8-4a8a-a535-805cb059ade8"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"VaultCheckout": null,
				"Workflow": {
					"WorkflowApprover": ...
					"WorkflowApprovers": ...
					"WorkflowEnabled": true,
					"PasswordResettable": false,
					"AccountUnlockable": false
				},
				"VaultAccount": {
					"Row": {
						"IsQuickStartAccount": null,
						"ID": "aff00452-71d8-4a8a-a535-805cb059ade8",
						"ActiveSessions": 0,
						"Host": "b2888c5d-6ec8-43c9-a0f6-918f9c43c0e5",
						"LastChange": "/Date(1588075662725)/",
						"DefaultCheckoutTime": null,
						"EffectiveWorkflowEnabled": null,
						"DueBack": null,
						"ProxyUser": "admin",
						"CanManage": false,
						"Mode": null,
						"Status": "",
						"User": "adminuser",
						"PasswordResetRetryCount": 0,
						"DomainID": null,
						"WorkflowApproversList": null,
						"KmipId": null,
						"PasswordResetLastError": "",
						"AccountRights": "0000000000000000000000000000000000000000000000000000011101111011",
						"WorkflowApprover": ...
						"NeedsPasswordReset": "NotNeeded",
						"HealthError": "OK",
						"SessionType": "Ssh",
						"MissingPassword": false,
						"CredentialId": null,
						"MPParent": null,
						"IsPrivileged": null,
						"WorkflowApprovers": ...
						"Description": "",
						"FQDN": "192.168.18.53",
						"IsManaged": true,
						"EffectiveWorkflowApprover": null,
						"Name": "FortiManager",
						"Rights": "Login, Naked, Manage, Owner, Delete, UpdatePassword, RotatePassword, View, FileTransfer",
						"OwnerId": null,
						"UseWheel": false,
						"WorkflowDefaultOptions": null,
						"DeviceID": null,
						"LastHealthCheck": "/Date(1588075662725)/",
						"DesktopApps": [],
						"CredentialType": "Password",
						"ComputerClass": "CustomSsh",
						"ActiveCheckouts": 0,
						"WorkflowEnabled": null,
						"DiscoveredTime": null,
						"_MatchFilter": null,
						"Healthy": "OK",
						"OwnerName": null,
						"EffectiveWorkflowApprovers": null,
						"UserDisplayName": "adminuser (FortiManager)",
						"IsFavorite": false,
						"DatabaseID": null,
						"NeedPassword": false
					},
					"Entities": [
						{
							"Type": "Server",
							"Key": "b2888c5d-6ec8-43c9-a0f6-918f9c43c0e5",
							"IsForeignKey": true
						},
						{
							"Type": "VaultAccount",
							"Key": "aff00452-71d8-4a8a-a535-805cb059ade8",
							"IsForeignKey": false
						}
					]
				},
				"RelatedResource": {
					"Row": {
						...
					},
					"Entities": [
						...
					]
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

	Create Account
	https://developer.centrify.com/reference#post_servermanage-addaccount

		Request body format
		{
			"User": "testaccount",
			"CredentialType": "Password",
			"IsManaged": false,
			"Password": "xxxxxxxxxx",
			"UseWheel": false,
			"Description": "afdasdfasfd",
			"Host": "ab6d8e59-39e5-4050-ac38-4b77122a817c"
		}
		or
		{
			"User": "account2",
			"CredentialType": "SshKey",
			"IsManaged": false,
			"undefined": "false",
			"SshKeyName": "my_ami_key",
			"SshKeyId": "56ec631d-635a-48f7-b100-6b97b6eb6f47",
			"Host": "ab6d8e59-39e5-4050-ac38-4b77122a817c",
			"Description": ""
		}
		or
		{
			"User": "sa",
			"Password": "xxxxxxxxx",
			"UseWheel": false,
			"Description": "sa account",
			"IsManaged": false,
			"DatabaseID": "6bbc81ca-36a5-4a77-8f41-fd79f4af11af"
		}

		Respond result
		{
			"success": true,
			"Result": "c9772d3d-02b1-411e-ba77-58795586907d",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update Account
	https://developer.centrify.com/reference#post_servermanage-updateaccount

		Request body format
		{
			"ID": "c9772d3d-02b1-411e-ba77-58795586907d",
			"ActiveSessions": 0,
			"Host": "80b5847a-fcff-471f-8278-990c620eed24",
			"CanManage": true,
			"Status": "Unknown",
			"User": "testaccount",
			"AccountRights": "0000000000000000000000000000000000000000000000000000011101111011",
			"WorkflowApprover": ...
			"NeedsPasswordReset": "NotNeeded",
			"HealthError": "NoManagementChannelAvailable",
			"SessionType": "Rdp",
			"MissingPassword": false,
			"WorkflowApprovers": ...
			"Description": "Test account",
			"FQDN": "192.168.2.3",
			"IsManaged": false,
			"Name": "Windows 01",
			"Rights": "Login, Naked, Manage, Owner, Delete, UpdatePassword, RotatePassword, View, FileTransfer",
			"DesktopApps": [],
			"CredentialType": "Password",
			"ComputerClass": "Windows",
			"ActiveCheckouts": 0,
			"Healthy": "Unknown",
			"UserDisplayName": "testaccount (Windows 01)",
			"IsFavorite": false,
			"NeedPassword": false,
			"IsAdministrativeAccount": false,
			"_RowKey": "c9772d3d-02b1-411e-ba77-58795586907d",
			"updateChallenges": true,
			"WorkflowEnabled": null
		}

		Respond result
		{
			"success": true,
			"Result": {
				"PVID": "c9772d3d-02b1-411e-ba77-58795586907d"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete Account
	https://developer.centrify.com/reference#post_servermanage-deleteaccount

		Check if Account can be deleted
		Request body format
		{
			"ID": "c9772d3d-02b1-411e-ba77-58795586907d",
			"RRFormat": true
		}

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

		Deleting account
		Request body format
		{
			"ID": "c9772d3d-02b1-411e-ba77-58795586907d"
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

		Upate password
		https://developer.centrify.com/reference#post_servermanage-updatepassword

		Request body format
		{
			"ID": "da96a581-22a8-4aa5-a077-7cf43c78666f"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"PasswordCheckoutDefaultProfile": "30804754-3b87-4862-a39e-0f042825a3a0",
				"PasswordCheckoutRules": {
					"_UniqueKey": "Condition",
					"_Value": [
						{
							"Conditions": [
								{
									"Prop": "IdentityCookie",
									"Op": "OpNotExists"
								}
							],
							"ProfileId": "bdfced4d-e2d3-4dae-a84f-dc8282faf8eb"
						}
					],
					"Enabled": true,
					"_Type": "RowSet"
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

	Set Admin Account
		Request body format
		{
			"Systems": [
				"949cd9c7-2949-446a-8e78-1fbc3d58b4ee"
			],
			"PVID": "60164260-a7b1-4a1e-9f5f-63f88d99c682"
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

/*
{
    "ID": "0714835b-2a52-423a-9b9e-c1934dadb84c",
    "ActiveSessions": 0,
    "Host": "0b07c7b1-5481-4bbc-8ef8-ea7874ea68f2",
    "ProxyUser": "root",
    "CanManage": false,
    "User": "user1",
    "PasswordResetRetryCount": 0,
    "WorkflowApproversList": "[
	{
		\"DisplayName\":\"mspadmin\",
		\"ObjectType\":\"User\",
		\"DistinguishedName\":\"admin@example.com\",
		\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",
		\"SystemName\":\"madmin@example.com\",
		\"ServiceInstance\":\"CDS\",
		\"Locked\":false,
		\"InternalName\":\"c2c7bcc6-9560-44e0-8dff-5be221cd37ee\",
		\"StatusEnum\":\"Active\",
		\"ServiceInstanceLocalized\":\"Centrify Directory\",
		\"ServiceType\":\"CDS\",
		\"EMail\":\"admin@demo.lab\",
		\"Status\":\"Active\",
		\"Enabled\":true,\
		"Name\":\"admin@example.com\",
		\"Guid\":\"c2c7bcc6-9560-44e0-8dff-5be221cd37ee\",
		\"Type\":\"User\",
		\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
		\"Principal\":\"admin@example.com\",
		\"PType\":\"User\",
		\"OptionsSelector\":true
	},
	{
		\"ReadOnly\":false,
		\"Description\":\"AD accounts that are granted local administrator access to non-domain joined machines.\",
		\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",
		\"RoleType\":\"PrincipalList\",
		\"_ID\":\"d958fad8_f90a_4c40_b986_f6fa31713bba\",
		\"Name\":\"LAB Cloud Local Admins\",
		\"Guid\":\"d958fad8_f90a_4c40_b986_f6fa31713bba\",
		\"Type\":\"Role\",
		\"ObjectType\":\"Role\",
		\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
		\"Principal\":\"LAB Cloud Local Admins\",
		\"PType\":\"Role\"
	},
	{
		\"Type\":\"Manager\",
		\"NoManagerAction\":\"useBackup\",
		\"BackupApprover\":
		{
			\"ReadOnly\":false,
			\"Description\":\"Machines and users who are enforced MFA for direct access without going through PAS.\",
			\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",
			\"RoleType\":\"PrincipalList\",
			\"_ID\":\"d06fc516_8c9b_4f76_a08d_797ca6fd0a55\",
			\"Name\":\"LAB MFA Machines & Users\",
			\"Guid\":\"d06fc516_8c9b_4f76_a08d_797ca6fd0a55\",
			\"Type\":\"Role\",
			\"ObjectType\":\"Role\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
			\"Principal\":\"LAB MFA Machines & Users\",
			\"PType\":\"Role\"
		},
		\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\"
	}
	]",
    "AccountRights": "0000000000000000000000000000000000000000000000000000011101111011",
    "NeedsPasswordReset": "NotNeeded",
    "HealthError": "OK",
    "SessionType": "Ssh",
    "MissingPassword": false,
    "WorkflowApprovers": "[
		{
			\"DisplayName\":\"mspadmin\",
			\"ObjectType\":\"User\",
			\"DistinguishedName\":\"admin@example.com\",
			\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",
			\"SystemName\":\"admin@example.com\",
			\"ServiceInstance\":\"CDS\",
			\"Locked\":false,
			\"InternalName\":\"c2c7bcc6-9560-44e0-8dff-5be221cd37ee\",\
			"StatusEnum\":\"Active\",
			\"ServiceInstanceLocalized\":\"Centrify Directory\",
			\"ServiceType\":\"CDS\",
			\"EMail\":\"admin@demo.lab\",
			\"Status\":\"Active\",
			\"Enabled\":true,
			\"Name\":\"admin@example.com\",
			\"Guid\":\"c2c7bcc6-9560-44e0-8dff-5be221cd37ee\",
			\"Type\":\"User\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
			\"Principal\":\"admin@example.com\",
			\"PType\":\"User\",
			\"OptionsSelector\":true
		},
		{
			\"ReadOnly\":false,
			\"Description\":\"AD accounts that are granted local administrator access to non-domain joined machines.\",
			\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",
			\"RoleType\":\"PrincipalList\",
			\"_ID\":\"d958fad8_f90a_4c40_b986_f6fa31713bba\",
			\"Name\":\"LAB Cloud Local Admins\",
			\"Guid\":\"d958fad8_f90a_4c40_b986_f6fa31713bba\",
			\"Type\":\"Role\",
			\"ObjectType\":\"Role\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
			\"Principal\":\"LAB Cloud Local Admins\",
			\"PType\":\"Role\"
		},
		{
			\"Type\":\"Manager\",
			\"NoManagerAction\":\"useBackup\",
			\"BackupApprover\":
			{
				\"ReadOnly\":false,
				\"Description\":\"Machines and users who are enforced MFA for direct access without going through PAS.\",
				\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",
				\"RoleType\":\"PrincipalList\",
				\"_ID\":\"d06fc516_8c9b_4f76_a08d_797ca6fd0a55\",
				\"Name\":\"LAB MFA Machines & Users\",
				\"Guid\":\"d06fc516_8c9b_4f76_a08d_797ca6fd0a55\",
				\"Type\":\"Role\",
				\"ObjectType\":\"Role\",
				\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
				\"Principal\":\"LAB MFA Machines & Users\",
				\"PType\":\"Role\"
			},
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\"
		}]",
	"WorkflowApprovers": "[
		{
			\"Type\":\"Manager\",
			\"NoManagerAction\":\"approve\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
			\"OptionsSelector\":
		}
	]",
	"WorkflowApprovers": "[
		{
			\"Type\":\"Manager\",
			\"NoManagerAction\":\"deny\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
			\"OptionsSelector\":true
		}
	]",
    "FQDN": "192.168.18.15",
    "IsManaged": true,
    "Name": "XML File",
    "Rights": "Login, Naked, Manage, Owner, Delete, UpdatePassword, RotatePassword, View, FileTransfer",
    "UseWheel": true,
    "WorkflowDefaultOptions": "{\"GrantMin\":60}",
    "DesktopApps": [],
    "CredentialType": "Password",
    "ComputerClass": "CustomSsh",
    "ActiveCheckouts": 0,
    "WorkflowEnabled": true,
    "Healthy": "OK",
    "UserDisplayName": "user1 (XML File)",
    "IsFavorite": false,
    "NeedPassword": false,
    "IsAdministrativeAccount": false,
    "NoManagerAction": "useBackup",
    "BackupApprover": {
        "ReadOnly": false,
        "Description": "Machines and users who are enforced MFA for direct access without going through PAS.",
        "DirectoryServiceUuid": "09B9A9B0-6CE8-465F-AB03-65766D33B05E",
        "RoleType": "PrincipalList",
        "_ID": "d06fc516_8c9b_4f76_a08d_797ca6fd0a55",
        "Name": "LAB MFA Machines & Users",
        "Guid": "d06fc516_8c9b_4f76_a08d_797ca6fd0a55",
        "Type": "Role",
        "ObjectType": "Role",
        "Type-generated-field": "/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128",
        "Principal": "LAB MFA Machines & Users",
        "PType": "Role"
    },
    "WorkflowSent": true,
    "_RowKey": "0714835b-2a52-423a-9b9e-c1934dadb84c",
    "updateChallenges": true,
    "Description": ""
}
*/
