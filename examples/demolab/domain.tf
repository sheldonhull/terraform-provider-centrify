output "domain_id" {
    value = data.centrifyvault_vaultdomain.centrify_lab.id
}

// Import domain
resource "centrifyvault_vaultdomain" "centrify_lab" {
    name = "centrify.lab"
    administrative_account_id = data.centrifyvault_directoryobject.ad_admin.id
    administrative_account_domain = data.centrifyvault_directoryobject.ad_admin.forest
    administrative_account_password = var.PASSWORD
    administrative_account_name = data.centrifyvault_directoryobject.ad_admin.system_name
    auto_domain_account_maintenance = true
    auto_local_account_maintenance = true
    manual_domain_account_unlock = true
    manual_local_account_unlock = true
    enable_password_rotation = true
    password_rotate_interval = 180
    
    permission {
        principal_id = data.centrifyvault_role.sysadmin.id
        principal_name = data.centrifyvault_role.sysadmin.name
        principal_type = "Role"
        rights = [
          "AddAccount",
        ]
    }
}


// Vault helpdesk domain account
resource "centrifyvault_vaultaccount" "helpdesk" {
    name = "helpdesk"
    credential_type = "Password"
    password = var.PASSWORD
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
    managed = true
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.tier1_accounts.id
    ]
}

// Vault service domain account1
resource "centrifyvault_vaultaccount" "svc_acct1" {
    name = "svc_acct1"
    credential_type = "Password"
    password = var.PASSWORD
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
    managed = true
}

// Vault service domain account2
resource "centrifyvault_vaultaccount" "svc_acct2" {
    name = "svc_acct2"
    credential_type = "Password"
    password = var.PASSWORD
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
    managed = true
}
