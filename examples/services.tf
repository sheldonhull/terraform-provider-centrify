/*
data "centrifyvault_vaultaccount" "ad_admin" {
    name = "ad_admin"
    domain_id = centrifyvault_vaultdomain.centrify_lab.id
}

resource "centrifyvault_multiplexedaccount" "testwinservice" {
    name = "Account for TestWindowsService"
    description = "Multiplexed account for TestWindowsService"
    accounts = [
        centrifyvault_vaultaccount.svc_acct1.id,
        centrifyvault_vaultaccount.svc_acct2.id,
    ]
}

data "centrifyvault_manualset" "test" {
    type = "Subscriptions"
    name = "test"
}

resource "centrifyvault_service" "testservice" {
    service_name = "TestWindowsService"
    description = "Test Windows Service in member1"
    system_id = centrifyvault_vaultsystem.member1.id
    service_type = "WindowsService"
    enable_management = false
    admin_account_id = data.centrifyvault_vaultaccount.ad_admin.id
    multiplexed_account_id = centrifyvault_multiplexedaccount.testwinservice.id
    restart_service = false
    restart_time_restriction = false
    days_of_week = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
    restart_start_time = "09:00"
    restart_end_time = "10:00"
    use_utc_time = false
    permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "Delete",
        ]
    }
    sets = [
        data.centrifyvault_manualset.test.id,
    ]
}
*/