package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConnector() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Connector",
			},
		},
	}
}

func dataSourceConnectorRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding connector")
	client := m.(*restapi.RestClient)
	object := NewConnector(client)
	object.Name = d.Get("name").(string)

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving vault object: %s", err)
	}

	//LogD.Printf("Found connector: %+v", result)
	d.SetId(result["ID"].(string))
	d.Set("name", result["Name"].(string))

	return nil
}
