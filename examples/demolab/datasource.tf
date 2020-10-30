
// data source for System Administrator role
data "centrifyvault_role" "sysadmin" {
    name = "System Administrator"
}

// data source for centrify.lab domain
data "centrifyvault_directoryservice" "centrify_lab" {
    name = "centrify.lab"
    type = "Active Directory"
}

// data source Centrify Direftory
data "centrifyvault_directoryservice" "centrifydir" {
    name = "Centrify Directory"
    type = "Centrify Directory"
}

// data source for AD accounnt ad_admin@centrify.lab
data "centrifyvault_directoryobject" "ad_admin" {
    directory_services = [
        data.centrifyvault_directoryservice.centrify_lab.id
    ]
    name = "ad_admin"
    object_type = "User"
}

// data source for AD group LAB_GRP_InfraAdmins@centrify.lab
data "centrifyvault_directoryobject" "grp_infra_admins" {
    directory_services = [
        data.centrifyvault_directoryservice.centrify_lab.id
    ]
    name = "LAB_GRP_InfraAdmins"
    object_type = "Group"
}

// data source for AD group LAB_GRP_InfraOwners@centrify.lab
data "centrifyvault_directoryobject" "grp_infra_owners" {
    directory_services = [
        data.centrifyvault_directoryservice.centrify_lab.id
    ]
    name = "LAB_GRP_InfraOwners"
    object_type = "Group"
}

// data source for Default Policy
data "centrifyvault_policy" "default_policy" {
    name = "Default Policy"
}

// Existing authentication profile - "Default Other Login Profile"
data "centrifyvault_authenticationprofile" "default_other_login" {
    name = "Default Other Login Profile"
}

// Existing domain account ad_admin
data "centrifyvault_vaultaccount" "ad_admin" {
    name = "ad_admin"
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
}

// Existing Centrify Directory user admin
data "centrifyvault_user" "admin" {
    username = "admin@centrify.lab"
}

// Existing domain
data "centrifyvault_vaultdomain" "centrify_lab" {
    name = "centrify.lab"
}

/*
data "centrifyvault_vaultaccount" "centos1_localaccount" {
    name = "local_account"
    host_id = centrifyvault_vaultsystem.centos1.id
    checkout = true
}

output "centos1_localaccount" {
    value = data.centrifyvault_vaultaccount.centos1_localaccount.password
}

data "centrifyvault_vaultsecret" "pas_admin_credential" {
    secret_name = "Centrify PAS Admin Credential"
    checkout = true
}

output "pas_admin_credential" {
    value = data.centrifyvault_vaultsecret.pas_admin_credential.secret_text
}
*/