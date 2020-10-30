/*
resource "centrifyvault_role" "infra_admin" {
    name = "${var.PREFIX}Infrastructure Admins"
    description = "Requester who can request access to all lab systems."
    adminrights = [
        "Privileged Access Service User",
    ]
    member {
        id = "8443cca3-a038-4f77-ae82-a2574e94f347"
        type = "Group"
    }
    member {
        id = "b9654e32-c10e-4e8e-9f7b-14955511d0e6"
        type = "User"
    }
    member {
        id = data.centrifyvault_user.admin.id
        type = "User"
    }
}

resource "centrifyvault_role" "infra_owner" {
    name = "${var.PREFIX}Infrastructure Owners"
    description = "Approver who can approve access to access lab systems."
    adminrights = [
        "Privileged Access Service User",
    ]
}

resource "centrifyvault_role" "mfa_machines_users" {
    name = "${var.PREFIX}MFA Machines & Users"
    description = "Machines and users who are enforced MFA for direct access without going through PAS."
    adminrights = [
        "Computer Login and Privilege Elevation",
    ]
}

resource "centrifyvault_role" "cloud_normal_user" {
    name = "${var.PREFIX}Cloud Normal Users"
    description = "AD accounts who can login to non-domain joined machines but without any privileges."
    adminrights = [
        "Computer Login and Privilege Elevation",
    ]
}

resource "centrifyvault_role" "cloud_local_admin" {
    name = "${var.PREFIX}Cloud Local Admins"
    description = "AD accounts that are granted local administrator access to non-domain joined machines."
    adminrights = [
        "Computer Login and Privilege Elevation",
    ]
}
*/

/*
resource "centrifyvault_role" "testrole1" {
    name = "test role1"
    description = "test role 1 changed"
    adminrights = [
        "Privileged Access Service User",
        "Application Management",
        "MFA Unlock"
    ]
}
*/