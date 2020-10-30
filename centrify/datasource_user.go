package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username in loginid@suffix format",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address",
			},
			"displayname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name",
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding user")
	client := m.(*restapi.RestClient)
	object := NewUser(client)
	object.Name = d.Get("username").(string)

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving vault object: %s", err)
	}

	//LogD.Printf("Found user: %+v", result)
	d.SetId(result["ID"].(string))
	d.Set("username", result["Username"].(string))
	d.Set("email", result["Email"].(string))
	d.Set("displayname", result["DisplayName"].(string))

	return nil
}
