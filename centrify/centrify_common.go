package centrify

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
	jsoniter "github.com/json-iterator/go"
)

// Define this for convience usage
type keyValue map[string]interface{}

var setPermissions = map[string]string{"Grant": "Grant", "View": "View", "Edit": "Edit", "Delete": "Delete"}
var winnixPermissions = map[string]string{"Grant": "Grant", "View": "View", "ManageSession": "ManageSession", "Edit": "Edit", "Delete": "Delete", "AgentAuth": "AgentAuth", "OfflineRescue": "OfflineRescue", "AddAccount": "AddAccount", "UnlockAccount": "UnlockAccount", "RequestZoneRole": "RequestZoneRole"}
var systemPermissions = map[string]string{"Grant": "Grant", "View": "View", "ManageSession": "ManageSession", "Edit": "Edit", "Delete": "Delete", "AgentAuth": "AgentAuth", "OfflineRescue": "OfflineRescue", "AddAccount": "AddAccount", "UnlockAccount": "UnlockAccount", "RequestZoneRole": "RequestZoneRole"}
var databasePermissions = map[string]string{"Grant": "Grant", "View": "View", "Edit": "Edit", "Delete": "Delete"}
var domainPermissions = map[string]string{"Grant": "Grant", "View": "View", "Edit": "Edit", "Delete": "Delete", "UnlockAccount": "UnlockAccount", "AddAccount": "AddAccount"}
var accountPermissions = map[string]string{"Grant": "Owner", "View": "View", "Checkout": "Naked", "Login": "Login", "FileTransfer": "FileTransfer", "Edit": "Manage", "Delete": "Delete", "UpdatePassword": "UpdatePassword", "WorkspaceLogin": "UserPortalLogin", "RotatePassword": "RotatePassword"}
var dbaccountPermissions = map[string]string{"Grant": "Owner", "View": "View", "Checkout": "Naked", "Edit": "Manage", "Delete": "Delete", "UpdatePassword": "UpdatePassword", "RotatePassword": "RotatePassword"}
var domainaccountPermissions = map[string]string{"Grant": "Owner", "View": "View", "Checkout": "Naked", "Login": "Login", "FileTransfer": "FileTransfer", "Edit": "Manage", "Delete": "Delete", "UpdatePassword": "UpdatePassword", "RotatePassword": "RotatePassword"}
var secretPermissions = map[string]string{"Grant": "Grant", "View": "View", "Edit": "Edit", "Delete": "Delete", "RetrieveSecret": "Retrieve"}
var sshkeyPermissions = map[string]string{"Grant": "Owner", "View": "View", "Retrieve": "Checkout", "Edit": "Manage", "Delete": "Delete"}
var servicePermissions = map[string]string{"Grant": "Grant", "Edit": "Edit", "Delete": "Delete"}
var appPermissions = map[string]string{"Grant": "Grant", "View": "View", "Run": "Execute"}
var folderPermissions = map[string]string{"Grant": "Grant", "View": "View", "Edit": "Edit", "Delete": "Delete", "Add": "Add"}

var (
	// LogE logs error message
	LogE = log.New(LogWriter{}, "[ERROR] ", 0)
	// LogW logs warning message
	LogW = log.New(LogWriter{}, "[WARN] ", 0)
	// LogI logs information message
	LogI = log.New(LogWriter{}, "[INFO] ", 0)
	// LogD logs debug message
	LogD = log.New(LogWriter{}, "[DEBUG] ", 0)

	subArgs = make(map[string]interface{})
)

// LogWriter struct
type LogWriter struct{}

func init() {
	subArgs["Caching"] = -1
	//subArgs["PageSize"] = 10000
	//subArgs["Limit"] = 10000
}

func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	if logPath != "" {
		logf, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer logf.Close()
		log.SetOutput(logf)
	}

	log.Printf("%s:%d %s: %s", filepath.Base(file), line, fnName, p)
	return len(p), nil
}

// fillWithMap takes map as input and populate struct attribute accordingly
//func fillWithMap(i vaultObjectInterface, m map[string]interface{}) error {
func fillWithMap(i interface{}, m map[string]interface{}) error {
	jsonString, _ := json.Marshal(m)
	err := json.Unmarshal(jsonString, i)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal map: %v", err)
	}
	//LogD.Printf("Unmarshal object %+v", i)

	return nil
}

//func generateRequestMap(i vaultObjectInterface) (map[string]interface{}, error) {
func generateRequestMap(i interface{}) (map[string]interface{}, error) {
	var mapData = make(map[string]interface{})
	dataBytes, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, err
	}

	return mapData, nil
}

//func generateSchemaMap(i vaultObjectInterface) (map[string]interface{}, error) {
func generateSchemaMap(i interface{}) (map[string]interface{}, error) {
	var mapData = make(map[string]interface{})
	schemaJSON := jsoniter.Config{TagKey: "schema", OnlyTaggedField: true}.Froze()
	dataBytes, err := schemaJSON.Marshal(i)
	if err != nil {
		panic(err)
	}
	err = schemaJSON.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, err
	}

	return mapData, nil
}

// flattenNestedMap converts nested map to flat map. It is used by Policy object.
// It is assumed that json tag of each nested struct element is unique
func flattenNestedMap(flatMap map[string]interface{}, nestedMap interface{}) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}:
			if err := flattenNestedMap(flatMap, v); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch nestedMap.(type) {
	case map[string]interface{}:
		for k, v := range nestedMap.(map[string]interface{}) {
			assign(k, v)
		}
	default:
		return errors.New("Not a valid input, must be a map")
	}

	return nil
}

func flattenSettings(flatten map[string]interface{}, nestedMap interface{}) error {
	if nestedMap != nil {
		for k1, v1 := range nestedMap.(map[string]interface{}) {
			// This is first level that deals with CentrifyServices, CentrifyClient, CentrifyCSSServer, etc.
			switch v1.(type) {
			case map[string]interface{}:
				// this is second level that deals with actual attributes but also may be map such as
				// 		/Core/Css/WinClient/AuthenticationRules
				// 		/Core/PasswordReset/ADAdminPass
				for k2, v2 := range v1.(map[string]interface{}) {
					flatten[k2] = v2
				}
			default:
				flatten[k1] = v1
			}
		}
	}
	return nil
}

func queryVaultObject(client *restapi.RestClient, query string) (map[string]interface{}, error) {
	var queryArg = make(map[string]interface{})
	queryArg["Script"] = query
	queryArg["Args"] = subArgs

	//LogD.Printf("Query arguments: %+v", queryArg)
	//fmt.Printf("Query arguments: %+v\n", queryArg)
	resp, err := client.CallGenericMapAPI("/RedRock/query", queryArg)
	//LogD.Printf("Query response from tenant: %v", resp)
	//fmt.Printf("Query response from tenant: %v\n", resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	// Loop through respond results and grab the first record
	var results = resp.Result["Results"].([]interface{})

	if len(results) == 0 {
		return nil, errors.New("Query returns 0 object")
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("Query returns too many objects (found %d, expected 1)", len(results))
	}
	var result = results[0].(map[string]interface{})
	var row = result["Row"].(map[string]interface{})
	//LogD.Printf("Retrieved row: %v", row)

	return row, nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

// Insert an element into slice at position i
func insert(a []map[string]interface{}, c map[string]interface{}, i int) []map[string]interface{} {
	return append(a[:i], append([]map[string]interface{}{c}, a[i:]...)...)
}

// Find the intersection of two iterable values.
func intersect(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		if contains(b, el) {
			set = append(set, el)
		}
	}

	return set
}

func contains(a interface{}, e interface{}) bool {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true
		}
	}
	return false
}
