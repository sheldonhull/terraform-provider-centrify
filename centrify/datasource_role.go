package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the role",
			},
		},
	}
}

func dataSourceRoleRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding role")
	client := m.(*restapi.RestClient)
	object := NewRole(client)
	object.Name = d.Get("name").(string)

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving vault object: %s", err)
	}

	//LogD.Printf("Found role: %+v", result)
	d.SetId(result["ID"].(string))
	d.Set("name", result["Name"].(string))

	return nil
}
