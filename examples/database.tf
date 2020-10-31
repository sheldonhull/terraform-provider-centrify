/*
resource "centrifyvault_vaultdatabase" "testdatabase1" {
    # Database -> Settings menu related settings
    name = "Database #1"
    hostname = "db1.example.com"
    database_class = "SQLServer"
    instance_name = "INSTANCE1"
    description = "Database system #1"
    port = 1433
    
    # Database -> Policy menu related settings
    //checkout_lifetime = 60

    # Database -> Advanced menu related settings
    //allow_multiple_checkouts = true
    //enable_password_rotation = true
    //password_rotate_interval = 60
    //enable_password_rotation_after_checkin = true
    //minimum_password_age = 90
    //password_profile_id = data.centrifyvault_passwordprofile.test_pw_pf1.id
    //enable_password_history_cleanup = true
    //password_historycleanup_duration = 100

	# System -> Connectors menu related settings
	//connector_list = [
    //    data.centrifyvault_connector.XXXXX-XXXXXX.id,
    //    data.centrifyvault_connector.dc01.id
    //]

    //sets = [
    //    centrifyvault_manualset.all_databases.id,
    //    data.centrifyvault_manualset.lab_databases.id
    //]

    permission {
        principal_id = data.centrifyvault_role.sso_role.id
        principal_name = data.centrifyvault_role.sso_role.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Edit",
        ]
    }

    lifecycle {
      ignore_changes = [
        database_class
        ]
    }
}
*/