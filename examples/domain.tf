/*
data "centrifyvault_passwordprofile" "domain_pw_pf" {
    name = "Domain Profile"
}

data "centrifyvault_role" "sysadmin" {
    name = "System Administrator"
}

resource "centrifyvault_vaultdomain" "centrify_lab" {
    name = "centrify.lab"
    //description = "centrify.lab domain"
    //checkout_lifetime = 50
    // Advanced menu -> Administrative Account Settings
    //administrative_account_id = data.centrifyvault_directoryobject.ad_admin.id
    //administrative_account_domain = data.centrifyvault_directoryobject.ad_admin.forest
    //administrative_account_password = "xxxxxxxxxx"
    //administrative_account_name = data.centrifyvault_directoryobject.ad_admin.system_name
    //auto_domain_account_maintenance = true
    //auto_local_account_maintenance = true
    //manual_domain_account_unlock = true
    //manual_local_account_unlock = true
    
    
    // Advanced -> Security Settings
    allow_multiple_checkouts = true
    enable_password_rotation = true
    password_rotate_interval = 90
    enable_password_rotation_after_checkin = true
    minimum_password_age = 120
    password_profile_id = data.centrifyvault_passwordprofile.domain_pw_pf.id
    // Advanced -> Maintenance Settings
    enable_password_history_cleanup = true
    password_historycleanup_duration = 100
    // Advanced -> Domain/Zone Tasks
    enable_zone_joined_check = true
    zone_joined_check_interval = 90
    enable_zone_role_cleanup = true
    zone_role_cleanup_interval = 6
    
    permission {
        principal_id = data.centrifyvault_role.sysadmin.id
        principal_name = data.centrifyvault_role.sysadmin.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Edit",
          "Delete",
          "UnlockAccount",
          "AddAccount",
        ]
    }
}
*/