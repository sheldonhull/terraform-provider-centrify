package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceVaultSystem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVaultSystemRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the system",
			},
			"fqdn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Hostname or IP address of the system",
				ValidateFunc: validation.NoZeroValues,
			},
			"computer_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the system",
				ValidateFunc: validation.StringInSlice([]string{
					"Windows",
					"Unix",
					"CiscoAsyncOS",
					"CiscoIOS",
					"CiscoNXOS",
					"JuniperJunos",
					"HpNonStopOS",
					"IBMi",
					"CheckPointGaia",
					"PaloAltoNetworksPANOS",
					"F5NetworksBIGIP",
					"VMwareVMkernel",
					"GenericSsh",
					"CustomSsh",
				}, false),
			},
		},
	}
}

func dataSourceVaultSystemRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding system")
	client := m.(*restapi.RestClient)
	object := NewVaultSystem(client)
	object.Name = d.Get("name").(string)
	object.FQDN = d.Get("fqdn").(string)
	if v, ok := d.GetOk("computer_class"); ok {
		object.ComputerClass = v.(string)
	}

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving vault object: %s", err)
	}

	//LogD.Printf("Found system: %+v", result)
	d.SetId(result["ID"].(string))
	d.Set("name", result["Name"].(string))
	d.Set("fqdn", result["FQDN"].(string))
	d.Set("computer_class", result["ComputerClass"].(string))

	return nil
}
