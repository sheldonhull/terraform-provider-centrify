/*
data "centrifyvault_vaultsystem" "member2" {
    name = "member2"
    fqdn = "member2.demo.lab"
    computer_class = "Windows"
}

data "centrifyvault_vaultsystem" "aws_console" {
    name = "AWS Console"
    fqdn = "192.168.18.15"
    computer_class = "CustomSsh"
}

data "centrifyvault_vaultdomain" "demo_lab" {
    name = "demo.lab"
}

data "centrifyvault_vaultaccount" "shared_account" {
    name = "shared_account"
    domain_id = data.centrifyvault_vaultdomain.demo_lab.id
}

data "centrifyvault_vaultaccount" "selabuser" {
    name = "selabuser"
    host_id = data.centrifyvault_vaultsystem.aws_console.id
}

resource "centrifyvault_desktopapp" "desktopapp1" {
    name = "Test Desktop App 1"
    description = "First Test Desktop Application"
    application_host_id = data.centrifyvault_vaultsystem.member2.id
    login_credential_type = "SharedAccount"
    application_account_id = data.centrifyvault_vaultaccount.shared_account.id
    application_alias = "pas_desktopapp"
    
    command_line = "--ini=ini\\web_aws_iamuser_webdriver.ini --aws_account={system.Description} --username={user.User} --password={user.Password}"
    command_parameter {
        name = "system"
        type = "Server"
        target_object_id = data.centrifyvault_vaultsystem.aws_console.id
    }
    command_parameter {
        name = "user"
        type = "VaultAccount"
        target_object_id = data.centrifyvault_vaultaccount.selabuser.id
    }
    
    default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "IpAddress"
        condition = "OpInCorpIpRange"
      }
    }
    
    sets = [
        centrifyvault_manualset.all_desktopapps.id
    ]
    permission {
        principal_id = data.centrifyvault_role.sso_role.id
        principal_name = data.centrifyvault_role.sso_role.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Run",
        ]
    }
    
}
*/