package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// DirectoryObjects -
type DirectoryObjects struct {
	client            *restapi.RestClient
	ObjectType        string // Either user or group
	queryName         string
	DirectoryServices []string          `json:"DirectoryServices,omitempty" schema:"directory_services,omitempty"`
	DirectoryObjects  []DirectoryObject `json:"DirectoryObjects,omitempty" schema:"directory_object,omitempty"`

	apiRead string
}

// DirectoryObject -
type DirectoryObject struct {
	ID                string `json:"InternalName,omitempty" schema:"id,omitempty"`
	Name              string `json:"Name,omitempty" schema:"name,omitempty"`
	SystemName        string `json:"SystemName,omitempty" schema:"system_name,omitempty"`
	DisplayName       string `json:"DisplayName,omitempty" schema:"display_name,omitempty"`
	DistinguishedName string `json:"DistinguishedName,omitempty" schema:"distinguished_name,omitempty"`
	ObjectType        string `json:"ObjectType,omitempty" schema:"object_type,omitempty"`
	Forest            string `json:"Forest,omitempty" schema:"forest,omitempty"`
}

// NewDirectoryObjects is a DirectoryObjects constructor
func NewDirectoryObjects(c *restapi.RestClient) *DirectoryObjects {
	s := DirectoryObjects{}
	s.client = c
	s.apiRead = "/UserMgmt/DirectoryServiceQuery"

	return &s
}

// Read function fetches directory objects from source
func (o *DirectoryObjects) Read() error {
	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	queryArg["directoryServices"] = o.DirectoryServices
	switch o.ObjectType {
	case "User":
		queryArg["user"] = "{\"_and\":[{\"SystemName\":{\"_like\":\"" + o.queryName + "\"}},{\"ObjectType\":\"User\"}]}"
	case "Group":
		queryArg["group"] = "{\"SystemName\":{\"_like\":\"" + o.queryName + "\"}}"
	}

	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}

	var rs map[string]interface{}
	switch o.ObjectType {
	case "User":
		rs = resp.Result["User"].(map[string]interface{})
	case "Group":
		rs = resp.Result["Group"].(map[string]interface{})
	}
	var results = rs["Results"].([]interface{})
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		obj := &DirectoryObject{}
		fillWithMap(obj, row)
		o.DirectoryObjects = append(o.DirectoryObjects, *obj)
	}

	return nil
}

/*
https://developer.centrify.com/reference#post_usermgmt-directoryservicequery

Request body
{
    "user": "{\"_and\":[{\"_or\":[{\"DisplayName\":{\"_like\":\"LAB\"}},{\"givenName\":{\"_like\":\"LAB\"}},{\"sn\":{\"_like\":\"LAB\"}},{\"SystemName\":{\"_like\":\"LAB\"}}]},{\"ObjectType\":\"user\"}]}",
    "directoryServices": [
        "09B9A9B0-6CE8-465F-AB03-65766D33B05E",
        "e09def65-17c3-0f40-c475-a6ee8825611f",
        "C30B30B1-0B46-49AC-8D99-F6279EED7999"
    ],
    "group": "{\"_or\":[{\"DisplayName\":{\"_like\":\"LAB\"}},{\"SystemName\":{\"_like\":\"LAB\"}}]}",
    "roles": "{\"_or\":[{\"_ID\":{\"_like\":\"LAB\"}},{\"Name\":{\"_like\":\"LAB\"}}]}",
    "Args": {
        "PageNumber": 1,
        "PageSize": 100000,
        "Limit": 100000,
        "SortBy": "",
        "direction": "False",
        "Caching": -1
    }
}

Respond result
{
    "success": true,
    "Result": {
        "User": {
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
                            "Type": "b44b1f2d-196c-48f5-825e-c4e40812d20a",
                            "Key": "User",
                            "IsForeignKey": false
                        }
                    ],
                    "Row": {
                        "Description": null,
                        "DisplayName": "Centrify AD Admin",
                        "ObjectType": "User",
                        "DistinguishedName": "CN=Centrify AD Admin,OU=Lab Service Accounts,DC=demo,DC=lab",
                        "DirectoryServiceUuid": "e09def65-17c3-0f40-c475-a6ee8825611f",
                        "SystemName": "ad_admin@demo.lab",
                        "ServiceInstance": "AdProxy_demo.lab",
                        "Locked": false,
                        "InternalName": "b44b1f2d-196c-48f5-825e-c4e40812d20a",
                        "StatusEnum": "Created",
                        "ServiceInstanceLocalized": "Active Directory (demo.lab)",
                        "ServiceType": "AdProxy",
                        "Forest": "demo.lab",
                        "EMail": "ad_admin@demo.lab",
                        "Status": "Not Invited",
                        "Enabled": true
                    }
                }
            ],
            "ReturnID": ""
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

*/
