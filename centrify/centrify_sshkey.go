package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// SSHKey - Encapsulates a single generic SSHKey
type SSHKey struct {
	vaultObject
	apiGetChallenge string
	apiRetrieve     string

	SSHKeysDefaultProfileID string          `json:"SshKeysDefaultProfile,omitempty" schema:"default_profile_id,omitempty"` // Default SSH Key Challenge Profile
	ChallengeRules          *ChallengeRules `json:"SshKeysRules,omitempty" schema:"challenge_rule,omitempty"`
	KeyFormat               string          `json:"KeyFormat,omitempty" schema:"key_format,omitempty"`
	KeyLength               int             `json:"KeyLength,omitempty" schema:"key_length,omitempty"`
	KeyType                 string          `json:"KeyType,omitempty" schema:"key_type,omitempty"`
	IsManaged               bool            `json:"IsManaged,omitempty" schema:"is_managed,omitempty"`
	Description             string          `json:"Comment,omitempty" schema:"description,omitempty"`
	PrivateKey              string          `json:"PrivateKey,omitempty" schema:"private_key,omitempty"`
	Passphrase              string          `json:"Passphrase,omitempty" schema:"passphrase,omitempty"`
	KeyPairType             string          `json:"KeyPairType,omitempty" schema:"key_pair_type,omitempty"` // Which key to retrieve from the pair, must be either PublicKey, PrivateKey, or PPK
}

// NewSSHKey is a SSHKey constructor
func NewSSHKey(c *restapi.RestClient) *SSHKey {
	s := SSHKey{}
	s.client = c
	s.apiRead = "/ServerManage/GetSshKeyInfo"
	s.apiCreate = "/ServerManage/AddSshKey"
	s.apiDelete = "/ServerManage/DeleteSshKey"
	s.apiUpdate = "/ServerManage/UpdateSshKey"
	s.apiRetrieve = "/ServerManage/RetrieveSshKey"
	s.apiGetChallenge = "/ServerManage/GetSshKeyRightsAndChallenges"
	s.apiPermissions = "/ServerManage/SetSSHKeyPermissions"

	return &s
}

// Read function fetches a SSHKey from source
func (o *SSHKey) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	LogD.Printf("Response for SSHKey from tenant: %v", resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}

	fillWithMap(o, resp.Result)

	// Get SSH Key challenge profile information
	resp, err = o.client.CallGenericMapAPI(o.apiGetChallenge, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	if v, ok := resp.Result["Challenges"]; ok {
		challenges := v.(map[string]interface{})
		if p, ok := challenges["SshKeysDefaultProfile"]; ok {
			o.SSHKeysDefaultProfileID = p.(string)
		}
		// Fill challenge rules
		if r, ok := challenges["SshKeysRules"]; ok {
			challengerules := &ChallengeRules{}
			fillWithMap(challengerules, r.(map[string]interface{}))
			o.ChallengeRules = challengerules
		}
	}

	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Create function creates a new SSHKey and returns a map that contains creation result
func (o *SSHKey) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	// Special handling of challenge checkout profile
	queryArg["updateChallenges"] = true

	LogD.Printf("Generated Map for Create(): %+v", queryArg)

	resp, err := o.client.CallStringAPI(o.apiCreate, queryArg)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

// Delete function deletes a SSHKey and returns a map that contains deletion result
func (o *SSHKey) Delete() (*restapi.StringResponse, error) {
	return o.deleteObjectStringAPI("")
}

// Update function updates an existing SSHKey and returns a map that contains update result
func (o *SSHKey) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	// Special handling of password checkout profile
	queryArg["updateChallenges"] = true

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

// Query function returns a single SSHKey object in map format
func (o *SSHKey) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM SshKeys WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// RetriveSSHKey retrieves SSH Key from vault
func (o *SSHKey) RetriveSSHKey() (*restapi.StringResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["KeyPairType"] = o.KeyPairType
	if o.Passphrase != "" {
		queryArg["Passphrase"] = o.Passphrase
	}
	if o.KeyFormat != "" {
		queryArg["KeyFormat"] = o.KeyFormat
	}

	resp, err := o.client.CallStringAPI(o.apiRetrieve, queryArg)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

/*
	Fetch SSH Key
	https://developer.centrify.com/reference#post_servermanage-getsshkeyinfo

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond Result
		{
			"success": true,
			"Result": {
				"_entitycontext": "W/\"datetime'2020-08-24T12%3A16%3A42.6342428Z'\"",
				"KeyFormat": "PEM",
				"LastUpdated": "/Date(1598271402593)/",
				"IsManaged": false,
				"_metadata": {
					"Version": 1,
					"IndexingVersion": 1
				},
				"_TableName": "sshkeys",
				"_encryptkeyid": "XXXXX",
				"_PartitionKey": "XXXXX",
				"_RowKey": "5xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Revision": 0,
				"_Timestamp": "/Date(1582552067897)/",
				"CreatedBy": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Name": "my_ami_key",
				"KeyType": "RSA",
				"Comment": "my AWS AMI key",
				"Created": "/Date(1582552067684)/",
				"ACL": "true",
				"KeyLength": 2048,
				"State": "Active",
				"LastUpdatedBy": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Add SSH Key
		https://developer.centrify.com/reference#post_servermanage-addsshkey

		Request body format
		{
			"Name": "Test Key",
			"Comment": "Test Key",
			"jsutil-enhancecheckbox-34778-inputEl": false,
			"PrivateKey": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCA...AwOQ==\n-----END RSA PRIVATE KEY-----",
			"Type": "Manual"
		}

		Respond Result
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

	Update SSH Key
	https://developer.centrify.com/reference#post_servermanage-updatesshkey

		Request body format
		{
			"IsManaged": false,
			"IsFavorite": false,
			"Revision": 0,
			"CreatedBy": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Test Key",
			"KeyType": "RSA",
			"Comment": "Test Key 1",
			"KeyLength": 2048,
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"State": "Active",
			"jsutil-enhancecheckbox-35013-inputEl": false
		}

		Respond Result
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

	Delete SSH Key
	https://developer.centrify.com/reference#post_servermanage-deletesshkey

		Reqeust body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond Result
		{
			"success": true,
			"Result": "Successfully deleted SSH Key Test Key",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Retrieve SSH Key
	https://developer.centrify.com/reference#post_servermanage-retrievesshkey

	Reqeust body format
		{
			"KeyPairType": "PrivateKey",
			"KeyFormat": "PEM",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond Result
		{
			"success": true,
			"Result": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCA...AwOQ==\n-----END RSA PRIVATE KEY-----",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}
*/
