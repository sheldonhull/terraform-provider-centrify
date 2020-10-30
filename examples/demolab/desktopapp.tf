
resource "centrifyvault_desktopapp" "centrify_pas" {
    template_name = "GenericDesktopApplication"
    name = "Centrify PAS Portal"
    description = "Web access for Centrify PAS portal"
    application_host_id = centrifyvault_vaultsystem.member2.id
    login_credential_type = "SharedAccount"
    application_account_id = centrifyvault_vaultaccount.helpdesk.id
    application_alias = "pas_desktopapp"
    
    command_line = "--ini=ini\\web_centrify_pas_webdriver.ini --url=https://pas.centrify.lab --username={login.Description} --password={login.SecretText}"
    command_parameter {
        name = "login"
        type = "DataVault"
        target_object_id = centrifyvault_vaultsecret.pas_admin_credential.id
    }
    
    sets = [
        centrifyvault_manualset.all_desktopapps.id
    ]
}

resource "centrifyvault_desktopapp" "ssms_demosa" {
    template_name = "Ssms"
    name = "SQL Server Management Studio (demo_sa)"
    description = "SQL Server Management Studio (SSMS) is an integrated environment for managing any SQL infrastructure."
    application_host_id = centrifyvault_vaultsystem.member2.id
    login_credential_type = "SharedAccount"
    application_account_id = centrifyvault_vaultaccount.helpdesk.id
    application_alias = "Ssms"
    
    command_line = "-S {database.FQDN}\\{database.InstanceName} -U {user.User} -P {user.Password}"
    command_parameter {
        name = "database"
        type = "VaultDatabase"
        target_object_id = centrifyvault_vaultdatabase.centrifysuite.id
    }
    command_parameter {
        name = "user"
        type = "VaultAccount"
        target_object_id = centrifyvault_vaultaccount.centrifysuite_demosa.id
    }
    
    sets = [
        centrifyvault_manualset.all_desktopapps.id
    ]
}
