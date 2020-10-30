package centrify

import (
	"fmt"
	"strings"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVaultSecretFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceVaultSecretFolderCreate,
		Read:   resourceVaultSecretFolderRead,
		Update: resourceVaultSecretFolderUpdate,
		Delete: resourceVaultSecretFolderDelete,
		Exists: resourceVaultSecretFolderExists,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the secret folder",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of an secret folder",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parent folder ID of an secret folder",
			},
			"parent_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Parent folder path of an secret folder",
			},
			"default_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default Secret Challenge Profile (used if no conditions matched)",
			},
			"permission":        getPermissionSchema(),
			"member_permission": getPermissionSchema(),
			"challenge_rule":    getChallengeRulesSchema(),
		},
	}
}

func resourceVaultSecretFolderExists(d *schema.ResourceData, m interface{}) (bool, error) {
	LogD.Printf("Checking VaultSecretFolder exist: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewVaultSecretFolder(client)
	object.ID = d.Id()
	err := object.Read()

	if err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	LogD.Printf("VaultSecretFolder exists in tenant: %s", object.ID)
	return true, nil
}

func resourceVaultSecretFolderRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Reading VaultSecretFolder: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	// Create a NewVaultSecretFolder object and populate ID attribute
	object := NewVaultSecretFolder(client)
	object.ID = d.Id()
	err := object.Read()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading VaultSecretFolder: %v", err)
	}
	//LogD.Printf("VaultSecretFolder from tenant: %+v", object)
	schemamap, err := generateSchemaMap(object)
	if err != nil {
		return err
	}
	LogD.Printf("Generated Map for resourceVaultSecretFolderRead(): %+v", schemamap)
	for k, v := range schemamap {
		d.Set(k, v)
	}

	LogD.Printf("Completed reading VaultSecretFolder: %s", object.Name)
	return nil
}

func resourceVaultSecretFolderCreate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning VaultSecretFolder creation: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)

	// Create a VaultSecretFolder object and populate all attributes
	object := NewVaultSecretFolder(client)
	err := getCreateSecretFolderData(d, object)
	if err != nil {
		return err
	}
	resp, err := object.Create()
	if err != nil {
		return fmt.Errorf("Error creating VaultSecretFolder: %v", err)
	}

	id := resp.Result
	if id == "" {
		return fmt.Errorf("VaultSecretFolder ID is not set")
	}
	d.SetId(id)
	// Need to populate ID attribute for subsequence processes
	object.ID = id

	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("folder_id")

	// 2nd step to update VaultSecretFolder login profile
	// Create API call doesn't set VaultSecretFolder login profile so need to run update again
	err = getUpdateSecretFolderData(d, object)
	if err != nil {
		return err
	}

	if object.CollectionMembersDefaultProfile != "" || object.ChallengeRules != nil {
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating VaultSecretFolder attribute: %v", err)
		}
		d.SetPartial("default_profile_id")
		d.SetPartial("challenge_rule")
	}

	// Handle Set permissions
	if _, ok := d.GetOk("permission"); ok {
		_, err = object.SetPermissions(false)
		if err != nil {
			return fmt.Errorf("Error setting VaultSecretFolder permissions: %v", err)
		}
		d.SetPartial("permission")
	}

	// Handle Set member permissions
	if _, ok := d.GetOk("member_permission"); ok {
		_, err = object.SetMemberPermissions(false)
		if err != nil {
			return fmt.Errorf("Error setting VaultSecretFolder member permissions: %v", err)
		}
		d.SetPartial("member_permission")
	}

	// Creation completed
	d.Partial(false)
	LogD.Printf("Creation of VaultSecretFolder completed: %s", object.Name)
	return resourceVaultSecretFolderRead(d, m)
}

func resourceVaultSecretFolderUpdate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning VaultSecretFolder update: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)
	object := NewVaultSecretFolder(client)
	object.ID = d.Id()
	err := getUpdateSecretFolderData(d, object)
	if err != nil {
		return err
	}

	// Deal with normal attribute changes first
	if d.HasChanges("name", "description", "default_profile_id", "challenge_rule") {
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating VaultSecret attribute: %v", err)
		}
		LogD.Printf("Updated attributes to: %v", object)
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("default_profile_id")
		d.SetPartial("challenge_rule")
	}

	if d.HasChange("parent_id") {
		_, new := d.GetChange("parent_id")
		object.ParentID = new.(string)
		resp, err := object.MoveFolder()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating VaultSecretFolder attribute: %v", err)
		}
		d.SetPartial("folder_id")
	}

	// Deal with permission changes
	if d.HasChange("permission") {
		old, new := d.GetChange("permission")
		// We don't want to care the details of changes
		// So, let's first remove the old permissions
		var err error
		if old != nil {
			// do not validate old values
			object.Permissions, err = expandPermissions(old, folderPermissions, false)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(true)
			if err != nil {
				return fmt.Errorf("Error removing VaultSecretFolder permissions: %v", err)
			}
		}

		if new != nil {
			object.Permissions, err = expandPermissions(new, folderPermissions, true)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(false)
			if err != nil {
				return fmt.Errorf("Error adding VaultSecretFolder permissions: %v", err)
			}
		}
		d.SetPartial("permission")
	}

	// Deal with member permission changes
	if d.HasChange("member_permission") {
		old, new := d.GetChange("member_permission")
		// We don't want to care the details of changes
		// So, let's first remove the old permissions
		if old != nil {
			var err error
			object.MemberPermissions, err = expandPermissions(old, secretPermissions, false)
			if err != nil {
				return err
			}
			_, err = object.SetMemberPermissions(true)
			if err != nil {
				return fmt.Errorf("Error removing VaultSecretFolder member permissions: %v", err)
			}
		}

		if new != nil {
			var err error
			object.MemberPermissions, err = expandPermissions(new, secretPermissions, true)
			if err != nil {
				return err
			}
			_, err = object.SetMemberPermissions(false)
			if err != nil {
				return fmt.Errorf("Error adding VaultSecretFolder member permissions: %v", err)
			}
		}
		d.SetPartial("member_permission")
	}

	// We succeeded, disable partial mode. This causes Terraform to save all fields again.
	d.Partial(false)
	LogD.Printf("Updating of VaultSecretFolder completed: %s", object.Name)
	return resourceVaultSecretFolderRead(d, m)
}

func resourceVaultSecretFolderDelete(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning deletion of VaultSecretFolder: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewVaultSecretFolder(client)
	object.ID = d.Id()

	// Remove challenge profile first otherwise deletion will fail
	err := getUpdateSecretFolderData(d, object)
	if err != nil {
		return err
	}
	if object.CollectionMembersDefaultProfile != "" || object.ChallengeRules != nil {
		object.CollectionMembersDefaultProfile = ""
		object.ChallengeRules = nil
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating VaultSecretFolder attribute: %v", err)
		}
	}

	resp, err := object.Delete()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		return fmt.Errorf("Error deleting VaultSecretFolder: %v", err)
	}

	if resp.Success {
		d.SetId("")
	}

	LogD.Printf("Deletion of VaultSecretFolder completed: %s", ResourceIDString(d))
	return nil
}

func getCreateSecretFolderData(d *schema.ResourceData, object *VaultSecretFolder) error {
	object.Name = d.Get("name").(string)
	if v, ok := d.GetOk("description"); ok {
		object.Description = v.(string)
	}
	if v, ok := d.GetOk("parent_id"); ok {
		object.ParentID = v.(string)
	}

	return nil
}

func getUpdateSecretFolderData(d *schema.ResourceData, object *VaultSecretFolder) error {
	getCreateSecretFolderData(d, object)

	if v, ok := d.GetOk("default_profile_id"); ok {
		object.CollectionMembersDefaultProfile = v.(string)
	}
	if v, ok := d.GetOk("permission"); ok {
		var err error
		object.Permissions, err = expandPermissions(v, folderPermissions, true)
		if err != nil {
			return err
		}
	}
	if v, ok := d.GetOk("member_permission"); ok {
		var err error
		object.MemberPermissions, err = expandPermissions(v, secretPermissions, true)
		if err != nil {
			return err
		}
	}
	// Challenge rules
	if v, ok := d.GetOk("challenge_rule"); ok {
		object.ChallengeRules = expandChallengeRules(v.([]interface{}))
		// Perform validations
		if err := validateChallengeRules(object.ChallengeRules); err != nil {
			return fmt.Errorf("Schema setting error: %s", err)
		}
	}

	return nil
}
