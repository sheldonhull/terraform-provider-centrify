
// Multiplex account for TestWindowsService
resource "centrifyvault_multiplexedaccount" "testwinservice" {
    name = "Account for TestWindowsService"
    description = "Multiplexed account for TestWindowsService"
    accounts = [
        centrifyvault_vaultaccount.svc_acct1.id,
        centrifyvault_vaultaccount.svc_acct2.id,
    ]
}

// Windows Service for TestWindowsService
resource "centrifyvault_service" "testservice" {
    service_name = "TestWindowsService"
    description = "Test Windows Service in member1"
    system_id = centrifyvault_vaultsystem.member1.id
    service_type = "WindowsService"
    enable_management = true
    admin_account_id = data.centrifyvault_vaultaccount.ad_admin.id
    multiplexed_account_id = centrifyvault_multiplexedaccount.testwinservice.id
}
