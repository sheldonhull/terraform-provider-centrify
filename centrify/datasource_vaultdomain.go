package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVaultDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVaultDomainRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the domain",
			},
		},
	}
}

func dataSourceVaultDomainRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding domain")
	client := m.(*restapi.RestClient)
	object := NewVaultDomain(client)
	object.Name = d.Get("name").(string)

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving vault object: %s", err)
	}

	//LogD.Printf("Found domain: %+v", result)
	d.SetId(result["ID"].(string))
	d.Set("name", result["Name"].(string))

	return nil
}
