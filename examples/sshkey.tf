/*
data "centrifyvault_manualset" "sshkey_set" {
    type = "SshKeys"
    name = "SSHKey Set"
}

data "centrifyvault_sshkey" "my_test_key" {
    name = "My Test Key"
    checkout = true
    key_pair_type = "PrivateKey"
}

resource "centrifyvault_sshkey" "test_key1" {
    name = "Test Key 1"
    description = "Test key 1-"
    private_key = data.centrifyvault_sshkey.my_test_key.ssh_key
    passphrase = "xxxxxxxxxx"
    default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
    
    //sets = [
    //    data.centrifyvault_manualset.sshkey_set.id,
    //]
    
    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "IpAddress"
        condition = "OpInCorpIpRange"
      }
      rule {
        filter = "Browser"
        condition = "OpEqual"
        value = "Chrome"
      }
    }
    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "Browser"
        condition = "OpEqual"
        value = "Chrome"
      }
    }

    permission {
        principal_id = data.centrifyvault_role.sso_role.id
        principal_name = data.centrifyvault_role.sso_role.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Edit",
          "Delete",
          "Retrieve",
        ]
    }
    
}
*/