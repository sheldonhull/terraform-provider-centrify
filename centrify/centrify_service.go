package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// Service - Encapsulates a single Service
type Service struct {
	vaultObject

	SystemID               string `json:"ComputerID,omitempty" schema:"system_id,omitempty"`
	ServiceType            string `json:"Type,omitempty" schema:"service_type,omitempty"`
	Name                   string `json:"WindowsServiceName,omitempty" schema:"service_name,omitempty"`
	EnableManagement       bool   `json:"IsActive" schema:"enable_management"`
	AdminAccountID         string `json:"PushCreds,omitempty" schema:"admin_account_id,omitempty"`
	MultiplexedAccountID   string `json:"AccountID,omitempty" schema:"multiplexed_account_id,omitempty"`
	RestartService         bool   `json:"RestartService" schema:"restart_service"`
	RestartTimeRestriction bool   `json:"RestartTimeRestriction" schema:"restart_time_restriction"`
	DaysOfWeek             string `json:"DaysOfWeek,omitempty" schema:"days_of_week,omitempty"`
	RestartStartTime       string `json:"RestartStartTime,omitempty" schema:"restart_start_time,omitempty"`
	RestartEndTime         string `json:"RestartEndTime,omitempty" schema:"restart_end_time,omitempty"`
	UseUTCTime             bool   `json:"RestartTimeInUtc" schema:"use_utc_time"`
}

// NewService is a Service constructor
func NewService(c *restapi.RestClient) *Service {
	s := Service{}
	s.client = c
	s.MyPermissionList = map[string]string{"Grant": "Grant", "Edit": "Edit", "Delete": "Delete"}
	s.apiRead = "/Subscriptions/GetSubscription"
	s.apiCreate = "/Subscriptions/AddSubscription"
	s.apiDelete = "/Subscriptions/DeleteSubscription"
	s.apiUpdate = "/Subscriptions/UpdateSubscription"
	s.apiPermissions = "/Subscriptions/SetSubscriptionPermissions"

	return &s
}

// Read function fetches a Service from source, including attribute values. Returns error if any
func (o *Service) Read() error {
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

// Create function creates a new Service
func (o *Service) Create() (*restapi.StringResponse, error) {
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

// Delete function deletes a Service
func (o *Service) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing Service
func (o *Service) Update() (*restapi.StringResponse, error) {
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

// Query function returns a single Service object in map format
func (o *Service) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Subscriptions WHERE 1=1"
	if o.Name != "" {
		query += " AND WindowsServiceName='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

/*
Get Service
https://developer.centrify.com/reference#post_subscriptions-getsubscription

	Request body
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

	Responde Result
	{
		"success": true,
		"Result": {
			"ComputerID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Description": "",
			"WindowsServiceName": "TestWindowsService",
			"SecretID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"_STAMP": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RestartService": true,
			"LogonCreds": "ad_admin (example.com)",
			"_TableName": "cpssubscriptions",
			"_PartitionKey": "XXXXXX",
			"IsActive": true,
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RestartTimeRestriction": false,
			"Type": "WindowsService",
			"Account": "Account for TestWindowsService",
			"PushCreds": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"LastUpdate": "/Date(1595301043898)/",
			"AccountID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Status": "OK",
			"Resource": "member1.example.com",
			"ACL": "true",
			"Mode": "Push",
			"CurrentMultiplexAccount": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RestartTimeInUtc": true
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Create Service
https://developer.centrify.com/reference#post_subscriptions-addsubscription

	Request body
	{
		"RestartTimeInUtc": false,
		"ComputerID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"AccountID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"PushCreds": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Resource": "member2",
		"Description": "TestWinService",
		"Type": "WindowsService",
		"WindowsServiceName": "testwinsrv",
		"IsActive": true,
		"LogonCreds": "ad_admin (example.com)",
		"Account": "test",
		"RestartService": true,
		"RestartTimeRestriction": true,
		"DaysOfWeek": "Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday",
		"RestartStartTime": "10:00",
		"RestartEndTime": "09:00",
		"Mode": "Push"
	}

	Responde Result
	{
		"success":true,
		"Result":"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Message":null,
		"MessageID":null,
		"Exception":null,
		"ErrorID":null,
		"ErrorCode":null,"
		IsSoftError":false,
		"InnerExceptions":null
	}

Update Service
https://developer.centrify.com/reference#post_subscriptions-updatesubscription

	Request body
	{
		"ComputerID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Description": "Test Windows Service",
		"WindowsServiceName": "TestWindowsService",
		"SecretID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"_STAMP": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"RestartService": true,
		"LogonCreds": "ad_admin (example.com)",
		"_TableName": "cpssubscriptions",
		"_PartitionKey": "XXXXX",
		"IsActive": true,
		"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"RestartTimeRestriction": false,
		"Type": "WindowsService",
		"Account": "Account for TestWindowsService",
		"PushCreds": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"AccountID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Status": "OK",
		"Resource": "member1.example.com",
		"ACL": "true",
		"Mode": "Push",
		"CurrentMultiplexAccount": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Responde Result
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

Delete Service
https://developer.centrify.com/reference#post_subscriptions-deletesubscription

	Request body
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Responde Result
	{
		"success":true,
		"Result":null,
		"Message":null,
		"MessageID":null,
		"Exception":null,
		"ErrorID":null,
		"ErrorCode":null,
		"IsSoftError":false,
		"InnerExceptions":null
	}
*/
