package centrify

import (
	"errors"
	"fmt"
	"strings"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// Generic object
type vaultObject struct {
	// Rest client
	client *restapi.RestClient
	// Standard attributes
	ID               string `json:"ID,omitempty" schema:"id,omitempty"`
	Name             string `json:"Name,omitempty" schema:"name,omitempty"`
	Description      string `json:"Description,omitempty" schema:"description,omitempty"`
	MyPermissionList map[string]string
	// Sets
	Sets        []string `json:"Sets,omitempty" schema:"sets,omitempty"`
	Permissions []Permission

	// API endpoints
	apiRead        string //`json:"-"` // Ignoring this JSON field
	apiCreate      string //`json:"-"` // Ignoring this JSON field
	apiDelete      string //`json:"-"` // Ignoring this JSON field
	apiUpdate      string //`json:"-"` // Ignoring this JSON field
	apiPermissions string
}

// Permission represents object permission
type Permission struct {
	PrincipalID   string `json:"PrincipalId,omitempty" schema:"principal_id,omitempty"` // Uuid of the principal
	PrincipalName string `json:"Principal,omitempty" schema:"principal_name,omitempty"` // User name or role name
	PrincipalType string `json:"PType,omitempty" schema:"principal_type,omitempty"`     // Principal type: User, Role etc..
	Rights        string `json:"Rights,omitempty" schema:"rights,omitempty"`            // Permissions: Grant,View,Edit,Delete or None to remove this item
}

// ChallengeRules represents list of login rule set
type ChallengeRules struct {
	Enabled   bool            `json:"Enabled,omitempty" schema:"enabled,omitempty"`
	UniqueKey string          `json:"_UniqueKey,omitempty" schema:"unique_key,omitempty"`
	Type      string          `json:"_Type,omitempty" schema:"type,omitempty"`
	Rules     []ChallengeRule `json:"_Value,omitempty" schema:"rule,omitempty"`
}

// ChallengeRule represents a set of login rule
type ChallengeRule struct {
	ChallengeCondition []ChallengeCondition `json:"Conditions,omitempty" schema:"rule,omitempty"`
	AuthProfileID      string               `json:"ProfileId,omitempty" schema:"authentication_profile_id,omitempty"` // "-1" means Not Allowed
}

// ChallengeCondition represents a single challenge rule
type ChallengeCondition struct {
	Filter    string `json:"Prop,omitempty" schema:"filter,omitempty"`
	Condition string `json:"Op,omitempty" schema:"condition,omitempty"`
	Value     string `json:"Val,omitempty" schema:"value,omitempty"`
}

// deleteObjectBoolAPI a object and returns a map that contains deletion result
func (o *vaultObject) deleteObjectBoolAPI(idfield string) (*restapi.BoolResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var funcArg = make(map[string]interface{})
	if idfield == "" {
		funcArg["ID"] = o.ID
	} else {
		funcArg[idfield] = o.ID
	}

	reply, err := o.client.CallBoolAPI(o.apiDelete, funcArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// deleteObjectMapAPI a object and returns a map that contains deletion result
func (o *vaultObject) deleteObjectMapAPI(idfield string) (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var funcArg = make(map[string]interface{})
	if idfield == "" {
		funcArg["ID"] = o.ID
	} else {
		funcArg[idfield] = o.ID
	}

	reply, err := o.client.CallGenericMapAPI(o.apiDelete, funcArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// deleteObjectStringAPI a object and returns a map that contains deletion result
func (o *vaultObject) deleteObjectStringAPI(idfield string) (*restapi.StringResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var funcArg = make(map[string]interface{})
	if idfield == "" {
		funcArg["ID"] = o.ID
	} else {
		funcArg[idfield] = o.ID
	}

	reply, err := o.client.CallStringAPI(o.apiDelete, funcArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// SetPermissions sets permissions. isRemove indicates whether to remove all permissions instead of setting permissions
func (o *vaultObject) SetPermissions(isRemove bool) (*restapi.BaseAPIResponse, error) {
	//func (o *vaultObject) SetPermissions(isRemove bool) (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}

	var permissions []map[string]interface{}
	for _, v := range o.Permissions {
		var permission = make(map[string]interface{})
		permission, err := generateRequestMap(v)
		if isRemove {
			permission["Rights"] = "None"
		}
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if len(permissions) > 0 {
		var queryArg = make(map[string]interface{})
		queryArg["ID"] = o.ID
		queryArg["PVID"] = o.ID
		queryArg["Grants"] = permissions
		LogD.Printf("Generated Map for SetPermissions(): %+v", queryArg)
		resp, err := o.client.CallBaseAPI(o.apiPermissions, queryArg)
		if err != nil {
			return nil, err
		}
		if !resp.Success {
			return nil, errors.New(resp.Message)
		}
		return resp, nil
	}
	return nil, nil
}

// FillStruct function fills a struct with map
func (o *vaultObject) FillStruct(m map[string]interface{}) error {
	LogD.Printf("Printing input map...")
	LogD.Printf("Input map: %v", m)
	for k, v := range m {
		LogD.Printf("Map key %v, map value: %v", k, v)
		err := setField(o, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func validatePermissions(permissions []Permission, valid []string) error {
	if permissions != nil {
		for _, v := range permissions {
			rights := strings.Split(v.Rights, ",")
			if len(intersect(rights, valid)) != len(rights) {
				return fmt.Errorf("%v can only contain %v", rights, valid)
			}
		}
	}

	return nil
}
