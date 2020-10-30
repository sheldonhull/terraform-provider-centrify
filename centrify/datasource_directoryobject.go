package centrify

import (
	"errors"
	"fmt"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceDirectoryObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDirectoryObjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the directory object",
			},
			"object_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the directory object",
				ValidateFunc: validation.StringInSlice([]string{
					"User",
					"Group",
				}, false),
			},
			"system_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "UPN of the directory object",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display name of the directory object",
			},
			"distinguished_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Distinguished name of the directory object",
			},
			"forest": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Forest name of the directory object",
			},
			"directory_services": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of UUID of directory services",
			},
		},
	}
}

func dataSourceDirectoryObjectRead(d *schema.ResourceData, m interface{}) error {
	LogD.Printf("Finding Directory Object")
	client := m.(*restapi.RestClient)
	object := NewDirectoryObjects(client)
	object.queryName = d.Get("name").(string)
	object.ObjectType = d.Get("object_type").(string)
	object.DirectoryServices = flattenSchemaSetToStringSlice(d.Get("directory_services"))

	err := object.Read()
	if err != nil {
		return fmt.Errorf("Error retrieving directory services: %s", err)
	}

	var results []DirectoryObject
	// Further narrow down with Distinguished if specified
	dn := d.Get("distinguished_name").(string)
	if dn != "" {
		for _, v := range object.DirectoryObjects {
			if dn == v.DistinguishedName {
				results = append(results, v)
			}
		}
	} else {
		results = object.DirectoryObjects
	}

	if len(results) == 0 {
		return errors.New("Query returns 0 object")
	}
	if len(results) > 1 {
		return fmt.Errorf("Query returns too many objects (found %d, expected 1)", len(results))
	}

	var result = results[0]
	d.SetId(result.ID)
	d.Set("ID", result.ID)
	d.Set("display_name", result.DisplayName)
	d.Set("name", result.Name)
	d.Set("forest", result.Forest)
	d.Set("system_name", result.SystemName)

	return nil
}
