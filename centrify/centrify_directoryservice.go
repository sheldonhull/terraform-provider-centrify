package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// DirectoryServices - Encapsulates Directory Services
type DirectoryServices struct {
	client      *restapi.RestClient
	DirServices []DirectoryService `json:"DirServices,omitempty" schema:"directory_service,omitempty"`

	apiRead string
}

// DirectoryService represents directory service
type DirectoryService struct {
	ID               string `json:"directoryServiceUuid,omitempty" schema:"id,omitempty"`
	Name             string `json:"Name,omitempty" schema:"name,omitempty"`
	Description      string `json:"Description,omitempty" schema:"description,omitempty"`
	DisplayName      string `json:"DisplayName,omitempty" schema:"displayName,omitempty"`
	DisplayNameShort string `json:"DisplayNameShort,omitempty" schema:"short_name,omitempty"`
	Service          string `json:"Service,omitempty" schema:"service,omitempty"`
	Status           string `json:"Status,omitempty" schema:"status,omitempty"`
	Config           string `json:"Config,omitempty" schema:"config,omitempty"`
	Forest           string `json:"Forest,omitempty" schema:"forest,omitempty"`
}

// NewDirectoryServices is a DirectoryServices constructor
func NewDirectoryServices(c *restapi.RestClient) *DirectoryServices {
	s := DirectoryServices{}
	s.client = c
	s.apiRead = "/core/GetDirectoryServices"

	return &s
}

// GetDirectorServices etches a DirectorServices from source and returns list of map
func (o *DirectoryServices) GetDirectorServices() ([]map[string]interface{}, error) {
	var dirs []map[string]interface{}

	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	var results = resp.Result["Results"].([]interface{})
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		dirs = append(dirs, row)
	}

	return dirs, nil
}

// Read function fetches a DirectorServices from source
func (o *DirectoryServices) Read() error {
	dirs, err := o.GetDirectorServices()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		obj := &DirectoryService{}
		fillWithMap(obj, dir)
		o.DirServices = append(o.DirServices, *obj)
	}
	LogD.Printf("Filled object: %+v", o)

	return nil
}

/*
https://developer.centrify.com/reference#post_core-getdirectoryservices

{
    "success": true,
    "Result": {
        "IsAggregate": false,
        "Count": 4,
        "Columns": [
            ...
        ],
        "FullCount": 4,
        "Results": [
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "CDS",
                    "DisplayName": "Centrify Centrify Directory",
                    "Tenant": "XXXXX",
                    "Name": "CDS",
                    "Status": "Active",
                    "Config": "Centrify Directory",
                    "StatusDisplay": "Online",
                    "Everybody": true,
                    "Description": "Centrify Directory",
                    "DisplayNameShort": "Centrify Directory",
                    "directoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "AdProxy",
                    "DisplayName": "Active Directory: example.com",
                    "Tenant": "XXXXX",
                    "Name": "AdProxy_example.com",
                    "Status": "Inactive",
                    "Config": "example.com",
                    "StatusDisplay": "Offline",
                    "Everybody": true,
                    "Description": "Active Directory",
                    "DisplayNameShort": "AD: example.com",
                    "directoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                    "Forest": "example.com"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "AdProxy",
                    "DisplayName": "Active Directory: example.com",
                    "Tenant": "XXXXX",
                    "Name": "AdProxy_example.com",
                    "Status": "Inactive",
                    "Config": "example.com",
                    "StatusDisplay": "Offline",
                    "Everybody": true,
                    "Description": "Active Directory",
                    "DisplayNameShort": "AD: example.com",
                    "directoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                    "Forest": "example.com"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "FDS",
                    "DisplayName": "Federated Directory Service",
                    "Tenant": "XXXXXX",
                    "Name": "FDS",
                    "Status": "Active",
                    "Config": "Federated Directory Service",
                    "StatusDisplay": "Online",
                    "Everybody": true,
                    "Description": "Federated Directory",
                    "DisplayNameShort": "FDS",
                    "directoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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
*/
