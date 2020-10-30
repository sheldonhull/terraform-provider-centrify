
// Windows system member1
resource "centrifyvault_vaultsystem" "member1" {
    name = "member1.centrify.lab"
    fqdn = "member1.centrify.lab"
    computer_class = "Windows"
    session_type = "Rdp"
    local_account_automatic_maintenance = true
    local_account_manual_unlock = true
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
    //default_profile_id = centrifyvault_authenticationprofile.stepup_authprofile.id
    sets = [
        centrifyvault_manualset.all_systems.id,
        centrifyvault_manualset.critical_systems.id,
    ]
}

// Windows system member2
resource "centrifyvault_vaultsystem" "member2" {
    name = "member2.centrify.lab"
    fqdn = "member2.centrify.lab"
    computer_class = "Windows"
    session_type = "Rdp"
    local_account_automatic_maintenance = true
    local_account_manual_unlock = true
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
    sets = [
        centrifyvault_manualset.all_systems.id,
    ]
}


// CentOS system centos1
// local_account_automatic_maintenance needs to be false during creation
resource "centrifyvault_vaultsystem" "centos1" {
    name = "centos1.centrify.lab"
    fqdn = "centos1.centrify.lab"
    computer_class = "Unix"
    session_type = "Ssh"
    use_my_account = true
    //local_account_automatic_maintenance = true
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
    sets = [
        centrifyvault_manualset.all_systems.id,
    ]
}
