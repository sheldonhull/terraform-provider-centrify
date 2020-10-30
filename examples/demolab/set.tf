
// All Systems Set
resource "centrifyvault_manualset" "all_systems" {
    name = "${var.PREFIX}Systems"
    type = "Server"
    description = "This Set contains all systems for ${var.PREFIX}."
    permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View"
        ]
    }
    member_permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View",
          "ManageSession",
          "RequestZoneRole"
        ]
    }
    
}

// Critical Systems Set
resource "centrifyvault_manualset" "critical_systems" {
    name = "${var.PREFIX}Critical Systems"
    type = "Server"
    description = "Critical systems that require step up authentication upon login."
}

// All Databases Set
resource "centrifyvault_manualset" "all_databases" {
    name = "${var.PREFIX}Databases"
    type = "VaultDatabase"
    description = "This Set contains all databases for ${var.PREFIX}."
    permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View"
        ]
    }
    member_permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View"
        ]
    }
}

// All Accounts Set
resource "centrifyvault_manualset" "all_accounts" {
    name = "${var.PREFIX}All Accounts"
    type = "VaultAccount"
    description = "This Set contains all vaulted accounts for ${var.PREFIX}."
    permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View"
        ]
    }
    member_permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View"
        ]
    }
}

// Tier-1 Accounts Set
resource "centrifyvault_manualset" "tier1_accounts" {
    name = "${var.PREFIX}Accounts Tier-1"
    type = "VaultAccount"
    description = "This Set contains vaulted accounts doesn't require approval for ${var.PREFIX}."
    member_permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View",
          "Checkout",
          "Login",
          "FileTransfer",
          "WorkspaceLogin"
        ]
    }
}

// Break Glass Accounts Set
resource "centrifyvault_manualset" "breakglass_accounts" {
    name = "${var.PREFIX}Break Glass Accounts"
    type = "VaultAccount"
    description = "Break glass accounts that require step up authentication upon password checkout."
}

// All Web Applications Set
resource "centrifyvault_manualset" "all_webapps" {
    name = "${var.PREFIX}Web Applications"
    type = "Application"
    subtype = "Web"
    description = "This Set contains all web applications for ${var.PREFIX}."
    permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View"
        ]
    }
    member_permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View",
          "Run"
        ]
    }
    lifecycle {
      ignore_changes = [
        type,
        subtype,
      ]
    }
}

// All Desktop Apps Set
resource "centrifyvault_manualset" "all_desktopapps" {
    name = "${var.PREFIX}Desktop Apps"
    type = "Application"
    subtype = "Desktop"
    description = "This Set contains all desktop applications for ${var.PREFIX}."
    permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View"
        ]
    }
    member_permission {
        principal_id = centrifyvault_role.infra_admin.id
        principal_name = centrifyvault_role.infra_admin.name
        principal_type = "Role"
        rights = [
          "View",
          "Run"
        ]
    }
    lifecycle {
      ignore_changes = [
        type,
        subtype,
      ]
    }
}
