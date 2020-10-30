package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePasswordProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePasswordProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of password profile",
			},
			"profile_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of password profile",
			},
			// computed attributes
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of password profile",
			},
		},
	}
}

func dataSourcePasswordProfileRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding password profile")
	client := m.(*restapi.RestClient)
	object := NewPasswordProfile(client)
	object.Name = d.Get("name").(string)
	object.ProfileType = d.Get("profile_type").(string)

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving vault object: %s", err)
	}

	if result["ID"] == nil {
		return fmt.Errorf("Password profile ID is not set")
	}
	//LogD.Printf("Found password profile: %+v", result)
	d.SetId(result["ID"].(string))
	d.Set("name", result["Name"].(string))
	d.Set("profile_type", result["ProfileType"].(string))
	if result["Description"] != nil {
		d.Set("description", result["Description"].(string))
	}

	return nil
}
