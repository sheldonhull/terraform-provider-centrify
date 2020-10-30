
// MS SQL database instance CENTRIFYSUITE
resource "centrifyvault_vaultdatabase" "centrifysuite" {
    name = "SQL-CENTRIFYSUITE"
    hostname = "member1.centrify.lab"
    database_class = "SQLServer"
    instance_name = "CENTRIFYSUITE"
    port = 1433

    sets = [
        centrifyvault_manualset.all_databases.id,
    ]
}

// demo_sa account for MS SQL
resource "centrifyvault_vaultaccount" "centrifysuite_demosa" {
    name = "demo_sa"
    credential_type = "Password"
    password = var.PASSWORD
    managed = true
    database_id = centrifyvault_vaultdatabase.centrifysuite.id
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.tier1_accounts.id
    ]
}

// sa account for MS SQL
resource "centrifyvault_vaultaccount" "centrifysuite_sa" {
    name = "sa"
    credential_type = "Password"
    password = var.PASSWORD
    managed = false
    database_id = centrifyvault_vaultdatabase.centrifysuite.id
    sets = [
        centrifyvault_manualset.all_accounts.id,
        centrifyvault_manualset.breakglass_accounts.id,
    ]
}
