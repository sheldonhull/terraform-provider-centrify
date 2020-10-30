package centrify

import (
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSSHKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the account",
			},
			"key_pair_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Which key to retrieve from the pair, must be either PublicKey, PrivateKey, or PPK",
				ValidateFunc: validation.StringInSlice([]string{
					"PublicKey",
					"PrivateKey",
					"PPK",
				}, false),
			},
			"key_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "KeyFormat to retrieve the key in - only works for PublicKey",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase to use for encrypting the PrivateKey",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"checkout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to retrieve SSH Key",
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Content of the SSH Key",
			},
		},
	}
}

func dataSourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding SSH Key")
	client := m.(*restapi.RestClient)
	object := NewSSHKey(client)
	object.Name = d.Get("name").(string)
	if v, ok := d.GetOk("key_pair_type"); ok {
		object.KeyPairType = v.(string)
	}
	if v, ok := d.GetOk("passphrase"); ok {
		object.Passphrase = v.(string)
	}
	if v, ok := d.GetOk("key_format"); ok {
		object.KeyFormat = v.(string)
	}

	result, err := object.Query()
	if err != nil {
		return fmt.Errorf("Error retrieving SSH Key: %s", err)
	}

	//LogD.Printf("Found account: %+v", result)
	object.ID = result["ID"].(string)
	d.SetId(object.ID)
	d.Set("name", result["Name"].(string))
	if result["Comment"] != nil {
		d.Set("description", result["Comment"].(string))
	}
	if result["KeyType"] != nil {
		d.Set("key_type", result["KeyType"].(string))
	}

	// Retrieve SSH Key
	if d.Get("checkout").(bool) {
		resp, err := object.RetriveSSHKey()
		if err != nil {
			return fmt.Errorf("Error retrieving SSH Key: %s", err)
		}
		d.Set("ssh_key", resp.Result)
	}

	return nil
}
