/*
resource "centrifyvault_manualset" "all_systems" {
    name = "${var.PREFIX}All Systems"
    type = "Server"
    description = "This Set contains all systems for ${var.PREFIX}."

    permission {
        principal_id = data.centrifyvault_role.lab_infra_admin.id
        principal_name = data.centrifyvault_role.lab_infra_admin.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View"
        ]
    }

    permission {
        principal_id = data.centrifyvault_role.sso_role.id
        principal_name = data.centrifyvault_role.sso_role.name
        principal_type = "Role"
        rights = [
          "Grant",
          "View",
          "Edit"
        ]
    }

    member_permission {
        principal_id = data.centrifyvault_role.lab_infra_admin.id
        principal_name = data.centrifyvault_role.lab_infra_admin.name
        principal_type = "Role"
        rights = [
          "Grant",
          "ManageSession"
        ]
    }
    
}

resource "centrifyvault_manualset" "all_databases" {
    name = "${var.PREFIX}All Databases"
    type = "VaultDatabase"
    description = "This Set contains all databases for ${var.PREFIX}."
}

resource "centrifyvault_manualset" "all_accounts" {
    name = "${var.PREFIX}All Accounts"
    type = "VaultAccount"
    description = "This Set contains all vaulted accounts for ${var.PREFIX}."
}

resource "centrifyvault_manualset" "tier1_accounts" {
    name = "${var.PREFIX}Accounts Tier-1"
    type = "VaultAccount"
    description = "This Set contains vaulted accounts doesn't require approval for ${var.PREFIX}."
}

resource "centrifyvault_manualset" "all_webapps" {
    name = "${var.PREFIX}All Web Applications"
    type = "Application"
    subtype = "Web"
    description = "This Set contains all web applications for ${var.PREFIX}."
    lifecycle {
    ignore_changes = [
      type,
      subtype,
    ]
  }
}


resource "centrifyvault_manualset" "all_desktopapps" {
    name = "${var.PREFIX}All Desktop Apps"
    type = "Application"
    subtype = "Desktop"
    description = "This Set contains all desktop applications for ${var.PREFIX}."
    lifecycle {
      ignore_changes = [
        type,
        subtype,
      ]
    }
}
*/