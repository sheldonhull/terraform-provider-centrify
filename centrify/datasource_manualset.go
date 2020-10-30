package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceManualSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManualSetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the manual set",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of set",
				ValidateFunc: validation.StringInSlice([]string{
					"Server",
					"VaultAccount",
					"VaultDatabase",
					"VaultDomain",
					"DataVault",
					"SshKeys",
					"Subscriptions",
					"Application",
					"ResourceProfiles",
				}, false),
			},
			"subtype": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SubObjectType for application.",
				ValidateFunc: validation.StringInSlice([]string{
					"Web",
					"Desktop",
				}, false),
			},
			// computed attributes
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of an manual set",
			},
		},
	}
}

func dataSourceManualSetRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding Manual Set")
	client := m.(*restapi.RestClient)
	object := NewManualSet(client)
	object.Name = d.Get("name").(string)
	object.ObjectType = d.Get("type").(string)
	object.SubObjectType = d.Get("subtype").(string)

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving vault object: %s", err)
	}

	if result["ID"] == nil {
		return fmt.Errorf("ManualSet ID is not set")
	}
	d.SetId(result["ID"].(string))
	d.Set("name", result["Name"].(string))
	if result["Description"] != nil {
		d.Set("description", result["Description"].(string))
	}

	return nil
}
