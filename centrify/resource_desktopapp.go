package centrify

import (
	"fmt"
	"strings"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDesktopApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceDesktopAppCreate,
		Read:   resourceDesktopAppRead,
		Update: resourceDesktopAppUpdate,
		Delete: resourceDesktopAppDelete,
		Exists: resourceDesktopAppExists,

		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the SSH Key",
				ValidateFunc: validation.StringInSlice([]string{
					"GenericDesktopApplication",
					"Ssms",
					"Toad",
					"VpxClient",
				}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the SSH Key",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the SSH Key",
			},
			"application_host_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application host",
			},
			"login_credential_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host login credential type",
				ValidateFunc: validation.StringInSlice([]string{
					"ADCredential",
					"SetByUser",
					"AlternativeAccount",
					"SharedAccount",
				}, false),
			},
			"application_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host login credential account",
			},
			"application_alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application alias",
			},
			"command_line": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Command line",
			},
			"command_parameter": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      customCommandParamHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the parameter",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the parameter",
							ValidateFunc: validation.StringInSlice([]string{
								"int",
								"date",
								"string",
								"User",
								"Role",
								"Device",
								"Server",
								"VaultAccount",
								"VaultDomain",
								"VaultDatabase",
								"Subscriptions",
								"DataVault",
								"SshKeys",
								"system_profile",
							}, false),
						},
						"target_object_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of selected parameter value",
						},
					},
				},
			},
			"default_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "AlwaysAllowed", // It must to be "--", "AlwaysAllowed", "-1" or UUID of authen profile
				Description: "Default authentication profile ID",
			},
			"sets": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Add to list of Sets",
			},
			"permission":     getPermissionSchema(),
			"challenge_rule": getChallengeRulesSchema(),
			"policy_script": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"challenge_rule"},
				Description:   "Use script to specify authentication rules (configured rules are ignored)",
			},
		},
	}
}

func resourceDesktopAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	LogD.Printf("Checking DesktopApp exist: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewDesktopApp(client)
	object.ID = d.Id()
	err := object.Read()

	if err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	LogD.Printf("DesktopApp exists in tenant: %s", object.ID)
	return true, nil
}

func resourceDesktopAppRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Reading DesktopApp: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	// Create a NewVaultSecret object and populate ID attribute
	object := NewDesktopApp(client)
	object.ID = d.Id()
	err := object.Read()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading DesktopApp: %v", err)
	}
	//LogD.Printf("DesktopApp from tenant: %+v", object)
	schemamap, err := generateSchemaMap(object)
	if err != nil {
		return err
	}
	LogD.Printf("Generated Map for resourceDesktopAppRead(): %+v", schemamap)
	for k, v := range schemamap {
		d.Set(k, v)
	}

	LogD.Printf("Completed reading DesktopApp: %s", object.Name)
	return nil
}

func resourceDesktopAppCreate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning DesktopApp creation: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)

	// Create a DesktopApp object
	object := NewDesktopApp(client)
	object.TemplateName = d.Get("template_name").(string)
	resp, err := object.Create()
	if err != nil {
		return fmt.Errorf("Error creating DesktopApp: %v", err)
	}
	if len(resp.Result) <= 0 {
		return fmt.Errorf("Import application template returns incorrect result")
	}

	id := resp.Result[0].(map[string]interface{})["_RowKey"].(string)

	if id == "" {
		return fmt.Errorf("DesktopApp ID is not set")
	}
	d.SetId(id)
	// Need to populate ID attribute for subsequence processes
	object.ID = id

	// Update attributes to complete creation
	err = getUpateGetDesktopAppData(d, object)
	if err != nil {
		return err
	}

	resp2, err2 := object.Update()
	if err2 != nil || !resp2.Success {
		return fmt.Errorf("Error updating DesktopApp attribute: %v", err2)
	}

	d.SetPartial("name")
	d.SetPartial("template_name")
	d.SetPartial("description")
	d.SetPartial("application_host_id")
	d.SetPartial("login_credential_type")
	d.SetPartial("application_account_id")
	d.SetPartial("application_alias")
	d.SetPartial("command_line")
	d.SetPartial("command_parameter")
	d.SetPartial("default_profile_id")
	d.SetPartial("challenge_rule")
	d.SetPartial("policy_script")

	if len(object.Sets) > 0 {
		for _, v := range object.Sets {
			setObj := NewManualSet(client)
			setObj.ID = v
			setObj.ObjectType = "Application"
			resp, err := setObj.UpdateSetMembers([]string{object.ID}, "add")
			if err != nil || !resp.Success {
				return fmt.Errorf("Error adding DesktopApp to Set: %v", err)
			}
		}
		d.SetPartial("sets")
	}

	// add permissions
	if _, ok := d.GetOk("permission"); ok {
		_, err = object.SetPermissions(false)
		if err != nil {
			return fmt.Errorf("Error setting DesktopApp permissions: %v", err)
		}
		d.SetPartial("permission")
	}

	// Creation completed
	d.Partial(false)
	LogD.Printf("Creation of DesktopApp completed: %s", object.Name)
	return resourceDesktopAppRead(d, m)
}

func resourceDesktopAppUpdate(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning DesktopApp update: %s", ResourceIDString(d))

	// Enable partial state mode
	d.Partial(true)

	client := m.(*restapi.RestClient)
	object := NewDesktopApp(client)
	object.ID = d.Id()
	err := getUpateGetDesktopAppData(d, object)
	if err != nil {
		return err
	}

	// Deal with normal attribute changes first
	if d.HasChanges("name", "template_name", "description", "application_host_id", "login_credential_type", "application_account_id", "application_alias",
		"command_line", "command_parameter", "default_profile_id", "challenge_rule", "policy_script") {
		resp, err := object.Update()
		if err != nil || !resp.Success {
			return fmt.Errorf("Error updating DesktopApp attribute: %v", err)
		}
		LogD.Printf("Updated attributes to: %v", object)
		d.SetPartial("name")
		d.SetPartial("template_name")
		d.SetPartial("description")
		d.SetPartial("application_host_id")
		d.SetPartial("login_credential_type")
		d.SetPartial("application_account_id")
		d.SetPartial("application_alias")
		d.SetPartial("command_line")
		d.SetPartial("command_parameter")
		d.SetPartial("default_profile_id")
		d.SetPartial("challenge_rule")
		d.SetPartial("policy_script")
	}

	if d.HasChange("sets") {
		old, new := d.GetChange("sets")
		// Remove old Sets
		for _, v := range flattenSchemaSetToStringSlice(old) {
			setObj := NewManualSet(client)
			setObj.ID = v
			setObj.ObjectType = "Application"
			resp, err := setObj.UpdateSetMembers([]string{object.ID}, "remove")
			if err != nil || !resp.Success {
				return fmt.Errorf("Error removing DesktopApp from Set: %v", err)
			}
		}
		// Add new Sets
		for _, v := range flattenSchemaSetToStringSlice(new) {
			setObj := NewManualSet(client)
			setObj.ID = v
			setObj.ObjectType = "Application"
			resp, err := setObj.UpdateSetMembers([]string{object.ID}, "add")
			if err != nil || !resp.Success {
				return fmt.Errorf("Error adding DesktopApp to Set: %v", err)
			}
		}
		d.SetPartial("sets")
	}

	// Deal with Permissions
	if d.HasChange("permission") {
		old, new := d.GetChange("permission")
		// We don't want to care the details of changes
		// So, let's first remove the old permissions
		var err error
		if old != nil {
			// do not validate old values
			object.Permissions, err = expandPermissions(old, appPermissions, false)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(true)
			if err != nil {
				return fmt.Errorf("Error removing DesktopApp permissions: %v", err)
			}
		}

		if new != nil {
			object.Permissions, err = expandPermissions(new, appPermissions, true)
			if err != nil {
				return err
			}
			_, err = object.SetPermissions(false)
			if err != nil {
				return fmt.Errorf("Error adding DesktopApp permissions: %v", err)
			}
		}
		d.SetPartial("permission")
	}

	// We succeeded, disable partial mode. This causes Terraform to save all fields again.
	d.Partial(false)
	LogD.Printf("Updating of DesktopApp completed: %s", object.Name)
	return resourceDesktopAppRead(d, m)
}

func resourceDesktopAppDelete(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Beginning deletion of DesktopApp: %s", ResourceIDString(d))
	client := m.(*restapi.RestClient)

	object := NewDesktopApp(client)
	object.ID = d.Id()
	resp, err := object.Delete()

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if err != nil {
		return fmt.Errorf("Error deleting DesktopApp: %v", err)
	}

	if resp.Success {
		d.SetId("")
	}

	LogD.Printf("Deletion of DesktopApp completed: %s", ResourceIDString(d))
	return nil
}

func getUpateGetDesktopAppData(d *schema.ResourceData, object *DesktopApp) error {
	object.Name = d.Get("name").(string)
	if v, ok := d.GetOk("description"); ok {
		object.Description = v.(string)
	}
	if v, ok := d.GetOk("application_host_id"); ok {
		object.DesktopAppRunHostID = v.(string)
	}
	if v, ok := d.GetOk("login_credential_type"); ok {
		object.DesktopAppRunAccountType = v.(string)
	}
	if v, ok := d.GetOk("application_account_id"); ok {
		object.DesktopAppRunAccountID = v.(string)
	}
	if v, ok := d.GetOk("application_alias"); ok {
		object.DesktopAppProgramName = v.(string)
	}
	if v, ok := d.GetOk("command_line"); ok {
		object.DesktopAppCmdline = v.(string)
	}
	if v, ok := d.GetOk("command_parameter"); ok {
		object.DesktopAppParams = expandCommandParams(v)
	}
	if v, ok := d.GetOk("default_profile_id"); ok {
		object.DefaultAuthProfile = v.(string)
	}
	if v, ok := d.GetOk("policy_script"); ok {
		object.PolicyScript = v.(string)
	}
	if v, ok := d.GetOk("sets"); ok {
		object.Sets = flattenSchemaSetToStringSlice(v)
	}
	// Permissions
	if v, ok := d.GetOk("permission"); ok {
		var err error
		object.Permissions, err = expandPermissions(v, appPermissions, true)
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
