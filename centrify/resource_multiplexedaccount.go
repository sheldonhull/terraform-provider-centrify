package centrify

import (
	"fmt"
	"strings"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMultiplexedAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceMultiplexedAccountCreate,
		Read:   resourceMultiplexedAccountRead,
		Update: resourceMultiplexedAccountUpdate,
		Delete: resourceMultiplexedAccountDelete,
		Exists: resourceMultiplexedAccountExists,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the multiplexed account",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the multiplexed account",
			},
			"account1_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account2_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"accounts": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 2,
				MaxItems: 2,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"active_account": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"permission": getPermissionSchema(),
		},
	}
}

func resourceMultiplexedAccountExists(d *schema.ResourceData, m interface{}) (bool, error) {
	LogD.Printf("Checking multiplexed account exist: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewMultiplexedAccount(client)
	object.ID = d.Id()
	err := object.Read()

	if err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	LogD.Printf("Multiplexed account exists in tenant: %s", object.ID)
	return true, nil
}

func resourceMultiplexedAccountRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Reading multiplexed account: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	// Create a NewMultiplexedAccount object and populate ID attribute
	object := NewMultiplexedAccount(client)
	object.ID = d.Id()
	err := object.Read()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading multiplexed account: %v", err)
	}
	//LogD.Printf("Multiplexed account from tenant: %+v", object)
	schemamap, err := generateSchemaMap(object)
	if err != nil {
		return err
	}
	LogD.Printf("Generated Map for resourceMultiplexedAccountRead(): %+v", schemamap)
	for k, v := range schemamap {
		d.Set(k, v)
	}

	LogD.Printf("Completed reading multiplexed account: %s", object.Name)
	return nil
}

func resourceMultiplexedAccountCreate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning multiplexed account creation: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)

	// Create a multiplexed account object and populate all attributes
	object := NewMultiplexedAccount(client)
	err := createUpateGetMultiplexedAccountData(d, object)
	if err != nil {
		return err
	}

	resp, err := object.Create()
	if err != nil {
		return fmt.Errorf("Error creating multiplexed account: %v", err)
	}

	id := resp.Result
	if id == "" {
		return fmt.Errorf("Multiplexed account ID is not set")
	}
	d.SetId(id)
	// Need to populate ID attribute for subsequence processes
	object.ID = id

	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("accounts")

	// add permissions
	if _, ok := d.GetOk("permission"); ok {
		_, err = object.SetPermissions(false)
		if err != nil {
			return fmt.Errorf("Error setting multiplexed account permissions: %v", err)
		}
		d.SetPartial("permission")
	}

	// Creation completed
	d.Partial(false)
	LogD.Printf("Creation of multiplexed account completed: %s", object.Name)
	return resourceMultiplexedAccountRead(d, m)
}

func resourceMultiplexedAccountUpdate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning multiplexed account update: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)
	object := NewMultiplexedAccount(client)
	object.ID = d.Id()
	err := createUpateGetMultiplexedAccountData(d, object)
	if err != nil {
		return err
	}

	// Deal with normal attribute changes first
	if d.HasChanges("name", "description", "accounts") {
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating multiplexed account attribute: %v", err)
		}
		LogD.Printf("Updated attributes to: %v", object)
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("accounts")
	}

	// Deal with Permissions
	if d.HasChange("permission") {
		old, new := d.GetChange("permission")
		// We don't want to care the details of changes
		// So, let's first remove the old permissions
		var err error
		if old != nil {
			// do not validate old values
			object.Permissions, err = expandPermissions(old, object.MyPermissionList, false)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(true)
			if err != nil {
				return fmt.Errorf("Error removing multiplexed account permissions: %v", err)
			}
		}

		if new != nil {
			object.Permissions, err = expandPermissions(new, object.MyPermissionList, true)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(false)
			if err != nil {
				return fmt.Errorf("Error adding multiplexed account permissions: %v", err)
			}
		}
		d.SetPartial("permission")
	}

	// We succeeded, disable partial mode. This causes Terraform to save all fields again.
	d.Partial(false)
	LogD.Printf("Updating of multiplexed account completed: %s", object.Name)
	return resourceMultiplexedAccountRead(d, m)
}

func resourceMultiplexedAccountDelete(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning deletion of multiplexed account: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewMultiplexedAccount(client)
	object.ID = d.Id()

	resp, err := object.Delete()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		return fmt.Errorf("Error deleting multiplexed account: %v", err)
	}

	if resp.Success {
		d.SetId("")
	}

	LogD.Printf("Deletion of multiplexed account completed: %s", ResourceIDString(d))
	return nil
}

func createUpateGetMultiplexedAccountData(d *schema.ResourceData, object *MultiplexedAccount) error {
	object.Name = d.Get("name").(string)
	if v, ok := d.GetOk("description"); ok {
		object.Description = v.(string)
	}
	object.RealAccounts = flattenSchemaSetToStringSlice(d.Get("accounts"))

	// Permissions
	if v, ok := d.GetOk("permission"); ok {
		var err error
		object.Permissions, err = expandPermissions(v, object.MyPermissionList, true)
		if err != nil {
			return err
		}
	}

	return nil
}
