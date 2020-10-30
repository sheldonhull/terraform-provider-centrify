package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// MultiplexedAccount - Encapsulates a single MultiplexedAccount
type MultiplexedAccount struct {
	vaultObject

	RealAccount1ID string   `json:"RealAccount1ID,omitempty" schema:"account1_id,omitempty"`
	RealAccount2ID string   `json:"RealAccount2ID,omitempty" schema:"account2_id,omitempty"`
	RealAccount1   string   `json:"RealAccount1,omitempty" schema:"account1,omitempty"`
	RealAccount2   string   `json:"RealAccount2,omitempty" schema:"account2,omitempty"`
	ActiveAccount  string   `json:"ActiveAccount,omitempty" schema:"active_account,omitempty"`
	RealAccounts   []string `json:"RealAccounts,omitempty" schema:"accounts,omitempty"`
}

// NewMultiplexedAccount is a MultiplexedAccount constructor
func NewMultiplexedAccount(c *restapi.RestClient) *MultiplexedAccount {
	s := MultiplexedAccount{}
	s.client = c
	s.MyPermissionList = map[string]string{"Grant": "Grant", "Edit": "Edit", "Delete": "Delete"}
	s.apiRead = "/Subscriptions/GetMultiplexedAccount"
	s.apiCreate = "/Subscriptions/CreateMPAccount"
	s.apiDelete = "/Subscriptions/DeleteMPAccount"
	s.apiUpdate = "/Subscriptions/UpdateMPAccount"
	s.apiPermissions = "/Subscriptions/SetMultiplexedAccountPermissions"

	return &s
}

// Read function fetches a MultiplexedAccount from source, including attribute values. Returns error if any
func (o *MultiplexedAccount) Read() error {
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

// Create function creates a new MultiplexedAccount
func (o *MultiplexedAccount) Create() (*restapi.StringResponse, error) {
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

// Delete function deletes a MultiplexedAccount
func (o *MultiplexedAccount) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing MultiplexedAccount
func (o *MultiplexedAccount) Update() (*restapi.StringResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}

	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}

	LogD.Printf("Generated Map for Update(): %+v", queryArg)
	reply, err := o.client.CallStringAPI(o.apiUpdate, queryArg)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Query function returns a single MultiplexedAccount object in map format
func (o *MultiplexedAccount) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM MultiplexedAccount WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

/*
Get Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-getmultiplexedaccount

	Request body
	{
		"ID": "399db3d4-473a-452f-bfdf-6f8c000c545d",
		"RRFormat": true,
		"Args": {
			"PageNumber": 1,
			"Limit": 1,
			"PageSize": 1,
			"Caching": -1
		}
	}

	Responde Result
	{
		"success": true,
		"Result": {
			"RealAccount2ID": "3ab3d4e3-a7b9-4992-b5a2-b2b983df86b0",
			"Name": "My Multiplexed Account",
			"_RowKey": "399db3d4-473a-452f-bfdf-6f8c000c545d",
			"ActiveAccount": "csvc_acct1 (demo.lab)",
			"RealAccount1ID": "8835f86f-36fd-481d-81ff-e28fc3079f1b",
			"Issues": "",
			"Description": "My Multiplexed Account -",
			"RealAccount2": "csvc_acct2 (demo.lab)",
			"RealAccount1": "csvc_acct1 (demo.lab)"
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Create Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-deletempaccount

	Request body
	{
		"RealAccount1ID": "8835f86f-36fd-481d-81ff-e28fc3079f1b",
		"RealAccount2ID": "3ab3d4e3-a7b9-4992-b5a2-b2b983df86b0",
		"Name": "My Multiplexed Account",
		"Description": "My Multiplexed Account",
		"RealAccount1": "csvc_acct1 (demo.lab)",
		"RealAccount2": "csvc_acct2 (demo.lab)",
		"RealAccounts": [
			"8835f86f-36fd-481d-81ff-e28fc3079f1b",
			"3ab3d4e3-a7b9-4992-b5a2-b2b983df86b0"
		]
	}

	Responde Result
	{
		"success": true,
		"Result": "399db3d4-473a-452f-bfdf-6f8c000c545d",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Update Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-updatempaccount

	Request body
	{
		"RealAccount2ID": "3ab3d4e3-a7b9-4992-b5a2-b2b983df86b0",
		"Name": "My Multiplexed Account",
		"_RowKey": "399db3d4-473a-452f-bfdf-6f8c000c545d",
		"ActiveAccount": "csvc_acct1 (demo.lab)",
		"RealAccount1ID": "8835f86f-36fd-481d-81ff-e28fc3079f1b",
		"Description": "My Multiplexed Account -",
		"RealAccount2": "csvc_acct2 (demo.lab)",
		"RealAccount1": "csvc_acct1 (demo.lab)",
		"ID": "399db3d4-473a-452f-bfdf-6f8c000c545d",
		"RealAccounts": [
			"8835f86f-36fd-481d-81ff-e28fc3079f1b",
			"3ab3d4e3-a7b9-4992-b5a2-b2b983df86b0"
		]
	}

	Responde Result
	{
		"success": true,
		"Result": "399db3d4-473a-452f-bfdf-6f8c000c545d",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Delete Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-updatempaccount

	Request body
	{
		"ID": "399db3d4-473a-452f-bfdf-6f8c000c545d"
	}

	Responde Result
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
