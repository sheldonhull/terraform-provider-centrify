package centrify

import (
	"errors"
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceDirectoryService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDirectoryServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Directory Service",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the Directory Service",
				ValidateFunc: validation.StringInSlice([]string{
					"Centrify Directory",
					"Active Directory",
					"Federated Directory",
					"Google Directory",
					"LDAP Directory",
				}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of the Directory Service",
			},
		},
	}
}

func dataSourceDirectoryServiceRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding DirectoryService")
	client := m.(*restapi.RestClient)
	object := NewDirectoryServices(client)

	err := object.Read()
	if err != nil {
		return fmt.Errorf("Error retrieving directory services: %s", err)
	}

	name := d.Get("name").(string)
	var dirtype string
	switch d.Get("type").(string) {
	case "Centrify Directory":
		dirtype = "CDS"
	case "Active Directory":
		dirtype = "AdProxy"
	case "Federated Directory":
		dirtype = "FDS"
	case "LDAP Directory":
		dirtype = "LdapProxy"
	}

	var results []DirectoryService
	for _, v := range object.DirServices {
		if dirtype == v.Service && name == v.Config {
			results = append(results, v)
		}
	}
	if len(results) == 0 {
		return errors.New("Query returns 0 object")
	}
	if len(results) > 1 {
		return fmt.Errorf("Query returns too many objects (found %d, expected 1)", len(results))
	}

	var result = results[0]
	//LogD.Printf("Found connector: %+v", result)
	d.SetId(result.ID)
	d.Set("name", result.Config)
	d.Set("status", result.Status)
	d.Set("type", name)

	return nil
}
