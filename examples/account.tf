/*
resource "centrifyvault_vaultaccount" "testsystem1_account1" {
    name = "account1"
    credential_type = "Password"
    password = "xxxxxxxxxxxxxx"
    host_id = centrifyvault_vaultsystem.testsystem1.id
    description = "Account 1 for Windows system 1"
    use_proxy_account = false
    checkout_lifetime = 70
    //managed = true
    default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "IpAddress"
        condition = "OpInCorpIpRange"
      }
    }
}


resource "centrifyvault_vaultaccount" "testsystem2_account1" {
    name = "account1"
    //credential_type = "Password"
    credential_type = "SshKey"
    //password = "xxxxxxxxxxxxxx"
    sshkey_id = "xxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    host_id = centrifyvault_vaultsystem.testsystem2.id
    description = "Account 1 for Linux system 1"
    //use_proxy_account = true
    //checkout_lifetime = 70
    //managed = true
    //default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id

    //sets = [
    //    data.centrifyvault_manualset.lab_all_accounts.id
    //]
    
    permission {
        principal_id = data.centrifyvault_role.sso_role.id
        principal_name = data.centrifyvault_role.sso_role.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Checkout",
          "Edit"
        ]
    }
    
}


resource "centrifyvault_vaultaccount" "testdatabase1_account1" {
    name = "account1"
    credential_type = "Password"
    password = "xxxxxxxxxxxxxx"
    database_id = centrifyvault_vaultdatabase.testdatabase1.id
    description = "Account 1 for database 1"
    checkout_lifetime = 70
    //managed = true
    default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
}
*/