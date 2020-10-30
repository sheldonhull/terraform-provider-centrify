package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePolicyLinks() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyLinksCreate,
		Read:   resourcePolicyLinksRead,
		Update: resourcePolicyLinksUpdate,
		Delete: resourcePolicyLinksDelete,

		Schema: map[string]*schema.Schema{
			"policy_order": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePolicyLinksRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Reading policy links: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	// Create policy links object
	object := NewPolicyLinks(client)
	err := object.Read()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading policy: %v", err)
	}

	d.Set("policy_order", object.Plinks)

	return nil
}

func resourcePolicyLinksCreate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning policy links creation: %s", ResourceIDString(d))

	d.SetId("centrifyvault_policy_links")

	client := m.(*restapi.RestClient)
	object := NewPolicyLinks(client)

	// Upon creating policy links in local state, update the order in tenant as well
	ids := d.Get("policy_order").([]interface{})
	for _, v := range ids {
		plink := PolicyLink{}
		plink.ID = v.(string)
		object.Plinks = append(object.Plinks, plink)
	}
	resp, err := object.Update()
	if err != nil || !resp.Success {
		return fmt.Errorf("Error updating policy links: %v", err)
	}

	// Creation completed
	LogD.Printf("Creation of policy links completed: %s", d.Id())
	return resourcePolicyLinksRead(d, m)
}

func resourcePolicyLinksUpdate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning policy links update: %s", ResourceIDString(d))

	client := m.(*restapi.RestClient)
	object := NewPolicyLinks(client)

	ids := d.Get("policy_order").([]interface{})
	for _, v := range ids {
		plink := PolicyLink{}
		plink.ID = v.(string)
		object.Plinks = append(object.Plinks, plink)
	}

	if d.HasChanges("policy_order") {
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating policy links: %v", err)
		}
	}

	return resourcePolicyLinksRead(d, m)
}

func resourcePolicyLinksDelete(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning deletion of policy links: %s", ResourceIDString(d))

	// We not actually deleting anything from the tenant
	d.SetId("")

	LogD.Printf("Deletion of policy links completed: %s", ResourceIDString(d))
	return nil
}
