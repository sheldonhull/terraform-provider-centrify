
// local_account for member1
resource "centrifyvault_vaultaccount" "member1_localaccount" {
    name = "local_account"
    credential_type = "Password"
    password = var.PASSWORD
    host_id = centrifyvault_vaultsystem.member1.id
    managed = true
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.tier1_accounts.id
    ]
}

// administrator account for member1
resource "centrifyvault_vaultaccount" "member1_administrator" {
    name = "administrator"
    credential_type = "Password"
    password = var.PASSWORD
    host_id = centrifyvault_vaultsystem.member1.id
    managed = false
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.breakglass_accounts.id,
    ]
}

// local_account for member2
resource "centrifyvault_vaultaccount" "member2_localaccount" {
    name = "local_account"
    credential_type = "Password"
    password = var.PASSWORD
    host_id = centrifyvault_vaultsystem.member2.id
    managed = true
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.tier1_accounts.id
    ]
}

// administrator account for member2
resource "centrifyvault_vaultaccount" "member2_administrator" {
    name = "administrator"
    credential_type = "Password"
    password = var.PASSWORD
    host_id = centrifyvault_vaultsystem.member2.id
    managed = false
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.breakglass_accounts.id,
    ]
}

// local_account for centos1
resource "centrifyvault_vaultaccount" "centos1_localaccount" {
    name = "local_account"
    credential_type = "Password"
    password = var.PASSWORD
    host_id = centrifyvault_vaultsystem.centos1.id
    managed = true
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.tier1_accounts.id
    ]
}

// root account for centos1
resource "centrifyvault_vaultaccount" "centos1_root" {
    name = "root"
    credential_type = "Password"
    password = var.PASSWORD
    host_id = centrifyvault_vaultsystem.centos1.id
    managed = false
    is_admin_account = true
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.breakglass_accounts.id,
    ]
}
