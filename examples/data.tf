
// Existing role
data "centrifyvault_role" "lab_infra_admin" {
    name = "LAB Infrastructure Admins"
}
data "centrifyvault_role" "sso_role" {
    name = "SSO_Role"
}
/*
//Existing user
data "centrifyvault_user" "admin" {
    username = "admin@example.com"
}

// Existing System Set "LAB Systems"
data "centrifyvault_manualset" "lab_systems" {
    type = "Server"
    name = "LAB Systems"
}

// Existing Database Set "LAB_Databases"
data "centrifyvault_manualset" "lab_databases" {
    type = "VaultDatabase"
    name = "LAB_Databases"
}

// Existing Account Set "LAB All Accounts"
data "centrifyvault_manualset" "lab_all_accounts" {
    type = "VaultAccount"
    name = "LAB All Accounts"
}

// Existing Secret Set "POC Secrets"
data "centrifyvault_manualset" "poc_secrets" {
    type = "DataVault"
    name = "POC Secrets"
}
data "centrifyvault_manualset" "test_secrets" {
    type = "DataVault"
    name = "Test Secrets"
}


// Existing password profile
data "centrifyvault_passwordprofile" "test_pw_pf1" {
    //name = "Windows Profile"
    name = "vCenter Profile"
    //profile_type = "UserDefined"
}
*/

// Existing authentication profile
data "centrifyvault_authenticationprofile" "step_up_auth_pf" {
    name = "LAB Step-up Authentication Profile"
    //name = "Default New Device Login Profile"
}

/*
// Connector
data "centrifyvault_connector" "XXXXX-XXXXXX" {
    name = "XXXXX-XXXXXX"
}

// Connector
data "centrifyvault_connector" "dc01" {
    name = "dc01"
}

// Existing domain
data "centrifyvault_vaultdomain" "demo_lab" {
    name = "demo.lab"
}

// Existing system
data "centrifyvault_vaultsystem" "centos1" {
    name = "centos1"
    fqdn = "centos1.demo.lab"
    computer_class = "Unix"
}

// Existing account
data "centrifyvault_vaultaccount" "centos1_clocal_account" {
    name = "clocal_account"
    host_id = data.centrifyvault_vaultsystem.centos1.id
}

// Existing secret
data "centrifyvault_vaultsecret" "test_secret" {
    secret_name = "Centrify PAS Admin Credential"
    checkout = true
}

// Existing secret folder at top level
data "centrifyvault_vaultsecretfolder" "lab_level1_folder" {
    name = "LAB Level 1 Folder"
}

// Existing secret folder at 2nd level
data "centrifyvault_vaultsecretfolder" "lab_level2_folder" {
    name = "LAB Level 2 Folder"
    parent_path = "LAB Level 1 Folder"
}

// Existing secret folder at 2nd level. There can be the same name in other folder so need to search parent_path too.
data "centrifyvault_vaultsecretfolder" "lab_level2_test_folder" {
    name = "Test Folder"
    parent_path = "LAB Level 1 Folder"
}

// Existing secret folder at 3rd level
data "centrifyvault_vaultsecretfolder" "lab_level3_folder" {
    name = "LAB Level 3 Folder"
    parent_path = "LAB Level 1 Folder\\LAB Level 2 Folder"
}

// Existing secret folder at 3rd level. There can be the same name in other folder so need to search parent_path too.
data "centrifyvault_vaultsecretfolder" "lab_level3_test_folder" {
    name = "Test Folder"
    parent_path = "LAB Level 1 Folder\\LAB Level 2 Folder"
}

*/