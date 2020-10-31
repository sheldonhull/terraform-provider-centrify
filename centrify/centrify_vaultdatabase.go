package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// VaultDatabase - Encapsulates a single Database
type VaultDatabase struct {
	// Database -> Settings menu related settings
	vaultObject
	apiGetChallenge string
	apiAddToSets    string
	setTable        string

	FQDN                 string `json:"FQDN,omitempty" schema:"hostname,omitempty"`
	DatabaseClass        string `json:"DatabaseClass,omitempty" schema:"database_class,omitempty"` // Valid values are: SQLServer, Oracle, SAPAse
	Port                 int    `json:"Port,omitempty" schema:"port,omitempty"`
	InstanceName         string `json:"InstanceName,omitempty" schema:"instance_name,omitempty"`
	SkipReachabilityTest bool   `json:"SkipReachabilityTest,omitempty" schema:"skip_reachability_test,omitempty"`

	// Database -> Policy menu related settings
	DefaultCheckoutTime int `json:"DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)

	// Database -> Advanced menu related settings
	AllowMultipleCheckouts            bool   `json:"AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts for related accounts
	AllowPasswordRotation             bool   `json:"AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration            int    `json:"PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin bool   `json:"AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                int    `json:"MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	PasswordProfileID                 string `json:"PasswordProfileID,omitempty" schema:"password_profile_id,omitempty"`                                    // Password Complexity Profile
	AllowPasswordHistoryCleanUp       bool   `json:"AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`              // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration    int    `json:"PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"`          // Password history cleanup (days)

	// Database -> Connectors menu related settings
	ProxyCollectionList string `json:"ProxyCollectionList,omitempty" schema:"connector_list,omitempty"` // List of Connectors used
}

// NewVaultDatabase is a Database constructor
func NewVaultDatabase(c *restapi.RestClient) *VaultDatabase {
	s := VaultDatabase{}
	s.client = c
	s.apiRead = "/RedRock/query"
	s.apiCreate = "/ServerManage/AddDatabase"
	s.apiDelete = "/ServerManage/DeleteDatabase"
	s.apiUpdate = "/ServerManage/UpdateDatabase"
	s.apiGetChallenge = "/ServerManage/GetComputerChallenges"
	s.apiAddToSets = "/Collection/UpdateMembersCollection"
	s.apiPermissions = "/ServerManage/SetDatabasePermissions"
	s.setTable = "VaultDatabase"

	return &s
}

// Read function fetches a VaultDatabase from source, including attribute values. Returns error if any
func (o *VaultDatabase) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Script"] = "SELECT * FROM VaultDatabase WHERE VaultDatabase.ID = '" + o.ID + "'"
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
		LogD.Printf("Returning error: VaultDatabase does not exist in tenant")
		return errors.New("VaultDatabase does not exist in tenant")
	} else if len(results) > 1 {
		// this should never happen
		return errors.New("There are more than one VaultDatabase with the same ID in tenant")
	}
	var result = results[0].(map[string]interface{})
	// Populate vaultObject struct with map from response
	var row = result["Row"].(map[string]interface{})
	//LogD.Printf("Input map: %+v", row)
	fillWithMap(o, row)

	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Create function creates a new VaultDatabase and returns a map that contains creation result
func (o *VaultDatabase) Create() (*restapi.StringResponse, error) {
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

// Delete function deletes a VaultDatabase and returns a map that contains deletion result
func (o *VaultDatabase) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing VaultDatabase and returns a map that contains update result
func (o *VaultDatabase) Update() (*restapi.GenericMapResponse, error) {
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

// Query function returns a single database object in map format
func (o *VaultDatabase) Query() (map[string]interface{}, error) {
	query := "SELECT ID, Name FROM VaultDatabase WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.FQDN != "" {
		query += " AND FQDN='" + o.FQDN + "'"
	}
	if o.DatabaseClass != "" {
		query += " AND DatabaseClass='" + o.DatabaseClass + "'"
	}

	return queryVaultObject(o.client, query)
}

/*
	API to manage system

	Fetch Database

		Request body format
		{
			"Script": "SELECT * FROM VaultDatabase WHERE VaultDatabase.ID = 'xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx'",
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
        "IsAggregate": false,
        "Count": 1,
        ...
        "FullCount": 1,
        "Results": [
            {
                "Entities": [
                    {
                        "Type": "VaultDatabase",
                        "Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Name": "Test DB",
                    "LastHealthCheck": "/Date(1597222428470)/",
                    "PasswordRotateInterval": null,
                    "HealthCheckInterval": null,
                    "ProxyCollectionList": null,
                    "PasswordHistoryCleanUpDuration": null,
                    "ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                    "MinimumPasswordAge": null,
                    "Description": "Test DB Instance 2",
                    "HealthStatus": "Unreachable",
                    "IPAddress": null,
                    "PasswordProfileID": null,
                    "InstanceName": "INSTANCE",
                    "FQDN": "127.0.0.1",
                    "DefaultCheckoutTime": null,
                    "AllowPasswordRotation": null,
                    "ReachableError": "_I18N_NoCloudConnectorsError",
                    "AllowPasswordHistoryCleanUp": null,
                    "AllowHealthCheck": null,
                    "ServiceName": null,
                    "LastState": "Unreachable",
                    "PasswordRotateDuration": null,
                    "HealthStatusError": "_I18N_NoCloudConnectorsError",
                    "Reachable": false,
                    "AllowMultipleCheckouts": null,
                    "_MatchFilter": null,
                    "DatabaseClass": "SQLServer",
                    "AllowPasswordRotationAfterCheckin": null,
                    "Port": 1433
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

	Create Database
	https://developer.centrify.com/reference#post_servermanage-adddatabase

		Request body format
		{
			"Name": "Test DB",
			"DatabaseClass": "SQLServer",
			"FQDN": "127.0.0.1",
			"Port": 1433,
			"InstanceName": "INSTANCE",
			"Description": "Test DB Instance",
			"SkipReachabilityTest": false
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

	Update Database
	https://developer.centrify.com/reference#post_servermanage-updatedatabase

		Request body format
		{
			"Name": "Test DB",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Description": "Test DB Instance 2",
			"HealthStatus": "Unreachable",
			"InstanceName": "INSTANCE",
			"FQDN": "127.0.0.1",
			"ReachableError": "_I18N_NoCloudConnectorsError",
			"LastState": "Unreachable",
			"HealthStatusError": "_I18N_NoCloudConnectorsError",
			"Reachable": false,
			"DatabaseClass": "SQLServer",
			"Port": 1433,
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Delete Database
	https://developer.centrify.com/reference#post_servermanage-deletedatabase

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
*/
