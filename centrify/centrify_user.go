package centrify

import (
	"errors"

	"github.com/centrify/terraform-provider/cloud-golang-sdk/restapi"
)

// User - Encapsulates a single user
type User struct {
	vaultObject
	apiUpdatePassword string

	Name                    string `json:"Name,omitempty" schema:"username,omitempty"`
	Mail                    string `json:"Mail,omitempty" schema:"email,omitempty"` // Email address
	DisplayName             string `json:"DisplayName,omitempty" schema:"displayname,omitempty"`
	Password                string `json:"Password,omitempty" schema:"password,omitempty"`
	ConfirmPassword         string `json:"confirmPassword,omitempty" schema:"confirm_password,omitempty"`
	PasswordNeverExpire     bool   `json:"PasswordNeverExpire,omitempty" schema:"password_never_expire,omitempty"`          // Password never expires
	ForcePasswordChangeNext bool   `json:"ForcePasswordChangeNext,omitempty" schema:"force_password_change_next,omitempty"` // Require password change at next login
	OauthClient             bool   `json:"OauthClient" schema:"oauth_client"`                                               // Is OAuth confidential client
	SendEmailInvite         bool   `json:"SendEmailInvite,omitempty" schema:"send_email_invite,omitempty"`                  // Send email invite for user profile setup
	OfficeNumber            string `json:"OfficeNumber,omitempty" schema:"office_number,omitempty"`
	HomeNumber              string `json:"HomeNumber,omitempty" schema:"home_number,omitempty"`
	MobileNumber            string `json:"MobileNumber,omitempty" schema:"mobile_number,omitempty"`
	//RedirectMFA             bool   `json:"jsutil-checkbox-2598-inputEl" schema:"redirect_mfa"` // Redirect multi factor authentication to a different user account
	RedirectMFAUserID string `json:"CmaRedirectedUserUuid" schema:"redirect_mfa_user_id"` // Redirect multi factor authentication to a different user account
	ReportsTo         string `json:"ReportsTo" schema:"manager_username"`                 // Manager

	// Roles
	Roles []string `json:"Roles,omitempty" schema:"roles,omitempty"`
}

// NewUser is a user constructor
func NewUser(c *restapi.RestClient) *User {
	s := User{}
	s.client = c
	s.apiRead = "/UserMgmt/GetUserAttributes"
	s.apiCreate = "/CDirectoryService/CreateUser"
	s.apiDelete = "/UserMgmt/RemoveUser"
	s.apiUpdate = "/CDirectoryService/ChangeUser"
	s.apiUpdatePassword = "/UserMgmt/ResetUserPassword"

	return &s
}

// Read function fetches a user from source, including attribute values. Returns error if any
func (o *User) Read() error {
	if o.ID == "" {
		return errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)

	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}

	fillWithMap(o, resp.Result)
	LogD.Printf("Filled object: %+v", o)

	return nil
}

// Delete function deletes a user and returns a map that contains deletion result
func (o *User) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Create function creates a new user and returns a map that contains creation result
func (o *User) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}

	LogD.Printf("Generated Map for Create(): %+v", queryArg)

	reply, err := o.client.CallStringAPI(o.apiCreate, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// Update function updates a existing user and returns a map that contains update result
func (o *User) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		return nil, err
	}

	LogD.Printf("Generated Map for Update(): %+v", queryArg)

	reply, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

// ChangePassword function changes user's password
func (o *User) ChangePassword() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		return nil, errors.New("error: ID is empty")
	}
	var queryArg = make(map[string]interface{})
	// Mandatory attributes
	queryArg["ID"] = o.ID
	queryArg["newPassword"] = o.Password

	reply, err := o.client.CallBoolAPI(o.apiUpdatePassword, queryArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	return reply, nil
}

/*
	{
		"DisplayName": "mspadmin",
		"SourceDsInstance": "CDS",
		"ServiceUser": null,
		"DirectoryServiceUuid": "09B9A9B0-6CE8-465F-AB03-65766D33B05E",
		"CloudState": null,
		"LastInvite": null,
		"LastLogin": "/Date(1597376561524)/",
		"UserType": "User",
		"SearchEmail": "ADMIN@DEMO.LAB",
		"SourceDsLocalized": "Centrify Directory",
		"StatusEnum": "Active",
		"SecurityQuestionSet": true,
		"_MatchFilter": null,
		"Email": "admin@demo.lab",
		"Username": "mspadmin@centrify.com.207",
		"Forest": null,
		"SourceDs": "CDS",
		"Status": "Active",
		"ID": "c2c7bcc6-9560-44e0-8dff-5be221cd37ee",
		"SecurityQuestionCount": 1,
		"SourceDsType": "CDS",
		"PhonePinLastChangeDate": "/Date(1582610444169)/"
    }
*/

// Query function returns a single user object in map format
func (o *User) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM User WHERE 1=1"
	if o.Name != "" {
		query += " AND Username='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

/*
	API to manage user
	https://centrify-dev.readme.io/docs/create-and-manage-cloud-directory-users-_new

	Fetch user
	https://developer.centrify.com/reference#post_cdirectoryservice-exemptuserfrommfa

		Request body format
		{
			"ID": "3e34ea2b-0c74-4a28-9307-4908fb94c3a3",
			"DirectoryServiceUuid": "09B9A9B0-6CE8-465F-AB03-65766D33B05E"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"Description": "Test user",
				"ForcePasswordChangeNext": "True",
				"directoryServiceUuid": "09B9A9B0-6CE8-465F-AB03-65766D33B05E",
				"DisplayName": "Test User",
				"PictureUri": "/UserMgmt/GetUserPicture?id=3e34ea2b-0c74-4a28-9307-4908fb94c3a3",
				"CloudState": "None",
				"InEverybodyRole": false,
				"OauthClient": false,
				"MobileNumber": "+65 97323323",
				"LastPasswordChangeDate": "/Date(-62135596800000)/",
				"CreateDate": "/Date(1597323830750)/",
				"OfficeNumber": "+65 97323323",
				"CmaRedirectedUser": "mspadmin@centrify.com.207",
				"SubjectToCloudLocks": true,
				"Alias": "centrify.com.207",
				"HomeNumber": "+65 97323323",
				"ReportsTo": "mspadmin@centrify.com.207",
				"Name": "testuser@centrify.com.207",
				"PreferredCulture": "",
				"Version": "1",
				"Mail": "testuser@centrify.lab",
				"CmaRedirectedUserUuid": "c2c7bcc6-9560-44e0-8dff-5be221cd37ee",
				"Uuid": "3e34ea2b-0c74-4a28-9307-4908fb94c3a3",
				"State": "None"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Create user
	https://developer.centrify.com/reference#post_cdirectoryservice-createuser

		Step 1: Create user
		Request body format
		{
			"LoginName": "testuser",
			"Mail": "testuser@centrify.lab",
			"DisplayName": "Test User",
			"Password": "xxxxxxxxx",
			"confirmPassword": "xxxxxxxxx",
			"PasswordNeverExpire": false,
			"ForcePasswordChangeNext": true,
			"InEverybodyRole": false,
			"OauthClient": false,
			"SendEmailInvite": true,
			"Description": "Test user",
			"OfficeNumber": "+65 97323323",
			"HomeNumber": "+65 97323323",
			"MobileNumber": "+65 97323323",
			"fileName": "centrify_logo.jpg",
			"ID": "",
			"state": "None",
			"jsutil-checkbox-2598-inputEl": true,
			"CmaRedirectedUserUuid": "c2c7bcc6-9560-44e0-8dff-5be221cd37ee",
			"jsutil-text-2601-inputEl": "mspadmin@centrify.com.207",
			"ReportsTo": "mspadmin@centrify.com.207",
			"Name": "testuser@centrify.com.207"
		}

		Respond result
		{
			"success": true,
			"Result": "3e34ea2b-0c74-4a28-9307-4908fb94c3a3",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

		Step 2: Set User Picture

	Update user
	https://developer.centrify.com/reference#post_cdirectoryservice-changeuser

		Request body format
		{
			"LoginName": "testuser",
			"Mail": "testuser@centrify.lab",
			"DisplayName": "Test User test",
			"CloudState": false,
			"PasswordNeverExpire": false,
			"ForcePasswordChangeNext": true,
			"InEverybodyRole": true,
			"OauthClient": false,
			"SendEmailInvite": false,
			"Description": "Test user",
			"OfficeNumber": "+65 97323323",
			"HomeNumber": "+65 97323323",
			"MobileNumber": "+65 97323323",
			"fileName": "",
			"ID": "3e34ea2b-0c74-4a28-9307-4908fb94c3a3",
			"state": "None",
			"jsutil-checkbox-2914-inputEl": true,
			"CmaRedirectedUserUuid": "c2c7bcc6-9560-44e0-8dff-5be221cd37ee",
			"jsutil-text-2917-inputEl": "mspadmin@centrify.com.207",
			"ReportsTo": "mspadmin@centrify.com.207",
			"Name": "testuser@centrify.com.207",
			"Password": "undefined"
		}

		Respond result
		{
			"success": true,
			"Result": null,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete user
	https://developer.centrify.com/reference#post_usermgmt-removeuser

		Request body format
		{
			"ID": "56ad7afe-0044-4b19-a52f-47b8685d048d"
		}

		Respond result
		{
			"success": true,
			"Result": null,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Change password

		Request body format
		{
			"ID": "7c062918-b1f6-4e1e-95ba-b3d0fc12b8f4",
			"newPassword": "xxxxxx"
		}
		Respond result
		{
			"success": true,
			"Result": null,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Assign rights to role
	https://developer.centrify.com/docs/manage-rolesnew#assigning-rights-to-the-role
	https://developer.centrify.com/reference#post_core-getassignedadministrativerights

		Request body format


		Respond result

*/
