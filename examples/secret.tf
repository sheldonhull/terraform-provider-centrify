/*
resource "centrifyvault_vaultsecret" "test_secret1" {
    secret_name = "Test Secret 2"
    description = "Test Secret 2"
    secret_text = "xxxxxxxxxxxxx
    type = "Text"
    //folder_id = centrifyvault_vaultsecretfolder.level2_folder.id
    default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
    //sets = [
    //    data.centrifyvault_manualset.poc_secrets.id,
    //    data.centrifyvault_manualset.test_secrets.id,
    //]
    permission {
        principal_id = data.centrifyvault_role.sso_role.id
        principal_name = data.centrifyvault_role.sso_role.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Edit",
          "Delete",
          "RetrieveSecret",
        ]
    }

    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "IpAddress"
        condition = "OpInCorpIpRange"
      }
    }
    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "DayOfWeek"
        condition = "OpIsDayOfWeek"
        value = "L,1,3,4,5"
      }
      rule {
        filter = "Browser"
				condition = "OpNotEqual"
        value = "Firefox"
      }
      rule {
        filter = "CountryCode"
				condition = "OpNotEqual"
        value = "GA"
      }
    }
}


// When parent folder is enforced with challenge authentication, deletion of subfolder will fail
resource "centrifyvault_vaultsecretfolder" "level1_folder" {
    name = "${var.PREFIX}Level 1 Folder"
    description = "Level 1 Folder"
    //default_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
    permission {
        principal_id = data.centrifyvault_role.sso_role.id
        principal_name = data.centrifyvault_role.sso_role.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Edit"
        ]
    }

    member_permission {
        principal_id = data.centrifyvault_role.lab_infra_admin.id
        principal_name = data.centrifyvault_role.lab_infra_admin.name
        principal_type = "Role"
        rights = [
          "Grant",
          "RetrieveSecret"
        ]
    }

    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "IpAddress"
        condition = "OpInCorpIpRange"
      }
    }
    challenge_rule {
      authentication_profile_id = data.centrifyvault_authenticationprofile.step_up_auth_pf.id
      rule {
        filter = "DayOfWeek"
        condition = "OpIsDayOfWeek"
        value = "L,1,3,4,5"
      }
      rule {
        filter = "Browser"
				condition = "OpNotEqual"
        value = "Firefox"
      }
      rule {
        filter = "CountryCode"
				condition = "OpNotEqual"
        value = "GA"
      }
    }
    
}


resource "centrifyvault_vaultsecretfolder" "level2_folder" {
    name = "${var.PREFIX}Level 2 Folder"
    description = "Level 2 Folder"
    parent_id = centrifyvault_vaultsecretfolder.level1_folder.id
}

resource "centrifyvault_vaultsecretfolder" "level3_folder" {
    name = "${var.PREFIX}Level 3 Folder"
    description = "Level 3 Folder"
    parent_id = centrifyvault_vaultsecretfolder.level2_folder.id
}
*/